package server

import (
	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/pkg/nacos"
)

type IDiscovery interface {
	RegisterInstance(int, float64, map[string]string) error
	UnregisterInstance(int) error
}

var discoverer IDiscovery

func initDiscoveryAdapter() {
	initNacos()
}

func initNacos() {
	nacos.InitNaming(config.GetCenterConfig().Nacos)
	discoverer = nacos.GetNamingClient()
}

func GetDiscoverer() IDiscovery {
	return discoverer
}
