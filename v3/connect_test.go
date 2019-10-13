package mysql

import (
	"github.com/vaughan0/go-ini"
	"testing"
)

func TestInit(t *testing.T) {
	testConfig := ini.File{
		"mysql_node_default": ini.Section{
			"dsn":   "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
			"slave": "slave1,slave2",
		},
		"mysql_node_slave1": ini.Section{
			"dsn": "root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
		},
		"mysql_node_slave2": ini.Section{
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
	
	//if res, err := Select().Engine().Exec("CREATE TABLE `test` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT,`value` int(11) DEFAULT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4"); err != nil {
	//	t.Fatal("create table fail,error:", err)
	//} else {
	//	if _, err := res.RowsAffected(); err != nil {
	//		t.Error("fetch affected num fail!error:", err)
	//	}
	//}
	
	type Test struct {
		Id    int    `xorm:"not null pk autoincr INT(11)"`
		Value string `xorm:"not null VARCHAR(11) 'name'"`
	}
	
	if _, err := Select().Slave().Insert(&Test{
		//Id:    1,
		Value: "scofield",
	}); err != nil {
		t.Error("insert fail,error:", err)
	}
	
	ts := Test{
		Id: 1,
	}
	
	if exist, err := Select().Slave("slave2").Get(&ts); err != nil {
		t.Error("get data fail!error:", err)
	} else if !exist {
		t.Error("data not exsit!")
	}
}
