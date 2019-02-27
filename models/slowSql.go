package models

import (
	"Athena/util/spilt"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
	"time"
)

//集群慢sql信息
type SlowSqlInfo struct {
	Id             int       `json:"id"`
	Node_addr      string    `json:"node_addr"`
	Node_port      uint16    `json:"node_port"`
	Cluster_name   string    `json:"cluster_name"`
	User_max       string    `json:"user"`
	Db_max         string    `json:"db"`
	Sample         string    `json:"sample"`
	TsCnt          float64   `json:"ts_cnt"`
	First_seen     time.Time `json:"first_seen"`
	Last_seen      time.Time `json:"last_seen"`
	Query_time_sum float64   `json:"query_time_sum"`
	Query_time_min float64   `json:"query_time_min"`
	Query_time_max float64   `json:"query_time_max"`
	Lock_time_sum  float64   `json:"lock_time_sum"`
}

//慢查询中 sql，db，host 及port信息
type SampleInfo struct {
	Sample      string
	ServeridMax string
	DbMax       string
}

//查找某个集群的慢sql
func GetOneCluserSlowSql(cluserName string) map[string]interface{} {
	var slowSqlInfo []SlowSqlInfo
	DB1 := orm.NewOrm()
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000

	num, err := DB1.Raw("select a.id, d.node_addr	,d.node_port,b.cluster_name	,a.user_max,a.db_max,a.sample,"+
		"a.ts_cnt, e.first_seen, e.last_seen, a.query_time_sum, a.query_time_min, a.query_time_max,"+
		"a.lock_time_sum "+
		"from mysql_slow_query_review_history a "+
		"join mysql_slow_query_review e ON a.CHECKSUM = e.CHECKSUM "+
		"join instance_info d on a.serverid_max=concat(d.node_addr,'-',d.node_port) "+
		"join cluster_instance c on c.instance_id=d.id "+
		"join cluster_info b  on c.cluster_info_id=b.id "+
		"where b.cluster_name=? order by ts_cnt desc, query_time_max desc", cluserName).QueryRows(&slowSqlInfo)
	if err != nil || num == 0 {
		resultJson["error"] = "no row found"
		resultJson["code"] = 1001
		return resultJson
	}
	resultJson["data"] = slowSqlInfo
	return resultJson
}

//根据id查找慢sql的表的索引信息
func GetSlowSqlTableInfoById(id int) map[string]interface{} {
	DB1 := orm.NewOrm()
	var Sampleinfo []SampleInfo
	tables := make(map[string]interface{})
	table := make(map[string]interface{})
	var results []interface{}
	var explain,o []orm.Params
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	//获取慢查询的sql，表，库，ip，端口
	num, err := DB1.Raw("select sample,serverid_max,db_max from mysql_slow_query_review_history "+
		"where id=?", id).QueryRows(&Sampleinfo)
	if err != nil || num == 0 {
		resultJson["error"] = "no row found"
		resultJson["code"] = 1001
		return resultJson
	}
	for _, v := range Sampleinfo {
		//根据关键字截取到表名
		s := spilt.SplitSQL(v.Sample)
		host := strings.Split(v.ServeridMax, "-")[0]
		port, _ := strconv.Atoi(strings.Split(v.ServeridMax, "-")[1])
		tables["sql"] = v.Sample
		//拿到sql直接去生产库查询获取执行计划
		db := fmt.Sprintf("%s_%s_%d", "mysql", host, port)
		DB := orm.NewOrm()
		DB.Using(db)
		s1 := fmt.Sprintf("use %s;", v.DbMax)
		DB.Raw(s1).Exec()
		sql := fmt.Sprintf("explain %s;", v.Sample)
		DB.Raw(sql).Values(&explain)
		for _, v := range explain {
			o=append(o,v)
		}
		tables["explain"] =o
		//根据拿到的表名，获取表信息，索引信息
		dbName := v.DbMax
		for k, va := range s {
			if k != "default_1athena1_db" {
				dbName = k
			}
			for _, tableName := range va {
				rr:=GetTableInfo(host, dbName, tableName, port)["data"]
				table[tableName]=rr
			}
			tables["table"]=table
		}
	}
	results=append(results,tables)
	resultJson["data"] = results
	return resultJson
}

