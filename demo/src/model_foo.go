package demo

import (
	"github.com/EvisuXiao/andrews-common/database"
)

type Foo struct {
	ModelInsertable
	Uid     int    `json:"uid"`
	Content string `json:"content"`
}

type Foos []*Foo

var fooModel = &Foo{}

func init() {
	database.RegisterModel(fooModel)
}

func NewFooModel() *Foo {
	return fooModel
}

func (m *Foo) TableName() string {
	return "foo"
}

func (m *Foo) GetRows(options *database.Options) (Foos, error) {
	var rows Foos
	err := m.GetAnyRows(options, &rows)
	return rows, err
}