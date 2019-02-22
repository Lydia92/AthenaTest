package controllers

import (
	"Athena/models"
	"github.com/astaxie/beego"
)

type InitController struct {
	beego.Controller
}

func (this *InitController) Get() {
	models.InitInfo()
	this.ServeJSON()
}
