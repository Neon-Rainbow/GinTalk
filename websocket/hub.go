package websocket

import (
	"go.uber.org/zap"
	"sync"
)

// Hub 负责保存活跃客户端并进行消息转发和存储
type Hub struct {
	Message      chan Message         // 消息通道
	NewClients   chan *Client         // 新连接的客户端
	CloseClients chan *Client         // 断开连接的客户端
	Clients      map[string]*Client   // 在线客户端
	Messages     map[string][]Message // 离线消息
	mu           sync.RWMutex         // 保护Clients和Messages的互斥锁
}

var hub = newHub()

// newHub 创建一个新的Hub实例
func newHub() *Hub {
	return &Hub{
		NewClients:   make(chan *Client),
		CloseClients: make(chan *Client),
		Message:      make(chan Message),
		Clients:      make(map[string]*Client),
		Messages:     make(map[string][]Message),
	}
}

func GetHub() *Hub {
	return hub
}

// run 启动Hub，处理连接、断开和消息
func (h *Hub) run() {
	for {
		select {
		case client := <-h.NewClients:
			h.mu.Lock()

			// 注册新客户端
			h.Clients[client.User] = client

			h.mu.Unlock()

			// 广播其他用户的在线状态给新用户
			//h.mu.RLock()
			//for k, v := range h.Clients {
			//	if k != client.User {
			//		client.Send <- Message{Kind: "online", From: k}
			//		v.Send <- Message{Kind: "online", From: client.User}
			//	}
			//}
			//h.mu.RUnlock()

			// 下发离线消息
			h.mu.Lock()
			if messages, ok := h.Messages[client.User]; ok {
				for _, msg := range messages {
					client.Send <- msg
				}
				delete(h.Messages, client.User) // 清除已发送的离线消息
			}
			h.mu.Unlock()

			zap.L().Info("Client connected", zap.String("user", client.User))

		case client := <-h.CloseClients:
			h.mu.Lock()
			close(client.Send)             // 关闭通道
			delete(h.Clients, client.User) // 从Clients中移除
			h.mu.Unlock()

			// 广播下线状态
			//h.mu.RLock()
			//for k, v := range h.Clients {
			//	if k != client.User {
			//		v.Send <- Message{Kind: "offline", From: client.User}
			//	}
			//}
			//h.mu.RUnlock()

			zap.L().Info("Client disconnected", zap.String("user", client.User))

		case msg := <-h.Message:
			h.mu.RLock()
			client, ok := h.Clients[msg.To]
			h.mu.RUnlock()

			if ok {
				// 用户在线，直接转发消息
				client.Send <- msg
			} else {
				// 用户不在线，存储离线消息
				h.mu.Lock()
				h.Messages[msg.To] = append(h.Messages[msg.To], msg)
				h.mu.Unlock()
			}
		}
	}
}
