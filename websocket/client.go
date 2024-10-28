package websocket

import (
	"encoding/json"
	"go.uber.org/zap"
	"log"
	"time"

	"github.com/gorilla/websocket"
)

const (
	// 写超时, 10s
	writeWait = 10 * time.Second
	// 读 pong 帧超时, 60s
	pongWait = 60 * time.Second
	// 写 ping 帧周期(必须小于 pongWait), 54s
	pingPeriod = (pongWait * 9) / 10
)

// http/1.1 -> websocket 协议升级
var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
}

func GetUpgrader() websocket.Upgrader {
	return upgrader
}

// Client 是 websocket 和 Hub 的中间人
type Client struct {
	// 用户
	User string
	Hub  *Hub
	// websocket 连接
	Conn *websocket.Conn
	// 发送消息通道
	Send chan Message
}

// ReadPump 从 websocket 连接读取消息并写入 Client.Send 通道
func (c *Client) ReadPump() {
	defer func() {
		c.Hub.CloseClients <- c
		c.Conn.Close()
	}()

	_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
	c.Conn.SetPongHandler(func(string) error {
		_ = c.Conn.SetReadDeadline(time.Now().Add(pongWait))
		return nil
	})

	for {
		_, bytes, err := c.Conn.ReadMessage()
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Printf("Unexpected close error: %v", err)
			}
			break
		}

		var message Message
		if err := json.Unmarshal(bytes, &message); err != nil {
			log.Printf("JSON unmarshal error: %v", err)
			continue
		}

		c.Hub.Message <- message
	}
}

// WritePump 从 Client.Send 通道读取消息并写入 websocket 连接
func (c *Client) WritePump() {
	ticker := time.NewTicker(pingPeriod)
	defer func() {
		ticker.Stop()
		err := c.Conn.Close()
		if err != nil {
			zap.L().Error("Close websocket connection error", zap.Error(err))
		}
	}()

	for {
		select {
		case message, ok := <-c.Send:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if !ok {
				err := c.Conn.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					zap.L().Error("Write close message error", zap.Error(err))
				}
				return
			}

			bytes, err := json.Marshal(&message)
			if err != nil {
				log.Printf("JSON marshal error: %v", err)
				continue
			}

			if err := c.Conn.WriteMessage(websocket.TextMessage, bytes); err != nil {
				log.Printf("WebSocket write error: %v", err)
				return
			}

		case <-ticker.C:
			c.Conn.SetWriteDeadline(time.Now().Add(writeWait))
			if err := c.Conn.WriteMessage(websocket.PingMessage, nil); err != nil {
				return
			}
		}
	}
}

// SendToUser 向特定用户发送消息
func (h *Hub) SendToUser(message Message) error {
	h.mu.RLock()
	userID := message.To
	client, exists := h.Clients[userID]
	h.mu.RUnlock()

	if !exists {
		log.Printf("用户 %s 不在线", userID)
		return nil
	}

	client.Send <- message // 将消息发送到客户端的Send通道
	return nil
}
