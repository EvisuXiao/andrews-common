package common

import (
	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/pkg/cron"
	"github.com/EvisuXiao/andrews-common/pkg/validator"
	"github.com/EvisuXiao/andrews-common/utils"
)

func Init(serviceName string, cfgs ...config.IConfig) {
	validator.Init()
	config.Init(serviceName, cfgs...)
	cron.Start()
}

func Stop() {
	cron.Stop()
	if err := config.Stop(); utils.HasErr(err) {
		logging.Error("stop process err: %+v", err)
	}
}
