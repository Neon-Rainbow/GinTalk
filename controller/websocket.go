package controller

import (
	"GinTalk/websocket"
	"fmt"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"log"
)

func ServeWs(hub *websocket.Hub, c *gin.Context) {
	_userID, err := getCurrentUserID(c)
	if err != nil {
		zap.L().Error("Get user id from context failed")
		return
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
