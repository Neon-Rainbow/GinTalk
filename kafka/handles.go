package kafka

import (
	"GinTalk/DTO"
	"GinTalk/cache"
	"GinTalk/dao"
	"GinTalk/model"
	"GinTalk/websocket"
	"context"
	"encoding/json"
	"strconv"

	"github.com/segmentio/kafka-go"
	"go.uber.org/zap"
)

// handleLikeMessage 处理点赞消息
func handleLikeMessage(msg kafka.Message) {
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
	err = dao.AddPostVoteWithTx(context.Background(), postID, userID, voteMsg.Vote)
	if err != nil {
		zap.L().Error("添加投票记录失败", zap.Error(err))
		return
	}

	// 更新 Redis 热度
	oldUp, err := dao.GetPostVoteCount(context.Background(), postID)
	if err != nil {
		zap.L().Error("获取帖子投票数失败", zap.Error(err))
		return
	}
	err = cache.AddPostHot(context.Background(), postID, int(oldUp.Vote), int(oldUp.Vote)+voteMsg.Vote)
	if err != nil {
		zap.L().Error("更新 Redis 热度失败", zap.Error(err))
		return
	}

	zap.L().Info("更新 Redis 热度成功", zap.Int64("post_id", postID), zap.Int("vote", voteMsg.Vote))

	// 如果是取消点赞，不发送通知
	if voteMsg.Vote == 0 {
		return
	}

	// 如果是点赞,发送通知
	notificationMsg := websocket.Message{
		Kind: websocket.MessageKindNotificationVote,
		From: strconv.FormatInt(userID, 10),
		To:   strconv.FormatInt(postID, 10),
	}

	// 获取帖子作者
	post, err := dao.GetPostDetail(context.Background(), postID)
	if err != nil {
		zap.L().Error("获取帖子详情失败", zap.Error(err))
		return
	}
	userID = post.AuthorId

	err = websocket.GetHub().SendToUser(notificationMsg)
	if err != nil {
		zap.L().Error("发送通知失败", zap.Error(err))
		return
	}

	zap.L().Info("发送通知成功", zap.Int64("post_id", postID), zap.Int64("user_id", userID))

	return
}

// handleCommentMessage 处理评论消息
func handleCommentMessage(msg kafka.Message) {
	var commentMsg DTO.CommentDetail
	if err := json.Unmarshal(msg.Value, &commentMsg); err != nil {
		zap.L().Error("序列化消息失败", zap.Error(err))
		return
	}
	commentModel := model.Comment{
		CommentID:  commentMsg.Comment.CommentID,
		PostID:     commentMsg.Comment.PostID,
		AuthorID:   commentMsg.Comment.AuthorID,
		AuthorName: commentMsg.Comment.AuthorName,
		Content:    commentMsg.Comment.Content,
	}
	err := dao.CreateComment(context.Background(), &commentModel, commentMsg.CommentRelation.ReplyID, commentMsg.CommentRelation.ParentID)
	if err != nil {
		zap.L().Error("保存评论到数据库失败", zap.Error(err))
		return
	}
	zap.L().Info("保存评论成功", zap.Int64("comment_id", commentMsg.Comment.CommentID))

	parentComment, err := dao.GetCommentByID(context.Background(), commentMsg.CommentRelation.ReplyID)
	if err != nil {
		zap.L().Error("获取父评论失败", zap.Error(err))
		return
	}
	if parentComment.AuthorID == commentMsg.Comment.AuthorID {
		// 评论自己的帖子，不发送通知
		return
	}

	// 发送通知
	notificationMsg := websocket.Message{
		Kind: websocket.MessageKindNotificationComment,
		From: strconv.FormatInt(commentMsg.Comment.AuthorID, 10),
		To:   strconv.FormatInt(parentComment.AuthorID, 10),
	}
	err = websocket.GetHub().SendToUser(notificationMsg)
	if err != nil {
		zap.L().Error("发送通知失败", zap.Error(err))
		return
	}

	return
}

// handleCreatePostMessage 处理帖子消息
func handleCreatePostMessage(msg kafka.Message) {
	var postMsg DTO.PostDetail
	if err := json.Unmarshal(msg.Value, &postMsg); err != nil {
		zap.L().Error("序列化消息失败", zap.Error(err))
		return
	}

	// 保存帖子到数据库
	err := dao.CreatePost(context.Background(), &postMsg)
	if err != nil {
		zap.L().Error("保存帖子到数据库失败", zap.Error(err))
		return
	}
	// 保存帖子到 Redis
	err = cache.SavePost(context.Background(), postMsg.ConvertToSummary())
	if err != nil {
		zap.L().Error("保存帖子到 Redis 失败", zap.Error(err))
		return
	}
	zap.L().Info("保存帖子成功", zap.Int64("post_id", postMsg.PostID))
}
