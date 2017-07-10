package db

import (
	"allsum_bi/db/conn"
	"allsum_bi/models"
	"allsum_bi/util"

	"github.com/astaxie/beego"
)

func InitDb() {
	//链接BI数据库
	InitBIConn()
	//按需初始化数据库
	InitDatabase()

	//连接源数据库
	InitSourceDBConn()

}

func InitDatabase() {
	//TODO
	db, err := conn.GetBIConn()
	if err != nil {
		beego.Error("InitDatabase | get conn err :%v", err)
		return
	}
	//auto migrate table
	beego.Debug("create manager schema")
	CreateManagerSchema()
	beego.Debug("create system schema")
	CreateSystemSchema()

	beego.Debug("create system schema")
	db.Table(models.GetDatabaseManagerTableName()).AutoMigrate(&models.DatabaseManager{})
	db.Table(models.GetDemandTableName()).AutoMigrate(&models.Demand{})
	db.Table(models.GetReportTableName()).AutoMigrate(&models.Report{})
	db.Table(models.GetSynchronousTableName()).AutoMigrate(&models.Synchronous{})
	db.Table(models.GetSynchronousLogTableName()).AutoMigrate(&models.SynchronousLog{})
	db.Table(models.GetDataLoadTableName()).AutoMigrate(&models.DataLoad{})
	db.Table(models.GetAggregateOpsTableName()).AutoMigrate(&models.AggregateOps{})
	beego.Debug("testinfo table: ", models.GetTestInfoTableName())
	//	db.Table(models.GetTestInfoTableName()).AutoMigrate(&models.TestInfo{})
	//	db.Table(models.GetAggregateOpsLogTableName()).AutoMigrate(&models.AggregateOps{})
}

func InitBIConn() (err error) {
	host := beego.AppConfig.String("bi_base_db::host")
	port, err := beego.AppConfig.Int("bi_base_db::port")
	if err != nil {
		port = 5432
	}
	user := beego.AppConfig.String("bi_base_db::user")
	password := beego.AppConfig.String("bi_base_db::password")
	dbname := beego.AppConfig.String("bi_base_db::dbname")
	conninfo := conn.Conn{
		Id:     util.BASEDB_CONNID,
		Dbtype: util.PG_DB_TYPE,
		Name:   "BIBaseDB",
		Host:   host,
		Port:   port,
		DbUser: user,
		Passwd: password,
		Dbname: dbname,
	}
	err = conn.CreateConn(conninfo)
	if err != nil {
		return err
	}
	return
}

func InitSourceDBConn() (err error) {
	//TODO
	dbinfos, err := models.ListDatabaseManager()
	if err != nil {
		beego.Error("InitSourceDBConn ListDatabaseManager err :", dbinfos, err)
		return
	}
	beego.Debug("dbinfos :", dbinfos)
	for _, dbinfo := range dbinfos {
		connStruct := conn.Conn{
			Id:     dbinfo.Dbid,
			Dbtype: dbinfo.Dbtype,
			Name:   dbinfo.Name,
			Host:   dbinfo.Host,
			Port:   dbinfo.Port,
			DbUser: dbinfo.Dbuser,
			Passwd: dbinfo.Password,
			Dbname: dbinfo.Dbname,
		}
		beego.Debug("init db", connStruct.Id)
		err = conn.CreateConn(connStruct)
		if err != nil {
			beego.Error("db:  conn err:", dbinfo.Name, err)
			continue
		}
	}
	return
}
