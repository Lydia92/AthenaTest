package models

import (
	"Athena/util/spilt"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

//集群慢sql信息
type SlowSqlInfo struct {
	Id                        int     `json:"id"`
	Node_addr                 string  `json:"node_addr"`
	Node_port                 uint16  `json:"node_port"`
	Cluster_name              string  `json:"cluster_name"`
	User_max                  string  `json:"user"`
	Db_max                    string  `json:"db"`
	Sample                    string  `json:"sample"`
	TsCnt                     float64 `json:"ts_cnt"`
	Query_time_pct            float64 `json:"query_time_pct"`
	Lock_time_pct             float64 `json:"lock_time_pct"`
	Rows_sent_pct             float64 `json:"rows_sent_pct"`
	Rows_examined_pct         float64 `json:"row_examined_pct"`
	Innodb_io_r_bytes_pct     float64 `json:"innodb_io_r_bytes_pct"`
	Innodb_io_r_ops_pct       float64 `json:"innodb_io_r_ops_pct"`
	Innodb_io_r_wait_pct      float64 `json:"innodb_io_r_wait_pctt"`
	Innodb_pages_distinct_pct float64 `json:"innodb_pages_distinct_pct"`
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

	num, err := DB1.Raw("select a.id,d.node_addr	,d.node_port,b.cluster_name	,a.user_max,a.db_max,a.sample,"+
		"a.ts_cnt,a.Query_time_pct_95  as query_time_pct ,a.Lock_time_pct_95 as lock_time_pct ,"+
		"a.Rows_sent_pct_95 as rows_sent_pct ,a.Rows_examined_pct_95 as rows_examined_pct,"+
		"a.InnoDB_IO_r_bytes_pct_95 as innodb_io_r_bytes_pct,a.innoDB_IO_r_ops_pct_95 as innodb_io_r_ops_pct,"+
		"a.InnoDB_IO_r_wait_pct_95 as Innodb_io_r_wait_pct,a.innoDB_pages_distinct_pct_95 as  innodb_pages_distinct_pct "+
		"from mysql_slow_query_review_history a "+
		"join instance_info d on a.serverid_max=concat(d.node_addr,'-',d.node_port) "+
		"join cluster_instance c on c.instance_id=d.id "+
		"join cluster_info b  on c.cluster_info_id=b.id "+
		" where b.cluster_name=? order by ts_cnt desc,Query_time_pct desc", cluserName).QueryRows(&slowSqlInfo)
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
	var results []interface{}
	var explain []orm.Params
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
		db_account_info := GetAccountForDBA()
		for _, v1 := range db_account_info {
			if v1.NodeAddr == host && int(v1.NodePort) == port {
				//因为之前方法是默认连接了mysql库，所以在此重新写注册数据库，后面可以优化此方法
				dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", v1.Account, v1.Passwd, v1.NodeAddr, port, v.DbMax)
				//注册数据库，连接对应的生产环境

				orm.RegisterDataBase(v.DbMax, "mysql", dns)
				DB := orm.NewOrm()
				DB.Using(v.DbMax)
				sql := fmt.Sprintf("explain %s;", v.Sample)
				DB.Raw(sql).Values(&explain)
				for _, v := range explain {
					tables["explain"] = v
				}
			}
		}
		//根据拿到的表名，获取表信息，索引信息
		for _, v1 := range s {
			v1 = strings.TrimSpace(v1)
			res := GetIndexInfoByTable(host, v.DbMax, v1, port)
			tables["index"] = res["index"]
			res1 := GetTableInfoBySchema(host, v.DbMax, v1, port)
			tables["tableinfo"] = res1["tableInfo"]
			results = append(results, tables)
		}
	}

	resultJson["data"] = results
	return resultJson

}
