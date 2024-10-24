package kafka

import (
	"GinTalk/dao"
	"context"
	"encoding/json"
	"fmt"
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
	MessageTypeAddPostVote = "add_post_vote"
	MessageTypeSubPostVote = "sub_post_vote"
)

type KafkaInterface interface {
	ProduceMessage(ctx context.Context, messageType string, message any) error
	ConsumeMessages(ctx context.Context) error
}

type Kafka struct {
	writer  *kafka.Writer
	reader  *kafka.Reader
	voteDao dao.PostVoteDaoInterface
	postDao dao.PostDaoInterface
}

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

func NewKafka(voteDao dao.PostVoteDaoInterface) KafkaInterface {
	w, r := InitKafka()
	k := &Kafka{
		writer:  w,
		reader:  r,
		voteDao: voteDao,
	}
	go func() {
		err := k.ConsumeMessages(context.Background())
		if err != nil {
			zap.L().Error("消费消息失败", zap.Error(err))
		}
	}()
	return k
}

// ProduceMessage Kafka Producer：生产者函数
func (k *Kafka) ProduceMessage(ctx context.Context, messageType string, message any) error {
	writer := kafka.Writer{
		Addr:     kafka.TCP(kafkaBrokerAddress),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
	}

	defer func(writer *kafka.Writer) {
		err := writer.Close()
		if err != nil {
			zap.L().Error("关闭 Kafka Writer 失败", zap.Error(err))
		}
	}(&writer)

	msg := KafkaMessage{
		MessageType: messageType,
		Message:     message,
	}
	msgBytes, err := json.Marshal(msg)
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
	zap.L().Info("发送消息成功, 消息内容: ", zap.Any("message", msg))
	return nil
}

// ConsumeMessages Kafka Consumer：消费者函数
func (k *Kafka) ConsumeMessages(ctx context.Context) error {
	reader := kafka.NewReader(kafka.ReaderConfig{
		Brokers:   []string{kafkaBrokerAddress},
		Topic:     topic,
		GroupID:   "my-group",
		Partition: 0,
		MinBytes:  10e3, // 10KB
		MaxBytes:  10e6, // 10MB
	})

	defer func(reader *kafka.Reader) {
		err := reader.Close()
		if err != nil {
			zap.L().Error("关闭 Kafka Reader 失败", zap.Error(err))
		}
	}(reader)

	zap.L().Info("开始消费消息")
	for {
		msg, err := reader.ReadMessage(ctx)
		if err != nil {
			zap.L().Error("读取消息失败", zap.Error(err))
			return err
		}
		var kafkaMsg KafkaMessage
		err = json.Unmarshal(msg.Value, &kafkaMsg)
		if err != nil {
			zap.L().Error("解析消息失败", zap.Error(err))
			return err
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

// 捕获系统中断信号
func handleInterrupt(ctx context.Context) {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, syscall.SIGTERM)
	<-c
	fmt.Println("\n收到中断信号，准备退出...")
	cancelFunc, ok := ctx.Value("cancel").(context.CancelFunc)
	if ok {
		cancelFunc()
	}
	os.Exit(0)
}
