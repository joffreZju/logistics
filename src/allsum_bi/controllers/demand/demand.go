package demand

import (
	"allsum_bi/models"
	"allsum_bi/services/demandsvs"
	"allsum_bi/services/oa"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	base "common/lib/baseController"
	"common/lib/errcode"
	"encoding/json"
	_ "fmt"
	"io/ioutil"
	"time"

	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListDemand() {
	roleType, err := c.GetInt("type")
	if err != nil {
		beego.Error("get type err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 20
	}
	index, err := c.GetInt("index")
	if err != nil {
		index = 0
	}
	index = limit * index
	var fields []string
	var values []interface{}
	var action string
	if roleType == util.ROLETYPE_ASSIGNER {
		fields = []string{}
		values = []interface{}{}
		action = util.ACTION_LISTDEMAND_ASSIGNER
	} else if roleType == util.ROLETYPE_PROJECTOR {
		fields = []string{"handlerid"}
		values = []interface{}{c.UserID}
		if c.UserID == 0 {
			values = []interface{}{-1}
		}
		action = util.ACTION_LISTDEMAND_PROJECTOR
	} else if roleType == util.ROLETYPE_TESTER {
		fields = []string{"status"}
		values = []interface{}{util.DEMAND_STATUS_TESTING}
		action = util.ACTION_LISTDEMAND_TESTER
	} else {
		beego.Error("error roleType:", roleType)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	if !oa.CheckActionEnable(c.UserID, action) {
		beego.Error("no have authority", c.UserID)
		c.ReplyErr(errcode.ErrActionNoAuthority)
		return
	}

	demands, err := models.ListDemandByField(fields, values, limit, index)
	if err != nil {
		beego.Error("list demands error :", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}
	count, err := models.CountDemandByField(fields, values)
	if err != nil {
		beego.Error("list demands error :", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}
	var res map[string]interface{}
	var data []map[string]interface{}
	for _, demand := range demands {
		report, err := models.GetReport(demand.Reportid)
		if err != nil {
			continue
		}
		mapdemand := map[string]interface{}{
			"index":      demand.Id,
			"uuid":       demand.Uuid,
			"reportuuid": report.Uuid,
			//			"owner":             demand.Owner,
			"owner_name": demand.OwnerName,
			"exhibitor":  demand.Exhibitor,
			//			"reportid":          demand.Reportid,
			"description": demand.Description,
			//			"handleid":          demand.Handleid,
			"handler_name": demand.HandlerName,
			//			"assignerid":        demand.Assignerid,
			"assigner_name": demand.AssignerName,
			"init_time":     demand.Inittime,
			"deadline":      demand.Deadline,
			"assigne_time":  demand.Assignetime,
			"complettime":   demand.Complettime,
			"status":        demand.Status,
		}
		data = append(data, mapdemand)
	}
	res = map[string]interface{}{
		"count": count,
		"data":  data,
	}
	c.ReplySucc(res)
	return
}

//工单过来的需求
func (c *Controller) AddDemand() {
	owner := c.GetString("ownerid")

	ownername := c.GetString("owner_name")
	exhibitor := c.GetString("exhibitor")
	description := c.GetString("description")
	if owner == "" || ownername == "" || exhibitor == "" || description == "" {
		beego.Error("miss param ")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	inittime := time.Now()
	status := util.DEMAND_STATUS_NO_ASSIGN
	demand := models.Demand{
		Owner:       owner,
		OwnerName:   ownername,
		Description: description,
		Inittime:    inittime,
		Exhibitor:   exhibitor,
		Status:      status,
	}
	demandres, err := models.InsertDemand(demand)
	if err != nil {
		beego.Error("insert demand err : ", err)
		c.ReplyErr(errcode.ErrActionPutDemand)
		return
	}
	report_create := models.Report{
		Demandid:    demandres.Id,
		Name:        "",
		Owner:       demandres.Owner,
		Status:      util.REPORT_STATUS_ANALYS,
		Description: demandres.Description,
	}
	report, err := models.InsertReport(report_create)
	if err != nil {
		beego.Error("insert report err : ", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	demandres.Reportid = report.Id
	err = models.UpdateDemand(demandres, "reportid")
	if err != nil {
		beego.Error("update demand err", err)
		c.ReplyErr(errcode.ErrActionPutDemand)
		return
	}

	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

//发布需求
func (c *Controller) PublishDemand() {
	demanduuid := c.GetString("demanduuid")
	err := demandsvs.ChangeStatus(demanduuid, util.DEMAND_STATUS_RELEASE, util.REPORT_STATUS_RELEASE, c.UserComp)
	if err != nil {
		beego.Error("update demand err : ", err)
		c.ReplyErr(errcode.ErrActionPutReport)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

//提交审核需求 fortest
func (c *Controller) ReviewDemand() {
	demanduuid := c.GetString("demanduuid")
	err := demandsvs.ChangeStatus(demanduuid, util.DEMAND_STATUS_REVIEW, util.REPORT_STATUS_REVIEW, c.UserComp)
	if err != nil {
		beego.Error("update demand err : ", err)
		c.ReplyErr(errcode.ErrActionPutReport)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

//需求分析 注: 将生成报表数据
func (c *Controller) AnalyzeDemand() {
	demanduuid := c.GetString("demanduuid")
	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		return
	}
	var report models.Report
	if demand.Reportid == 0 {
		report_create := models.Report{
			Demandid:    demand.Id,
			Name:        "",
			Owner:       demand.Owner,
			Status:      util.REPORT_STATUS_ANALYS,
			Description: demand.Description,
		}
		report, err = models.InsertReport(report_create)
		if err != nil {
			beego.Error("insert report err : ", err)
			c.ReplyErr(errcode.ErrActionGetReport)
			return
		}
		demand.Reportid = report.Id
		err = models.UpdateDemand(demand, "reportid")
		if err != nil {
			beego.Error("update demand err", err)
			c.ReplyErr(errcode.ErrActionPutDemand)
			return
		}
	} else {
		report, err = models.GetReport(demand.Reportid)
		if err != nil {
			beego.Error("get report err : ", err)
			c.ReplyErr(errcode.ErrActionGetReport)
			return
		}
	}
	res := map[string]interface{}{
		"reportuuid":   report.Uuid,
		"demand_owner": demand.OwnerName,
		"exhibitor":    demand.Exhibitor,
		"inittime":     demand.Inittime,
		"contactid":    demand.Contactid,
		"description":  demand.Description,
	}
	c.ReplySucc(res)
}

//获取需求分析数据
func (c *Controller) GetAnalyzeReport() {
	demanduuid := c.GetString("demanduuid")
	//	reportuuid := c.GetString("reportuuid")

	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		beego.Error("get demand err :", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}
	report, err := models.GetReport(demand.Reportid)
	if err != nil {
		beego.Error("get report err :", err)
		c.ReplyErr(errcode.ErrActionGetReport)
	}
	var assigner_authority []map[string]string
	if demand.AssignerAuthority == "" {
		assigner_authority = []map[string]string{}
	} else {
		err = json.Unmarshal([]byte(demand.AssignerAuthority), &assigner_authority)
		if err != nil {
			beego.Error("unmarshal assigner_authority err :", err)
			c.ReplyErr(errcode.ErrServerError)
		}
	}

	beego.Debug("authority: ", assigner_authority)
	authoritymap := map[string]interface{}{}
	for _, v := range assigner_authority {
		dbid := v["dbid"]
		dbname := v["dbname"]
		data, ok := authoritymap[dbid]
		if !ok {
			data = map[string]interface{}{
				"dbid":    dbid,
				"dbname":  dbname,
				"schemas": []interface{}{v["schema"]},
			}
		} else {
			datamap := data.(map[string]interface{})
			datamap["schemas"] = append(datamap["schemas"].([]interface{}), v["schema"])
			data = datamap
		}
		authoritymap[dbid] = data
	}
	var authorityArray []interface{}
	for _, v := range authoritymap {
		authorityArray = append(authorityArray, v)
	}
	beego.Debug("authoritymap:", authoritymap)

	res := map[string]interface{}{
		"reportuuid":         report.Uuid,
		"demanduuid":         demand.Uuid,
		"demand_owner":       demand.OwnerName,
		"exhibitor":          demand.Exhibitor,
		"contactid":          demand.Contactid,
		"description":        demand.Description,
		"inittime":           demand.Inittime,
		"price":              demand.Price,
		"report_type":        report.Reporttype,
		"doc_name":           demand.DocName,
		"assignetime":        demand.Assignetime,
		"handler_name":       demand.HandlerName,
		"deadline":           demand.Deadline,
		"assigner_authority": authorityArray,
	}
	c.ReplySucc(res)
}

//需求分析设置
func (c *Controller) SetDemand() {
	var reqbody map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqbody)

	if err != nil {
		beego.Error("json unmarshal fail err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	//	reportuuid := reqbody["reportuuid"].(string)
	demanduuid := reqbody["demanduuid"].(string)
	description := reqbody["description"].(string)
	price := reqbody["price"].(float64)
	report_type := int(reqbody["report_type"].(float64))
	handlerid := int(reqbody["handlerid"].(float64))
	var deadline time.Time
	deadlinestr := reqbody["deadline"].(string)
	deadline, err = time.Parse("2015-01-01 00:00:00", deadlinestr)
	if err != nil {
		deadline = time.Now().AddDate(0, 0, 7)
	}
	assigner_authority_bytes, err := json.Marshal(reqbody["assigner_authority"])
	if err != nil {
		beego.Error("json marshal assigner_authority err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	var assigner_authority_maps []map[string]string
	err = json.Unmarshal(assigner_authority_bytes, &assigner_authority_maps)
	if err != nil {
		beego.Error("json unmarshal assigner_authority err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	err = demandsvs.AddAuthority(assigner_authority_maps, c.UserComp, c.UserID)
	if err != nil {
		beego.Error("add assigner_authority err", err)
		c.ReplyErr(errcode.ErrActionAddAuthority)
		return
	}
	assigner_authority := string(assigner_authority_bytes)

	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		beego.Error("get demand err : ", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}
	report, err := models.GetReport(demand.Reportid)
	if err != nil {
		beego.Error("get report err : ", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	demand.Description = description
	demand.Handlerid = handlerid
	demand.Deadline = deadline
	demand.AssignerAuthority = assigner_authority
	demand.Status = util.DEMAND_STATUS_BUILDING
	demand.Price = price

	err = models.UpdateDemand(demand, "description", "price", "handleid", "deadline", "assigner_authority", "status")
	if err != nil {
		beego.Error("update demand err :", err)
		c.ReplyErr(errcode.ErrActionPutDemand)
		return
	}

	report.Reporttype = report_type
	report.Status = util.REPORT_STATUS_DEVELOP

	err = models.UpdateReport(report, "reporttype", "status")
	if err != nil {
		beego.Error("update report err :", err)
		c.ReplyErr(errcode.ErrActionPutReport)
		return
	}

	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)

}

//下载需求文档
func (c *Controller) GetDemandDoc() {
	demanduuid := c.GetString("demanduuid")
	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		beego.Error("get demand err :", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}
	filedata, err := ossfile.GetFile(demand.DocUrl)
	if err != nil {
		beego.Error("get file err :", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	c.ReplyFile("application/octet-stream", demand.DocName, filedata)
	return
}

func (c *Controller) UploadDemandDoc() {
	demanduuid := c.Ctx.Request.Header.Get("demanduuid")
	f, h, err := c.GetFile("uploadfile")
	if err != nil {
		beego.Error("get file err : ", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	defer f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Error("read filehandle err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	filename := uuid.NewV4().String() + "-" + h.Filename
	uri, err := ossfile.PutFile("demand", filename, data)
	if err != nil {
		beego.Error("put file to oss err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}

	demand, err := models.GetDemandByUuid(demanduuid)
	demand.DocUrl = uri
	demand.DocName = h.Filename
	err = models.UpdateDemand(demand, "doc_url", "doc_name")
	if err != nil {
		beego.Error("update demand err :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}
