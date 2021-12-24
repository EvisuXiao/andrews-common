package ding

import (
	"fmt"
	"net/http"

	"github.com/EvisuXiao/andrews-common/curl"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Robot struct {
	token string
}

func RobotClient(token string) *Robot {
	return &Robot{token}
}

func (r *Robot) Send(title, message string) {
	dingUrl := fmt.Sprintf("%s?access_token=%s", host+funcSendMsg, r.token)
	if !utils.IsEmpty(title) {
		message = fmt.Sprintf("%s\r\n%s", title, message)
	}
	data := map[string]interface{}{
		"msgtype": "text",
		"text":    map[string]interface{}{"content": message},
	}
	var result ErrorResult
	err := curl.Request(dingUrl, http.MethodPost, data, &result)
	if utils.HasErr(err) {
		logging.Error("Ding send message err: %+v", err)
		return
	}
	if result.GetErrCode() != successCode {
		logging.Error("Ding send message err: %s", result.GetErrMsg())
	}
}
