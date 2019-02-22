package models

import (
	"time"
)

// 集群信息表
type ClusterInfo struct {
	Id          uint64             `orm:"auto" description:"主键"`
	Cluster     []*ClusterInstance `orm:"reverse(many)"`
	ClusterName string             `orm:"size(32);unique" description:"集群名,全局唯一,推荐使用业务名+别名+序号"`
	WriteVip    string             `orm:"size(15);default(---)" description:"表示表示集群写 vip,可能没有值"`
	WritePort   uint16             `orm:"default(0)" description:"表示集群写端口号,可能没有值"`
	ReadPort    uint16             `orm:"default(0)" description:"表示集群读端口号,可能没有值"`
	ReadVip     string             `orm:"size(15);default(---)" description:"表示表示集群读 vip，可能没有值"`
	ClusterType uint8              `orm:"default(1)" description:"0:unknown,1:主从,2:MHA,3:MGR,4:单节点。默认为1"`
	Created     time.Time          `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated     time.Time          `orm:"auto_now;type(datetime);auto_now"`
}

// 表示集群与实例的关系, 一对多关系
type ClusterInstance struct {
	Id          uint64        `orm:"auto" description:"主键"`
	ClusterInfo *ClusterInfo  `orm:"rel(fk)"`
	Instance    *InstanceInfo `orm:"rel(one);index"`
	Role        uint8         `orm:"default(1)" description:"表示实例角色：1:primary, 2:secondary"`
	Created     time.Time     `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated     time.Time     `orm:"auto_now;type(datetime);auto_now"`
}

// 建立唯一索引
func (self *ClusterInstance) TableUnique() [][]string {
	return [][]string{
		[]string{"ClusterInfo", "Instance"},
	}
}

// 业务与集群的对应关系， 一对一关系
type BusinessCluster struct {
	Id           uint64       `orm:"auto" description:"主键"`
	ClusterInfo  *ClusterInfo `orm:"rel(one)" description:"集群信息表主键，唯一"`
	BusinessName string       `orm:"size(32)" description:"业务名称"`
	Leader       string       `orm:"size(32)" description:"业务负责人"`
	Created      time.Time    `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated      time.Time    `orm:"auto_now;type(datetime);auto_now"`
}

// 实例信息表， 实例对集群为多对一关系
type InstanceInfo struct {
	Id              uint64           `orm:"auto" description:"主键"`
	ClusterInstance *ClusterInstance `orm:"reverse(one)"`
	NodeAddr        string           `orm:"size(15)" description:"表示实例地址"`
	NodePort        uint16           `orm:"default(3001)" description:"表示实例端口"`
	BaseDir         string           `orm:"size(50);default(/usr/local/mysql)" description:"表示实例安装目录存放位置" `
	DataDir         string           `orm:"size(50);default(/data/mysql/mysqldata3001)" description:"表示实例数据目录存放位置" `
	ConfPath        string           `orm:"size(50);default(/etc/my3001.cnf)" description:"表示实例配置文件全路径" `
	InstanceVersion string           `orm:"size(15);default(mysql-5.7.24)" description:"表示实例版本信息"`
	InstanceType    uint8            `orm:"default(1)" description:"表示实例是人工录入还是平台部署，1: 平台部署，2: 人工录入"`
	InstanceStatus  uint8            `orm:"default(3)" description:"实例状态,1:部署中，2:部署失败，3:待分配，4:已分配，5:下线"`
	InstanceTime    time.Time        `orm:"type(datetime);null" description:"表示实例分配给集群使用的时间"`
	Created         time.Time        `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated         time.Time        `orm:"auto_now;type(datetime);auto_now"`
}

func (self *InstanceInfo) TableUnique() [][]string {
	return [][]string{
		[]string{"NodeAddr", "NodePort"},
	}
}

// 数据库用户信息， 与实例为 多对一关系
type DbAccountInfo struct {
	Id        uint64        `orm:"auto" description:"主键"`
	Instance  *InstanceInfo `orm:"rel(fk);index"`
	NodeAddr  string        `orm:"size(15)" description:"表示实例地址"`
	NodePort  uint16        `orm:"default(3001)" description:"表示实例端口, node_addr+node_port 表示唯一一个实例"`
	Ownership string        `orm:"size(30)" description:"账户拥有者"`
	Account   string        `orm:"size(30)" description:"登录到数据库的用户名"`
	Passwd    string        `orm:"size(128)" description:"登录数据时使用的密码，加密存储"`
	Priv      string        `orm:"size(530)" description:"账户所拥有的权限"`
	DbName    string        `orm:"size(64)" description:"权限作用于的库名，如果是所有库，使用*表示"`
	TableName string        `orm:"size(64)" description:"权限作用于的表名，如果是所有表，使用*表示"`
	LoginAddr string        `orm:"size(15)" description:"允许登录数据库的地址，可以是%，网段"`
	IsGrant   uint8         `orm:"default(0)" description:"表示此账户是否有 with grant option 权限，0:表示没有，1: 表示有"`
	Validity  uint          `orm:"default(0)" description:"表示此账户的生命周期，0表示永久，大于0 表示天数，到期时间为create_time + validity"`
	Role      uint8         `orm:"default(4)" description:"1:管理员账号,2:应用账户,3:只读账号*表示,4:手动创建,5:人工录入"`
	Md5       string        `orm:"type(char)" description:"account+priv+db_name+tab_name+login_addr+is_grant 生成MD5,判断是否变化过"`
	Created   time.Time     `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated   time.Time     `orm:"auto_now;type(datetime);auto_now"`
}

func (u *DbAccountInfo) TableUnique() [][]string {
	return [][]string{
		[]string{"NodeAddr", "NodePort", "Account", "LoginAddr"},
	}
}

//主机基础信息
type HostInfo struct {
	Id              uint64             `orm:"auto" description:"主键"`
	HostAccountInfo []*HostAccountInfo `orm:"reverse(many)"`
	HostAppType     []*HostAppType     `orm:"reverse(many)"`
	Host            string             `orm:"size(15);unique" description:"表示主机地址"`
	Port            uint16             `orm:"default(22)" description:"表示主机ssh端口"`
	Disk            string             `orm:"size(255);null" description:"磁盘信息,{目录:{总大小，空闲量，使用比}}"`
	Cpu             string             `orm:"size(32);null" description:"cpu信息,{物理id:[核心数,频率],}"`
	Mem             uint16             `orm:"null" description:"内存信息,单位是MB"`
	Environment     uint8              `orm:"index;default(0)" description:"0:unknown,1:dev,2:sit,3:uat,4:vis,5:prd"`
	Created         time.Time          `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated         time.Time          `orm:"auto_now;type(datetime);auto_now"`
}

//主机账号信息
type HostAccountInfo struct {
	Id       uint64    `orm:"auto" description:"主键"`
	HostInfo *HostInfo `orm:"rel(fk);index"`
	User     string    `orm:"size(30)" description:"登录到主机的用户名"`
	Passwd   string    `orm:"size(128)" description:"登录主机时使用的密码，加密存储"`
	Validity uint      `orm:"default(0)" description:"此账户的生命周期，0表示永久，大于0表示天数，到期时间为create_time + validity"`
	Created  time.Time `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated  time.Time `orm:"auto_now;type(datetime);auto_now"`
}

//应用类型， 表示一台机器上跑了哪些服务
type HostAppType struct {
	Id       uint64    `orm:"auto" description:"主键"`
	HostInfo *HostInfo `orm:"rel(fk)"`
	AppType  *AppType  `orm:"rel(fk)"`
	Created  time.Time `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated  time.Time `orm:"auto_now;type(datetime);auto_now"`
}

func (self *HostAppType) TableUnique() [][]string {
	return [][]string{
		[]string{"HostInfo", "AppType"},
	}
}

//应用名称
type AppType struct {
	Id      uint64    `orm:"auto" description:"主键"`
	AppType string    `orm:"size(32);default(os);unique" description:"应用名，mongo/mysql/redis/java..."`
	Created time.Time `orm:"auto_now_add;type(datetime);auto_now_add"`
	Updated time.Time `orm:"auto_now;type(datetime);auto_now"`
}

//系统元表信息
type MysqlMetadataTables struct {
	Id             uint64    `orm:"auto" description:"主键"`
	NodeAddr       string    `orm:"size(15)" description:"表示实例地址"`
	NodePort       uint16    `orm:"default(3001)" description:"表示实例端口, node_addr+node_port 表示唯一一个实例"`
	TableSchema    string    `orm:"size(64)" description:"库名"`
	TableName      string    `orm:"size(64)" description:"表名"`
	DbEngine       string    `orm:"size(64)" description:"引擎"`
	RowFormat      string    `orm:"size(10)" description:"行格式"`
	TableRows      uint64    `orm:"size(21)" description:"行数"`
	AvgRowLength   uint64    `orm:"size(21)" description:"平均行长度"`
	MaxDataLength  uint64    `orm:"size(21)" description:"最大行长度的"`
	DataLength     uint64    `orm:"size(21)" description:"数据长度"`
	IndexLength    uint64    `orm:"size(21)" description:"索引长度"`
	DataFree       uint64    `orm:"size(21)" description:"空闲"`
	ChipSize       uint64    `orm:"size(21)" description:"碎片"`
	AutoIncrement  uint64    `orm:"size(21)" description:"下一个自增的值"`
	TableCollation string    `orm:"size(32)" description:"字符集"`
	CreateTime     time.Time `orm:"type(datetime)" description:"表创建时间"`
	UpdateTime     time.Time `orm:"null; type(datetime)" description:"表更新时间"`
	CheckTime      time.Time `orm:"null; type(datetime)" description:"表检测时间"`
	TableComment   string    `orm:"size(2048)" description:"表注释"`
	TableMd5       string    `orm:"type(char)" description:"HOST+PORT+db_name+tab_name生成MD5,判断是否变化过"`
}

func (u *MysqlMetadataTables) TableUnique() [][]string {
	return [][]string{
		[]string{"NodeAddr", "NodePort", "TableSchema", "TableName"},
	}
}

//系统列信息
type MysqlMetadataColumns struct {
	Id            uint64 `orm:"auto" description:"主键"`
	NodeAddr      string `orm:"size(15)" description:"表示实例地址"`
	NodePort      uint16 `orm:"default(3001)" description:"表示实例端口, node_addr+node_port 表示唯一一个实例"`
	TableSchema   string `orm:"size(64)" description:"库名"`
	TableName     string `orm:"size(64)" description:"表名"`
	ColumnName    string `orm:"size(64)" description:"列名"`
	ColumnType    string `orm:"size(64)" description:"列类型"`
	CollationName string `orm:"size(32)" description:"字符集"`
	IsNullable    string `orm:"size(3)" description:"是否为空"`
	ColumnKey     string `orm:"size(64)" description:"键类型"`
	ColumnDefault string `orm:"type(text)" description:"列默认值"`
	Extra         string `orm:"size(30)" description:"说明"`
	ColPrivileges string `orm:"size(80)" description:"权限"`
	ColumnComment string `orm:"size(1024)" description:"列备注"`
	ColumnMd5     string `orm:"size(64)" description:"md5值 node_addr+node_port+table_schema+table_name_column_name"`
}

//系统索引信息
type MysqlMetadataIndexs struct {
	Id           uint64 `orm:"auto" description:"主键"`
	NodeAddr     string `orm:"size(15)" description:"表示实例地址"`
	NodePort     uint16 `orm:"default(3001)" description:"表示实例端口, node_addr+node_port 表示唯一一个实例"`
	TableSchema  string `orm:"size(64)" description:"库名"`
	TableName    string `orm:"size(64)" description:"表名"`
	ColumnName   string `orm:"size(64)" description:"列名"`
	NonUnique    uint64 `orm:"default(0)" description:"是否可包含重复项 0表示不能重复，索引为唯一索引"`
	IndexName    string `orm:"size(64)" description:"索引名"`
	SeqInIndex   uint64 `orm:"default(0)" description:"索引的序列号，从1开始"`
	Cardinality  uint64 `orm:"size(21)" description:"索引中唯一值的数量"`
	Nullable     string `orm:"size(3)" description:"能否为空"`
	IndexType    string `orm:"size(16)" description:"索引类型"`
	IndexComment string `orm:"size(16)" description:"索引备注"`
	IndexMd5     string `orm:"size(32)" description:"md5值 node_addr+node_port+table_schema+table_name_column_name+index_name"`
}

//慢查询表信息
type MysqlSlowQueryReviewHistory struct {
	Id                        uint64    `orm:"auto" description:"主键"`
	ServeridMax               string    `orm:"size(20)" description:"主机+port"`
	DbMax                     string    `orm:"size(100);null" description:"数据库名"`
	UserMax                   string    `orm:"size(100)" description:"用户"`
	Checksum                  string    `orm:"size(32)" description:"MD5"`
	Sample                    string    `orm:"type(text)" description:"sql"`
	TsMin                     time.Time `orm:"auto_now;type(datetime);auto_now"`
	TsMax                     time.Time `orm:"auto_now;type(datetime);auto_now"`
	TsCnt                     float64   `orm:"null"`
	QueryTimeSum              float64   `orm:"null"`
	QueryTimeMin              float64   `orm:"null"`
	QueryTimeMax              float64   `orm:"null"`
	QueryTimePct_95           float64   `orm:"null"`
	QueryTimeStddev           float64   `orm:"null"`
	QueryTimeMedian           float64   `orm:"null"`
	LockTimeSum               float64   `orm:"null"`
	LockTimeMin               float64   `orm:"null"`
	LockTimeMax               float64   `orm:"null"`
	LockTimePct_95            float64   `orm:"null"`
	LockTimeStddev            float64   `orm:"null"`
	LockTimeMedian            float64   `orm:"null"`
	RowsSentSum               float64   `orm:"null"`
	RowsSentMin               float64   `orm:"null"`
	RowsSentMax               float64   `orm:"null"`
	RowsSentPct_95            float64   `orm:"null"`
	RowsSentStddev            float64   `orm:"null"`
	RowsSentMedian            float64   `orm:"null"`
	RowsExaminedSum           float64   `orm:"null"`
	RowsExaminedMin           float64   `orm:"null"`
	RowsExaminedMax           float64   `orm:"null"`
	RowsExaminedPct_95        float64   `orm:"null"`
	RowsExaminedStddev        float64   `orm:"null"`
	RowsExaminedMedian        float64   `orm:"null"`
	RowsAffectedSum           float64   `orm:"null"`
	RowsAffectedMin           float64   `orm:"null"`
	RowsAffectedMax           float64   `orm:"null"`
	RowsAffectedPct_95        float64   `orm:"null"`
	RowsAffectedStddev        float64   `orm:"null"`
	RowsAffectedMedian        float64   `orm:"null"`
	RowsReadSum               float64   `orm:"null"`
	RowsReadMin               float64   `orm:"null"`
	RowsReadMax               float64   `orm:"null"`
	RowsReadPct_95            float64   `orm:"null"`
	RowsReadStddev            float64   `orm:"null"`
	RowsReadMedian            float64   `orm:"null"`
	MergePassesSum            float64   `orm:"null"`
	MergePassesMin            float64   `orm:"null"`
	MergePassesMax            float64   `orm:"null"`
	MergePassesPct_95         float64   `orm:"null"`
	MergePassesStddev         float64   `orm:"null"`
	MergePassesMedian         float64   `orm:"null"`
	InnodbIoROpsMin           float64   `orm:"null"`
	InnodbIoROpsMax           float64   `orm:"null"`
	InnodbIoROpsPct_95        float64   `orm:"null"`
	InnodbIoROpsStddev        float64   `orm:"null"`
	InnodbIoROpsMedian        float64   `orm:"null"`
	InnodbIoRBytesMin         float64   `orm:"null"`
	InnodbIoRBytesMax         float64   `orm:"null"`
	InnodbIoRBytesPct_95      float64   `orm:"null"`
	InnodbIoRBytesStddev      float64   `orm:"null"`
	InnodbIoRBytesMedian      float64   `orm:"null"`
	InnodbIoRWaitMin          float64   `orm:"null"`
	InnodbIoRWaitMax          float64   `orm:"null"`
	InnodbIoRWaitPct_95       float64   `orm:"null"`
	InnodbIoRWaitStddev       float64   `orm:"null"`
	InnodbIoRWaitMedian       float64   `orm:"null"`
	InnodbRecLockWaitMin      float64   `orm:"null"`
	InnodbRecLockWaitMax      float64   `orm:"null"`
	InnodbRecLockWaitPct_95   float64   `orm:"null"`
	InnodbRecLockWaitStddev   float64   `orm:"null"`
	InnodbRecLockWaitMedian   float64   `orm:"null"`
	InnodbQueueWaitMin        float64   `orm:"null"`
	InnodbQueueWaitMax        float64   `orm:"null"`
	InnodbQueueWaitPct_95     float64   `orm:"null"`
	InnodbQueueWaitStddev     float64   `orm:"null"`
	InnodbQueueWaitMedian     float64   `orm:"null"`
	InnodbPagesDistinctMin    float64   `orm:"null"`
	InnodbPagesDistinctMax    float64   `orm:"null"`
	InnodbPagesDistinctPct_95 float64   `orm:"null"`
	InnodbPagesDistinctStddev float64   `orm:"null"`
	InnodbPagesDistinctMedian float64   `orm:"null"`
	QxHitCnt                  float64   `orm:"null"`
	QcHitSum                  float64   `orm:"null"`
	FullScanCnt               float64   `orm:"null"`
	FullScanSum               float64   `orm:"null"`
	FullJoinCnt               float64   `orm:"null"`
	FullJoinSum               float64   `orm:"null"`
	TmpTableCnt               float64   `orm:"null"`
	TmpTableSum               float64   `orm:"null"`
	TmpTableOnDisk_cnt        float64   `orm:"null"`
	TmpTableOnDiskSum         float64   `orm:"null"`
	FilesortCnt               float64   `orm:"null"`
	FilesortSum               float64   `orm:"null"`
	FilesortOnDiskCnt         float64   `orm:"null"`
	FilesortOnDisksum         float64   `orm:"null"`
}

func (u *MysqlSlowQueryReviewHistory) TableIndex() [][]string {
	return [][]string{
		[]string{"ServeridMax"},
	}
}
func (u *MysqlSlowQueryReviewHistory) TableUnique() [][]string {
	return [][]string{
		[]string{"Checksum"},
	}
}

type MysqlSlowQueryReview struct {
	Id          uint64    `orm:"auto" description:"主键"`
	Checksum    string    `orm:"size(32)"`
	Fingerprint string    `orm:"type(text)" `
	Sample      string    `orm:"type(text)" `
	FirstSeen   time.Time `orm:"null"`
	LastSeen    time.Time `orm:"null"`
	ReviewedBy  string    `orm:"size(32)"`
	ReviewedOn  time.Time `orm:"null"`
	Comments    string    `orm:"type(text);null"`
}
