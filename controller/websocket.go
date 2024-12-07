package controller

import (
	"GinTalk/websocket"
	"fmt"
	"log"

	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// WebsocketHandle 处理给定 Gin 上下文的 WebSocket 连接。
// 它从上下文中检索当前用户 ID，将 HTTP 连接升级为 WebSocket 连接，
// 创建一个新的客户端实例，将客户端注册到 hub，并启动客户端的读写泵。
//
// 参数:
//   - c: 请求的 Gin 上下文。
//
// 该函数执行以下步骤:
//  1. 检索 WebSocket hub 实例。
//  2. 从上下文中获取当前用户 ID。
//  3. 使用用户 ID 记录 WebSocket 连接。
//  4. 将 HTTP 连接升级为 WebSocket 连接。
//  5. 创建一个新的 WebSocket 客户端实例。
//  6. 将新客户端注册到 hub。
//  7. 启动客户端的读写泵。
func WebsocketHandle(c *gin.Context) {
	hub := websocket.GetHub()
	_userID, exist := getCurrentUserID(c)
	if !exist {
		zap.L().Error("获取用户 ID 失败")
	}

	userID := fmt.Sprintf("%v", _userID)
	zap.L().Info("WebSocket connected", zap.String("user_id", userID))
	// 升级HTTP到WebSocket协议
	upgrader := websocket.GetUpgrader()
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Printf("WebSocket upgrade error: %v", err)
		return
	}

	// 创建客户端实例
	client := &websocket.Client{
		User: userID,
		Hub:  hub,
		Conn: conn,
		Send: make(chan websocket.Message, 256),
	}

	// 注册新客户端
	client.Hub.NewClients <- client

	// 启动消息读写协程
	go client.WritePump()
	go client.ReadPump()
}
