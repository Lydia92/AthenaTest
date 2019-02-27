package controllers

import (
	"Athena/models"
	"github.com/astaxie/beego"
)

type GoalController struct {
	beego.Controller
}

//添加元数据信息
func (this *GoalController) AddAllMetaInfo() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	models.InertAllInfo()
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()
}
