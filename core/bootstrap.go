package core

import (
	"os"

	"github.com/jxcia/go-donkey/core/drives/etcd"
	"github.com/jxcia/go-donkey/core/log"
)

func (g *Garden) bootstrap(bootstrap, env string) {
	g.cfg.Bootstrap = bootstrap
	g.cfg.Env = env
	g.bootConfig("yml")
	g.checkConfig()
	log.Setup(g.cfg.Service.LogPath, g.cfg.Service.ServiceName)
	log.Info("bootstrap", "", g.cfg.Service.ServiceName+" running")
	g.bootEtcd()
	//g.bootOpenTracing() 暂时不需要链路
}

// 检查关键字日志
func (g *Garden) checkConfig() {
	if g.cfg.Bootstrap == "" {
		os.Exit(-1)
		log.Fatal("config", "", "empty option serviceName")
	}

}
func (g *Garden) bootEtcd() {
	etcdC, err := etcd.Connect(g.cfg.Service.EtcdAddress, log.GetLogger().Desugar())
	if err != nil {
		log.Fatal("etcd", "", err)
	}
	g.setSafe("etcd", etcdC)
}
