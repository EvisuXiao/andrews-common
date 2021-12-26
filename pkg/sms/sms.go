package sms

import (
	"encoding/json"
	"fmt"

	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/profile"
	"github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/common/regions"
	sms "github.com/tencentcloud/tencentcloud-sdk-go/tencentcloud/sms/v20210111"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

var client *sms.Client

// Init 需启用cloud配置, config.RegisterConfig(config.CloudConfig)或config.Init("serviceName", config.CloudConfig)
func Init() {
	var err error
	cfg := config.GetCloudConfig().Tencent
	credential := common.NewCredential(cfg.SecretId, cfg.SecretKey)
	client, err = sms.NewClient(credential, regions.Beijing, profile.NewClientProfile())
	if utils.HasErr(err) {
		logging.Fatal("SMS client init err: %+v", err)
	}
}

func SendCapchaMessage(phone, capcha string) error {
	cfg := config.GetCloudConfig().Tencent.Sms
	req := buildSendRequest([]string{phone}, cfg.Templates.Capcha.Id, capcha, fmt.Sprint(cfg.Templates.Capcha.Expired.Minutes()))
	res, err := client.SendSms(req)
	b, _ := json.Marshal(res)
	logging.Debug("SMS resp: %s", string(b))
	return err
}

func buildSendRequest(phones []string, templateId string, templateParams ...string) *sms.SendSmsRequest {
	cfg := config.GetCloudConfig().Tencent.Sms
	req := sms.NewSendSmsRequest()
	req.SignName = common.StringPtr(cfg.Sign)
	req.SmsSdkAppId = common.StringPtr(cfg.AppId)
	req.TemplateId = common.StringPtr(templateId)
	if !utils.IsEmpty(templateParams) {
		req.TemplateParamSet = common.StringPtrs(templateParams)
	}
	req.PhoneNumberSet = common.StringPtrs(phones)
	return req
}
