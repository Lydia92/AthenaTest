package models

import (
	"Athena/util/numberutil"
	"github.com/astaxie/beego/orm"
)

func AddInstance(ip string, node_port int, base_dir, data_dir, conf_path string, instance_type,
	instance_status int, db_version string) {
	o := orm.NewOrm()
	i := InstanceInfo{
		NodeAddr:        ip,
		NodePort:        numberutil.IntToUint16(node_port),
		BaseDir:         base_dir,
		DataDir:         data_dir,
		ConfPath:        conf_path,
		InstanceType:    numberutil.IntToUint8(instance_type),
		InstanceStatus:  numberutil.IntToUint8(instance_status),
		InstanceVersion: db_version,
	}

	o.Insert(&i)

}

type instanceInfo struct {
	Id             string
	NodeAddr       string
	NodePort       string
	InstanceType   string
	InstanceStatus string
	InstanceTime   string
	ClusterName    string
	BusinessName   string
	Leader         string
}

//TODO:添加分页功能， 目前是将所有结果全部查出返回
//获取所有实例信息
func GetInstance() map[string]interface{} {
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	var instanceAllRet []instanceInfo
	o := orm.NewOrm()
	num, err := o.Raw("select a.id, a.node_addr, a.node_port, a.instance_type, a.instance_status, " +
		"a.instance_time, c.cluster_name, b.business_name, b.leader " +
		"from instance_info a, business_cluster b, cluster_info c, cluster_instance d " +
		"where a.id = d.instance_id and c.id = d.cluster_info_id " +
		"and b.cluster_info_id = c.id").QueryRows(&instanceAllRet)
	if err != nil {
		resultJson["error"] = err.Error()
		resultJson["code"] = 1001
		return resultJson
	}
	//无结果
	if num == 0 {
		resultJson["code"] = 1002
		return resultJson
	}

	resultJson["data"] = instanceAllRet
	return resultJson
}

//返回单个实例详情
func GetOneInstanceInfo(id int) map[string]interface{} {
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	var instanceOneRet = InstanceInfo{Id: uint64(id)}
	o := orm.NewOrm()
	err := o.Read(&instanceOneRet, "id")
	if err != nil {
		resultJson["error"] = "no row found"
		resultJson["code"] = 1001
		return resultJson
	}

	resultJson["data"] = instanceOneRet
	return resultJson
}
