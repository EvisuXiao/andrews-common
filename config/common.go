package config

import (
	"github.com/EvisuXiao/andrews-common/utils"
)

type Common struct {
	PageSize             int    `json:"page_size" default:"20"`
	MaxPageSize          int    `json:"max_page_size" default:"5000"`
	DateFormat           string `json:"date_format" default:"2006-01-02"`
	DatetimeFormat       string `json:"datetime_format" default:"2006-01-02 15:04:05"`
	SerialDatetimeFormat string `json:"serial_datetime_format" default:"20060102150405"`
	TempPath             string `json:"temp_path" default:"temp/"`
}

var CommonConfig = &Common{}

func init() {
	RegisterConfig(CommonConfig)
}

func GetCommonConfig() *Common {
	return CommonConfig
}

func (c *Common) Name() string {
	return "common"
}

func (c *Common) Source() string {
	return SourceDefault
}

func (c *Common) FileType() string {
	return TypeJson
}

func (c *Common) Init() {
	utils.SetDateFormat(c.DateFormat)
	utils.SetDatetimeFormat(c.DatetimeFormat)
	utils.SetSerialDatetimeFormat(c.SerialDatetimeFormat)
}
