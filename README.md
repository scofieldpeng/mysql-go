# mysql-go

a simple mysql package for golang projects base on xorm,support connection pools and master&slave operation 

## install

```go
go get github.com/scofieldpeng/mysql-go
```

## Usage

First, create a mysql.ini(mysql_debug.ini when debug mode) file under your `$ProjectPath/config` directory like below:

```ini
# connection pool configurations(optional)
[pool]
# the idle connection num in the pool 
idleNum=5

# max open connection in the pool
maxOpenConn=10

# node_* sections are the mysql nodes configurations,the node name is the `*` part,ie node_default,the node name is default
[node_default]
# mysql dsn driver configuration,more detail see here https://github.com/go-sql-driver/mysql#dsn-data-source-name
dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4

# if the node have slaves nodes,add the node name,each nodes devide by `,`
# note make every slave have configuration in the ini file!!!
slave=slave1,slave2,slave3
```

Second, initialize the mysql engine

```go
// initialize the mysql ini file,more document see the package https://github.com/scofieldpeng/config-go
mysqlFileConfig := config.Select("mysql")
if err := mysql.Init(mysqlFileConfig);err != nil {
    log.Fatalln(err)
}
```

Now, you can use as you like

```go
// return a xorm engine object,now you can use mysql operation freely
xormEngine := mysql.Select("default").XormEngine()

// most of time we always choose the master node, but you want to use slave,you can use the Slave method,
// when you set the slave name param, you can operation the specific slave node, or the system will return a random slave node
xormSlaveEngine := mysql.Select("default").Slave("slave1").XormEngine()
```

if you are not familiar with xorm, you can see the xorm document [http://xorm.io](http://xorm.io)

## Licence

MIT Licence

## Thanks

[https://github.com/go-xorm/xorm](https://github.com/go-xorm/xorm)
