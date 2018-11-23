package mysql

import (
	"errors"
	"github.com/go-xorm/xorm"
	"strings"
)

type Table interface {
	// 获取表节点
	TableNode() string
}

type TableFactory struct {
	tableNode string      `json:"-" xorm:"-" xml:"-"`
	tableObj  interface{} `json:"-" xorm:"-" xml:"-"`
}

func NewTableFactory() *TableFactory {
	tf := &TableFactory{}
	tf.tableNode = "default"
	tf.tableObj = tf

	return tf
}

func (tf TableFactory) checkNode() error {
	if tf.tableNode == "" {
		return errors.New("table_node not set")
	}

	return nil
}

func (tf TableFactory) checkObj() error {
	if tf.tableObj == nil {
		return errors.New("pls use NewXXX function to create a table struct")
	}
	return nil
}

// 校验
func (tf TableFactory) check() error {
	if err := tf.checkNode(); err != nil {
		return err
	}
	if err := tf.checkObj(); err != nil {
		return err
	}

	return nil
}

func (tf TableFactory) TableNode() string {
	return tf.tableNode
}

func (tf *TableFactory) Insert() (affectRows int64, err error) {
	if err := tf.check(); err != nil {

	}
	affectRows, err = Select(tf.tableNode).Master().Insert(tf.tableObj)
	return
}

// 删除表结构
func (tf *TableFactory) Delete() (int64, error) {
	return Select(tf.tableNode).Master().Delete(tf.tableObj)
}

// 查询一条数据
func (tf *TableFactory) Get(fromMaster ...bool) (bool, error) {
	if len(fromMaster) > 0 && fromMaster[0] {
		return Select(tf.tableNode).Master().Get(tf.tableObj)
	}

	return Select(tf.tableNode).Slave().Get(tf.tableObj)
}

// 查询列表
func (tf *TableFactory) Find(where string, fields string, fromMaster ...bool) (res []Table, err error) {
	res = make([]Table, 0)
	if where == "" {
		return res, errors.New("where required")
	}

	master := false
	if len(fromMaster) > 0 && fromMaster[0] {
		master = true
	}

	var (
		engine *xorm.Engine
	)

	if master {
		engine = Select(tf.tableNode).Master()
	} else {
		engine = Select(tf.tableNode).Slave()
	}

	engine.Where(where)
	if len(fields) == 0 {
		engine.AllCols()
	} else {
		engine.Cols(strings.Split(fields,",")...)
	}


}

// 更新表
func (tf *TableFactory) Update(where string, whereData []interface{}, updateFields ...interface{}) (int64, error) {}

// 单独定制一个
func (tf *TableFactory) Exec(sql string, fromMaster ...bool) (int64, error) {}
