package ws

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"net/http"
)

func Handler(c *gin.Context) {
	uid := c.Query("uid") // 客户端的id
	conn, err := (&websocket.Upgrader{
		CheckOrigin: func(r *http.Request) bool { // CheckOrigin解决跨域问题
			return true
		}}).Upgrade(c.Writer, c.Request, nil) // 升级成ws协议
	if err != nil {
		http.NotFound(c.Writer, c.Request)
		return
	}
	// 创建一个客户端实例
	client := NewClient(uid, conn)
	// 用户注册到用户管理上
	Manager.Register(client)
	go client.Read()
	go client.Write()
}
