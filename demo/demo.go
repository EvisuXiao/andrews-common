package main

import (
	"fmt"
	"time"

	common "github.com/EvisuXiao/andrews-common"
	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/database"
	demo "github.com/EvisuXiao/andrews-common/demo/src"
	"github.com/EvisuXiao/andrews-common/http"
	"github.com/EvisuXiao/andrews-common/logging"
)

func init() {
	common.Init("Andrews-common", config.DatabaseConfigs)
}

// 公司执行: go run demo/demo.go -dir=/Users/xiaowenbin/Workspace/code/andrews-project -env=local -source=center
// 家里执行: go run demo/demo.go -dir=/Users/evisu/Workspace/Andrews/app -env=local -source=center
func main() {
	logging.Info("All Environment initialized successfully!")
	//startHttpServer()
	modelTest()
}

func startHttpServer() {
	http.StartServer(http.NewOption(demo.InitRouter()).WithQuitHandler(quitTest))
}

func quitTest() {
	logging.Debug("quitting start")
	time.Sleep(10 * time.Second)
	logging.Debug("quitting end")
}

func modelTest() {
	database.Init()
	m := demo.NewUserModel()
	rows, err := m.GetRows(database.NewOptions())
	fmt.Println(rows[0].Oauth, err)
}
