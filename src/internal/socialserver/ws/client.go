package ws

import (
	"context"
	"encoding/json"
	"github.com/davecgh/go-spew/spew"
	"github.com/gorilla/websocket"
	"github.com/marmotedu/iam/pkg/log"
	whatsappbase "go-socialapp/internal/pkg/third-party/whatsapp"
	"go-socialapp/internal/socialserver/client/whatsapp"
)

type Client struct {
	id        string
	socket    *websocket.Conn
	send      chan []byte
	waService whatsapp.Factory
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		id:        id,
		socket:    conn,
		send:      make(chan []byte),
		waService: whatsapp.Client(),
	}
}

func (c *Client) Read() {
	defer func() { // 避免忘记关闭，所以要加上close
		Manager.unRegister(c)
		_ = c.socket.Close()
	}()
	for {
		c.socket.PongHandler()

		messageType, message, err := c.socket.ReadMessage()
		//sendMsg := new(SendMsg)
		//err := c.Socket.ReadJSON(&sendMsg) // 读取json格式，如果不是json格式，会报错
		//if err != nil {
		//	log.Println("数据格式不正确", err)
		//	Manager.Unregister <- c
		//	_ = c.Socket.Close()
		//	break
		//}
		if err != nil {
			if websocket.IsUnexpectedCloseError(err, websocket.CloseGoingAway, websocket.CloseAbnormalClosure) {
				log.Errorf("read error: %s", err.Error())
			}
			return // Calls the deferred function, i.e. closes the connection on error
		}

		if messageType == websocket.TextMessage {
			// Broadcast the received message
			var messageData whatsappbase.BroadcastMessage
			err := json.Unmarshal(message, &messageData)
			if err != nil {
				log.Errorf("error unmarshal message: %s", err.Error())
				return
			}
			if messageData.Code == "FETCH_DEVICES" {
				devices, err := c.waService.App().FetchDevices(context.Background())
				if err != nil {
					log.Errorf("FETCH_DEVICES err: %s", err.Error())
				}
				bmsg := whatsappbase.BroadcastMessage{
					Code:   "LIST_DEVICES",
					Result: devices,
				}
				Manager.broadcastMsg(bmsg)
			}

			if messageData.Code == "QRCODE" {
				log.Infof("收到请求 qrcode")
				ch, err := c.waService.App().GetQrCode(context.Background())
				if err != nil {
					log.Errorf("QRCODE err: %s", err.Error())
					return
				}
				go func() {
					log.Infof("遍历获取到的 qrcode")
					for evt := range ch {
						spew.Dump(evt)
						if evt.Event == "code" {
							replyMsg := whatsappbase.BroadcastMessage{
								Code:   "QRCODE",
								Result: evt.Code,
							}
							msg, _ := json.Marshal(replyMsg)
							_ = c.socket.WriteMessage(websocket.TextMessage, msg)

						} else {
							log.Errorf("error when get qrCode ,%v", evt.Event)
						}
					}
				}()

			}

		} else {
			log.Errorf("websocket message received of type %s", messageType)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		Manager.unRegister(c)
		_ = c.socket.Close()
	}()
	for {
		select {
		//case message, ok := <-c.send:
		//	if !ok {
		//		_ = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
		//		return
		//	}
		//	log.Infof(c.id, "接受消息: %s", string(message))
		//	replyMsg := wsMessage{
		//		Code:    e.WebsocketSuccessMessage,
		//		Content: fmt.Sprintf("%s", string(message)),
		//	}
		//	msg, _ := json.Marshal(replyMsg)
		//	_ = c.socket.WriteMessage(websocket.TextMessage, msg)
		}
	}
}
