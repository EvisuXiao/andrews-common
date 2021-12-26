package database

import (
	"fmt"

	"gorm.io/driver/mysql"
	"gorm.io/driver/postgres"
	"gorm.io/driver/sqlserver"
	"gorm.io/gorm"
	"gorm.io/gorm/schema"
	"gorm.io/plugin/dbresolver"

	"github.com/EvisuXiao/andrews-common/config"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

// 数据库驱动类型
const (
	DriverPostgres  = "postgres"
	DriverMySQL     = "mysql"
	DriverSQLServer = "mssql"
)

type database struct {
	name string
	db   *gorm.DB
}

var databases = make(map[string]*database)

// Init 初始化已注册的数据库及数据模型
func Init() {
	for _, db := range databases {
		db.setup()
	}
	for _, m := range models {
		m.MountDb()
		m.InitSchema(m)
	}
}

/**
 * 初始化数据库
 * @receiver *Database
 */
func (db *database) setup() {
	cnf, ok := (*config.GetDatabaseConfigs())[db.name]
	if !ok {
		logging.Fatal("Init: database %s connection name not found", db.name)
	}
	if utils.IsEmpty(cnf.Master) {
		logging.Fatal("Init:database %s master database must be valid", db.name)
	}
	db.db = conn(cnf)
	logging.Info("Database %s setup successfully!", db.name)
}

/**
 * 数据库连接
 * @param  *setting.Database cnf DB配置
 * @return *gorm.DB
 */
func conn(cnf *config.Database) *gorm.DB {
	// 连接主库
	db, err := gorm.Open(dbDialer(cnf.Driver, cnf.Master), &gorm.Config{
		NamingStrategy:         schema.NamingStrategy{TablePrefix: cnf.TablePrefix, SingularTable: true},
		SkipDefaultTransaction: true,
		NowFunc:                utils.LocalTime,
	})
	if utils.HasErr(err) {
		logging.Fatal("Init: database connection err: %+v", err)
	}
	// 是否开启读写分离
	if cnf.Separation {
		if utils.IsEmpty(cnf.Slave) {
			logging.Fatal("Setup: slave database must be valid when separated")
		}
		// 注册从库
		resolver := dbresolver.Register(dbresolver.Config{
			Replicas: []gorm.Dialector{dbDialer(cnf.Driver, cnf.Slave)},
		})
		err = db.Use(resolver)
		if utils.IsEmpty(cnf.Slave) {
			logging.Fatal("Init: slave database registered err: %+v", err)
		}
	}
	// 本地开启调试模式
	if config.IsLocalEnv() {
		db = db.Debug()
	}
	// 设置连接池
	if !utils.IsEmpty(cnf.PoolSize) {
		sqlDB, _ := db.DB()
		sqlDB.SetMaxOpenConns(cnf.PoolSize)
		sqlDB.SetMaxIdleConns(utils.CeilInt(cnf.PoolSize, 2))
		sqlDB.SetConnMaxLifetime(cnf.PoolLifeTime)
	}
	return db
}

/**
 * 获取DB连接器
 * @param  string driver DB驱动类型
 * @param  *setting.DatabaseConnection DB连接配置
 * @return gorm.Dialector
 */
func dbDialer(driver string, cnf *config.DatabaseConnection) gorm.Dialector {
	var dialer gorm.Dialector
	switch driver {
	case DriverMySQL:
		dialer = mysql.Open(mysqlDSN(cnf))
	case DriverPostgres:
		dialer = postgres.Open(postgresDSN(cnf))
	case DriverSQLServer:
		dialer = sqlserver.Open(mssqlDSN(cnf))
	default:
		logging.Fatal("Setup: unknown database driver: %s", driver)
	}
	return dialer
}

/**
 * 获取MySQL连接DSN
 * @return string
 */
func mysqlDSN(cnf *config.DatabaseConnection) string {
	return fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?parseTime=True&charset=utf8mb4&loc=Local", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Database)
}

/**
 * 获取PostgresSQL连接DSN
 * @return string
 */
func postgresDSN(cnf *config.DatabaseConnection) string {
	return fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable statement_cache_mode=describe", cnf.Host, cnf.Port, cnf.User, cnf.Password, cnf.Database)
}

/**
 * 获取SQLServer连接DSN
 * @return string
 */
func mssqlDSN(cnf *config.DatabaseConnection) string {
	return fmt.Sprintf("sqlserver://%s:%s@%s:%d?database=%s", cnf.User, cnf.Password, cnf.Host, cnf.Port, cnf.Database)
}

/**
 * 注册数据库
 * 初始化数据库模型时需调用此方法
 * @param  string name 数据库标识名
 */
func RegisterDatabase(dbName string) {
	databases[dbName] = &database{name: dbName}
}

/**
 * 根据数据库标识名获取ORM实例
 * @param  string dbName 数据库标识名
 * @return *gorm.DB
 */
func GetDbByName(dbName string) *gorm.DB {
	db := getDatabaseByName(dbName)
	if !utils.IsEmpty(db) {
		return db.GetDb()
	}
	return nil
}

/**
 * 根据数据库标识名获取数据库实例
 * @param  string dbName 数据库标识名
 * @return *database
 */
func getDatabaseByName(dbName string) *database {
	if db, ok := databases[dbName]; ok {
		return db
	}
	logging.Error("cannot find database(%s)", dbName)
	return nil
}

/**
 * 获取ORM实例
 * @receiver *Database
 * @return *gorm.DB
 */
func (db *database) GetDb() *gorm.DB {
	return db.db.Session(&gorm.Session{})
}

/**
 * 获取数据库标识名
 * @receiver *Database
 * @return string
 */
func (db *database) GetDbName() string {
	return db.name
}
