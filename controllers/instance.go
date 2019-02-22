package controllers

import (
	"Athena/models"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strconv"
	"time"
)

type testJson struct {
	Id   int
	Name string
	Age  int
	D    time.Time
}

type InstanceController struct {
	beego.Controller
}

//获取所有实例概述信息
func (this *InstanceController) Get() {
	//models.AddInstance()
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	ret := models.GetInstance()
	this.Data["json"] = ret
	this.ServeJSON()
}

//获取单个实例详情
func (this *InstanceController) GetOneInstanceInfo() {
	var ret map[string]interface{}
	ret = make(map[string]interface{})
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	idValue := this.Ctx.Input.Param(":id")
	id, err := strconv.Atoi(idValue)
	if err != nil {
		ret["error"] = fmt.Sprintf("invalid id, id is:%s", idValue)
		ret["code"] = "1001"
	} else {
		ret = models.GetOneInstanceInfo(id)
	}
	this.Data["json"] = ret
	this.ServeJSON()
}

/*
1. 通过平台创建实例时， 页面传入 ip, instance_num, db_version
2. 将已有实例录入数据库
	node_addr,node_port,base_dir,data_dir,conf_path,instance_version,

	instance_type =2,instance_status=2
	instance_time,create_time,update_time 这三个时间为 now

*/

func (this *InstanceController) AddInstance() {
	this.Ctx.Output.Header("Access-Control-Allow-Origin", "*")
	//fmt.Println(self.Input().Get("a"))
	//abc(self.Input().Get("ip"),self.Input().Get("num"), self.Input().Get("version") )
	j := new(testJson)
	fmt.Println(json.Unmarshal(this.Ctx.Input.RequestBody, &j))
	fmt.Println(j)

	this.Data["json"] = "asdfasfsafasdfasf"
	this.ServeJSON()
}

func abc(ip, instance_num, db_version string) (err error) {
	num, err := strconv.Atoi(instance_num)
	if err != nil {
		return err
	}
	port := 3000
	// 生成多条记录，每一条表示一个实例
	for i := 1; i <= num; i++ {
		node_port := port + i
		base_dir := "/usr/local/mysql"
		data_dir := fmt.Sprintf("/data/mysql/data/mysqldata%d", node_port)
		conf_path := fmt.Sprintf("/etc/my%d.conf", node_port)
		instance_type := 1
		instance_status := 3
		models.AddInstance(ip, node_port, base_dir, data_dir, conf_path, instance_type, instance_status, db_version)
		//fmt.Println(ip, node_port, base_dir, data_dir, conf_path, instance_type, instance_status, db_version)
	}
	return nil
}
