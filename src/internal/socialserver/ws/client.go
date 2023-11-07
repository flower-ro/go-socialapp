package ws

import (
	"encoding/json"
	"github.com/gorilla/websocket"
	"github.com/marmotedu/iam/pkg/log"
	"go-socialapp/internal/socialserver/cache/notlogin"
)

type Client struct {
	id      string
	socket  *websocket.Conn
	send    chan Message
	isClose bool
}

func NewClient(id string, conn *websocket.Conn) *Client {
	return &Client{
		id:     id,
		socket: conn,
		send:   make(chan Message, 10),
	}
}

func (c *Client) sendMsg(replyMsg Message) {
	c.send <- replyMsg
}

func (c *Client) Read() {
	defer func() { // 避免忘记关闭，所以要加上close
		if err := recover(); err != nil {
			//	log.Errorf("%v", err)
		}
		//log.Infof("来自 read unRegister")
		Manager.unRegister(c)
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
			var messageData Message
			err := json.Unmarshal(message, &messageData)
			if err != nil {
				log.Errorf("error unmarshal message: %s", err.Error())
				continue
			}
			//if messageData.Code == "FETCH_DEVICES" {
			//	devices, err := notlogin.TmpWaClientCache.GetQrCodeByNewWaClient()
			//	if err != nil {
			//		log.Errorf("FETCH_DEVICES err: %s", err.Error())
			//	}
			//	bmsg := BroadcastMessage{
			//		Code:   "LIST_DEVICES",
			//		Result: devices,
			//	}
			//	Manager.broadcastMsg(bmsg)
			//}

			if messageData.Code == "QRCODE" {
				//log.Infof("收到请求 qrcode")
				ch, err := notlogin.TmpWaClientCache.GetQrCodeByNewWaClient()
				if err != nil {
					log.Errorf("QRCODE err: %s", err.Error())
					continue
				}
				go func() {
					defer func() {
						if r := recover(); r != nil {
							//		log.Errorf("%v", err)
						}
					}()
					//log.Infof("遍历获取到的 qrcode")
					var i int
					for evt := range ch {
						//spew.Dump(evt)
						if evt.Event == "code" && i == 0 {
							replyMsg := Message{
								Code:   "QRCODE",
								Result: evt.Code,
							}
							c.sendMsg(replyMsg)

						} else {
							log.Errorf("error when get qrCode ,%v", evt.Event)
						}
						i++
					}
					//log.Infof("遍历结束")
				}()

			}

		} else {
			log.Errorf("websocket message received of type %s", messageType)
		}
	}
}

func (c *Client) Write() {
	defer func() {
		//log.Infof("来自 write unRegister")
		Manager.unRegister(c)
	}()
	for {
		select {
		case message, ok := <-c.send:
			//log.Infof("准备发送消息")
			if !ok {
				//log.Infof("通道没有数据，且关闭了，，如果通道没有数据且通道未关闭就会堵塞")
				return
			}
			msg, err := json.Marshal(message)
			if err != nil {
				log.Errorf("error Marshal message: %s", err.Error())
				continue
			}

			//log.Infof("开始发送消息了")
			err = c.socket.WriteMessage(websocket.TextMessage, msg)

			if err != nil {
				//log.Infof("发送消息失败")
				log.Errorf("write message close error: %v", err)
				err = c.socket.WriteMessage(websocket.CloseMessage, []byte{})
				if err != nil {
					log.Errorf(" empty message write error:%v", err)
					return
				}
			}
			//log.Infof("发送消息成功")
		}
	}
}
