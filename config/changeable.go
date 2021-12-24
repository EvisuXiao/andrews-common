package config

import (
	"log"

	"github.com/EvisuXiao/andrews-common/utils"
)

type Changeable struct {
	raw IConfig
}

func (c *Changeable) Name() string {
	return ""
}

func (c *Changeable) Source() string {
	// must source center!
	return SourceCenter
}

func (c *Changeable) FileType() string {
	return TypeJson
}

func (c *Changeable) Init() {
	// 尽量不要进行初始化修改原始配置, 会导致持久化时将默认值写入配置文件
}

func (c *Changeable) Persist(cfg IConfig) error {
	b, err := putContent(cfg)
	if utils.HasErr(err) {
		return err
	}
	name := cfg.Name()
	err = centerClient.PublishConfig(name, string(b))
	if utils.HasErr(err) {
		return err
	}
	log.Printf("[INFO] %s configuration published successfully!\n", name)
	return nil
}

//func (c *Changeable) SetXXX() error {
//	c.XXX = XXX
//	return c.Persist()
//}
