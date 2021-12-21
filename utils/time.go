package utils

import (
	"time"
)

var (
	DateFormat           = "2006-01-02"
	DatetimeFormat       = "2006-01-02 15:04:05"
	SerialDatetimeFormat = "20060102150405"
)

func LocalTime() time.Time {
	return time.Now().Local()
}

func InvalidTime(t time.Time) bool {
	if t.IsZero() {
		return true
	}
	valid, _ := time.Parse(DateFormat, "1000-01-02")
	return t.Before(valid)
}

func SetDateFormat(format string) {
	DateFormat = format
}

func SetDatetimeFormat(format string) {
	DatetimeFormat = format
}

func SetSerialDatetimeFormat(format string) {
	SerialDatetimeFormat = format
}

func LocalDateStr(t time.Time) string {
	if InvalidTime(t) {
		return ""
	}
	return t.Local().Format(DateFormat)
}

func LocalTimeStr(t time.Time) string {
	if InvalidTime(t) {
		return ""
	}
	return t.Local().Format(DatetimeFormat)
}

func SerialTimeStr() string {
	return time.Now().Format(SerialDatetimeFormat)
}

func StrToTime(str string) time.Time {
	t, _ := time.Parse(DatetimeFormat, str)
	return t.Local()
}
