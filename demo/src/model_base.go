package demo

import (
	"gorm.io/gorm"

	"github.com/EvisuXiao/andrews-common/database"
)

var dbName = "andrews_office"

func init() {
	database.RegisterDatabase(dbName)
}

// Model 只需要ID的模型嵌套此结构
type Model struct {
	database.Model
}

// ModelInsertable 需要ID, 创建时间的模型嵌套此结构
type ModelInsertable struct {
	database.ModelInsertable
}

// ModelUpdatable 需要ID, 创建时间, 创建人, 更新时间, 更新人的模型嵌套此结构
type ModelUpdatable struct {
	database.ModelUpdatable
}

/**
 * 关联数据库
 * @receiver *Manager
 */
func (m *Model) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerInsertable
 */
func (m *ModelInsertable) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerUpdatable
 */
func (m *ModelUpdatable) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 获取moderation数据库实例
 * @return *gorm.DB
 */
func GetDb() *gorm.DB {
	return database.GetDbByName(dbName)
}
