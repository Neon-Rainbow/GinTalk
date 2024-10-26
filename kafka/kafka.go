package kafka

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"os/signal"
	"syscall"
	"time"

	"go.uber.org/zap"

	"github.com/segmentio/kafka-go"
)

// Kafka 配置
const (
	kafkaBrokerAddress = "localhost:9092"
	topic              = "test-topic"
)

type KafkaMessage struct {
	MessageType string `json:"message_type"`
	Message     any    `json:"message"`
}

const (
	MessageTypeAddPostVote     = "add_post_vote"
	MessageTypeSubPostVote     = "sub_post_vote"
	MessageTypeSavePostToRedis = "save_post_to_redis"
)

var _ KafkaInterface = (*Kafka)(nil)

type KafkaInterface interface {
	produceMessage(ctx context.Context, msg KafkaMessage) error
	consumeMessages(ctx context.Context) error
	SavePostSummaryToRedis(ctx context.Context, postSummary *DTO.PostSummary) error
	AddContentVote(ctx context.Context, postID int64) error
	SubContentVote(ctx context.Context, postID int64) error
	Close()
}

type Kafka struct {
	writer    *kafka.Writer
	reader    *kafka.Reader
	voteDao   dao.PostVoteDaoInterface
	voteCache cache.VoteCacheInterface
	postCache cache.PostCacheInterface
}

// InitKafka 初始化 Kafka
func InitKafka() (writer *kafka.Writer, reader *kafka.Reader) {
	return &kafka.Writer{
			Addr:     kafka.TCP(kafkaBrokerAddress),
			Topic:    topic,
			Balancer: &kafka.LeastBytes{},
		}, kafka.NewReader(kafka.ReaderConfig{
			Brokers:   []string{kafkaBrokerAddress},
			Topic:     topic,
			GroupID:   "my-group",
			Partition: 0,
			MinBytes:  10e3, // 10KB
			MaxBytes:  10e6, // 10MB
		})
}

// NewKafka 初始化 Kafka
func NewKafka(voteDao dao.PostVoteDaoInterface, voteCache cache.VoteCacheInterface, postCache cache.PostCacheInterface) KafkaInterface {
	w, r := InitKafka()

	k := &Kafka{
		writer:    w,
		reader:    r,
		voteDao:   voteDao,
		voteCache: voteCache,
		postCache: postCache,
	}
	go func() {
		err := k.consumeMessages(context.Background())
		if err != nil {
			zap.L().Error("消费消息失败", zap.Error(err))
		}
	}()
	return k
}

// ProduceMessage Kafka Producer：生产者函数
func (k *Kafka) produceMessage(ctx context.Context, message KafkaMessage) error {
	writer := k.writer

	msgBytes, err := json.Marshal(message)
	if err != nil {
		return err
	}
	// 发送消息到 Kafka
	err = writer.WriteMessages(ctx, kafka.Message{
		Key:   []byte(fmt.Sprintf("key-%d", time.Now().Unix())),
		Value: msgBytes,
	})
	if err != nil {
		zap.L().Error("发送消息失败", zap.Error(err))
		return err
	}
	zap.L().Info("发送消息成功, 消息内容: ", zap.Any("message", message))
	return nil
}

// ConsumeMessages Kafka Consumer：消费者函数
func (k *Kafka) consumeMessages(ctx context.Context) error {
	reader := k.reader

	zap.L().Info("开始消费消息")
	for {
		msg, err := reader.ReadMessage(ctx)
		zap.L().Info("读取消息", zap.Any("message", msg))
		if err != nil {
			zap.L().Error("读取消息失败", zap.Error(err))
			return err
		}
		var kafkaMsg KafkaMessage
		err = json.Unmarshal(msg.Value, &kafkaMsg)
		if err != nil {
			zap.L().Error("解析消息失败", zap.Error(err))
		}
		zap.L().Info("接收到消息", zap.Any("message", kafkaMsg))
		k.handleKafkaMessage(ctx, kafkaMsg)
	}
}

// 处理 Kafka 消息
func (k *Kafka) handleKafkaMessage(ctx context.Context, kafkaMsg KafkaMessage) {
	handleFunc, ok := kafkaHandleMap[kafkaMsg.MessageType]
	if !ok {
		zap.L().Error("未找到消息处理函数", zap.Any("message", kafkaMsg))
		return
	}

	zap.L().Info("开始处理消息", zap.Any("kafkaMsg", kafkaMsg))

	err := handleFunc(ctx, k, kafkaMsg)
	if err != nil {
		zap.L().Error("消息处理失败", zap.Error(err))
	}
}

func (k *Kafka) Close() {
	err := k.writer.Close()
	if err != nil {
		zap.L().Error("关闭 Kafka Writer 失败", zap.Error(err))
	}
	zap.L().Info("关闭 Kafka Writer 成功")
	err = k.reader.Close()
	if err != nil {
		zap.L().Error("关闭 Kafka Reader 失败", zap.Error(err))
	}
	zap.L().Info("关闭 Kafka Reader 成功")
}

// 捕获系统中断信号
func HandleInterrupt(ctx context.Context, k KafkaInterface) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c

	fmt.Println("\n收到中断信号，准备退出...")
	k.Close()
	os.Exit(0)
}

// ResetKafkaTopic 删除并重建 Topic
func ResetKafkaTopic() error {
	conn, err := kafka.Dial("tcp", kafkaBrokerAddress)
	if err != nil {
		return err
	}
	defer conn.Close()

	// 删除 Topic
	err = conn.DeleteTopics(topic)
	if err != nil {
		log.Printf("Failed to delete topic: %v", err)
		return err
	}
	zap.L().Info("Topic deleted successfully")

	// 重建 Topic
	topicConfig := kafka.TopicConfig{
		Topic:             topic,
		NumPartitions:     1,
		ReplicationFactor: 1,
	}
	err = conn.CreateTopics(topicConfig)
	if err != nil {
		zap.L().Error("Failed to create topic", zap.Error(err))
		return err
	}
	zap.L().Info("Topic created successfully")
	return nil
}

func (k *Kafka) SavePostSummaryToRedis(ctx context.Context, postSummary *DTO.PostSummary) error {
	msg := KafkaMessage{
		MessageType: MessageTypeSavePostToRedis,
		Message:     *postSummary,
	}
	return k.produceMessage(ctx, msg)
}

func (k *Kafka) AddContentVote(ctx context.Context, postID int64) error {
	msg := KafkaMessage{
		MessageType: MessageTypeAddPostVote,
		Message:     postID,
	}
	return k.produceMessage(ctx, msg)
}

func (k *Kafka) SubContentVote(ctx context.Context, postID int64) error {
	msg := KafkaMessage{
		MessageType: MessageTypeSubPostVote,
		Message:     postID,
	}
	return k.produceMessage(ctx, msg)
}
