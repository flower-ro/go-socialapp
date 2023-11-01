package socialserver

import (
	"github.com/gin-gonic/gin"
	"github.com/marmotedu/component-base/pkg/core"
	"github.com/marmotedu/errors"
	"go-socialapp/internal/pkg/code"
	"go-socialapp/internal/pkg/middleware"
	"go-socialapp/internal/socialserver/controller/v1/account"

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

	v1 := g.Group("/api/v1")
	routeGroup(v1)

	return g
}

func routeGroup(group *gin.RouterGroup) {

	// account RESTful resource
	accountRoute := group.Group("/account")
	{
		accountController := account.NewAccountController()

		accountRoute.POST("", middleware.DealHanlder(accountController.Create))
		//accountRoute.GET(":chainId", middleware.DealHanlder(accountController.GetLasted)) // admin api
	}

}
