package config

import (
	"log"

	"github.com/EvisuXiao/andrews-common/constants"
	"github.com/EvisuXiao/andrews-common/pkg/nacos"
	"github.com/EvisuXiao/andrews-common/utils"
)

const (
	CenterNacos  = "nacos"
	CenterApollo = "apollo"
)

type ICenter interface {
	GetConfig(string) (string, error)
	PublishConfig(string, string) error
	ListenConfig(string, func(string)) error
	CancelListenConfig(string) error
}

type Center struct {
	Nacos *constants.Nacos `json:"nacos"`
}

var (
	center       string
	centerClient ICenter
	centerConfig = &Center{}
)

func GetCenterConfig() *Center {
	return centerConfig
}

func (c *Center) Name() string {
	return "center"
}

func (c *Center) Source() string {
	return SourceFile
}

func (c *Center) FileType() string {
	return TypeJson
}

func (c *Center) Init() {
	c.initNacos()
}

func (c *Center) initNacos() {
	c.Nacos.ServiceName = GetServiceName()
	if utils.IsEmpty(c.Nacos.TempPath) {
		c.Nacos.TempPath = utils.AddDirSuffixSlash(AppFilePath("temp/nacos"))
	}
}

func initCenter() {
	if source != SourceCenter {
		return
	}
	MapTo(centerConfig)
	if center == CenterNacos {
		initNacosCenter()
	}
}

func initNacosCenter() {
	nacos.InitConfig(GetCenterConfig().Nacos)
	centerClient = nacos.GetConfigClient()
}

func readFromCenter(cfg IConfig) ([]byte, error) {
	name := cfg.Name()
	content, err := centerClient.GetConfig(name)
	if utils.HasErr(err) {
		return nil, err
	}
	err = centerClient.ListenConfig(name, func(content string) {
		_ = mapCfg([]byte(content), cfg)
		log.Printf("[INFO] %s configuration reloaded successfully!\n", name)
	})
	if utils.HasErr(err) {
		return nil, err
	}
	log.Printf("[INFO] listening %s configuration successfully!\n", name)
	return []byte(content), nil
}
