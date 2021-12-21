package mail

import (
	"fmt"
	"strings"

	"gopkg.in/gomail.v2"

	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

type Mail struct {
	AddressSuffix    string
	Server           string
	Port             int
	Password         string
	Sender           string
	DefaultReceivers []string
}

var message *gomail.Message
var dialer *gomail.Dialer
var mailConfig *Mail

func Init(cfg *Mail) {
	mailConfig = cfg
	message = gomail.NewMessage()
	message.SetAddressHeader("From", GetAddress(mailConfig.Sender), mailConfig.Sender)
	dialer = gomail.NewDialer(mailConfig.Server, mailConfig.Port, GetAddress(mailConfig.Sender), mailConfig.Password)
}

func Send(title, content string, receivers ...string) {
	if utils.IsEmpty(receivers) {
		receivers = mailConfig.DefaultReceivers
	}
	for idx, receiver := range receivers {
		receivers[idx] = GetAddress(receiver)
	}
	message.SetHeader("To", receivers...)
	message.SetHeader("Subject", title)
	message.SetBody("text/html", content)
	go func() {
		err := dialer.DialAndSend(message)
		if utils.HasErr(err) {
			logging.Error("Mail send err: %+v", err)
		}
	}()
}

func GetAddress(name string) string {
	if strings.Contains(name, "@") {
		return name
	}
	return fmt.Sprintf("%s@%s", name, mailConfig.AddressSuffix)
}
