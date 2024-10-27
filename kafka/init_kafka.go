package kafka

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"context"
	"encoding/json"
	"log"
	"net"
	"os"
	"os/signal"
	"strconv"
	"sync"
	"syscall"
	"time"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

const (
	TopicComments      = "comments"
	TopicNotifications = "notifications"
	TopicVote          = "vote"
	TopicPost          = "post"
)

// InitProducers 初始化所有生产者
func InitProducers() map[string]*kafka.Writer {
	brokers := []string{"localhost:9092"}
	return map[string]*kafka.Writer{
		TopicComments:      initKafkaWriter(brokers, "comments"),
		TopicNotifications: initKafkaWriter(brokers, "notifications"),
		TopicVote:          initKafkaWriter(brokers, "vote"),
		TopicPost:          initKafkaWriter(brokers, "post"),
	}
}

// InitConsumers 启动多个消费者
func InitConsumers(brokers []string, topics []string) {
	var wg sync.WaitGroup

	for _, topic := range topics {
		wg.Add(1)
		go StartSingleConsumer(topic, brokers, &wg)
	}

	// 捕获系统中断信号，优雅退出
	sigChan := make(chan os.Signal, 1)
	signal.Notify(sigChan, syscall.SIGINT, syscall.SIGTERM)
	<-sigChan

	log.Println("接收到中断信号，正在关闭消费者...")
	wg.Wait() // 等待所有消费者退出
	log.Println("所有消费者已停止")
}

// initKafkaWriter 初始化 Kafka 生产者
func initKafkaWriter(brokers []string, topic string) *kafka.Writer {
	return &kafka.Writer{
		Addr:     kafka.TCP(brokers...),
		Topic:    topic,
		Balancer: &kafka.LeastBytes{},
		Transport: &kafka.Transport{
			Dial: (&net.Dialer{
				Timeout: 10 * time.Second,
			}).DialContext,
		},
	}
}

// initKafkaConsumer 初始化 Kafka 消费者
func initKafkaConsumer(brokers []string, topic, groupID string) *kafka.Reader {
	return kafka.NewReader(kafka.ReaderConfig{
		Brokers: brokers,
		Topic:   topic,
		GroupID: groupID,
	})
}

// StartSingleConsumer 启动单个消费者
func StartSingleConsumer(topic string, brokers []string, wg *sync.WaitGroup) {
	defer wg.Done() // 确保消费者退出时减少计数

	reader := initKafkaConsumer(brokers, topic, topic+"_group")
	defer reader.Close() // 关闭消费i

	log.Printf("启动消费者，监听 topic: %s", topic)
	for {
		msg, err := reader.ReadMessage(context.Background())
		if err != nil {
			log.Printf("读取消息失败: %v", err)
			continue
		}
		// 异步处理消息，避免阻塞
		go handleMessage(msg)
	}
}

type MessageHandle interface {
	HandleMessage(msg kafka.Message)
}

var handles map[string]MessageHandle

var _ MessageHandle = (*VotePostHandle)(nil)

type VotePostHandle struct {
	voteDao   dao.PostVoteDaoInterface
	voteCache cache.VoteCacheInterface
}

func NewVotePostHandle(voteDao dao.PostVoteDaoInterface, voteCache cache.VoteCacheInterface) MessageHandle {
	return &VotePostHandle{
		voteDao:   voteDao,
		voteCache: voteCache,
	}
}

var _ MessageHandle = (*PostHandle)(nil)

type PostHandle struct {
	postDao   dao.PostDaoInterface
	postCache cache.PostCacheInterface
}

func NewPostHandle(postDao dao.PostDaoInterface, postCache cache.PostCacheInterface) MessageHandle {
	return &PostHandle{
		postDao:   postDao,
		postCache: postCache,
	}
}

// HandleMessage 处理帖子消息
// 帖子消息格式：{"post_id": "1", "author_id": "1", "community_id": "1", "title": "标题", "content": "内容"}
// 处理消息时，需要将帖子保存到数据库和 Redis 中
// 保存到数据库时，需要生成帖子摘要
// 该接口可以将发布帖子的过程通过消息队列实现异步化
func (p *PostHandle) HandleMessage(msg kafka.Message) {
	var postMsg DTO.PostDetail
	if err := json.Unmarshal(msg.Value, &postMsg); err != nil {
		zap.L().Error("序列化消息失败", zap.Error(err))
		return
	}

	summary := postMsg.GenerateSummary()

	// 保存帖子到数据库
	err := p.postDao.CreatePost(context.Background(), &postMsg, summary)
	if err != nil {
		zap.L().Error("保存帖子到数据库失败", zap.Error(err))
		return
	}
	// 保存帖子到 Redis
	err = p.postCache.SavePostToRedis(context.Background(), postMsg.ConvertToSummary())
	if err != nil {
		zap.L().Error("保存帖子到 Redis 失败", zap.Error(err))
		return
	}
	zap.L().Info("保存帖子成功", zap.Int64("post_id", postMsg.PostID))
}

// HandleMessage 处理投票帖子消息
// 投票消息格式：{"post_id": "1", "user_id": "1", "vote": 1}
// 投票数为 1 时表示赞成，为 -1 时表示反对,为 0 时表示取消投票
// 投票消息会增加投票记录,更新帖子的投票数,并且重新计算 Redis 热度
func (v *VotePostHandle) HandleMessage(msg kafka.Message) {
	// 处理消息
	var voteMsg Vote
	if err := json.Unmarshal(msg.Value, &voteMsg); err != nil {
		zap.L().Error("序列化消息失败", zap.Error(err))
		return
	}
	postID, err := strconv.ParseInt(voteMsg.PostID, 10, 64)
	if err != nil {
		zap.L().Error("转换 post id 失败", zap.Error(err))
		return
	}
	userID, err := strconv.ParseInt(voteMsg.UserID, 10, 64)
	if err != nil {
		zap.L().Error("转换 user id 失败", zap.Error(err))
		return
	}
	// 向数据库中添加投票记录和更新投票数
	err = v.voteDao.AddPostVoteWithTx(context.Background(), postID, userID, voteMsg.Vote)
	if err != nil {
		zap.L().Error("添加投票记录失败", zap.Error(err))
		return
	}

	// 更新 Redis 热度
	oldUp, err := v.voteDao.GetPostVoteCount(context.Background(), postID)
	if err != nil {
		zap.L().Error("获取帖子投票数失败", zap.Error(err))
		return
	}
	err = v.voteCache.AddPostHot(context.Background(), postID, int(oldUp.Vote), int(oldUp.Vote)+voteMsg.Vote)
	if err != nil {
		zap.L().Error("更新 Redis 热度失败", zap.Error(err))
		return
	}
	zap.L().Info("更新 Redis 热度成功", zap.Int64("post_id", postID), zap.Int("vote", voteMsg.Vote))
	return
}

// 消息处理逻辑
func handleMessage(msg kafka.Message) {
	zap.L().Info("接收到消息", zap.String("key", string(msg.Key)), zap.String("value", string(msg.Value)))
	handles[msg.Topic].HandleMessage(msg)
}
