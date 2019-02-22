package models

import "github.com/astaxie/beego/orm"

//获取所有集群信息
type clusterInfo struct {
	Id           string
	ClusterName  string
	WriteVip     string
	WritePort    string
	ReadVip      string
	ReadPort     string
	ClusterType  string
	BusinessName string
	Leader       string
}

//获取所有集群信息
func GetClusterInfo() map[string]interface{} {
	resultJson := make(map[string]interface{}, 5)
	resultJson["code"] = 1000
	var clusterAllRet []clusterInfo
	o := orm.NewOrm()
	num, err := o.Raw("select c.id, c.cluster_name, c.write_vip, c.write_port, c.read_vip, c.read_port," +
		" c.cluster_type, b.business_name, b.leader " +
		"from business_cluster b, cluster_info c where b.cluster_info_id = c.id").QueryRows(&clusterAllRet)
	if err != nil {
		resultJson["error"] = err.Error()
		resultJson["code"] = 1001
		return resultJson
	}

	if num == 0 {
		resultJson["code"] = 1002
		return resultJson
	}

	resultJson["data"] = clusterAllRet
	return resultJson
}

//获取集群详情
type clusterOneInfo struct {
	Id             string
	NodeAddr       string
	NodePort       string
	InstanceType   string
	InstanceStatus string
	InstanceTime   string
	ClusterName    string
	BusinessName   string
	Leader         string
	Role           string
	WriteVip       string
	WritePort      string
	ReadVip        string
	ReadPort       string
}

//获取集群详情
func GetOneClusterInfo(id int) map[string]interface{} {
	resultJson := make(map[string]interface{}, 5)
	resultJson["code"] = 1000
	var clusterOneRet []clusterOneInfo
	o := orm.NewOrm()
	num, err := o.Raw("select a.id, a.node_addr, a.node_port, a.instance_type, a.instance_status, d.role, "+
		"a.instance_time, c.cluster_name, b.business_name, b.leader, c.write_vip, c.write_port, c.read_vip, c.read_port "+
		"from instance_info a, business_cluster b, cluster_info c, cluster_instance d "+
		"where a.id = d.instance_id and c.id = d.cluster_info_id "+
		"and b.cluster_info_id = c.id and c.id=?", id).QueryRows(&clusterOneRet)

	if err != nil {
		resultJson["error"] = err.Error()
		resultJson["code"] = 1001
		return resultJson
	}
	if num == 0 {
		resultJson["code"] = 1002
		return resultJson
	}
	resultJson["data"] = clusterOneRet
	return resultJson
}

//获取集群名
func GetClusterName() map[string]interface{} {
	resultJson := make(map[string]interface{}, 5)
	resultJson["code"] = 1000
	var clusterName []string
	o := orm.NewOrm()
	num, err := o.Raw("select cluster_name from cluster_info").QueryRows(&clusterName)
	if err != nil {
		resultJson["error"] = err.Error()
		resultJson["code"] = 1001
		return resultJson
	}
	if num == 0 {
		resultJson["code"] = 1002
		return resultJson
	}
	resultJson["data"] = clusterName
	return resultJson
}
