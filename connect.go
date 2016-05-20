package mysql

import (
	"github.com/go-xorm/xorm"
	"math/rand"
	"github.com/vaughan0/go-ini"
	_ "github.com/go-sql-driver/mysql"
	"github.com/scofieldpeng/config-go"
	"strings"
	"errors"
)

// engine mysql的连接对象
type engine struct {
	xormEngine *xorm.Engine // 当前的xorm链接engine
	slaveNamesMap     map[string]string // 该engine下属的slave的名称map,key和value都为slave的名称
	slaveNames []string // 该engine下属的slave的名称slice
}

const(
	DefaultIdleNum = 5 // 默认的连接池空闲数大小
	DefaultMaxOpenConn = 10 // 默认的最大打开连接数
)

var (
	engines map[string]*engine // 所有的节点配置项
)

// Slave 获取slave节点,如果不指定slave名称,则随机返回slave中的一个节点
func (e *engine) Slave(slave ...string) *engine {
	selectSlave := false
	if len(slave) != 0 {
		selectSlave = true
	}

	if !selectSlave {
		slave = make([]string, 1)
		slave[0] = e.slaveNames[rand.Intn(len(e.slaveNamesMap))]
	}

	return engines[e.slaveNamesMap[slave[0]]]
}

// XormEngine 返回该节点的xorm的engine对象
func (e *engine) XormEngine() *xorm.Engine {
	return e.xormEngine
}

// Select 选择某个节点,如果没有选择节点名称,默认选择default节点
func Select(node ...string) *engine {
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = "default"
	}

	return engines[node[0]]
}

// Init 初始化mysql,传入ini.File类型的值,将会解析所有的配置项
func Init(iniConfigs ini.File) error {
	engines = make(map[string]*engine)

	idleNum:= config.Int(iniConfigs.Get("pool","idleNum"))
	if idleNum < 1 {
		idleNum = DefaultIdleNum
	}

	maxOpenConn := config.Int(iniConfigs.Get("pool","maxOpenConn"))
	if maxOpenConn < 1 {
		maxOpenConn = DefaultMaxOpenConn
	}

	var (
		err error
		xormEngine *xorm.Engine
		e *engine
		findDsn bool =false
	)

	// 遍历,初始化配置项
	for sectionName,section := range iniConfigs {
		if strings.Index(sectionName,"node_") == 0 {
			findDsn = false
			for confName,conf := range section {
				if confName == "dsn" {
					xormEngine,err = xorm.NewEngine("mysql",conf)
					if err != nil {
						continue
					}

					xormEngine.SetMaxIdleConns(idleNum)
					xormEngine.SetMaxOpenConns(maxOpenConn)

					e = &engine{
						xormEngine:xormEngine,
						slaveNames:make([]string,0),
						slaveNamesMap:make(map[string]string),
					}
					findDsn = true
				}
				if confName == "slave" {
					slaveSlice := strings.Split(conf,",")
					for _,slaveName := range slaveSlice {
						e.slaveNamesMap[slaveName] = slaveName
						e.slaveNames = append(e.slaveNames,slaveName)
					}
				}
			}
			if findDsn {
				nodeNameSlice := string([]byte(sectionName)[5:])
				engines[nodeNameSlice] = e
			} else {
				return errors.New("节点[" + sectionName + "]没有找到dsn配置项")
			}
		}
	}

	return nil
}