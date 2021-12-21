package database

import (
	"time"

	"github.com/EvisuXiao/andrews-common/utils"
)

// ModelUpdatable 预定义带有自动添加更新时间戳模型
type ModelUpdatable struct {
	ModelInsertable
	CreatedBy int       `json:"created_by"`
	UpdatedBy int       `json:"updated_by"`
	UpdatedAt time.Time `gorm:"autoUpdateTime" json:"updated_at"`
}

/**
 * 获取更新时间键名
 * @receiver *ModelUpdatable
 * @return string
 */
func (m *ModelUpdatable) getUpdatedTimeKey() string {
	if utils.IsEmpty(m.schema) {
		return ""
	}
	for _, f := range m.schema.Fields {
		if _, ok := f.TagSettings["AUTOUPDATETIME"]; ok {
			return f.DBName
		}
	}
	return "updated_time"
}

/**
 * 更新数据
 * 由于是以map形式更新, autoUpdateTime失效, 需手动添加时间戳
 * @receiver *ModelUpdatable
 * @param  Conditions conditions 条件
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *ModelUpdatable) UpdateRows(conditions Conditions, data map[string]interface{}) error {
	data[m.getUpdatedTimeKey()] = utils.LocalTime()
	return m.Model.UpdateRows(conditions, data)
}

/**
 * 根据ID更新数据
 * @receiver *ModelUpdatable
 * @param  int64 id ID
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *ModelUpdatable) UpdateRowById(id int64, data map[string]interface{}) error {
	return m.UpdateRows(NewSingleEqConditions(m.GetPk(), id), data)
}

/**
 * 根据批量ID更新数据
 * @receiver *ModelUpdatable
 * @param  []int64 ids 批量ID
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *ModelUpdatable) UpdateRowsByIds(ids []int64, data map[string]interface{}) error {
	return m.UpdateRows(NewSingleConditions(m.GetPk(), OP_IN, ids), data)
}
