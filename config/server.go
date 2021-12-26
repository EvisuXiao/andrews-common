package config

import (
	"time"
)

const (
	EnvLocal   = "local"
	EnvTesting = "testing"
	EnvProd    = "production"
)

type Server struct {
	Port      int     `json:"port" binding:"gt=0,lt=65536"`
	Discovery bool    `json:"discovery" default:"true"`
	Weight    float64 `json:"weight" default:"100"`
	Timeout   Timeout `json:"timeout"`
	RateLimit int     `json:"rate_limit"`
}
type Timeout struct {
	Read  time.Duration `json:"read" default:"60"`
	Write time.Duration `json:"write" default:"60"`
	Exit  time.Duration `json:"exit" default:"3"`
}

var (
	env          string
	ServerConfig = &Server{}
)

func init() {
	RegisterConfig(ServerConfig)
}

func GetServerConfig() *Server {
	return ServerConfig
}

func (c *Server) Name() string {
	return "server"
}

func (c *Server) Source() string {
	return SourceDefault
}

func (c *Server) FileType() string {
	return TypeJson
}

func (c *Server) Init() {
	c.Timeout.Read = c.Timeout.Read * time.Second
	c.Timeout.Write = c.Timeout.Write * time.Second
	c.Timeout.Exit = c.Timeout.Exit * time.Second
}

func GetEnv() string {
	return env
}

func IsLocalEnv() bool {
	return env == EnvLocal
}

func IsTestingEnv() bool {
	return env == EnvTesting
}

func IsProdEnv() bool {
	return env == EnvProd
}
