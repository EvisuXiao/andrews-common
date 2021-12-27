package config

import (
	"time"
)

type Caches struct {
	Redis Redises
}

type Redis struct {
	Driver     string
	MasterName string
	Host       string
	Port       string
	Hosts      []string
	Password   string
	Database   int
	Prefix     string
	Timeout    Timeout
}

type Redises map[string]*Redis

var CacheConfigs = &Caches{}

func GetCacheConfigs() *Caches {
	return CacheConfigs
}

func GetRedisConfigs() Redises {
	return GetCacheConfigs().Redis
}

func (c *Caches) Name() string {
	return "cache"
}

func (c *Caches) Source() string {
	return SourceDefault
}

func (c *Caches) FileType() string {
	return TypeJson
}

func (c *Caches) Init() {
	for _, cache := range CacheConfigs.Redis {
		cache.Timeout.Read = cache.Timeout.Read * time.Second
		cache.Timeout.Write = cache.Timeout.Write * time.Second
	}
}
