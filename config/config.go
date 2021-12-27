package config

import (
	"flag"
	"fmt"
	"log"
	"strings"

	"github.com/EvisuXiao/andrews-common/utils"
)

const (
	SourceDefault = ""
	SourceFile    = "file"
	SourceCenter  = "center"

	TypeJson = "json"
	TypeYaml = "yaml"
)

type IConfig interface {
	Name() string
	Source() string
	FileType() string
	Init()
}

var (
	ServiceName string
	dir         string
	source      string
	configs     []IConfig
	inited      bool
)

// Init 默认加载server, common配置
// 其他配置请在可选参数中加载, 或手动调用RegisterConfig
func Init(serviceName string, cfgs ...IConfig) {
	if inited {
		return
	}
	configs = append(configs, cfgs...)
	setServiceName(serviceName)
	parseFlag()
	initCenter()
	loadConf()
	log.Println("[INFO] All configuration loaded successfully!")
	inited = true
}

func GetServiceName() string {
	return ServiceName
}

func setServiceName(name string) {
	ServiceName = name
}

func RegisterConfig(cfg IConfig) {
	configs = append(configs, cfg)
}

func parseFlag() {
	flag.StringVar(&dir, "dir", "./", "The application root directory")
	flag.StringVar(&env, "env", "testing", fmt.Sprintf("The application environment. %s, %s, %s is available. Default: %s", EnvLocal, EnvTesting, EnvProd, EnvTesting))
	flag.StringVar(&source, "source", SourceFile, fmt.Sprintf("The source of config file. %s, %s is available. Default: %s", SourceFile, SourceCenter, SourceFile))
	flag.StringVar(&center, "center", CenterNacos, fmt.Sprintf("The config center adapter. %s is supported, %s is in the todo list. Default: %s", CenterNacos, CenterApollo, CenterNacos))
	flag.Parse()
	dir = utils.AddDirSuffixSlash(dir)
	if env != EnvLocal && env != EnvProd {
		env = EnvTesting
	}
	source = strings.ToLower(source)
}

func loadConf() {
	log.Println("[INFO] Load configuration")
	var err error
	for _, cfg := range configs {
		if err = MapTo(cfg); utils.HasErr(err) {
			log.Fatalf("[FATAL] Init fatal: map conf %s err: %+v\n", cfg.Name(), err)
		}
		log.Printf("[INFO] %s configuration loaded successfully!\n", cfg.Name())
	}
}

func MapTo(cfg IConfig) error {
	read, err := readContent(cfg)
	if utils.HasErr(err) {
		return err
	}
	if err = mapCfg(read, cfg); utils.HasErr(err) {
		return err
	}
	return nil
}

func Stop() error {
	if source != SourceCenter {
		return nil
	}
	for _, cfg := range configs {
		name := cfg.Name()
		err := centerClient.CancelListenConfig(name)
		if utils.HasErr(err) {
			return err
		}
		log.Printf("[INFO] cancel listening %s configuration successfully!\n", name)
	}
	return nil
}

func AppFilePath(filename string) string {
	return dir + filename
}

func TempFilePath(filename string) string {
	tempPath := GetCommonConfig().TempPath
	if strings.HasPrefix(tempPath, "/") {
		return tempPath + filename
	}
	return AppFilePath(tempPath + filename)
}
