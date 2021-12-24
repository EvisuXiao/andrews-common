package nacos

import (
	"log"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/naming_client"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/EvisuXiao/andrews-common/constants"
	"github.com/EvisuXiao/andrews-common/utils"
)

type NamingClient struct {
	client      naming_client.INamingClient
	groupName   string
	serviceName string
}

var (
	namingClient = &NamingClient{}
)

func InitNaming(cfg *constants.Nacos) {
	var err error
	namingClient.client, err = clients.NewNamingClient(buildClientParam(cfg))
	if utils.HasErr(err) {
		log.Fatalf("[FATAL] Init fatal: init nacos naming client error: %+v\n", err)
	}
	namingClient.groupName = cfg.GroupName
	namingClient.serviceName = cfg.ServiceName
	log.Println("[INFO] Init nacos naming client successfully")
}

func GetNamingClient() *NamingClient {
	return namingClient
}

func (c *NamingClient) RegisterInstance(port int, weight float64, meta map[string]string) error {
	param := vo.RegisterInstanceParam{
		Ip:          utils.GetLocalIP(),
		Port:        uint64(port),
		ServiceName: c.serviceName,
		GroupName:   c.groupName,
		Weight:      weight,
		Enable:      true,
		Healthy:     true,
		Ephemeral:   true,
		Metadata:    meta,
	}
	_, err := namingClient.client.RegisterInstance(param)
	return err
}

func (c *NamingClient) UnregisterInstance(port int) error {
	param := vo.DeregisterInstanceParam{
		Ip:          utils.GetLocalIP(),
		Port:        uint64(port),
		ServiceName: c.serviceName,
		GroupName:   c.groupName,
		Ephemeral:   true,
	}
	_, err := namingClient.client.DeregisterInstance(param)
	return err
}
