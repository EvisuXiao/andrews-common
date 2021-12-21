package config

import (
	"encoding/json"
	"errors"
	"flag"
	"fmt"
	"io/ioutil"
	"log"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/EvisuXiao/andrews-common/utils"
)

const (
	SOURCE_JSON   = "json"
	SOURCE_YAML   = "yaml"
	SOURCE_CENTER = "center"
)

type IConfig interface {
	Name() string
	Source() string
	Check() error
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
	loadConf()
	log.Print("[INFO] All configuration loaded successfully!")
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
	flag.StringVar(&dir, "dir", "./", "The application directory")
	flag.StringVar(&source, "source", "json", "The source of config file. json, yaml, center is available")
	flag.Parse()
	dir = utils.AddDirSuffixSlash(dir)
	source = strings.ToLower(source)
}

func loadConf() {
	log.Println("[INFO] Load configuration")
	for _, cfg := range configs {
		MapTo(cfg)
	}
}

func MapTo(cfg IConfig) {
	var err error
	s := utils.Or(cfg.Source(), source).(string)
	switch s {
	case SOURCE_JSON:
		err = LoadFromJson(cfg)
	case SOURCE_YAML:
		err = LoadFromYaml(cfg)
	case SOURCE_CENTER:
		err = LoadFromCenter(cfg)
	default:
		log.Fatalf("[FATAL] Init fatal: invalid conf source(%s)\n", source)
	}
	name := cfg.Name()
	if utils.HasErr(err) {
		log.Fatalf("[FATAL] Init fatal: load %s error: %+v\n", name, err)
	}
	utils.SetStructDefaultValue(cfg)
	if err = cfg.Check(); utils.HasErr(err) {
		log.Fatalf("[FATAL] Init fatal: check %s error: %+v\n", name, err)
	}
	cfg.Init()
	log.Printf("[INFO] %s configuration loaded successfully!", name)
}

func LoadFromYaml(cfg IConfig) error {
	filename := AppFilePath(fmt.Sprintf("conf/%s.json", cfg.Name()))
	f, err := os.Open(filename)
	defer f.Close()
	if utils.HasErr(err) {
		return err
	}
	read, err := ioutil.ReadAll(f)
	if utils.HasErr(err) {
		return err
	}
	return yaml.Unmarshal(read, cfg)
}

func LoadFromJson(cfg IConfig) error {
	filename := AppFilePath(fmt.Sprintf("conf/%s.json", cfg.Name()))
	f, err := os.Open(filename)
	defer f.Close()
	if utils.HasErr(err) {
		return err
	}
	read, err := ioutil.ReadAll(f)
	if utils.HasErr(err) {
		return err
	}
	return json.Unmarshal(read, cfg)
}

func LoadFromCenter(cfg IConfig) error {
	return errors.New("todo function")
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
