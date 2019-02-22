package controllers

import (
	"Athena/models"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
)

type ClusterController struct {
	beego.Controller
}

//获取所有集群列表
func (this *ClusterController) Get() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ret := models.GetClusterInfo()
	this.Data["json"] = ret
	this.ServeJSON()
}

//获取集群详情
func (this *ClusterController) GetOneClusterInfo() {
	var ret map[string]interface{}
	ret = make(map[string]interface{})
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	idValue := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		ret["error"] = fmt.Sprintf("invalid id, id is:%s", idValue)
		ret["code"] = "1001"
	} else {
		ret = models.GetOneClusterInfo(id)
	}
	this.Data["json"] = ret
	this.ServeJSON()
}

//获取集群名
func (this *ClusterController) GetClusterName() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ret := models.GetClusterName()
	this.Data["json"] = ret
	this.ServeJSON()
}
