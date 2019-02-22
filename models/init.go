package models

import (
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
)

func init() {
	// "username:password@tcp(127.0.0.1:3306)/db_name?charset=utf8"

	mysqlHost := beego.AppConfig.String("mysql::host")
	mysqlPort := beego.AppConfig.String("mysql::port")
	mysqlUser := beego.AppConfig.String("mysql::user")
	mysqlPass := beego.AppConfig.String("mysql::password")
	mysqlDb := beego.AppConfig.String("mysql::database")
	mysqlChar := beego.AppConfig.String("mysql::charset")

	maxIdleConns, err := beego.AppConfig.Int("mysql::MaxIdleConns")
	if err != nil {
		beego.BeeLogger.Warning("%s error message:[%s]",
			"MaxIdleConns parameter get error, set value is 10", err)
		maxIdleConns = 10
	}

	maxOpenConns, err := beego.AppConfig.Int("mysql::MaxOpenConns")
	if err != nil {
		beego.BeeLogger.Warning("%s error message:[%s]",
			"MaxOpenConns parameter get error, set value is 30", err)
		maxOpenConns = 30
	}

	dsn := fmt.Sprintf("%s:%s@tcp(%s:%s)/%s?%s&parseTime=true&loc=Local",
		mysqlUser, mysqlPass, mysqlHost, mysqlPort, mysqlDb, mysqlChar)

	orm.RegisterDataBase("default", "mysql", dsn, maxIdleConns, maxOpenConns)
	// register model
	orm.RegisterModel(
		new(InstanceInfo),
		new(ClusterInfo),
		new(ClusterInstance),
		new(BusinessCluster),
		new(DbAccountInfo),
		new(HostInfo),
		new(HostAccountInfo),
		new(HostAppType),
		new(AppType),
		new(MysqlMetadataTables),
		new(MysqlMetadataColumns),
		new(MysqlMetadataIndexs),
		new(MysqlSlowQueryReviewHistory),
		new(MysqlSlowQueryReview),
	)

	// create table
	orm.RunSyncdb("default", false, true)
}
