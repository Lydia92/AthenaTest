package controllers

import (
	"Athena/models"
	"github.com/astaxie/beego"
)

type DBAccountController struct {
	beego.Controller
}

//获取所有的用户信息
func (this *DBAccountController) Geta() {

	user := models.GetUserHost()
	this.Data["user"] = user
	models.ComparAccout()
	//fmt.Println(models.GetAccountForDBA())
	this.TplName = "hello.tpl"

}

//添加用户信息到cmdb
func (this *DBAccountController) AddAccountInfo() {
	node_addr := "192.168.160.132"
	node_port := 3388
	ownership := "dba"
	account := "dba"
	passwd := "123456"
	priv := "all"
	db_name := "aa"
	table_name := "aa"
	login_addr := "%"
	md5 := "aaaa"
	//node_addr, ownership, account, passwd, priv, db_name, table_name, login_addr,	md5 string, node_port int
	models.InserAccountInfo(node_addr, ownership, account, passwd, priv, db_name, table_name, login_addr, md5, uint16(node_port))
	this.TplName = "hello1.tpl"
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()
}

//根据集群名获取集群内的所有账户信息
func (this *DBAccountController) Get() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	clusterName := this.Ctx.Input.Param(":clusterName")
	ret := models.GetAccountInfo(clusterName)
	this.Data["json"] = ret
	this.ServeJSON()
}

func (this *DBAccountController) GetPassword() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	id := this.Ctx.Input.Param(":id")

	ret := models.GetPassword(id)
	this.Data["json"] = ret
	this.ServeJSON()
}
