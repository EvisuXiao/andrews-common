package nacos

import (
	"errors"
	"log"
	"net/url"
	"strconv"

	"github.com/nacos-group/nacos-sdk-go/clients"
	"github.com/nacos-group/nacos-sdk-go/clients/config_client"
	"github.com/nacos-group/nacos-sdk-go/common/constant"
	"github.com/nacos-group/nacos-sdk-go/vo"

	"github.com/EvisuXiao/andrews-common/constants"
	"github.com/EvisuXiao/andrews-common/utils"
)

type ConfigClient struct {
	client      config_client.IConfigClient
	groupName   string
	serviceName string
}

var (
	configClient = &ConfigClient{}

	ErrConfigNotFound = errors.New("config is not found")
)

func InitConfig(cfg *constants.Nacos) {
	var err error
	configClient.client, err = clients.NewConfigClient(buildClientParam(cfg))
	if utils.HasErr(err) {
		log.Fatalf("[FATAL] Init fatal: init nacos config client error: %+v\n", err)
	}
	configClient.groupName = cfg.GroupName
	configClient.serviceName = cfg.ServiceName
	log.Println("[INFO] Init nacos config client successfully")
}

func GetConfigClient() *ConfigClient {
	return configClient
}

func buildClientParam(cfg *constants.Nacos) vo.NacosClientParam {
	cCfg := constant.NewClientConfig(
		constant.WithNamespaceId(cfg.Namespace),
		constant.WithUsername(cfg.Username),
		constant.WithPassword(cfg.Password),
		constant.WithLogLevel("error"),
		constant.WithLogDir(cfg.TempPath+"log"),
		constant.WithCacheDir(cfg.TempPath+"cache"),
	)
	var sCfg []constant.ServerConfig
	for _, host := range cfg.Hosts {
		u, err := url.Parse(host)
		if utils.HasErr(err) {
			log.Fatalf("[FATAL] Init fatal: parse nacos host error: %+v\n", err)
		}
		p, _ := strconv.ParseUint(u.Port(), 10, 32)
		if utils.IsEmpty(p) {
			if u.Scheme == "http" {
				p = 80
			}
			if u.Scheme == "https" {
				p = 443
			}
		}
		sc := constant.NewServerConfig(u.Hostname(), p, constant.WithScheme(u.Scheme))
		sCfg = append(sCfg, *sc)
	}
	return vo.NacosClientParam{
		ClientConfig:  cCfg,
		ServerConfigs: sCfg,
	}
}

func (c *ConfigClient) GetConfig(dataId string) (string, error) {
	group, err := c.getGroupName(dataId)
	if utils.HasErr(err) {
		return "", err
	}
	param := vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	}
	content, err := configClient.client.GetConfig(param)
	if utils.HasErr(err) {
		return "", err
	}
	if utils.IsEmpty(content) {
		return "", ErrConfigNotFound
	}
	return content, nil
}

func (c *ConfigClient) PublishConfig(dataId, content string) error {
	group, err := c.getGroupName(dataId)
	if utils.HasErr(err) {
		return err
	}
	_, err = configClient.client.PublishConfig(vo.ConfigParam{
		DataId:  dataId,
		Group:   group,
		Content: content,
	})
	return err
}

func (c *ConfigClient) ListenConfig(dataId string, reload func(string)) error {
	group, err := c.getGroupName(dataId)
	if utils.HasErr(err) {
		return err
	}
	return configClient.client.ListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
		OnChange: func(namespace, group, dataId, data string) {
			reload(data)
		},
	})
}

func (c *ConfigClient) CancelListenConfig(dataId string) error {
	group, err := c.getGroupName(dataId)
	if utils.HasErr(err) {
		return err
	}
	return configClient.client.CancelListenConfig(vo.ConfigParam{
		DataId: dataId,
		Group:  group,
	})
}

func (c *ConfigClient) getGroupName(dataId string) (string, error) {
	param := vo.SearchConfigParam{
		Search:  "accurate",
		DataId:  dataId,
		Group:   c.groupName,
		AppName: c.serviceName,
	}
	res, err := configClient.client.SearchConfig(param)
	if utils.HasErr(err) {
		return "", err
	}
	if utils.IsEmpty(res.PageItems) {
		param.Group = constant.DEFAULT_GROUP
		param.AppName = ""
		res, err = configClient.client.SearchConfig(param)
		if utils.HasErr(err) {
			return "", err
		}
		if utils.IsEmpty(res.PageItems) {
			return "", ErrConfigNotFound
		}
	}
	return res.PageItems[0].Group, nil
}
