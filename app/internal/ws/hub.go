package ws

import (
	"commmunity/app/zlog"
	"encoding/json"
	"sync"

	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

type Client struct {
	Manager *Manager
	UserID  uint
	Socket  *websocket.Conn
	Send    chan []byte
}

type Manager struct {
	Clients    map[uint]*Client
	Register   chan *Client
	Unregister chan *Client
	Broadcast  chan []byte
	Lock       sync.RWMutex
}

var GlobalManager = Manager{
	Clients:    make(map[uint]*Client),
	Register:   make(chan *Client),
	Unregister: make(chan *Client),
	Broadcast:  make(chan []byte),
}

func (manager *Manager) Start() {
	zlog.Info("WebSocket 管理器启动...")
	for {
		select {
		case client := <-manager.Register:
			manager.Lock.Lock()
			manager.Clients[client.UserID] = client
			zlog.Info("用户上线", zap.Uint("user_id", client.UserID))
			manager.Lock.Unlock()
		case client := <-manager.Unregister:
			manager.Lock.Lock()
			if _, ok := manager.Clients[client.UserID]; ok {
				delete(manager.Clients, client.UserID)
				close(client.Send)
				zlog.Info("用户下线", zap.Uint("user_id", client.UserID))
			}
			manager.Lock.Unlock()
		}
	}
}

func (manager *Manager) SendToUser(userID uint, message interface{}) {
	manager.Lock.RLock()
	client, ok := manager.Clients[userID]
	manager.Lock.RUnlock()
	if ok {
		jsonMessage, _ := json.Marshal(message)
		select {
		case client.Send <- jsonMessage:
		default:
			close(client.Send)
			delete(manager.Clients, userID)
		}
	} else {
		zlog.Warn("用户不在线，消息未发送")
		//存库
	}
}
