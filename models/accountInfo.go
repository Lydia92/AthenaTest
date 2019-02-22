package models

import (
	"encoding/hex"
	"fmt"
	"github.com/astaxie/beego/orm"
	_ "github.com/go-sql-driver/mysql"
	"time"
	//"gormTest/util"
	"Athena/util/encrypt"
	"Athena/util/spilt"
	"strings"
)

type MysqlUser struct {
	User string
	Host string
}

type dbAccountInfo struct {
	Account string
	Passwd string
	NodeAddr string
	NodePort uint16
}

/*
type Grants struct {
	grant string
}
*/

//获取所有实例的账户信息
func GetUserHost() []*DbAccountInfo {
	DB := orm.NewOrm()
	var dbaccountInfo []*DbAccountInfo
	DB.QueryTable("DbAccountInfo").All(&dbaccountInfo)
	return dbaccountInfo
}

//获得dba账户的信息。以便后面方便登录对应服务器做操作
func GetAccountForDBA() []*dbAccountInfo {

	accountInfoList := []*dbAccountInfo{}
	DB := orm.NewOrm()
	var db_account_info []*DbAccountInfo
	DB.QueryTable("DbAccountInfo").Filter("Account", "dba").All(&db_account_info,
		"Account", "Passwd", "NodeAddr", "NodePort")

	for _, item := range db_account_info {
		accountInfo := dbAccountInfo{}
		accountInfo.Passwd = encrypt.DecryptStr(item.Passwd)
		accountInfo.Account = item.Account
		accountInfo.NodeAddr = item.NodeAddr
		accountInfo.NodePort = item.NodePort
		accountInfoList = append(accountInfoList, &accountInfo)
	}

	return accountInfoList
}

//连接新的数据库，注册新的DB
func ConnectMysql(user, pass, address string, port uint16) orm.Ormer {
	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8", user, pass, address, port, "mysql")
	db := fmt.Sprintf("%s_%s_%d", "mysql", address, port)
	//注册数据库，连接对应的生产环境
	orm.RegisterDataBase(db, "mysql", dns)
	DB1 := orm.NewOrm()
	DB1.Using(db)
	return DB1
}

//获取账户权限
func SelectGrant(DB1 orm.Ormer) []map[string]string {
	var mysqlUser []MysqlUser
	var res []map[string]string
	var grant []orm.Params
	DB1.Raw("select user,host from mysql.user where user not like 'mysql.%' and user!=' ' and host !='::1'" +
		" and host not like '%localdomain%'").QueryRows(&mysqlUser)
	for _, v := range mysqlUser {
		sql := fmt.Sprintf("show grants for '%s'@'%s';", v.User, v.Host)
		DB1.Raw(sql).Values(&grant)
		for _, v1 := range grant {
			for _, v2 := range v1 {
				if v2.(string) != "" {
					u := spilt.Split(v2.(string))
					res = append(res, u)
				}
			}
		}

	}

	return res
}

//比较原生产是否有新的账户添加，若有，则加入到cmdb中，若有账户权限及密码有改变则更新
func ComparAccout() {
	db_account_info := GetAccountForDBA()
	DB := orm.NewOrm()
	for _, v3 := range db_account_info {
		//连接生产DB
		DB1 := ConnectMysql(v3.Account, v3.Passwd, v3.NodeAddr, v3.NodePort)
		//获取生产库中账户权限
		res := SelectGrant(DB1)
		for _, v := range res {
			ss := fmt.Sprintf("%s%s%s%s%s%s", v["user"], v["grants"], v["database"],
				v["table"], v["hosts"], v["option"])
			ss = strings.TrimSpace(ss)
			ss = strings.Replace(ss, "`", "", -1)
			ss = strings.Replace(ss, "'", "", -1)
			md5 := encrypt.GenerateMD5(ss)
			//fmt.Println(md5)
			account1 := strings.Replace(v["user"], "'", "", -1)
			account1 = strings.Replace(account1, " ", "", -1)
			user := DbAccountInfo{NodeAddr: v3.NodeAddr, NodePort: v3.NodePort, Account: account1, LoginAddr: v["hosts"]}
			err := DB.Read(&user, "NodeAddr", "NodePort", "Account", "LoginAddr")
			//判断是否有此账户或者md5值是否相同
			if err == orm.ErrNoRows {
				InserAccountInfo(v3.NodeAddr, "admin", account1, "123456", v["grants"], v["database"],
					v["table"], v["hosts"], md5, v3.NodePort)
			} else {
				if user.Md5 == md5 {
					fmt.Println("do nothing")
				} else {
					DB.QueryTable("DbAccountInfo").Filter("Id", user.Id).Update(orm.Params{
						"Md5": md5, "Priv": v["grants"], "DbName": v["database"], "TableName": v["table"],
						"Updated": time.Now()})
				}

			}
		}

	}

}

//插入新增的账户
func InserAccountInfo(node_addr, ownership, account, passwd, priv, db_name, table_name, login_addr,
	md5 string, node_port uint16) {
	//user=User{}
	DB := orm.NewOrm()
	instanceInfo := InstanceInfo{NodeAddr: node_addr, NodePort: node_port}
	DB.Read(&instanceInfo, "NodeAddr", "NodePort")
	var accountInfo = DbAccountInfo{
		NodeAddr:  node_addr,
		Instance:  &instanceInfo,
		Ownership: ownership,
		Account:   account,
		Passwd:    passwd,
		Priv:      priv,
		DbName:    db_name,
		TableName: table_name,
		LoginAddr: login_addr,
		Md5:       md5,
		NodePort:  node_port,
		IsGrant:   0,
		Validity:  0,
		Role:      4}
	//fmt.Println(accountInfo)
	DB.Insert(&accountInfo)
	//fmt.Println(err)
}

/*func (User) TableName() string {
	return "user"
}*/

//获取指定集群下面的账户信息
func GetAccountInfo(clusterName string) (resultJson map[string]interface{}) {
	var id uint64

	resultJson = make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	var accountRet []DbAccountInfo
	o := orm.NewOrm()
	err := o.Raw("select b.id "+
		"from cluster_instance a, instance_info b, cluster_info c "+
		"where a.instance_id = b.id and a.cluster_info_id = 1 "+
		"and c.cluster_name = ? limit 1", clusterName).QueryRow(&id)

	if err != nil {
		resultJson["error"] = fmt.Sprintf("invalid clusterName, clusterName is:%s", clusterName)
		resultJson["code"] = 1001
		return
	}

	o.Raw("select id, node_addr, node_port, ownership, account, priv, db_name, table_name, "+
		"login_addr, is_grant, validity, role, created "+
		"from db_account_info where instance_id = ?", id).QueryRows(&accountRet)

	resultJson["data"] = accountRet
	return
}

//获取账户密码
func GetPassword(id string) (resultJson map[string]interface{}) {

	resultJson = make(map[string]interface{}, 1)
	resultJson["code"] = 1000
	var password string
	o := orm.NewOrm()
	o.Raw("select passwd from db_account_info where id = ?", id).QueryRow(&password)

	if password == "" {
		resultJson["code"] = 1002
		resultJson["error"] = "Not account"
		return
	}
	cryptstr, _ := hex.DecodeString(password)
	passwd, _ := encrypt.AESDecryptWithECB([]byte(cryptstr), []byte(encryptKey))

	resultJson["data"] = string(passwd)
	return
}
