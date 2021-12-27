package server

import (
	"sync"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/pkg/nacos"
	"github.com/EvisuXiao/andrews-common/utils"
)

type IDiscovery interface {
	RegisterInstance(int, float64, map[string]string) error
	UnregisterInstance(int) error
}

type emptyDiscovery struct{}

var (
	once       sync.Once
	discoverer IDiscovery = &emptyDiscovery{}
)

func initDiscoverer() {
	once.Do(initNacos)
}

func initNacos() {
	cfg := config.GetNacosConfig()
	if utils.IsEmpty(cfg) {
		return
	}
	nacos.InitNaming(cfg)
	discoverer = nacos.GetNamingClient()
}

func GetDiscoverer() IDiscovery {
	return discoverer
}

func (d *emptyDiscovery) RegisterInstance(port int, weight float64, meta map[string]string) error {
	return nil
}

func (d *emptyDiscovery) UnregisterInstance(port int) error {
	return nil
}
