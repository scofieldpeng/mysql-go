package mysql

import (
	"strings"
)

// WhereBuilder用来构建xorm的engine所需的where模块
// 使用方法很easy
// mysql.NewWhereBuilder()
// mysql.Add("id",1)
// 然后使用的时候直接mysql.Select()
type WhereBuilder struct {
	data map[string]interface{}
}

func NewWhereBuilder(args ...map[string]interface{}) *WhereBuilder {
	wb := &WhereBuilder{
		data: make(map[string]interface{}),
	}
	if len(args) > 0 {
		wb.data = args[0]
	}

	return wb
}

// 添加参数,需要注意的是，原来的xorm中where条件是id>? AND id
func (wb *WhereBuilder) Add(condition string, args interface{}) *WhereBuilder {
	wb.data[condition] = args
	return wb
}

// Encode用来生成xorm的engine.Where()条件的两个参数
func (wb *WhereBuilder) Encode() (whereStr string, beans []interface{}) {
	whereArr := make([]string, 0, len(wb.data))
	beans = make([]interface{}, 0, len(wb.data))
	for k, v := range wb.data {
		whereArr = append(whereArr, k+" ?")
		beans = append(beans, v)
	}

	return strings.Join(whereArr, " AND "), beans
}
