package controllers

import (
	"Athena/models"
	"github.com/astaxie/beego"
)

type MetaController struct {
	beego.Controller
}
//添加元数据信息
func (this *MetaController) AddMetadata(){
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	models.CompareMetaTable()
	models.CompareMetaColumns()
	models.CompareIndex()
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()
}
//添加元表信息
/*func (this *MetaController) AddMetadataTable() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	models.CompareMetaTable()
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()

}

//添加元表 列信息
func (this *MetaController) AddMetadataColumn() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	models.CompareMetaColumns()
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()

}

//添加元表 索引信息
func (this *MetaController) AddMetadataIndex() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	models.CompareIndex()
	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()

}*/

//根据集群名获取实例信息
func (this *MetaController) GetHostByCluster() {

	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	clusterName := this.Ctx.Input.Param(":clusterName")
	ret := models.GetHostByCluster(clusterName)
	this.Data["json"] = ret
	this.ServeJSON()
}

//根据host和端口获取库信息
func (this *MetaController) GetDatabaseByHost() {

	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	clusterName := this.GetString("host")
	clusterPort, _ := this.GetInt("clusterPort")
	ret := models.GetDatabaseByHost(clusterName, clusterPort)
	this.Data["json"] = ret
	this.ServeJSON()
}

//根据host等值查询表的详细信息
func (this *MetaController) GetTableInfoBySchema() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	host := this.GetString("host")
	schema := this.GetString("schema")
	tableName := this.GetString("tableName")
	clusterPort, _ := this.GetInt("clusterPort")
	ret := models.GetTableInfo(host, schema, tableName, clusterPort)
	this.Data["json"] = ret
	this.ServeJSON()
}
