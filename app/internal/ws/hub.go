package ws

import (
	"commmunity/app/zlog"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

const (
	Chat         = 1
	Notification = 2
)

type Response struct {
	Code int         `json:"code"`
	Data interface{} `json:"data"`
}

type ChatData struct {
	FromId    uint   `json:"from_id"`
	ToId      uint   `json:"to_id"`
	Content   string `json:"content"`
	Type      int    `json:"type"`
	CreatedAt string `json:"created_at"`
	IsMine    bool   `json:"is_mine"`
}

type NoticeData struct {
	Type      int    `json:"type"` // 1点赞 2评论 3系统
	SenderId  uint   `json:"sender_id"`
	Content   string `json:"content"`
	PostId    uint   `json:"post_id"`
	CreatedAt string `json:"created_at"`
}

type Client struct {
	Manager *Manager
	UserId  uint
	Socket  *websocket.Conn
	Send    chan []byte
}

type Manager struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Lock       sync.RWMutex
}

var GlobalManager = Manager{
	Clients:    make(map[uint]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
}

func (manager *Manager) Start() {
	zlog.Info("WebSocket 管理器启动...")
	for {
		select {
		case client := <-manager.Register:
			manager.Lock.Lock()
			manager.Clients[client.UserId] = client
			zlog.Info("用户上线", zap.Any("client", client.UserId))
			manager.Lock.Unlock()
		case client := <-manager.Unregister:
			manager.Lock.Lock()
			if _, ok := manager.Clients[client.UserId]; ok {
				delete(manager.Clients, client.UserId)
				close(client.Send)
				zlog.Info("用户下线", zap.Any("client", client.UserId))
			}
			manager.Lock.Unlock()
		}
	}
}

func (manager *Manager) SendToUser(userId uint, message interface{}) {
	manager.Lock.RLock()
	client, ok := manager.Clients[userId]
	manager.Lock.RUnlock()
	if ok {
		jsonMessage, _ := json.Marshal(message)
		select {
		case client.Send <- jsonMessage:
		default:
			close(client.Send)
			delete(manager.Clients, userId)
		}
	} else {
		zlog.Info("用户不在线，消息未发送")
	}
}
