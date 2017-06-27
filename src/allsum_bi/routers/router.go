package routers

import (
	"allsum_bi/controllers/base"
	"allsum_bi/controllers/demand"

	"github.com/astaxie/beego"
)

func init() {
	//router
	beego.Router("/list_demand", &demand.Controller{}, "get:ListDemand")
	beego.Router("/add_demand", &demand.Controller{}, "*:AddDemand")
	beego.Router("/analyze_demand", &demand.Controller{}, "post:AnalyzeDemand")
	beego.Router("/get_analyze_report", &demand.Controller{}, "get:GetAnalyzeReport")
	beego.Router("/set_demand", &demand.Controller{}, "post:SetDemand")
	beego.Router("/get_demand_doc", &demand.Controller{}, "get:GetDemandDoc")
	beego.Router("/upload_demand_doc", &demand.Controller{}, "post:UploadDemandDoc")
	beego.Router("/*", &base.Controller{}, "*:Index")
}
