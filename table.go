package mysql

import (
	"errors"
	"fmt"
	"github.com/go-xorm/xorm"
	"strings"
)

type Table interface {
	SelfObj() *xorm.TableName
}

type TableFactory struct {
	tableNode string `json:"-" xorm:"-" xml:"-"`
	myself    func() interface{}
}

func (tf TableFactory) checkNode() error {
	if tf.tableNode == "" {
		return errors.New("table_node not set")
	}

	return nil
}

// 校验
func (tf TableFactory) check() error {
	if err := tf.checkNode(); err != nil {
		return err
	}

	return nil
}

func (tf TableFactory) TableNode() string {
	return tf.tableNode
}

func (tf TableFactory) Myself() interface{} {
	return tf.myself()
}

func (tf *TableFactory) Insert() (affectRows int64, err error) {
	if err := tf.check(); err != nil {
		return affectRows, err
	}

	affectRows, err = Select(tf.TableNode()).Master().Insert(tf.Myself())
	return
}

// 删除表结构
func (tf *TableFactory) Delete() (int64, error) {
	return Select(tf.TableNode()).Master().Delete(tf.Myself())
}

// 查询一条数据
func (tf *TableFactory) Get(fromMaster ...bool) (bool, error) {
	if len(fromMaster) > 0 && fromMaster[0] {
		return Select(tf.TableNode()).Master().Get(tf.Myself())
	}

	return Select(tf.TableNode()).Slave().Get(tf.Myself())
}

// 查询列表
func (tf *TableFactory) Find(whereBuilder *WhereBuilder, fields string, orderBy string, fromMaster ...bool) (res []interface{}, err error) {
	var (
		session   *xorm.Session
		master        = false
		whereData  = ""
		whereArgs []interface{}
	)

	whereData, whereArgs = whereBuilder.Encode()
	if len(whereData) == 0 {
		return res, errors.New("where condition required")
	}
	fmt.Printf("where: %s\n",whereData)
	fmt.Printf("args: %v\n",whereArgs)

	if len(fromMaster) > 0 && fromMaster[0] {
		master = true
	}

	if master {
		session = Select(tf.TableNode()).Master().NewSession()
	} else {
		session = Select(tf.TableNode()).Slave().NewSession()
	}

	if len(fields) == 0 {
		session.AllCols()
	} else {
		session.Cols(strings.Split(fields, ",")...)
	}

	session.Where(whereData, whereArgs...)

	if orderBy != "" {
		session.OrderBy(orderBy)
	}

	list := make([]interface{}, 0)
	if err := session.Table(tf.Myself()).Find(&list); err != nil {
		return res, err
	}

	return list, err
}

// 更新表
func (tf *TableFactory) Update(where string, whereData []interface{}, updateFields ...string) (int64, error) {
	if err := tf.check(); err != nil {
		return 0, err
	}

	if where == "" {
		return 0, errors.New("where condition required")
	}
	if len(whereData) == 0 {
		return 0, errors.New("whereData required")
	}

	session := Select(tf.TableNode()).Master().NewSession()
	session.Where(where, whereData...)
	if len(updateFields) == 0 {
		session.AllCols()
	} else {
		session.Cols(updateFields...)
	}

	return session.Update(tf.myself())
}
