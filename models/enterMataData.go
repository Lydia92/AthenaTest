package models

import (
	"Athena/util/encrypt"
	"Athena/util/numberutil"
	"fmt"
	"github.com/360EntSecGroup-Skylar/excelize"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/orm"
	"strconv"
	"strings"
)

var file = "/Users/zs/go/src/Athena/tests/testMetadata.xlsx"
var encryptKey = beego.AppConfig.String("encrypt_key::key")

func ReadExcel(fileName, sheet string) (rows [][]string, err error) {
	xlsx, err := excelize.OpenFile(fileName)
	if err != nil {
		return nil, err
	}
	rows = xlsx.GetRows(sheet)
	return rows, nil
}

func InitInfo() {
	InitOsInfo("os")
	enterInstanceInfo("instance")
	enterDbAccount("dbAccount")
	enterCluster("cluster")
	enterHostApp("hostApp")
}

//TODO:需要处理单元格为空的情况, 写入前先检测是否存在，存在则不再插入
//录入机器信息
func InitOsInfo(sheet string) {
	rows, err := ReadExcel(file, sheet)
	if err != nil {
		beego.Error(err)
	}

	for i, row := range rows {
		// 读每一行数据
		if i == 0 {
			continue
		}
		hostInfo := HostInfo{}
		accountInfo := HostAccountInfo{}
	ExitNull:
		for j, colCell := range row {
			switch j {
			case 0:
				host := strings.TrimSpace(colCell)
				if len(host) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty host in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				hostInfo.Host = host
			case 1:
				user := strings.TrimSpace(colCell)
				if len(user) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty username in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				accountInfo.User = user
			case 2:
				passwd := strings.TrimSpace(colCell)
				if len(passwd) == 0 {
					beego.Warning(
						fmt.Sprintf("Have a empty password in file:[%d] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				//cryptoStr, _ := encrypt.AESEncryptWithECB([]byte(passwd), []byte(encryptKey))
				//accountInfo.Passwd = strings.ToUpper(hex.EncodeToString(cryptoStr))

				accountInfo.Passwd = encrypt.EncryptStr(passwd)

			case 3:
				p, _ := strconv.Atoi(strings.TrimSpace(colCell))
				if p == 0 {
					beego.Error(
						fmt.Sprintf("Have a invalid port, file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
				}
				hostInfo.Port = numberutil.IntToUint16(p)
			case 4:
				env := strings.TrimSpace(colCell)
				if len(env) == 0 {
					beego.Error(fmt.Sprintf("Have a empty environment in file:[%s] sheet[%s] at the line[%d].",
						file, sheet, i+1))
				}
				hostInfo.Environment = getEvent(env)
			}
		}

		// 进行数据插入
		o := orm.NewOrm()
		o.Begin()
		_, err1 := o.Insert(&hostInfo)
		if err1 != nil {
			o.Rollback()
			beego.Error(
				fmt.Sprintf("Insert host_info error, err:[%s] file:[%s] sheet[%s] line:[%d].",
					err1.Error(), file, sheet, i+1))
			continue
		}
		accountInfo.HostInfo = &hostInfo
		_, err2 := o.Insert(&accountInfo)
		if err2 != nil {
			o.Rollback()
			beego.Error(
				fmt.Sprintf("Insert account_info error, err:[%s] file:[%s] sheet[%s] line:[%d].",
					err2.Error(), file, sheet, i+1))
			continue
		}
		o.Commit()
	}
}

//录入实例信息
func enterInstanceInfo(sheet string) {

	rows, err := ReadExcel(file, sheet)
	if err != nil {
		beego.Error(err)
	}
	o := orm.NewOrm()
	for i, row := range rows {
		// 读每一行数据
		if i == 0 {
			continue
		}

		instanceInfo := InstanceInfo{}
	ExitNull:
		for j, colCell := range row {
			switch j {
			case 0:
				host := strings.TrimSpace(colCell)
				if len(host) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty host in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.NodeAddr = host
			case 1:
				p, _ := strconv.Atoi(strings.TrimSpace(colCell))
				if p == 0 {
					beego.Error(
						fmt.Sprintf("Have a invalid port in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.NodePort = numberutil.IntToUint16(p)

			case 2:
				baseDir := strings.TrimSpace(colCell)
				if len(baseDir) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty baseDir in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.BaseDir = baseDir

			case 3:
				dataDir := strings.TrimSpace(colCell)
				if len(dataDir) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty dataDir in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.DataDir = dataDir
			case 4:
				confPath := strings.TrimSpace(colCell)
				if len(confPath) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty confPath in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.ConfPath = confPath
			case 5:
				version := strings.TrimSpace(colCell)
				if len(version) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty version in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instanceInfo.InstanceVersion = version
			}
		}
		instanceInfo.InstanceType = 2
		instanceInfo.InstanceStatus = 2

		//进行数据插入
		//o.Insert(&instanceInfo)

		//如果查询后不为空，则表示没有结果，需要进行插入
		err := o.Read(&instanceInfo, "NodeAddr", "NodePort")
		if err != nil {
			//进行数据插入
			o.Insert(&instanceInfo)
		}

	}
}

//录入db账户信息
func enterDbAccount(sheet string) {
	rows, err := ReadExcel(file, sheet)
	if err != nil {
		beego.Error(err)
	}
	o := orm.NewOrm()
	for i, row := range rows {
		// 读每一行数据
		if i == 0 {
			continue
		}

		dbAccountInfo := DbAccountInfo{}
		instanceInfo := InstanceInfo{}
	ExitNull:
		for j, colCell := range row {
			switch j {
			case 0:
				host := strings.TrimSpace(colCell)
				if len(host) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty host in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.NodeAddr = host
				instanceInfo.NodeAddr = host
			case 1:
				dbName := strings.TrimSpace(colCell)
				if len(dbName) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty dbName in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.DbName = dbName
			case 2:
				p, _ := strconv.Atoi(colCell)
				if p == 0 {
					beego.Error(
						fmt.Sprintf("Have a invalid port in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.NodePort = uint16(p)
				instanceInfo.NodePort = uint16(p)
			case 3:
				account := strings.TrimSpace(colCell)
				if len(account) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty account in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.Account = account
				dbAccountInfo.Ownership = account
			case 4:
				passwd := strings.TrimSpace(colCell)
				if len(passwd) == 0 {
					beego.Warning(fmt.Sprintf("Have a empty password on file:[%d] sheet[%s] at the line[%d].",
						file, sheet, i+1))
					break ExitNull
				}
				//cryptoStr, _ := encrypt.AESEncryptWithECB([]byte(passwd), []byte(encryptKey))
				//dbAccountInfo.Passwd = strings.ToUpper(hex.EncodeToString(cryptoStr))
				dbAccountInfo.Passwd = encrypt.EncryptStr(passwd)
			case 5:
				priv := strings.TrimSpace(colCell)
				if len(priv) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty priv in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.Priv = priv
			case 6:
				tableName := strings.TrimSpace(colCell)
				if len(tableName) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty tableName in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.TableName = strings.TrimSpace(colCell)
			case 7:
				loginAddr := strings.TrimSpace(colCell)
				if len(loginAddr) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty loginAddr in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.LoginAddr = strings.TrimSpace(colCell)
			case 8:
				grant, _ := strconv.Atoi(colCell)
				if grant < 0 || grant > 1 {
					beego.Error(
						fmt.Sprintf("Have a invalid isGrant in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				dbAccountInfo.IsGrant = uint8(grant)
			}
		}
		dbAccountInfo.Validity = 0
		dbAccountInfo.Role = 5

		//account+priv+db_name+tab_name+login_addr+is_grant
		md5Str := fmt.Sprintf("%s%s%s%s%s%d", dbAccountInfo.Account, dbAccountInfo.Priv, dbAccountInfo.DbName,
			dbAccountInfo.LoginAddr, dbAccountInfo.IsGrant)

		md5s := encrypt.GenerateMD5(md5Str)
		dbAccountInfo.Md5 = md5s
		o.Read(&instanceInfo, "NodeAddr", "NodePort")
		//fmt.Println(&instanceInfo)
		dbAccountInfo.Instance = &instanceInfo

		//如果查询后不为空，则表示没有结果，需要进行插入
		err := o.Read(&dbAccountInfo, "NodeAddr", "NodePort", "Md5")
		//fmt.Println(err)
		if err != nil {
			//进行数据插入
			_, err = o.Insert(&dbAccountInfo)
			//fmt.Println(err)
		}

	}
}

//录入集群与业务信息
func enterCluster(sheet string) {
	rows, err := ReadExcel(file, sheet)
	if err != nil {
		beego.Error(err)
	}
	o := orm.NewOrm()

	for i, row := range rows {
		// 读每一行数据
		if i == 0 {
			continue
		}

		var instanceList []*InstanceInfo
		clusterInfo := ClusterInfo{}
		instancePrimary := InstanceInfo{}
		instanceInfo2 := InstanceInfo{}
		instanceInfo3 := InstanceInfo{}
		business := BusinessCluster{}
	ExitNull:
		for j, colCell := range row {
			switch j {
			case 0:
				clusterName := strings.TrimSpace(colCell)
				if len(clusterName) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty clusterName in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				clusterInfo.ClusterName = strings.TrimSpace(colCell)
			case 1:
				businessName := strings.TrimSpace(colCell)
				if len(businessName) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty businessName in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				business.BusinessName = businessName
			case 2:
				leader := strings.TrimSpace(colCell)
				if len(leader) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty leader in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				business.Leader = leader
			case 3:
				masterNode := strings.TrimSpace(colCell)
				if len(masterNode) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty masterNode in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instancePrimary.NodeAddr = masterNode
			case 4:
				p, _ := strconv.Atoi(strings.TrimSpace(colCell))
				if p == 0 {
					beego.Error(
						fmt.Sprintf("Have a invalid port in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				instancePrimary.NodePort = uint16(p)
			case 5:
				slaveNode1 := strings.TrimSpace(colCell)
				slavePort1 := instancePrimary.NodePort
				if len(slaveNode1) == 0 {
					slaveNode1 = "---"
					slavePort1 = 0
				}
				instanceInfo2.NodeAddr = slaveNode1
				instanceInfo2.NodePort = slavePort1
			case 6:
				slaveNode2 := strings.TrimSpace(colCell)
				slavePort2 := instancePrimary.NodePort
				if len(slaveNode2) == 0 {
					slaveNode2 = "---"
					slavePort2 = 0
				}
				instanceInfo3.NodeAddr = slaveNode2
				instanceInfo3.NodePort = slavePort2
			case 7:
				writeVip := strings.TrimSpace(colCell)
				if len(writeVip) == 0 {
					writeVip = "---"
				}
				clusterInfo.WriteVip = writeVip
			case 8:
				p, _ := strconv.Atoi(strings.TrimSpace(colCell))
				if clusterInfo.WriteVip == "---" {
					p = 0
				}
				clusterInfo.WritePort = uint16(p)
			case 9:
				readVip := strings.TrimSpace(colCell)
				if len(readVip) == 0 {
					readVip = "---"
				}
				clusterInfo.ReadVip = readVip
			case 10:
				p, _ := strconv.Atoi(strings.TrimSpace(colCell))
				if clusterInfo.ReadVip == "---" {
					p = 0
				}
				clusterInfo.ReadPort = uint16(p)
			case 11:
				clusterType := strings.TrimSpace(colCell)
				if len(clusterType) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty clusterType in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				clusterInfo.ClusterType = getEvent(clusterType)
			}
		}
		o.Read(&instancePrimary, "NodeAddr", "NodePort")
		o.Read(&instanceInfo2, "NodeAddr", "NodePort")
		o.Read(&instanceInfo3, "NodeAddr", "NodePort")

		instanceList = append(instanceList, &instancePrimary)
		instanceList = append(instanceList, &instanceInfo2)
		instanceList = append(instanceList, &instanceInfo3)

		o.Insert(&clusterInfo)

		business.ClusterInfo = &clusterInfo

		o.Insert(&business)

		for i, item := range instanceList {
			clusterInstance := ClusterInstance{}
			clusterInstance.ClusterInfo = &clusterInfo
			clusterInstance.Instance = item
			if i == 0 {
				clusterInstance.Role = 1
			} else {
				clusterInstance.Role = 2
			}
			if item.NodePort != 0 {
				o.Insert(&clusterInstance)
			}
		}

	}
}

//录入主机与应用关系
func enterHostApp(sheet string) {

	rows, err := ReadExcel(file, sheet)
	if err != nil {
		beego.Error(err)
	}
	o := orm.NewOrm()

	for i, row := range rows {
		// 读每一行数据
		if i == 0 {
			continue
		}

		hostInfo := HostInfo{}
		hostAppType := HostAppType{}
		hostApp := AppType{}
	ExitNull:
		for j, colCell := range row {
			switch j {
			case 0:
				host := strings.TrimSpace(colCell)
				if len(host) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty ip in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				hostInfo.Host = host
			case 1:
				appName := strings.TrimSpace(colCell)
				if len(appName) == 0 {
					beego.Error(
						fmt.Sprintf("Have a empty appName in file:[%s] sheet[%s] at the line[%d].",
							file, sheet, i+1))
					break ExitNull
				}
				hostApp.AppType = appName
			}
		}
		o.Read(&hostInfo, "Host")
		hostAppType.HostInfo = &hostInfo

		err := o.Read(&hostApp, "AppType")
		if err != nil {
			o.Insert(&hostApp)
		}
		hostAppType.AppType = &hostApp
		o.Insert(&hostAppType)
	}
}

func getEvent(colCell string) uint8 {
	var id uint8
	switch colCell {
	case "VIS", "vis", "单节点":
		id = 4
	case "dev", "DEV", "主从":
		id = 1
	case "sit", "SIT", "mha", "MHA":
		id = 2
	case "uat", "UAT", "MGR", "mgr":
		id = 3
	case "prd", "PRD":
		id = 5
	}
	return id
}

/*
select a.host, c.app_type
from host_info a, host_app_type b, app_type c
where a.id = b.host_info_id
and b.app_type_id = c.id

*/
