package middleware

import (
	"net/http"
	"time"

	"github.com/gin-gonic/gin"
	gindump "github.com/tpkeeper/gin-dump"
)

// Middlewares store registered middlewares.
var Middlewares = defaultMiddlewares()

// NoCache is a middleware function that appends headers
// to prevent the client from caching the HTTP response.
func NoCache(c *gin.Context) {
	c.Header("Cache-Control", "no-cache, no-store, max-age=0, must-revalidate, value")
	c.Header("Expires", "Thu, 01 Jan 1970 00:00:00 GMT")
	c.Header("Last-Modified", time.Now().UTC().Format(http.TimeFormat))
	c.Next()
}

// Options is a middleware function that appends headers
// for options requests and aborts then exits the middleware
// account and ends the network.
func Options(c *gin.Context) {
	if c.Request.Method != "OPTIONS" {
		c.Next()
	} else {
		c.Header("Access-Control-Allow-Origin", "*")
		c.Header("Access-Control-Allow-Methods", "GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Access-Control-Allow-Headers", "authorization, origin, content-type, accept")
		c.Header("Allow", "HEAD,GET,POST,PUT,PATCH,DELETE,OPTIONS")
		c.Header("Content-Type", "application/json")
		c.AbortWithStatus(http.StatusOK)
	}
}

// Secure is a middleware function that appends security
// and resource access headers.
func Secure(c *gin.Context) {
	c.Header("Access-Control-Allow-Origin", "*")
	c.Header("X-Frame-Options", "DENY") //DENY
	c.Header("X-Content-Type-Options", "nosniff")
	c.Header("X-XSS-Protection", "1; mode=block")

	if c.Request.TLS != nil {
		c.Header("Strict-Transport-Security", "max-age=31536000")
	}
}

func defaultMiddlewares() map[string]gin.HandlerFunc {
	return map[string]gin.HandlerFunc{
		"recovery":  gin.Recovery(), //recovery：捕获任何 panic，并恢复。
		"secure":    Secure,         //添加一些安全和资源访问相关的 HTTP 头
		"options":   Options,
		"nocache":   NoCache, //禁止客户端缓存 HTTP 请求的返回结果。
		"cors":      Cors(),  //HTTP 请求跨域中间件。
		"requestid": RequestID(),
		"logger":    Logger(),
		"dump":      gindump.Dump(), //打印出 HTTP 请求包和返回包的内容，方便 debug。注意，生产环境禁止加载该中间件。
	}
}
