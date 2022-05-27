package core

import (
	"errors"
	"fmt"

	clientV3 "go.etcd.io/etcd/client/v3"
)

// setSafe keys
func (g *Garden) setSafe(name string, val interface{}) {
	g.container.Store(name, val)
}
func (g *Garden) GetEtcd() (*clientV3.Client, error) {
	res, err := g.Get("etcd")
	if err != nil {
		return nil, err
	}
	return res.(*clientV3.Client), nil
}

//Get instance by name
func (g *Garden) Get(name string) (interface{}, error) {
	if res, ok := g.container.Load(name); ok {
		return res, nil
	}
	return nil, errors.New(fmt.Sprintf("Not found %s from container! ", name))
}
