package utils

import (
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"io"
	"os"
)

func EncodeMd5Str(value string) string {
	m := md5.New()
	m.Write([]byte(value))
	return hex.EncodeToString(m.Sum(nil))
}

func EncodeMd5File(filename string) (string, error) {
	file, err := os.Open(filename)
	if HasErr(err) {
		return "", err
	}
	m := md5.New()
	if _, err = io.Copy(m, file); HasErr(err) {
		return "", err
	}
	return hex.EncodeToString(m.Sum(nil)), nil
}

func EncodeJsonValue(value interface{}) []byte {
	b, _ := json.Marshal(value)
	return b
}
