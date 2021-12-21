package database

import (
	"fmt"

	"gorm.io/gorm"
	"gorm.io/plugin/dbresolver"

	"github.com/EvisuXiao/andrews-common/utils"
)

type Options struct {
	Fields     []string   // 查询字段
	Conditions Conditions // 查询条件
	Orders     Orders     // 排序
	Groups     []string   // 分组
	Page       int        // 页码
	PageSize   int        // 每页数量
	Force      bool       // 是否强制使用主库
}

// Condition 条件结构
type Condition struct {
	Column string
	Op     string
	Value  interface{}
}
type Conditions []Condition

// Order 排序结构
type Order struct {
	Column string
	Sort   string
}
type Orders []Order

const (
	OP_EQ     = "="
	OP_NEQ    = "!="
	OP_GT     = ">"
	OP_LT     = "<"
	OP_GEQ    = ">="
	OP_LEQ    = "<="
	OP_LIKE   = "LIKE"
	OP_IS     = "IS"
	OP_NIS    = "IS NOT"
	OP_IN     = "IN"
	OP_NIN    = "NOT IN"
	OP_NULL   = "NULL"
	OP_NNULL  = "NOT NULL"
	OP_RAW    = "RAW"
	SORT_ASC  = "ASC"
	SORT_DESC = "DESC"
)

/**
 * 构建预处理模型
 * @receiver *Model
 * @param  *Options options 条件
 * @return *gorm.DB
 */
func (m *Model) BuildQuery(options *Options) *gorm.DB {
	// 必须用副本形式DB层层传递, 保证线程安全
	dbClone := m.GetDb()
	// 填充字段
	if !utils.IsEmpty(options.Fields) {
		dbClone = dbClone.Select(options.Fields)
	}
	// 是否强制主库
	if options.Force {
		dbClone = dbClone.Clauses(dbresolver.Write)
	}
	// 构建条件
	for _, cond := range options.Conditions {
		cond.Op = utils.Or(cond.Op, OP_EQ).(string)
		switch cond.Op {
		case OP_IN, OP_NIN:
			if utils.IsEmpty(cond.Value) {
				dbClone = dbClone.Where("1=0")
			} else {
				dbClone = dbClone.Where(fmt.Sprintf("%s %s (?)", cond.Column, cond.Op), cond.Value)
			}
		case OP_NULL, OP_NNULL:
			dbClone = dbClone.Where(fmt.Sprintf("%s IS %s", cond.Column, cond.Op))
		case OP_IS, OP_NIS:
			dbClone = dbClone.Where(fmt.Sprintf("%s %s %v", cond.Column, cond.Op, cond.Value))
		case OP_RAW:
			dbClone = dbClone.Where(cond.Value)
		default:
			dbClone = dbClone.Where(fmt.Sprintf("%s %s ?", cond.Column, cond.Op), cond.Value)
		}
	}
	// 构建分组
	for _, group := range options.Groups {
		dbClone = dbClone.Group(group)
	}
	// 构建排序
	for _, order := range options.Orders {
		dbClone = dbClone.Order(fmt.Sprintf("%s %s", order.Column, order.Sort))
	}
	return dbClone
}

func NewOptions() *Options {
	return &Options{}
}

/**
 * 设置查询字段
 * @receiver *Options
 * @param  []string fields 查询字段
 * @return *Options
 */
func (o *Options) WithFields(fields []string) *Options {
	o.Fields = fields
	return o
}

/**
 * 是否设置查询字段
 * @receiver *Options
 * @return bool
 */
func (o *Options) HasFields() bool {
	return !utils.IsEmpty(o.Fields)
}

/**
 * 添加查询字段
 * @receiver *Options
 * @param  string field 查询字段
 * @return *Options
 */
func (o *Options) AddField(field string) *Options {
	o.Fields = append(o.Fields, field)
	return o
}

/**
 * 设置查询条件
 * @receiver *Options
 * @param  Conditions conditions 查询条件
 * @return *Options
 */
func (o *Options) WithConditions(conditions Conditions) *Options {
	o.Conditions = conditions
	return o
}

/**
 * 添加查询条件
 * @receiver *Options
 * @param  string column 查询字段
 * @param  string op 查询操作
 * @param  interface{} value 查询值
 * @return *Options
 */
func (o *Options) AddCondition(column, op string, value interface{}) *Options {
	o.Conditions.AddCondition(column, op, value)
	return o
}

/**
 * 添加等值查询条件
 * @receiver *Options
 * @param  string column 查询字段
 * @param  interface{} value 查询值
 * @return *Options
 */
func (o *Options) AddEqCondition(column string, value interface{}) *Options {
	o.Conditions.AddEqCondition(column, value)
	return o
}

/**
 * 设置分组
 * @receiver *Options
 * @param  []string group 分组
 * @return *Options
 */
func (o *Options) WithGroups(group []string) *Options {
	o.Groups = group
	return o
}

/**
 * 设置排序
 * @receiver *Options
 * @param  Orders orders 排序
 * @return *Options
 */
func (o *Options) WithOrders(orders Orders) *Options {
	o.Orders = orders
	return o
}

/**
 * 添加排序
 * @receiver *Options
 * @param  string column 字段
 * @param  string sortBy 顺序
 * @return *Options
 */
func (o *Options) AddOrder(column, sortBy string) *Options {
	o.Orders.AddOrder(column, sortBy)
	return o
}

/**
 * 添加正序排序
 * @receiver *Options
 * @param  string column 字段
 * @return *Options
 */
func (o *Options) AddAscOrder(column string) *Options {
	o.Orders.AddAscOrder(column)
	return o
}

/**
 * 添加倒序排序
 * @receiver *Options
 * @param  string column 字段
 * @return *Options
 */
func (o *Options) AddDescOrder(column string) *Options {
	o.Orders.AddDescOrder(column)
	return o
}

/**
 * 是否设置排序
 * @receiver *Options
 * @return bool
 */
func (o *Options) HasOrder() bool {
	return !o.Orders.IsEmpty()
}

/**
 * 设置分页
 * @receiver *Options
 * @param  int page 页码
 * @param  int pageSize 每页数量
 * @return *Options
 */
func (o *Options) WithPagination(page, pageSize int) *Options {
	o.Page = page
	o.PageSize = pageSize
	return o
}

/**
 * 是否设置分页
 * @receiver *Options
 * @return bool
 */
func (o *Options) HasPagination() bool {
	return o.Page > 0 && o.PageSize > 0
}

/**
 * 设置强制主库
 * @receiver *Options
 * @param  bool master 是否强制主库
 * @return *Options
 */
func (o *Options) WithMaster(master bool) *Options {
	o.Force = master
	return o
}

func NewConditions() Conditions {
	return Conditions{}
}

/**
 * 新建查询条件
 * @param  string column 查询字段
 * @param  string op 查询操作
 * @param  interface{} value 查询值
 * @return Conditions
 */
func NewSingleConditions(column, op string, value interface{}) Conditions {
	conditions := NewConditions()
	conditions.AddCondition(column, op, value)
	return conditions
}

/**
 * 新建等值查询条件
 * @param  string column 查询字段
 * @param  interface{} value 查询值
 * @return Conditions
 */
func NewSingleEqConditions(column string, value interface{}) Conditions {
	conditions := NewConditions()
	conditions.AddEqCondition(column, value)
	return conditions
}

/**
 * 添加查询条件
 * @receiver *Conditions
 * @param  string column 查询字段
 * @param  string op 查询操作
 * @param  interface{} value 查询值
 * @return Conditions
 */
func (c *Conditions) AddCondition(column, op string, value interface{}) Conditions {
	*c = append(*c, Condition{column, op, value})
	return *c
}

/**
 * 添加等值查询条件
 * @receiver *Conditions
 * @param  string column 查询字段
 * @param  interface{} value 查询值
 * @return Conditions
 */
func (c *Conditions) AddEqCondition(column string, value interface{}) Conditions {
	*c = append(*c, Condition{Column: column, Value: value})
	return *c
}

/**
 * 添加排序
 * @receiver *Orders
 * @param  string column 字段
 * @param  string sortBy 顺序
 * @return Orders
 */
func (o *Orders) AddOrder(column, sortBy string) Orders {
	*o = append(*o, Order{column, sortBy})
	return *o
}

/**
 * 添加正序排序
 * @receiver *Orders
 * @param  string column 字段
 * @return Orders
 */
func (o *Orders) AddAscOrder(column string) Orders {
	*o = append(*o, Order{column, SORT_ASC})
	return *o
}

/**
 * 添加倒序排序
 * @receiver *Orders
 * @param  string column 字段
 * @return Orders
 */
func (o *Orders) AddDescOrder(column string) Orders {
	*o = append(*o, Order{column, SORT_DESC})
	return *o
}

/**
 * 是否设置排序
 * @receiver Orders
 * @return bool
 */
func (o Orders) IsEmpty() bool {
	return utils.IsEmpty(o)
}
