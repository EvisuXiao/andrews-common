package main

import (
	"time"

	"github.com/EvisuXiao/andrews-common/config"
	demo "github.com/EvisuXiao/andrews-common/demo/src"
	"github.com/EvisuXiao/andrews-common/http"
	"github.com/EvisuXiao/andrews-common/logging"
)

func init() {
	config.Init("Andrews-common")
}

// 执行: go run demo/demo.go -dir=demo/src
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
