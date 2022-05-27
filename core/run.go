package core

import (
	"github.com/gin-gonic/gin"
	"github.com/jxcia/go-donkey/core/log"
)

// Run service start
func (g *Garden) Run(route func(r *gin.Engine), rpc interface{}, service string) {
	if service == "http" {
		go g.runHttpServer(route)
	} else if service == "rpc" {
		go g.runRpcServer(rpc)
	}
	forever := make(chan int, 0)
	<-forever
}

func (g *Garden) runHttpServer(route func(r *gin.Engine)) {
	address := g.GetServiceIp()
	listenAddress := address + ":" + g.cfg.Service.Port
	if err := g.ginListen(listenAddress, route); err != nil {
		log.Fatal("gin", "", err)
	}
}

func (g *Garden) runRpcServer(rpc interface{}) {

	if err := g.grpcLinsten(g.cfg.Service.Port); err != nil {
		log.Fatal("rpcRun", "", err)
	}
}
