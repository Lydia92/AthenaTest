package routers

import (
	"Athena/controllers"
	_ "Athena/models"
	"github.com/astaxie/beego"
)

func init() {
	beego.Router("/api/v1/main", &controllers.MainController{})
	beego.Router("/api/v1/initInfo", &controllers.InitController{})

	//获取实例信息
	beego.Router("/api/v1/instance", &controllers.InstanceController{}, "get:Get")
	beego.Router("/api/v1/instance/:id", &controllers.InstanceController{}, "get:GetOneInstanceInfo")

	//获取集群信息
	beego.Router("/api/v1/cluster", &controllers.ClusterController{}, "get:Get")
	beego.Router("/api/v1/cluster/:id", &controllers.ClusterController{}, "get:GetOneClusterInfo")
	beego.Router("/api/v1/cluster_name", &controllers.ClusterController{}, "get:GetClusterName")

	//获取某个集群下的账户信息
	beego.Router("/api/v1/dbaccount/:clusterName", &controllers.DBAccountController{}, "get:Get")
	//获取账户密码
	beego.Router("/api/v1/dbpasswd/:id", &controllers.DBAccountController{}, "get:GetPassword")

	//添加账户信息
	beego.Router("/hello1", &controllers.DBAccountController{}, "get:AddAccountInfo")

	//添加元数据信息
	beego.Router("/api/v1/addmetadata", &controllers.MetaController{}, "get:AddMetadata")
	/*beego.Router("/api/v1/metadatacolumn", &controllers.MetaController{}, "get:AddMetadataColumn")
	beego.Router("/api/v1/metadataindex", &controllers.MetaController{}, "get:AddMetadataIndex")*/

	//获取元表详细信息
	beego.Router("/api/v1/hostbycluster/:clusterName/", &controllers.MetaController{}, "get:GetHostByCluster")
	beego.Router("/api/v1/hostbycluster/", &controllers.MetaController{}, "get:GetDatabaseByHost")
	beego.Router("/api/v1/tableInfo/", &controllers.MetaController{}, "get:GetTableInfoBySchema")
	//beego.Router("/api/v1/tableIndexInfo/", &controllers.MetaController{}, "get:GetIndexInfoByTable")

	//获取集群下的慢查询
	beego.Router("/api/v1/slowsql/:clusterName", &controllers.SlowSqlController{}, "get:GetSlowSqlByClustername")
	//beego.Router("/api/v1/slowsql/", &controllers.SlowSqlController{}, "get:GetAll")
	beego.Router("/api/v1/slowsqlbyid/:id", &controllers.SlowSqlController{}, "get:GetSlowById")


}
