package socialserver

import (
	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
	"go-socialapp/internal/socialserver/controller/v1/account"
	"go-socialapp/internal/socialserver/ws"
	"net/http"

	// custom gin validators.
	_ "github.com/marmotedu/iam/pkg/validator"
)

func initRouter(g *gin.Engine) {
	installMiddleware(g)
	installController(g)
}

func installMiddleware(g *gin.Engine) {
}

func installController(g *gin.Engine) *gin.Engine {

	g.NoRoute(func(c *gin.Context) {
		core.WriteResponse(c, errors.WithCode(code.ErrPageNotFound, "Page not found."), nil)
	})
	v1 := g.Group("/")
	{
		v2 := v1.Group("api/v1")
		routeGroup(v2)
		v1.GET("ws", wsHandler)
	}

	return g
}

func routeGroup(group *gin.RouterGroup) {

	// account RESTful resource
	accountRoute := group.Group("/account")
	{
		accountController := account.NewAccountController()

		accountRoute.POST("", middleware.DealHanlder(accountController.Create))
		accountRoute.GET("/login", middleware.DealHanlder(accountController.Login))
		accountRoute.POST("/logout/:phone", middleware.DealHanlder(accountController.Logout)) // admin api
	}

}

func wsHandler(c *gin.Context) {
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
	client := ws.NewClient(uid, conn)
	// 用户注册到用户管理上
	ws.Manager.Register(client)
	go client.Read()
	//	go client.Write()
}
