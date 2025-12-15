package server

import (
	"log"

	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
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
