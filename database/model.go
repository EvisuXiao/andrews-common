package database

import (
	"encoding/json"
	"errors"
	"strconv"
	"sync"

	"gorm.io/gorm"
	"gorm.io/gorm/schema"

	"github.com/EvisuXiao/andrews-common/exception"
	"github.com/EvisuXiao/andrews-common/logging"
	"github.com/EvisuXiao/andrews-common/utils"
)

type IModel interface {
	SetDb()
	InitSchema(IModel)
}

type Model struct {
	dbName   string
	schema   *schema.Schema
	tenantId int
	Id       int64 `gorm:"primaryKey" json:"id"`
}

var models []IModel

/**
 * 注册数据模型
 * 初始化表数据模型时需调用此方法
 * @param  IModel m 数据模型
 */
func RegisterModel(m IModel) {
	models = append(models, m)
}

/**
 * 初始化表数据模型属性
 * @receiver *Model
 * @param  IModel s 表数据模型
 */
func (m *Model) InitSchema(s IModel) {
	sc, err := schema.Parse(s, &sync.Map{}, schema.NamingStrategy{SingularTable: true})
	if utils.HasErr(err) {
		logging.Fatal("Init: init db schema err: %+v", err)
	}
	m.schema = sc
}

func (m *Model) SetTenantId(tenantId int) {
	m.tenantId = tenantId
}

func (m *Model) IsSAASMode() bool {
	return !utils.IsEmpty(m.tenantId)
}

func (m *Model) SetDbName(dbName string) {
	m.dbName = dbName
}

func (m *Model) GetDb() *gorm.DB {
	return TenantDbResolver(m.tenantId, m.dbName).GetDb()
}

/**
 * 获取数据表模型属性
 * @receiver *Model
 * @return string
 */
func (m *Model) GetSchema() *schema.Schema {
	return m.schema
}

/**
 * 获取数据表名
 * @receiver *Model
 * @return string
 */
func (m *Model) GetTableName() string {
	if utils.IsEmpty(m.schema) {
		return ""
	}
	return m.schema.Table
}

/**
 * 获取数据表主键名
 * @receiver *Model
 * @return string
 */
func (m *Model) GetPk() string {
	if utils.IsEmpty(m.schema) || utils.IsEmpty(m.schema.PrimaryFieldDBNames) {
		return ""
	}
	return m.schema.PrimaryFieldDBNames[0]
}

/**
 * 获取数据表字段名集合
 * @receiver *Model
 * @return []string
 */
func (m *Model) GetFieldsName() []string {
	if utils.IsEmpty(m.schema) {
		return []string{}
	}
	return m.schema.DBNames
}

/***************************数据库操作***************************/

/**
 * 插入数据, 支持多条
 * @receiver *Model
 * @param  interface{} row 待插入数据
 * @return int 影响条数
 * @return error
 */
func (m *Model) AddRow(row interface{}) (int, error) {
	res := m.GetDb().Create(row)
	if res.Error != nil {
		return 0, exception.DbErrWrapper(res.Error)
	}
	return int(res.RowsAffected), nil
}

/**
 * 更新数据
 * @receiver *Model
 * @param  Conditions conditions 条件
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *Model) UpdateRows(conditions Conditions, data map[string]interface{}) error {
	pk := m.GetPk()
	if id, ok := data[pk]; ok {
		conditions.AddEqCondition(pk, id)
		delete(data, pk)
	}
	for k, v := range data {
		if utils.IsArrayOrMap(v) {
			data[k], _ = json.Marshal(v)
		} else if fv, ok := v.(float64); ok {
			data[k] = strconv.FormatFloat(fv, 'f', -1, 64)
		}
	}
	dbClone := m.BuildQuery(NewOptions().WithConditions(conditions))
	dbClone = dbClone.Table(m.GetTableName()).UpdateColumns(data)
	if !utils.IsEmpty(dbClone.Error) {
		return exception.DbErrWrapper(dbClone.Error)
	}
	if utils.IsEmpty(dbClone.RowsAffected) {
		return exception.DbErrWrapper(errors.New("no rows affected"))
	}
	return nil
}

/**
 * 根据ID更新数据
 * @receiver *Model
 * @param  int64 id ID
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *Model) UpdateRowById(id int64, data map[string]interface{}) error {
	return m.UpdateRows(NewSingleEqConditions(m.GetPk(), id), data)
}

/**
 * 根据批量ID更新数据
 * @receiver *Model
 * @param  []int64 ids 批量ID
 * @return map[string]interface{} data 待更新数据
 * @return error
 */
func (m *Model) UpdateRowsByIds(ids []int64, data map[string]interface{}) error {
	return m.UpdateRows(NewSingleConditions(m.GetPk(), OP_IN, ids), data)
}

/**
 * 删除数据
 * @receiver *Model
 * @param  Conditions conditions 条件
 * @return error
 */
func (m *Model) DeleteRows(conditions Conditions) error {
	dbClone := m.BuildQuery(NewOptions().WithConditions(conditions))
	return exception.DbErrWrapper(dbClone.Table(m.GetTableName()).Delete(nil).Error)
}

/**
 * 根据ID删除数据
 * @receiver *Model
 * @param  int64 id ID
 * @return error
 */
func (m *Model) DeleteRowById(id int64) error {
	return m.DeleteRows(NewSingleEqConditions(m.GetPk(), id))
}

/**
 * 根据批量ID删除数据
 * @receiver *Model
 * @param  []int64 ids 批量ID
 * @return error
 */
func (m *Model) DeleteRowsByIds(ids []int64) error {
	return m.DeleteRows(NewSingleConditions(m.GetPk(), OP_IN, ids))
}

/**
 * 数据是否存在
 * @receiver *Model
 * @param  Conditions conditions 条件
 * @return bool
 */
func (m *Model) Exists(conditions Conditions) bool {
	return m.GetCount(conditions) > 0
}

/**
 * 查询数据条数
 * @receiver *Model
 * @param  Conditions conditions 条件
 * @return int
 */
func (m *Model) GetCount(conditions Conditions) int {
	var total int64
	options := NewOptions().WithConditions(conditions)
	dbClone := m.BuildQuery(options)
	dbClone.Table(m.GetTableName()).Count(&total)
	return int(total)
}

/**
 * 查询去重列数据条数
 * @receiver *Model
 * @param  []string   fields 查询字段
 * @param  Conditions conditions 条件
 * @return int
 */
func (m *Model) GetFieldCount(fields []string, conditions Conditions) int {
	var total int64
	options := NewOptions().WithFields(fields).WithConditions(conditions)
	dbClone := m.BuildQuery(options)
	dbClone.Table(m.GetTableName()).Distinct().Count(&total)
	return int(total)
}

/**
 * 查询数据
 * @receiver *Model
 * @param  *Options options 选项
 * @param  interface{} rows 数据模板
 * @return error
 */
func (m *Model) GetAnyRows(options *Options, rows interface{}) error {
	if !options.HasOrder() {
		options.AddAscOrder(m.GetPk())
	}
	if !options.HasFields() {
		options.WithFields(m.GetFieldsName())
	}
	dbClone := m.BuildQuery(options)
	if options.Page > 0 && options.PageSize > 0 {
		offset := options.PageSize * (options.Page - 1)
		dbClone = dbClone.Offset(offset).Limit(options.PageSize)
	}
	return exception.DbErrWrapper(dbClone.Find(rows).Error)
}

/**
 * 根据批量ID查询数据
 * @receiver *Model
 * @param  []int64 ids 批量ID
 * @param  []string fields 查询字段
 * @param  interface{} rows 数据模板
 * @return error
 */
func (m *Model) GetAnyRowsByIds(ids []int64, fields []string, rows interface{}) error {
	options := NewOptions().WithFields(fields).AddCondition(m.GetPk(), OP_IN, ids)
	return m.GetAnyRows(options, rows)
}

/**
 * 查询单条数据
 * @receiver *Model
 * @param  *Options options 选项
 * @param  interface{} row 数据模板
 * @return error
 */
func (m *Model) GetAnyRow(options *Options, row interface{}) error {
	if !options.HasOrder() {
		options.AddAscOrder(m.GetPk())
	}
	if !options.HasFields() {
		options.WithFields(m.GetFieldsName())
	}
	err := m.BuildQuery(options).First(row).Error
	return exception.DbErrWrapper(err)
}

/**
 * 根据ID查询数据
 * @receiver *Model
 * @param  int64 ids ID
 * @param  []string fields 查询字段
 * @param  interface{} rows 数据模板
 * @return error
 */
func (m *Model) GetAnyRowById(id int64, fields []string, row interface{}) error {
	options := NewOptions().WithFields(fields).AddEqCondition(m.GetPk(), id)
	return m.GetAnyRow(options, row)
}
