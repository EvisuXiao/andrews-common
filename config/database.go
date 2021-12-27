package config

import (
	"time"
)

type Databases map[string]*Database
type Database struct {
	Driver       string              `json:"driver" binding:"required"`
	Separation   bool                `json:"separation"`
	Master       *DatabaseConnection `json:"master"`
	Slave        *DatabaseConnection `json:"slave"`
	TablePrefix  string              `json:"table_prefix"`
	PoolSize     int                 `json:"pool_size "default:"50"`
	PoolLifeTime time.Duration       `json:"pool_life_time "default:"3600"`
}
type DatabaseConnection struct {
	Host     string `json:"host"`
	Port     int    `json:"port"`
	User     string `json:"user"`
	Password string `json:"password"`
	Database string `json:"database"`
}

var DatabaseConfigs = &Databases{}

func GetDatabaseConfigs() Databases {
	return *DatabaseConfigs
}

func (c *Databases) Name() string {
	return "database"
}

func (c *Databases) Source() string {
	return SourceDefault
}

func (c *Databases) FileType() string {
	return TypeJson
}

func (c *Databases) Init() {
	for _, db := range *DatabaseConfigs {
		db.PoolLifeTime = db.PoolLifeTime * time.Second
	}
}
