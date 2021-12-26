package config

import (
	"log"

	"github.com/EvisuXiao/andrews-common/utils"
)

// Changeable 可编辑的配置可继承此结构, 编辑后调用Persist进行持久化
// 只支持配置中心, 防止多实例下本地配置不同步. 不可进行对读取的值进行默认值, 转换单位等初始化操作, 会导致将初始化后的内容同步到配置中心, 如有需要请编写GetXXX方法调用
type Changeable struct{}

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
