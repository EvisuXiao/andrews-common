package config

import "time"

type Cloud struct {
	SecretId  string `json:"secret_id"`
	SecretKey string `json:"secret_key"`
	Sms       *Sms   `json:"sms"`
}

type Sms struct {
	AppId     string      `json:"app_id"`
	Sign      string      `json:"sign"`
	Templates SmsTemplate `json:"templates"`
}

type SmsTemplate struct {
	Capcha SmsTemplateCapcha `json:"capcha"`
}

type SmsTemplateCapcha struct {
	Id      string        `json:"id"`
	Expired time.Duration `json:"expired" default:"2"`
}

var CloudConfig = &Cloud{}

func GetCloudConfig() *Cloud {
	return CloudConfig
}

func (c *Cloud) Name() string {
	return "cloud"
}

func (c *Cloud) Source() string {
	return ""
}

func (c *Cloud) Check() error {
	return nil
}

func (c *Cloud) Init() {
	c.Sms.Templates.Capcha.Expired = c.Sms.Templates.Capcha.Expired * time.Minute
}
