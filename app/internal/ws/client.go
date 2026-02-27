package ws

import (
	"commmunity/app/internal/db/global"
	"commmunity/app/internal/model"
	"commmunity/app/zlog"
	"encoding/json"
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"go.uber.org/zap"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  1024,
	WriteBufferSize: 1024,
	CheckOrigin: func(r *http.Request) bool {
		return true
	},
}

func HandleWebSocket(c *gin.Context) {
	conn, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		zlog.Error("WebSocket 升级失败", zap.Error(err))
		return
	}
	userId := c.MustGet("userId").(uint)
	client := &Client{
		Manager: &GlobalManager,
		UserId:  userId,
		Socket:  conn,
		Send:    make(chan []byte, 256),
	}
	GlobalManager.Register <- client
	go client.ReadMessage()
	go client.WriteMessage()
}

func (client *Client) ReadMessage() {
	for {
		_, message, err := client.Socket.ReadMessage()
		if err != nil {
			zlog.Warn("接收信息异常", zap.Error(err))
			break
		}
		zlog.Info("成功接收信息")
		var MessageRequest model.MessageRequest
		err = json.Unmarshal(message, &MessageRequest)
		if err != nil {
			zlog.Error("json序列化失败", zap.Error(err))
			continue
		}
		global.Message.SaveMessage(client.UserId, MessageRequest.ToUserID, MessageRequest.Content, MessageRequest.Type)
		_ = global.MessageRedis.DelMessageCache(client.UserId, MessageRequest.ToUserID)
		res := Response{
			Code: Chat,
			Data: ChatData{
				FromId:    client.UserId,
				ToId:      MessageRequest.ToUserID,
				Content:   MessageRequest.Content,
				Type:      MessageRequest.Type,
				CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
				IsMine:    false,
			},
		}
		GlobalManager.SendToUser(MessageRequest.ToUserID, res)
		myRes := res
		tem := res.Data.(ChatData)
		tem.IsMine = true
		myRes.Data = tem
		GlobalManager.SendToUser(client.UserId, myRes)
	}
	defer func() {
		client.Manager.Unregister <- client
		client.Socket.Close()
	}()
}

func (client *Client) WriteMessage() {
	defer client.Socket.Close()
	for {
		message, ok := <-client.Send
		if !ok {
			return
		}
		err := client.Socket.WriteMessage(websocket.TextMessage, message)
		if err != nil {
			zlog.Error("写入信息异常", zap.Error(err))
			return
		}
	}
}

func SendNotice(userId uint, tp int, senderId uint, postId uint, content string) {
	global.Message.SaveNotice(userId, senderId, tp, content, postId)
	data := NoticeData{
		Type:      tp,
		SenderId:  senderId,
		Content:   content,
		PostId:    postId,
		CreatedAt: time.Now().Format("2006-01-02 15:04:05"),
	}
	response := Response{
		Code: Notification,
		Data: data,
	}
	GlobalManager.SendToUser(userId, response)
}
