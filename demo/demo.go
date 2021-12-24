package main

import (
	"time"

	common "github.com/EvisuXiao/andrews-common"
	demo "github.com/EvisuXiao/andrews-common/demo/src"
	"github.com/EvisuXiao/andrews-common/http"
	"github.com/EvisuXiao/andrews-common/logging"
)

func init() {
	common.Init("Andrews-common")
}

// 公司执行: go run demo/demo.go -dir=/Users/xiaowenbin/Workspace/code/andrews-project
// 家里执行: go run demo/demo.go -dir=/Users/evisu/Workspace/Andrews/app
func main() {
	logging.Info("All Environment initialized successfully!")
	http.StartServer(http.NewOption(demo.InitRouter()).WithQuitHandler(quitTest))
}

func quitTest() {
	logging.Debug("quitting start")
	time.Sleep(10 * time.Second)
	logging.Debug("quitting end")
}
