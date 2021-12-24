package config

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"gopkg.in/yaml.v3"

	"github.com/EvisuXiao/andrews-common/pkg/validator"
	"github.com/EvisuXiao/andrews-common/utils"
)

func mapCfg(content []byte, cfg IConfig) error {
	err := loadContent(content, cfg)
	if utils.HasErr(err) {
		return err
	}
	utils.SetStructDefaultValue(cfg)
	err = validator.Check(cfg)
	if utils.HasErr(err) {
		return err
	}
	cfg.Init()
	return nil
}

func readContent(cfg IConfig) ([]byte, error) {
	s := cfg.Source()
	if s == SourceDefault {
		s = source
	}
	if s == SourceCenter {
		return readFromCenter(cfg)
	}
	return readFromFile(cfg)
}

func readFromFile(cfg IConfig) ([]byte, error) {
	filename := AppFilePath(fmt.Sprintf("conf/%s.%s", cfg.Name(), strings.ToLower(cfg.FileType())))
	f, err := os.Open(filename)
	defer f.Close()
	if utils.HasErr(err) {
		return nil, err
	}
	return ioutil.ReadAll(f)
}

func loadContent(content []byte, cfg IConfig) error {
	switch cfg.FileType() {
	case TypeJson:
		return json.Unmarshal(content, cfg)
	case TypeYaml:
		return yaml.Unmarshal(content, cfg)
	default:
		return fmt.Errorf("invalid file type: %s", cfg.FileType())
	}
}

func putContent(cfg IConfig) ([]byte, error) {
	switch cfg.FileType() {
	case TypeJson:
		return json.Marshal(cfg)
	case TypeYaml:
		return yaml.Marshal(cfg)
	default:
		return nil, fmt.Errorf("invalid file type: %s", cfg.FileType())
	}
}
