package core

import (
	"io/ioutil"

	"git.tanzk.cn/om/tools/go-donkey/core/log"
	"github.com/gin-gonic/gin"
	"github.com/google/uuid"
)

func (g *Garden) ginListen(listenAddress string, route func(r *gin.Engine)) error {
	gin.SetMode(gin.ReleaseMode)
	engine := gin.New()
	engine.Use(gin.Recovery(), SetLoggerMiddleware())

	gin.DefaultWriter = ioutil.Discard
	route(engine)

	log.Infof("http", "listen on: %s", listenAddress)
	return engine.Run(listenAddress)
}

func SetLoggerMiddleware() gin.HandlerFunc {
	return func(c *gin.Context) {
		traceId := c.GetHeader("traceId")
		if traceId == "" {
			traceId = uuid.New().String()
		}
		c.Set("traceId", traceId)
		c.Next()
	}
}
