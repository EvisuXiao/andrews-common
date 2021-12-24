package config

import (
	"time"
)

type Databases map[string]*Database
type Database struct {
	Driver       string
	Separation   bool
	Master       *DatabaseConnection
	Slave        *DatabaseConnection
	TablePrefix  string
	PoolSize     int           `default:"50"`
	PoolLifeTime time.Duration `default:"3600"`
}
type DatabaseConnection struct {
	Host     string
	Port     int
	User     string
	Password string
	Database string
}

var DatabaseConfigs = &Databases{}

func GetDatabaseConfigs() *Databases {
	return DatabaseConfigs
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
