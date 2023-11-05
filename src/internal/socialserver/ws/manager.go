package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/marmotedu/iam/pkg/log"
	whatsappbase "go-socialapp/internal/pkg/third-party/whatsapp"
)

var Manager = ClientManager{
	clients:    make(map[string]*Client), // 参与连接的用户，出于性能的考虑，需要设置最大连接数
	broadcast:  whatsappbase.Broadcast,
	register:   make(chan *Client),
	unregister: make(chan *Client),
}

// 用户管理
type ClientManager struct {
	clients    map[string]*Client
	broadcast  chan whatsappbase.BroadcastMessage
	register   chan *Client
	unregister chan *Client
}

func (manager *ClientManager) Register(client *Client) {
	manager.register <- client
}

func (manager *ClientManager) unRegister(client *Client) {
	manager.unregister <- client
}

func (manager *ClientManager) broadcastMsg(msg whatsappbase.BroadcastMessage) {
	manager.broadcast <- msg
}

func (manager *ClientManager) Start() {
	for {
		select {
		case conn := <-manager.register: // 建立连接
			log.Infof("建立新连接: %v", conn.id)
			Manager.clients[conn.id] = conn
			replyMsg := &whatsappbase.BroadcastMessage{
				Code: "CONNECT_SUCCESS",
			}
			msg, err := json.Marshal(replyMsg)
			if err != nil {
				log.Errorf("Marshal error:", err)
				continue
			}
			err = conn.socket.WriteMessage(websocket.TextMessage, msg)
			if err != nil {
				err = conn.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Errorf("write empty message error: %v", err)
					manager.unRegister(conn)
				}
			}
		case conn := <-Manager.unregister: // 断开连接
			if conn.isClose {
				continue
			}
			log.Infof("----连接失败:%v", conn.id)
			if _, ok := Manager.clients[conn.id]; ok {
				replyMsg := &whatsappbase.BroadcastMessage{
					Code: "DISCONNECT",
				}
				msg, _ := json.Marshal(replyMsg)
				_ = conn.socket.WriteMessage(websocket.TextMessage, msg)
				close(conn.send)
				conn.socket.Close()
				delete(Manager.clients, conn.id)
				conn.isClose = true
			}
		//广播信息
		case message := <-Manager.broadcast:
			log.Infof("message received: %s", message)
			if message.Code == "LOGIN_SUCCESS" {
				continue
			}
			// Send the message to all clients
			for _, connection := range manager.clients {
				connection.sendMsg(message)
			}
		}
	}
}
