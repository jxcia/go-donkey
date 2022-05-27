package core

import (
	"context"
	"errors"
	"strings"
	"time"

	"github.com/jxcia/go-donkey/core/log"
	clientV3 "go.etcd.io/etcd/client/v3"
)

type service struct {
	Nodes []node
}
type node struct {
	Addr    string
	Waiting int64
	Finish  int64
}
type serviceOperate struct {
	operate     string
	serviceName string
	serviceAddr string
	nodeIndex   int
}

func (g *Garden) bootService() {
	var err error

	g.services = map[string]*service{}
	g.cfg.ServiceIp, err = getOutboundIP()
	if err != nil {
		log.Fatal("bootService", "", err)
	}
	g.serviceManager = make(chan serviceOperate, 0)
	go g.serviceManageWatch(g.serviceManager)

	if err = g.serviceRegister(true); err != nil {
		log.Fatal("serviceRegister", "", err)
	}
}

func (g *Garden) serviceManageWatch(ch chan serviceOperate) {
	for {
		select {
		case sm := <-ch:
			switch sm.operate {

			case "addNode":
				g.createServiceIndex(sm.serviceName)
				g.services[sm.serviceName].Nodes = append(g.services[sm.serviceName].Nodes, node{Addr: sm.serviceAddr})
				break

			case "delNode":
				if g.existsService(sm.serviceName) {
					for i := 0; i < len(g.services[sm.serviceName].Nodes); i++ {
						if g.services[sm.serviceName].Nodes[i].Addr == sm.serviceAddr {
							g.services[sm.serviceName].Nodes = append(g.services[sm.serviceName].Nodes[:i], g.services[sm.serviceName].Nodes[i+1:]...)
							i--
						}
					}
				}
				break
			}

		}
	}
}

//GetServiceId get this service union id
func (g *Garden) GetServiceId() string {
	return g.cfg.Service.EtcdKey + "_" + g.cfg.Service.ServiceName + "_" + g.cfg.ServiceIp + ":" + g.cfg.Service.Port
}

func (g *Garden) serviceRegister(isReconnect bool) error {
	client, err := g.GetEtcd()
	if err != nil {
		return err
	}
	// New lease
	resp, err := client.Grant(context.TODO(), 2)
	if err != nil {
		return err
	}
	// The lease was granted
	if err != nil {
		return err
	}
	_, err = client.Put(context.TODO(), g.GetServiceId(), "0", clientV3.WithLease(resp.ID))
	if err != nil {
		return err
	}
	// keep alive
	ch, err := client.KeepAlive(context.TODO(), resp.ID)
	if err != nil {
		return err
	}
	// monitor etcd connection
	go func() {
		for {
			select {
			case resp := <-ch:
				if resp == nil {
					go g.serviceRegister(false)
					return
				}
			}
		}
	}()

	if isReconnect {
		go g.serviceWatcher()
		go func() {
			for {
				g.getAllServices()
				time.Sleep(time.Second * 5)
			}
		}()
	}
	return nil
}

func (g *Garden) serviceWatcher() {
	client, err := g.GetEtcd()
	if err != nil {
		log.Error("getEtcd", "", err)
		return
	}

	rch := client.Watch(context.Background(), g.cfg.Service.EtcdKey+"_", clientV3.WithPrefix())
	for wresp := range rch {
		for _, ev := range wresp.Events {
			arr := strings.Split(string(ev.Kv.Key), "_")
			serviceName := arr[1]
			httpAddr := arr[2]
			serviceAddr := httpAddr
			switch ev.Type {
			case 0: //put
				g.addServiceNode(serviceName, serviceAddr)
				log.Infof("service", "%s node %s join", serviceName, serviceAddr)
			case 1: //delete
				g.delServiceNode(serviceName, serviceAddr)
				log.Infof("service", "%s node %s leave", serviceName, serviceAddr)
			}
		}
	}
}

func (g *Garden) getAllServices() ([]string, error) {
	client, err := g.GetEtcd()
	if err != nil {
		return []string{}, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := client.Get(ctx, g.cfg.Service.EtcdKey+"_", clientV3.WithPrefix())
	cancel()
	if err != nil {
		log.Error("getAllServices", "", err)
		return []string{}, nil
	}
	var services []string
	for _, ev := range resp.Kvs {
		arr := strings.Split(string(ev.Key), g.cfg.Service.EtcdKey+"_")
		service := arr[1]
		services = append(services, service)
	}

	for _, service := range services {
		arr := strings.Split(service, "_")
		serviceName := arr[0]
		serviceHttpAddr := arr[1]

		g.addServiceNode(serviceName, serviceHttpAddr)
	}

	return services, nil
}

func (g *Garden) getServicesByName(serviceName string) ([]string, error) {
	client, err := g.GetEtcd()
	if err != nil {
		return nil, err
	}
	ctx, cancel := context.WithTimeout(context.Background(), 1*time.Second)
	resp, err := client.Get(ctx, g.cfg.Service.EtcdKey+"_"+serviceName, clientV3.WithPrefix())
	cancel()
	if err != nil {
		log.Error("getServicesByName", "", err)
		return []string{}, nil
	}
	var services []string
	for _, ev := range resp.Kvs {
		arr := strings.Split(string(ev.Key), g.cfg.Service.EtcdKey+"_"+serviceName+"_")
		serviceAddr := arr[1]
		services = append(services, serviceAddr)
	}
	return services, nil
}

func (g *Garden) addServiceNode(name, addr string) {
	sm := serviceOperate{
		operate:     "addNode",
		serviceName: name,
		serviceAddr: addr,
	}
	g.serviceManager <- sm
}

func (g *Garden) delServiceNode(name, addr string) {
	sm := serviceOperate{
		operate:     "delNode",
		serviceName: name,
		serviceAddr: addr,
	}
	g.serviceManager <- sm
}

func (g *Garden) createServiceIndex(name string) {
	if !g.existsService(name) {
		g.services[name] = &service{
			Nodes: []node{},
		}
	}
}

func (g *Garden) existsService(name string) bool {
	_, ok := g.services[name]
	return ok
}

func (g *Garden) getServiceHttpAddr(name string, index int) (string, error) {
	if index > len(g.services[name].Nodes)-1 {
		return "", errors.New("service node not found")
	}
	arr := strings.Split(strings.Split(g.services[name].Nodes[index].Addr, "_")[0], ":")
	return arr[0] + ":" + arr[1], nil
}

func (g *Garden) getServiceRpcAddr(name string, index int) (string, error) {
	if index > len(g.services[name].Nodes)-1 {
		return "", errors.New("service node not found")
	}
	arr := strings.Split(strings.Split(g.services[name].Nodes[index].Addr, "_")[0], ":")
	return arr[0] + ":" + arr[2], nil
}

//GetServiceIp get this service ip
func (g *Garden) GetServiceIp() string {
	return g.cfg.ServiceIp
}
