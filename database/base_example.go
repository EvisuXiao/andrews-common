package database

import (
	"gorm.io/gorm"
)

var dbName = "foo"

//func init() {
//	RegisterDatabase(dbName)
//}

// Manager 只需要ID的模型嵌套此结构
type Manager struct {
	Model
}

// ManagerInsertable 需要ID, 创建时间的模型嵌套此结构
type ManagerInsertable struct {
	ModelInsertable
}

// ManagerUpdatable 需要ID, 创建时间, 创建人, 更新时间, 更新人的模型嵌套此结构
type ManagerUpdatable struct {
	ModelUpdatable
}

/**
 * 关联数据库
 * @receiver *Manager
 */
func (m *Manager) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerInsertable
 */
func (m *ManagerInsertable) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerUpdatable
 */
func (m *ManagerUpdatable) MountDb() {
	m.SetDatabaseByName(dbName)
}

/**
 * 获取moderation数据库实例
 * @return *gorm.DB
 */
func GetDb() *gorm.DB {
	return GetDbByName(dbName)
}
