package demo

import (
	"gorm.io/gorm"

	"github.com/EvisuXiao/andrews-common/database"
)

var dbName = "office"

// Manager 只需要ID的模型嵌套此结构
type Manager struct {
	database.Model
}

// ManagerInsertable 需要ID, 创建时间的模型嵌套此结构
type ManagerInsertable struct {
	database.ModelInsertable
}

// ManagerUpdatable 需要ID, 创建时间, 创建人, 更新时间, 更新人的模型嵌套此结构
type ManagerUpdatable struct {
	database.ModelUpdatable
}

/**
 * 关联数据库
 * @receiver *Manager
 */
func (m *Manager) SetDb() {
	m.SetDbName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerInsertable
 */
func (m *ManagerInsertable) SetDb() {
	m.SetDbName(dbName)
}

/**
 * 关联数据库
 * @receiver *ManagerUpdatable
 */
func (m *ManagerUpdatable) SetDb() {
	m.SetDbName(dbName)
}

/**
 * 获取moderation数据库实例
 * @return *gorm.DB
 */
func GetDb() *gorm.DB {
	return database.GetDbByName(dbName)
}

func GetTenantDb(tenantId int) *gorm.DB {
	return database.GetDbByName(database.GetTenantDbName(tenantId, dbName))
}
