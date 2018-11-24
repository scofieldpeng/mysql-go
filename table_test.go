package mysql

import (
	"fmt"
	"github.com/vaughan0/go-ini"
	"testing"
)

type testTable struct {
	TableFactory `xorm:"-"`
	Id           int    `xorm:"'id' not null pk autoincr INT(11)"`
	Name         string `xorm:"'name' not null VARCHAR(50)"`
}

func (t *testTable) self() interface{} {
	return t
}

func (t testTable) TableName() string {
	return "test"
}

func (t testTable) TableNode() string {
	return t.tableNode
}

func newTestTable() *testTable {
	t := &testTable{}
	t.tableNode = "default"
	t.myself = t.self
	return t
}

func TestTableFactory(t *testing.T) {
	testConfig := ini.File{
		"mysql_node_default": ini.Section{
			"dsn": "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
		},
	}
	mysqlConfig := Config{
		Debug:   true,
		MaxIdle: 5,
		MaxConn: 10,
	}

	if err := Init(mysqlConfig, testConfig); err != nil {
		t.Fatal("init fail!error:", err.Error())
	}

	tt := newTestTable()
	tt.Name = "scofield"
	if _, err := tt.Insert(); err != nil {
		t.Error(err.Error())
	}
	insertId := tt.Id

	tt = newTestTable()
	tt.Id = insertId
	if _, err := tt.Get(); err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("%v\n", tt)

	tt = newTestTable()
	list, err := tt.Find(NewWhereBuilder(map[string]interface{}{"id>": 1}), "", "id ASC")
	if err != nil {
		t.Error(err.Error())
	}
	fmt.Printf("res: %v\n", list)
}
