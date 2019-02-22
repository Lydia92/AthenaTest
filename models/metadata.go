package models

import (
	"Athena/util/encrypt"
	"fmt"
	"github.com/astaxie/beego/orm"
	"strings"
	"time"
)

//原表信息
type InformationTable struct {
	Table_schema    string
	Table_name      string
	Engine          string
	Row_format      string
	Table_rows      uint64
	Avg_row_length  uint64
	Data_length     uint64
	Max_data_length uint64
	Index_length    uint64
	Data_free       uint64
	Auto_increment  uint64
	Table_collation string
	Table_comment   string
	Create_time     time.Time
	Update_time     time.Time
	Check_time      time.Time
}

//元表 列信息
type InfomationColumn struct {
	Table_schema   string
	Table_name     string
	Column_name    string
	Column_type    string
	Collation_name string
	Is_nullable    string
	Column_key     string
	Column_default string
	Extra          string
	Privileges     string
	Column_comment string
}

//元表 索引信息
type InfomationIndex struct {
	Table_schema  string
	Table_name    string
	Column_name   string
	Non_unique    uint64
	Index_name    string
	Seq_in_index  uint64
	Cardinality   uint64
	Nullable      string
	Index_type    string
	Index_comment string
}

//集群信息
type InfomationCluster struct {
	NodeAddr    string
	NodePort    int
	ClusterName string
}

//库名
type Database struct {
	TableSchema string
}

//表，列信息
type TableCloumnInfo struct {
	TableSchema    string
	TableName      string
	RowFormat      string
	TableRows      string
	AvgRowLength   int
	MaxDataLength  int
	DataLength     int
	IndexLength    int
	DataFree       int
	TableCollation string
	CreateTime     time.Time
	UpdateTime     time.Time
	CheckTime      time.Time
	ColumnName     string
	ColumnType     string
	IsNullable     string
	ColumnKey      string
	ColPrivileges  string
}

//表索引信息
type TableIndexInfo struct {
	ColumnComment string
	NonUnique     int
	IndexName     string
	IndexType     string
	SeqInIndex    int
	Cardinality   int
	IndexComment  string
}

//获取系统元表信息
func GetMetaTable(DB1 orm.Ormer) []InformationTable {
	var informationTable []InformationTable
	//var res []map[string]string
	sql := "SELECT table_schema, table_name, engine, row_format, table_rows, avg_row_length,data_length," +
		" max_data_length,index_length, data_free, auto_increment,table_collation, table_comment," +
		" create_time, update_time, check_time FROM information_schema.tables " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	DB1.Raw(sql).QueryRows(&informationTable)

	return informationTable
}

//插入元表信息
func (infomationTable *MysqlMetadataTables)InsertTableInfo(Host, Table_schema, Table_name, Engine, Row_format string, Table_rows, Avg_row_length,
	Data_length, Max_data_length, Index_length, Data_free, Auto_increment uint64,
	Table_collation, Table_comment, Md5 string, Create_time, Update_time, check_time time.Time, port uint16) {
	DB := orm.NewOrm()
	 infomationTable = &MysqlMetadataTables{
		NodeAddr:       Host,
		NodePort:       port,
		TableSchema:    Table_schema,
		TableName:      Table_name,
		DbEngine:       Engine,
		RowFormat:      Row_format,
		TableRows:      Table_rows,
		AvgRowLength:   Avg_row_length,
		MaxDataLength:  Max_data_length,
		DataLength:     Data_length,
		IndexLength:    Index_length,
		DataFree:       Data_free,
		ChipSize:       0,
		AutoIncrement:  Auto_increment,
		TableCollation: Table_collation,
		CreateTime:     Create_time,
		UpdateTime:     Update_time,
		CheckTime:      check_time,
		TableComment:   Table_comment,
		TableMd5:       Md5,
	}
	DB.Insert(infomationTable)
}

//对比元表数据与平台数据，更新新数据
func CompareMetaTable() {
	db_account_info := GetAccountForDBA()
	var informationTable []InformationTable
	DB := orm.NewOrm()
	mysqlMetadataTables:=MysqlMetadataTables{}
	for _, v3 := range db_account_info {
		DB1 := ConnectMysql(v3.Account, v3.Passwd, v3.NodeAddr, v3.NodePort)

		informationTable = GetMetaTable(DB1)
		for _, v := range informationTable {
			md5 := fmt.Sprintf("%s%d%s%s", v3.NodeAddr, v3.NodePort, v.Table_schema, v.Table_name)
			md5 = strings.TrimSpace(md5)
			md5 = encrypt.GenerateMD5(md5)
			mysqlmetadatatables := MysqlMetadataTables{NodeAddr: v3.NodeAddr, NodePort: v3.NodePort,
				TableSchema: v.Table_schema, TableName: v.Table_name}
			err := DB.Read(&mysqlmetadatatables, "NodeAddr", "NodePort", "TableSchema", "TableName")
			if err == orm.ErrNoRows {
				mysqlMetadataTables.InsertTableInfo(v3.NodeAddr, v.Table_schema, v.Table_name, v.Engine, v.Row_format, v.Table_rows,
					v.Avg_row_length, v.Data_length, v.Max_data_length, v.Index_length, v.Data_free, v.Auto_increment,
					v.Table_collation, v.Table_comment, md5, v.Create_time, v.Update_time, v.Check_time, v3.NodePort)

			} else {
				if mysqlmetadatatables.TableMd5 != md5 {
					DB.QueryTable("mysqlmetadatatables").Filter("id", mysqlmetadatatables.Id).Update(orm.Params{
						"TableMd5": md5, "DbEngine": v.Engine, "RowFormat": v.Row_format, "TableRows": v.Table_rows,
						"AvgRowLength": v.Avg_row_length, "MaxDataLength": v.Max_data_length, "DataLength": v.Data_length,
						"IndexLength": v.Index_length, "DataFree": v.Data_free, "AutoIncrement": v.Auto_increment,
						"TableCollation": v.Table_collation, "UpdateTime": v.Update_time, "CheckTime": v.Check_time,
						"TableComment": v.Table_comment})
				} else {
					fmt.Println("do nothing")
				}

			}
		}

	}

}

//获取元表列信息
func GetMetaColumns(DB1 orm.Ormer) []InfomationColumn {
	var infomationColumn []InfomationColumn
	sql := "SELECT table_schema,table_name,column_name,column_type,collation_name,is_nullable," +
		"column_key,column_default,extra,privileges,column_comment FROM information_schema.columns " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	DB1.Raw(sql).QueryRows(&infomationColumn)
	return infomationColumn
}

//插入列信息到平台表
func (mysqlMetadataColumns *MysqlMetadataColumns)InserColumnInfo(Node_add, Table_schema, Table_name, Column_name, Column_type, Collation_name, Is_nullable,
	Column_key, Column_default, Extra, Col_privileges, Column_comment, md5 string, port uint16) {
	DB := orm.NewOrm()
	 mysqlMetadataColumns = &MysqlMetadataColumns{
		NodeAddr:      Node_add,
		NodePort:      port,
		TableSchema:   Table_schema,
		TableName:     Table_name,
		ColumnName:    Column_name,
		ColumnType:    Column_type,
		CollationName: Collation_name,
		IsNullable:    Is_nullable,
		ColumnKey:     Column_key,
		ColumnDefault: Column_default,
		Extra:         Extra,
		ColPrivileges: Col_privileges,
		ColumnComment: Column_comment,
		ColumnMd5:     md5,
	}
	DB.Insert(mysqlMetadataColumns)

}

//比较列信息与平台信息
func CompareMetaColumns() {
	db_account_info := GetAccountForDBA()
	var infomationColumn []InfomationColumn
	mysqlMetadataColumns:=MysqlMetadataColumns{}
	DB := orm.NewOrm()
	for _, v3 := range db_account_info {
		DB1 := ConnectMysql(v3.Account, v3.Passwd, v3.NodeAddr, v3.NodePort)
		infomationColumn = GetMetaColumns(DB1)
		for _, v := range infomationColumn {
			md5 := fmt.Sprintf("%s%d%s%s%s%s%s%s%s%s%s%s%s%", v3.NodeAddr, v3.NodePort, v.Table_schema,
				v.Table_name, v.Column_name, v.Column_type, v.Collation_name, v.Is_nullable, v.Column_key, v.Column_default,
				v.Extra, v.Privileges, v.Extra)
			md5 = strings.TrimSpace(md5)
			md5 = encrypt.GenerateMD5(md5)
			mysqlmetadatacolumns := MysqlMetadataColumns{NodeAddr: v3.NodeAddr, NodePort: v3.NodePort,
				TableSchema: v.Table_schema, TableName: v.Table_name, ColumnName: v.Column_name}
			err := DB.Read(&mysqlmetadatacolumns, "NodeAddr", "NodePort", "TableSchema", "TableName", "ColumnName")
			if err == orm.ErrNoRows {
				mysqlMetadataColumns.InserColumnInfo(v3.NodeAddr, v.Table_schema, v.Table_name, v.Column_name, v.Column_type, v.Collation_name,
					v.Is_nullable, v.Column_key, v.Column_default, v.Extra, v.Privileges, v.Column_comment, md5, v3.NodePort)
			} else {
				if mysqlmetadatacolumns.ColumnMd5 != md5 {
					DB.QueryTable("mysqlmetadatacolumns").Filter("id", mysqlmetadatacolumns.Id).
						Update(orm.Params{"ColumnName": v.Column_name, "ColumnType": v.Column_type,
							"CollationName ": v.Collation_name, "IsNullable": v.Is_nullable, "ColumnKey ": v.Column_key,
							"ColumnDefault ": v.Column_default, "Extra ": v.Extra, "ColPrivileges ": v.Privileges,
							"ColumnComment ": v.Column_comment, "ColumnMd5 ": md5})

				} else {
					fmt.Println("do nothing")
				}
			}

		}

	}
}

//获取索引信息
func GetMetaIndex(DB1 orm.Ormer) []InfomationIndex {
	var infomationIndex []InfomationIndex
	sql := "select table_schema,table_name,column_name,non_unique,index_name,seq_in_index,cardinality," +
		"nullable,index_type,index_comment from information_schema.statistics " +
		"where table_schema not in ('sys', 'test', 'information_schema', 'performance_schema', 'mysql')"
	DB1.Raw(sql).QueryRows(&infomationIndex)
	return infomationIndex
}

//插入索引信息
func (mysqlMetadataIndexs *MysqlMetadataIndexs)InsertIndexInfo(Node_addr, Table_schema, Table_name, Column_name, Index_name string, Non_unique, Seq_in_index,
	Cardinality uint64, Nullable, Index_type, Index_comment, md5 string, port uint16) {
	DB1 := orm.NewOrm()
	 mysqlMetadataIndexs = &MysqlMetadataIndexs{
		NodeAddr:     Node_addr,
		NodePort:     port,
		TableSchema:  Table_schema,
		TableName:    Table_name,
		ColumnName:   Column_name,
		NonUnique:    Non_unique,
		IndexName:    Index_name,
		SeqInIndex:   Seq_in_index,
		Cardinality:  Cardinality,
		Nullable:     Nullable,
		IndexType:    Index_type,
		IndexComment: Index_comment,
		IndexMd5:     md5,
	}
	DB1.Insert(mysqlMetadataIndexs)

}

//比较索引信息，更新
func CompareIndex() {
	db_account_info := GetAccountForDBA()
	var infomationIndex []InfomationIndex
	mysqlMetadataIndexs:=MysqlMetadataIndexs{}
	DB := orm.NewOrm()
	for _, v3 := range db_account_info {
		DB1 := ConnectMysql(v3.Account, v3.Passwd, v3.NodeAddr, v3.NodePort)
		infomationIndex = GetMetaIndex(DB1)
		for _, v := range infomationIndex {
			md5 := fmt.Sprintf("%s%d%s%s%s%s%s%s%s%s%s%s", v3.NodeAddr, v3.NodePort, v.Table_schema,
				v.Table_name, v.Column_name, v.Non_unique, v.Index_name, v.Seq_in_index, v.Cardinality,
				v.Nullable, v.Index_type, v.Index_comment)
			md5 = strings.TrimSpace(md5)
			md5 = encrypt.GenerateMD5(md5)
			mysqlmetadataindexs := MysqlMetadataIndexs{NodeAddr: v3.NodeAddr, NodePort: v3.NodePort,
				TableSchema: v.Table_schema, TableName: v.Table_name}
			err := DB.Read(&mysqlmetadataindexs, "NodeAddr", "NodePort", "TableSchema", "TableName")
			if err == orm.ErrNoRows {
				mysqlMetadataIndexs.InsertIndexInfo(v3.NodeAddr, v.Table_schema, v.Table_name, v.Column_name, v.Index_name, v.Non_unique,
					v.Seq_in_index, v.Cardinality, v.Nullable, v.Index_type, v.Index_comment, md5, v3.NodePort)
			} else {
				if mysqlmetadataindexs.IndexMd5 != md5 {
					DB.QueryTable("mysqlMetadataIndexs").Filter("id", mysqlmetadataindexs.Id).
						Update(orm.Params{"ColumnName": v.Column_name, "NonUnique": v.Non_unique, "IndexName": v.Index_name,
							"SeqInIndex": v.Seq_in_index, "Cardinality": v.Cardinality, "Nullable": v.Nullable,
							"IndexType": v.Index_type, "IndexComment": v.Index_comment, "IndexMd5": md5})
				} else {
					fmt.Println("index do nothing")
				}
			}
		}
	}
}

//根据集群名获取集群中实例及端口
func GetHostByCluster(cname string) map[string]interface{} {

	DB1 := orm.NewOrm()
	var infomationCluster []InfomationCluster
	// sql:=""
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	_, err := DB1.Raw("select  a.node_addr, a.node_port, c.cluster_name	from instance_info a,  cluster_info c, cluster_instance d	"+
		"where  a.id=d.instance_id and c.id=d.cluster_info_id and cluster_name=?", cname).QueryRows(&infomationCluster)
	if err != nil {
		resultJson["error"] = "no row found"
		resultJson["code"] = 1001
		return resultJson
	}
	resultJson["data"] = infomationCluster
	return resultJson

}

//根据实例地址，端口，获取库名
func GetDatabaseByHost(host string, port int) map[string]interface{} {
	DB1 := orm.NewOrm()
	var database []Database
	// sql:=""
	resultJson := make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	_, err := DB1.Raw("select distinct table_schema from mysql_metadata_tables where node_addr=? and node_port=?", host, port).QueryRows(&database)
	if err != nil {
		resultJson["error"] = "no row found"
		resultJson["code"] = 1001
		return resultJson
	}
	resultJson["data"] = database
	return resultJson
}

//根据实例地址，端口，库名获取表信息
func GetTableInfoBySchema(host, table_schema, table_name string, port int) map[string]interface{} {
	DB1 := orm.NewOrm()
	var tableCloumnInfo []TableCloumnInfo
	res := make(map[string]interface{}, 1)
	res["code"] = 1000
	_, err := DB1.Raw("select a.table_schema,a.table_name,a.row_format,a.table_rows,a.avg_row_length,"+
		"a.max_data_length,a.data_length,a.index_length,a.data_free,a.table_collation,a.create_time,a.update_time,"+
		"a.check_time,b.column_name,b.column_type,b.column_type,b.is_nullable,b.column_key,b.col_privileges,"+
		"b.column_comment from mysql_metadata_tables a"+
		" join mysql_metadata_columns b on a.node_addr=b.node_addr "+
		"and a.node_port=b.node_port and a.table_schema=b.table_schema and a.table_name=b.table_name "+
		" where a.node_addr=? and a.node_port=? and a.table_schema=?"+
		" and a.table_name=? ", host, port, table_schema, table_name).QueryRows(&tableCloumnInfo)


	if err != nil {
		res["error"] = "no row found"
		res["code"] = 1001
		return res
	}
	res["tableInfo"] = tableCloumnInfo
	return res
}

func GetIndexInfoByTable(host, table_schema, table_name string, port int) map[string]interface{} {
	DB1 := orm.NewOrm()
	var tableIndexInfo []TableIndexInfo
	res := make(map[string]interface{})
	var create []orm.Params
	res["code"] = 1000
	_, err := DB1.Raw("select c.non_unique,c.index_name,c.index_type,c.seq_in_index,c.cardinality,"+
		"c.index_comment,c.column_name "+
		"from mysql_metadata_indexs c where c.node_addr=? and c.node_port=? and c.table_schema=?"+
		" and c.table_name=? ", host, port, table_schema, table_name).QueryRows(&tableIndexInfo)

	db_account_info := GetAccountForDBA()
	for _, v := range db_account_info {
		if v.NodeAddr == host && int(v.NodePort) == port {
			DB := ConnectMysql(v.Account, v.Passwd, v.NodeAddr, v.NodePort)
			sql := fmt.Sprintf("show create table %s.%s;", table_schema, table_name)
			DB.Raw(sql).Values(&create)
			for _, v := range create {
				res["create"] = v["Create Table"]
			}
		}
	}
	res["index"] = tableIndexInfo

	if err != nil {
		res["error"] = "no row found"
		res["code"] = 1001
		return res
	}
	return res
}

func GetTableInfo(host, table_schema, table_name string, port int) map[string]interface{} {
	resultJson := make(map[string]interface{}, 1)
	result := make(map[string]interface{})
	resultJson["code"] = 1000
	res := GetTableInfoBySchema(host, table_schema, table_name, port)
	res1 := GetIndexInfoByTable(host, table_schema, table_name, port)
	if res["code"] == 1001 || res1["code"] == 1001 {
		resultJson["code"] = 1001
		resultJson["error"] = "something err"
	}
	result["tableInfo"] = res["tableInfo"]
	result["create"] = res1["create"]
	result["index"] = res1["index"]

	resultJson["data"] = result
	return resultJson

}
