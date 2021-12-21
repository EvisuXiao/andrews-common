package utils

import (
	"bytes"
	"math/rand"
	"net/url"
	"strings"
	"time"
)

func GenerateRandomStr(size int) string {
	str := "abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ0123456789"
	rand.NewSource(time.Now().UnixNano())
	var s bytes.Buffer
	for i := 0; i < size; i++ {
		s.WriteByte(str[rand.Int63()%int64(len(str))])
	}
	return s.String()
}

func AddDirSuffixSlash(path string) string {
	if path == "" {
		path = "."
	}
	if !strings.HasSuffix(path, "/") {
		path += "/"
	}
	return path
}

func IsLocalUrl(u string) bool {
	uObj, err := url.Parse(u)
	if HasErr(err) {
		return false
	}
	return InSlice(uObj.Host, []string{"localhost", "127.0.0.1", "0.0.0.0"})
}
