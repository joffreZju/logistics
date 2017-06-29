package routers

import (
	"allsum_bi/controllers/base"
	"allsum_bi/controllers/demand"
	"allsum_bi/controllers/etlmgr"
	"common/filter"

	"github.com/astaxie/beego"
)

const (
	DemandPrefix string = "/v1/demand"
	ETLPrefix    string = "/v1/etl"
)

func init() {
	//Demand
	beego.Router(DemandPrefix+"/list_demand", &demand.Controller{}, "get:ListDemand")
	beego.Router(DemandPrefix+"/add_demand", &demand.Controller{}, "*:AddDemand")
	beego.Router(DemandPrefix+"/analyze_demand", &demand.Controller{}, "post:AnalyzeDemand")
	beego.Router(DemandPrefix+"/get_analyze_report", &demand.Controller{}, "get:GetAnalyzeReport")
	beego.Router(DemandPrefix+"/set_demand", &demand.Controller{}, "post:SetDemand")
	beego.Router(DemandPrefix+"/get_demand_doc", &demand.Controller{}, "get:GetDemandDoc")
	beego.Router(DemandPrefix+"/upload_demand_doc", &demand.Controller{}, "post:UploadDemandDoc")
	beego.Router("/*", &base.Controller{}, "*:Index")

	//ETL
	beego.Router(ETLPrefix+"/show_sycn_list", &etlmgr.Controller{}, "get:ShowSycnList")
	beego.Router(ETLPrefix+"/data_calibration", &etlmgr.Controller{}, "post:DataCalibration")
	beego.Router(ETLPrefix+"/set_etl", &etlmgr.Controller{}, "post:SetEtl")
	beego.Router(ETLPrefix+"/stop_etl", &etlmgr.Controller{}, "post:StopEtl")

	//api auth white list

	//	notNeedAuthList := []string{}

	//filter
	//	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
