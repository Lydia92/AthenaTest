package controllers

import (
	"Athena/models"
	"github.com/astaxie/beego"
	"strconv"
)

type SlowSqlController struct {
	beego.Controller
}

//传入集群名，获取集群中的慢sql
func (this *SlowSqlController) GetSlowSqlByClustername() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	clusterName := this.Ctx.Input.Param(":clusterName")
	ret := models.GetOneCluserSlowSql(clusterName)
	this.Data["json"] = ret
	this.ServeJSON()
}

//传入某个慢sql的id，获取慢sql中表的索引信息
func (this *SlowSqlController) GetSlowById() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	idValue := this.Ctx.Input.Param(":id")
	id, _ := strconv.Atoi(idValue)
	//ret := models.GetCluserSlowSql()
	ret := models.GetSlowSqlTableInfoById(id)
	this.Data["json"] = ret
	this.ServeJSON()
}
