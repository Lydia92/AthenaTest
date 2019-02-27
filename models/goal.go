package models

import (
	"fmt"
	"github.com/astaxie/beego/orm"
)

type OrmEng struct {
	SourceDB orm.Ormer
	Goal     orm.Ormer
}

/*
func (this *OrmEng) InitDB() {

	db_account_info := GetAccountForDBA()

	for _, v := range db_account_info {
		dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", v.Account, v.Passwd, v.NodeAddr, v.NodePort, "mysql")
		db := fmt.Sprintf("%s_%s_%d", "mysql", v.NodeAddr, v.NodePort)
		fmt.Println("--------------", v.Passwd)
		//注册数据库，连接对应的生产环境
		orm.RegisterDataBase(db, "mysql", dns)
	}
}
*/
func (this *OrmEng) Compare() {
	this.SourceDB = orm.NewOrm()
	db_account_info := GetAccountForDBA()
	for _, v := range db_account_info {
		db := fmt.Sprintf("%s_%s_%d", "mysql", v.NodeAddr, v.NodePort)
		o := orm.NewOrm()
		o.Using(db)
		this.Goal = o
		this.ComparAccoutInfo()
		this.CompareMetaTableInfo()
		this.CompareMetaColumnsInfo()
		this.CompareIndexInfo()
	}

}

func InertAllInfo() {
	ormeng := OrmEng{}
	//ormeng.InitDB()
	ormeng.Compare()
}
