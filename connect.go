/*
 mysql-go is a mysql library based on xorm[http://xorm.io]

config

```ini
[mysql_node_default]
dsn=
slave=abc,abc1,abc2,abc3
# 最大空闲数量
maxIdle=10
# 最大连接数量
maxConn=100
```

usage:

```go
// 获取主
engine := mysql.Select().Master()

// 获取主
engine := mysql.Select().Engine()

// 获取slave节点
engine := mysql.Select("default").Slave()
```
*/
package mysql

import (
	"errors"
	"fmt"
	_ "github.com/go-sql-driver/mysql"
	"github.com/go-xorm/core"
	"github.com/go-xorm/xorm"
	appConfig "github.com/scofieldpeng/config"
	"github.com/vaughan0/go-ini"
	"io"
	"math/rand"
	"os"
	"strings"
	"time"
)

// engine mysql的连接对象
type (
	engine struct {
		xormEngine    *xorm.Engine // 当前的xorm链接engine
		slaveNames    []string
		slaveNamesMap map[string]bool // 该engine下属的slave的名称map,key和value都为slave的名称
		// 最大空闲数量
		maxIdle int
		// 最大连接数量
		maxConn int
	}
	Config struct {
		Debug bool
		// 最大空闲数量
		MaxIdle int
		// 最大连接数量
		MaxConn int
		// log writer
		LogWriter io.Writer
		// log prefix
		LogPrefix string
		// log flag
		LogFlag int
		// log level
		LogLevel core.LogLevel
	}
)

const (
	DefaultIdleNum     = 5  // 默认的连接池空闲数大小
	DefaultMaxOpenConn = 10 // 默认的最大打开连接数

	DEFAULT_CONN_PREFIX = "mysql_node_"
)

var (
	engines map[string]*engine // 所有的节点配置项
	config  = Config{
		Debug:     false,
		MaxIdle:   DefaultIdleNum,
		MaxConn:   DefaultMaxOpenConn,
		LogWriter: os.Stdout,
		LogPrefix: xorm.DEFAULT_LOG_PREFIX,
		LogFlag:   xorm.DEFAULT_LOG_FLAG,
		LogLevel:  xorm.DEFAULT_LOG_LEVEL,
	}

	// Engine没有找到
	ErrEngineNotFound = errors.New("not found engine")
	// mysql操作影响了0行
	ErrAffectedZeroRow = errors.New("affect 0 rows")
)

// 设置config，如果传入的值存在，则覆盖默认配置项目
func (c *Config) Set(config Config) {
	c.Debug = config.Debug
	if config.MaxIdle > 0 {
		c.MaxIdle = config.MaxIdle
	}
	if config.MaxConn > 0 {
		c.MaxConn = config.MaxConn
	}
	if config.LogWriter != nil {
		c.LogWriter = config.LogWriter
	}
	if int(config.LogLevel) > 0 {
		c.LogLevel = config.LogLevel
	}
	if int(config.LogFlag) > 0 {
		c.LogFlag = config.LogFlag
	}
}

// Slave 获取slave节点,如果不指定slave名称,则随机返回slave中的一个节点
func (e *engine) Slave(slave ...string) *xorm.Engine {
	selectSlave := false
	if len(slave) != 0 {
		selectSlave = true
	}
	if !selectSlave && len(e.slaveNames) > 0 {
		slave = make([]string, 1)
		rand.Seed(time.Now().UnixNano())
		slave[0] = e.slaveNames[rand.Intn(len(e.slaveNamesMap))]
		
		return engines[slave[0]].xormEngine
	}
	
	return e.Engine()
}

// 获取主节点
func (e *engine) Master() *xorm.Engine {
	return e.xormEngine
}

func (e *engine) Engine() *xorm.Engine {
	return e.Master()
}

// 选择某个节点,如果没有选择节点名称,默认选择default节点
// BUG: 当用户没有设置default节点配置时如果又通过default来获取，可能会获取不到结果
func Select(node ...string) (e *engine) {
	var exist bool
	if len(node) == 0 {
		node = make([]string, 1)
		node[0] = "default"
	}

	// TODO 这里如果连default都没有，那么会抛出异常，暂时不考虑这种情况，否则返回值就得加上
	if e, exist = engines[node[0]]; !exist {
		return nil
	}

	return e
}

// 初始化
func Init(mysqlConfig Config, connConfig ini.File) error {
	var (
		connNameMap = make(map[string]bool)
		findDefault = false
	)

	config.Set(mysqlConfig)
	engines = make(map[string]*engine)

	for k, _ := range connConfig {
		if strings.Index(k, DEFAULT_CONN_PREFIX) == 0 {
			sectionName := string([]byte(k)[len(DEFAULT_CONN_PREFIX):])
			if len(sectionName) > 0 && appConfig.String(connConfig.Get(k, "dsn")) != "" {
				connNameMap[sectionName] = true
			}
		}
	}
	if len(connConfig) == 0 {
		return errors.New("not found any mysql connect config")
	}

	// 遍历,初始化配置项
	for k, _ := range connNameMap {
		if k == "default" {
			findDefault = true
		}
		mysqlEngine := &engine{slaveNames: make([]string, 0), slaveNamesMap: make(map[string]bool), maxIdle: config.MaxIdle, maxConn: config.MaxConn}
		if maxConn := appConfig.Int(connConfig.Get(DEFAULT_CONN_PREFIX+k, "maxConn")); maxConn > 0 {
			mysqlEngine.maxConn = maxConn
		}
		if maxIdle := appConfig.Int(connConfig.Get(DEFAULT_CONN_PREFIX+k, "maxIdle")); maxIdle > 0 {
			mysqlEngine.maxIdle = maxIdle
		}

		// 初始化xorm.Engine
		xormEngine, err := xorm.NewEngine("mysql", appConfig.String(connConfig.Get(DEFAULT_CONN_PREFIX+k, "dsn")))
		if err != nil {
			return errors.New(fmt.Sprintf("found invalid mysql dsn,key:%s", DEFAULT_CONN_PREFIX+k))
		}
		xormEngine.ShowSQL(config.Debug)
		xormEngine.SetMaxIdleConns(mysqlEngine.maxIdle)
		xormEngine.SetMaxOpenConns(mysqlEngine.maxConn)
		xormEngine.SetLogger(xorm.NewSimpleLogger3(config.LogWriter, config.LogPrefix, config.LogFlag, config.LogLevel))

		mysqlEngine.xormEngine = xormEngine

		// 遍历slave
		if slaves := appConfig.String(connConfig.Get(DEFAULT_CONN_PREFIX+k, "slave")); slaves != "" {
			slaveKeys := strings.Split(slaves, ",")
			for _, v := range slaveKeys {
				if _, exist := connNameMap[v]; exist {
					mysqlEngine.slaveNamesMap[v] = true
					mysqlEngine.slaveNames = append(mysqlEngine.slaveNames, v)
				}
			}
		}

		engines[k] = mysqlEngine
	}

	if !findDefault {
		fmt.Println("[warning]config中没有找到default的配置项，请勿使用mysql.Select()或者mysql.Select(\"default\")")
	}

	return nil
}
