# mysql-go

a simple mysql package for golang projects base on xorm,support connection pools and master&slave operation 

## install

```go
go get github.com/scofieldpeng/mysql-go/v3
```

## Usage

First, create a ini file under your `$ProjectPath/config`(or any dir that you like) directory like below:

```ini
# node_* sections are the mysql nodes configurations,the node name is the `*` part,ie node_default,the node name is default
[mysql_node_default]
# mysql dsn driver configuration,more detail see here https://github.com/go-sql-driver/mysql#dsn-data-source-name
dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4

# if the node have slaves nodes,add the node name,each nodes devide by `,`
# note make every slave have configuration in the ini file!!!
slave=slave1,slave2,slave3

# connection pool max connect number, only affect to this node
maxConn=100
# max idle connect number, only affect to this node
maxIdle=5
```

# slave config
[mysql_node_slave1]
dsn=root:123456@tcp(127.0.0.1:3306)/test?charset=utf8mb4

Second, initialize the mysql engine

```go
// initialize the mysql ini file,more document see the package https://github.com/scofieldpeng/config
connConfig := config.Data("mysql")
mysqlConfig := mysql.Config{
    Debug: false,
    MaxIdle:5,
    MaxConn:10,
}
if err := mysql.Init(mysqlConfig,mysqlFileConfig);err != nil {
    fmt.Println(err)
}
```

Now, you can use as you like

```go
// return a xorm engine object,now you can use mysql operation freely
xormEngine := mysql.Select("default").Engine()

// or you can use the master node
xormMasterEngine := mysql.Select("default").Master()

// most of time we always choose the master node, but you want to use slave,you can use the Slave method,
// when you set the slave name param, you can operation the specific slave node, or the system will return a random slave node
xormSlaveEngine := mysql.Select("default").Slave("slave1")
```

if you are not familiar with xorm, you can see the xorm document [http://xorm.io](http://xorm.io)

## Licence

MIT Licence

## Thanks

[https://github.com/go-xorm/xorm](https://github.com/go-xorm/xorm)
