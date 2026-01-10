package server

import (
	"log"

	"github.com/getmelove/gorder2/internal/common/middleware"
	"github.com/gin-gonic/gin"
	"github.com/sirupsen/logrus"
	"github.com/spf13/viper"
	"go.opentelemetry.io/contrib/instrumentation/github.com/gin-gonic/gin/otelgin"
)

func RunHttpServer(serviceName string, wrapper func(router *gin.Engine)) {
	// wrapper 来自于各个服务想要对gin做哪些更改
	addr := viper.Sub(serviceName).GetString("http-addr")
	if addr == "" {
		// TODO: Warning log
		log.Fatalf("Order service address is empty")
	}
	RunHttpServerOnAddr(addr, wrapper)
}

func RunHttpServerOnAddr(addr string, wrapper func(router *gin.Engine)) {
	apiRouter := gin.New()
	// 添加中间件
	setMiddlewares(apiRouter)
	//
	wrapper(apiRouter)
	apiRouter.Group("/api", func(c *gin.Context) {

	})
	// Here ---- 注册路由
	//apiRouter.GET("/ping", func(c *gin.Context) {
	//	c.JSON(200, gin.H{
	//		"message": "pong -- 12.13",
	//	})
	//})
	if err := apiRouter.Run(addr); err != nil {
		panic(err)
	}
}

func setMiddlewares(r *gin.Engine) {
	r.Use(middleware.StructuredLog(logrus.NewEntry(logrus.StandardLogger())))
	r.Use(gin.Recovery())
	r.Use(otelgin.Middleware("default_server"))
}
