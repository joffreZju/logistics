package routers

import (
	"allsum_bi/controllers/aggregatemgr"
	"allsum_bi/controllers/base"
	"allsum_bi/controllers/dataloadmgr"
	"allsum_bi/controllers/dbmgr"
	"allsum_bi/controllers/demand"
	"allsum_bi/controllers/etlmgr"
	"allsum_bi/controllers/kettlemgr"
	"allsum_bi/controllers/reportsetmgr"
	"allsum_bi/controllers/testmgr"
	"common/filter"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

const (
	DemandPrefix    string = "/v1/demand"
	ETLPrefix       string = "/v1/etl"
	DataLoadPrefix  string = "/v1/dataload"
	AggregatePrefix string = "/v1/aggregate"
	ReportSetPrefix string = "/v1/reportset"
	DbPrefix        string = "/v1/dbmgr"
	TestPrefix      string = "/v1/testmgr"
	KettlePrefix    string = "/v1/kettle"
)

func init() {
	//Demand
	beego.Router(DemandPrefix+"/list_demand", &demand.Controller{}, "get:ListDemand")
	beego.Router(DemandPrefix+"/add_demand", &demand.Controller{}, "*:AddDemand")
	//	beego.Router(DemandPrefix+"/analyze_demand", &demand.Controller{}, "post:AnalyzeDemand")
	beego.Router(DemandPrefix+"/get_analyze_report", &demand.Controller{}, "get:GetAnalyzeReport")
	beego.Router(DemandPrefix+"/set_demand", &demand.Controller{}, "post:SetDemand")
	beego.Router(DemandPrefix+"/get_demand_doc", &demand.Controller{}, "get:GetDemandDoc")
	beego.Router(DemandPrefix+"/upload_demand_doc", &demand.Controller{}, "post:UploadDemandDoc")
	beego.Router(DemandPrefix+"/publish_demand", &demand.Controller{}, "*:PublishDemand")
	beego.Router(DemandPrefix+"/review_demand", &demand.Controller{}, "*:ReviewDemand")
	beego.Router("/*", &base.Controller{}, "*:Index")

	//ETL
	beego.Router(ETLPrefix+"/show_sycn_list", &etlmgr.Controller{}, "get:ShowSycnList")
	beego.Router(ETLPrefix+"/data_calibration", &etlmgr.Controller{}, "post:DataCalibration")
	beego.Router(ETLPrefix+"/set_etl", &etlmgr.Controller{}, "post:SetEtl")
	beego.Router(ETLPrefix+"/stop_etl", &etlmgr.Controller{}, "post:StopEtl")

	//DATALOAD
	beego.Router(DataLoadPrefix+"/list", &dataloadmgr.Controller{}, "get:ListDataload")
	beego.Router(DataLoadPrefix+"/get", &dataloadmgr.Controller{}, "get:GetDataload")
	beego.Router(DataLoadPrefix+"/save", &dataloadmgr.Controller{}, "post:SaveDataload")
	beego.Router(DataLoadPrefix+"/test_create_script", &dataloadmgr.Controller{}, "post:TestDataLoadCreateScript")
	beego.Router(DataLoadPrefix+"/test_alter_script", &dataloadmgr.Controller{}, "post:TestDataLoadAlterScript")
	beego.Router(DataLoadPrefix+"/test_aggregate", &dataloadmgr.Controller{}, "post:TestAggregate")
	beego.Router(DataLoadPrefix+"/upload_web_file", &dataloadmgr.Controller{}, "post:UploadDataLoadWeb")
	beego.Router(DataLoadPrefix+"/download_web_file", &dataloadmgr.Controller{}, "get:DownloadDataLoadWeb")
	//DATALOAD_USER
	beego.Router(DataLoadPrefix+"/list_data", &dataloadmgr.Controller{}, "post:ListData")
	beego.Router(DataLoadPrefix+"/input_data", &dataloadmgr.Controller{}, "post:InputData")

	//AGGREAGATE

	beego.Router(AggregatePrefix+"/list", &aggregatemgr.Controller{}, "get:ListAggregate")
	beego.Router(AggregatePrefix+"/get", &aggregatemgr.Controller{}, "get:GetAggregate")
	beego.Router(AggregatePrefix+"/save", &aggregatemgr.Controller{}, "post:SaveAggregate")
	beego.Router(AggregatePrefix+"/test_create_script", &aggregatemgr.Controller{}, "post:TestAggregateCreateScript")
	beego.Router(AggregatePrefix+"/test_alter_script", &aggregatemgr.Controller{}, "post:TestAggregateAlterScript")
	beego.Router(AggregatePrefix+"/test_flush_script", &aggregatemgr.Controller{}, "post:TestAggregateFlushScript")

	//ReportSet
	beego.Router(ReportSetPrefix+"/list", &reportsetmgr.Controller{}, "get:ListReportSet")
	beego.Router(ReportSetPrefix+"/get", &reportsetmgr.Controller{}, "get:GetReportSet")
	beego.Router(ReportSetPrefix+"/get_reportset_web_file", &reportsetmgr.Controller{}, "get:GetReportSetWebFile")
	beego.Router(ReportSetPrefix+"/upload_reportset_web_file", &reportsetmgr.Controller{}, "post:UploadReportSetWeb")
	beego.Router(ReportSetPrefix+"/get_data", &reportsetmgr.Controller{}, "post:GetReportData")
	beego.Router(ReportSetPrefix+"/save", &reportsetmgr.Controller{}, "post:SaveReportSet")

	//DBMGR
	beego.Router(DbPrefix+"/add", &dbmgr.Controller{}, "post:AddDb")
	beego.Router(DbPrefix+"/list", &dbmgr.Controller{}, "get:ListDbDetail")
	beego.Router(DbPrefix+"/update", &dbmgr.Controller{}, "post:UpdateDb")
	beego.Router(DbPrefix+"/list_schema", &dbmgr.Controller{}, "get:ListSchema")
	beego.Router(DbPrefix+"/list_schema_table", &dbmgr.Controller{}, "get:ListSchemaTable")
	beego.Router(DbPrefix+"/list_all_db_schema", &dbmgr.Controller{}, "get:ListAllDbSchema")
	beego.Router(DbPrefix+"/delete", &dbmgr.Controller{}, "delete:DelDb")

	//testmgr
	beego.Router(TestPrefix+"/get", &testmgr.Controller{}, "get:GetTestInfo")
	beego.Router(TestPrefix+"/add", &testmgr.Controller{}, "post:AddTest")
	beego.Router(TestPrefix+"/get_image", &testmgr.Controller{}, "get:GetTestFile")

	//kettlemgr
	beego.Router(KettlePrefix+"/add_kettle_job", &kettlemgr.Controller{}, "post:AddKJob")
	beego.Router(KettlePrefix+"/list_kettle_job", &kettlemgr.Controller{}, "get:ListKJob")
	beego.Router(KettlePrefix+"/download_kettle_job", &kettlemgr.Controller{}, "get:DownloadKJob")
	beego.Router(KettlePrefix+"/set_kettle_job_enable", &kettlemgr.Controller{}, "post:SetJobEnable")
	beego.Router(KettlePrefix+"/delete_kettle_job", &kettlemgr.Controller{}, "delete:DeleteKJob")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowOrigins:     []string{"http://localhost:8090", "http://www.suanpeizaix.comw", "http://www.suanpeizaix.com:8090"},
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"uid", "cno", "access_token", "Authorization", "X-Requested-With", "Content-Type", "Origin", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST", "OPTIONS"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
	}))

	//api auth white list

	notNeedAuthList := []string{}

	//filter
	beego.InsertFilter("/v1/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
