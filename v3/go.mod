module github.com/scofieldpeng/mysql-go/v3

go 1.13

require (
	github.com/go-sql-driver/mysql v1.4.1
	github.com/go-xorm/core v0.7.2
	github.com/go-xorm/xorm v0.7.9
	github.com/scofieldpeng/config-go/v3 v3.0.0
	github.com/vaughan0/go-ini v0.0.0-20130923145212-a98ad7ee00ec
)

replace github.com/go-xorm/core => xorm.io/core v0.7.2

replace github.com/go-xorm/xorm => xorm.io/xorm v0.7.9
