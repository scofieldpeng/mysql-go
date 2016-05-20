package mysql

import (
    "testing"
    "github.com/vaughan0/go-ini"
)

func TestInit(t *testing.T) {
    testConfig := ini.File{
        "node_default":ini.Section{
            "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
            "slave":"slave1,slave2",
        },
        "node_slave1":ini.Section{
            "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
        },
        "node_slave2":ini.Section{
            "dsn":"root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4",
        },
    }

    if err := Init(testConfig);err != nil {
        t.Fatal("init fail!error:",err.Error())
    }

    if res,err := Select().XormEngine().Exec("CREATE TABLE `test` (`id` int(11) unsigned NOT NULL AUTO_INCREMENT,`value` int(11) DEFAULT NULL,PRIMARY KEY (`id`)) ENGINE=InnoDB DEFAULT CHARSET=utf8mb4");err != nil {
        t.Fatal("create table fail,error:",err)
    } else {
        if _,err := res.RowsAffected();err != nil {
            t.Error("fetch affected num fail!error:",err)
        }
    }

    type Test struct {
        Id int `xorm:"not null pk autoincr INT(11)"`
        Value int `xorm:"not null INT(11)"`
    }

    if _,err := Select().Slave().XormEngine().Insert(&Test{
        Id:1,
        Value:2,
    });err != nil {
        t.Error("insert fail,error:",err)
    }

    ts := Test{
        Id:1,
    }

    if exist,err := Select().Slave("slave2").XormEngine().Get(&ts);err != nil {
        t.Error("get data fail!error:",err)
    } else if !exist {
        t.Error("data not exsit!")
    }

}
