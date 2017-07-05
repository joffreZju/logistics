package reportsetmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/models"
	"allsum_bi/services/reportset"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	"common/lib/errcode"
	"encoding/json"
	"io/ioutil"
	"strconv"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListReportSet() {
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 20
	}
	index, err := c.GetInt("index")
	if err != nil {
		index = 0
	}
	reportuuid := c.GetString("reportuuid")
	if reportuuid == "" {
		beego.Error("miss reportuuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		beego.Error("error reportuuid")
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	reportsets, err := models.ListReportSetByField([]string{"reportid"}, []interface{}{report.Id}, limit, index)
	if err != nil {
		beego.Error("list dataload err: ", err)
		c.ReplyErr(errcode.ErrActionGetReportSet)
		return
	}
	var reportsetres []map[string]interface{}
	for _, reportset := range reportsets {
		reportmap := map[string]interface{}{
			"index":   reportset.Id,
			"uuid":    reportset.Uuid,
			"name":    "reprot-" + strconv.Itoa(reportset.Reportid),
			"webpath": reportset.WebPath,
		}
		reportsetres = append(reportsetres, reportmap)
	}
	c.ReplySucc(reportsetres)
}

//type condition struct {
//	field   string      `json:"field"`
//	Type    string      `json:"type"`
//	Greater interface{} `json:"omitempty"`
//	Smaller interface{} `json:"omitempty"`
//	Enum    []string    `json:"omitempty"`
//}

func (c *Controller) GetReportSet() {
	uuid := c.GetString("uuid")
	if uuid == "" {
		beego.Error("get reportset miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportset, err := models.GetReportSetByUuid(uuid)
	if err != nil {
		beego.Error("get reportset db err ", err)
		c.ReplyErr(errcode.ErrActionGetReportSet)
		return
	}
	var conditions []map[string]interface{}
	err = json.Unmarshal([]byte(reportset.Conditions), &conditions)

	if err != nil {
		beego.Error("unmarshal condition err ", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	res := map[string]interface{}{
		"uuid":       reportset.Uuid,
		"name":       "report-" + strconv.Itoa(reportset.Reportid),
		"script":     reportset.Script,
		"conditions": conditions,
		"web_path":   reportset.WebPath,
		//		"webfile_name": reportset.WebfileName,
		"status": reportset.Status,
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) SaveReportSet() {
	var reqbody map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqbody)
	if err != nil {
		beego.Error("json unmarshal err:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	uuid, ok := reqbody["uuid"]
	if !ok {
		reportuuid, ok := reqbody["reportuuid"]
		if !ok {
			beego.Error("miss Reportid")
			c.ReplyErr(errcode.ErrParams)
			return
		}
		report, err := models.GetReportByUuid(reportuuid.(string))
		if err != nil {
			beego.Error("get Report by uuid err", err)
			c.ReplyErr(errcode.ErrActionGetReport)
			return

		}
		reportsetdb := models.ReportSet{
			Reportid: report.Id,
			Status:   util.REPORTSET_BUILDING,
		}
		uuidstr, err := models.InsertReportSet(reportsetdb)
		if err != nil {
			beego.Error("insert report set err: ", err)
			c.ReplyErr(errcode.ErrActionPutReportSet)
			return
		}
		res := map[string]string{
			"uuid": uuidstr,
		}
		c.ReplySucc(res)
		return
	}
	reportsetdb, err := models.GetReportSetByUuid(uuid.(string))
	if err != nil {
		beego.Error("get report set err:", err)
		c.ReplyErr(errcode.ErrActionGetReportSet)
		return
	}
	getsql, ok := reqbody["get_script"]
	if !ok {
		beego.Error("miss get_script")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportsetdb.Script = getsql.(string)
	conditioninterface := reqbody["conditions"]
	//	var conditions []map[string]interface{}
	jsoncondition, err := json.Marshal(conditioninterface)
	if err != nil {
		beego.Error("marshal condition json err", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportsetdb.Conditions = string(jsoncondition)
	webpath := reqbody["web_path"].(string)
	reportsetdb.WebPath = webpath
	checkres := reportset.CheckConditionFormat(reportsetdb.Conditions)
	if !checkres {
		beego.Error("check report condition format ")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportsetdb.Status = util.REPORTSET_STARTED
	err = models.UpdateReportSet(reportsetdb, "script", "conditions", "status", "web_path")
	if err != nil {
		beego.Error("update report set ")
		c.ReplyErr(errcode.ErrActionPutReportSet)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) UploadReportSetWeb() {
	reportsetuuid := c.GetString("uuid")
	f, h, err := c.GetFile("uploadfile")
	if err != nil {
		beego.Error("get file err : ", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	f.Close()
	data, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Error("read filehandle err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	filename := uuid.NewV4().String() + "-" + h.Filename
	uri, err := ossfile.PutFile("reportset", filename, data)
	if err != nil {
		beego.Error("put file to oss err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}

	reportset, err := models.GetReportSetByUuid(reportsetuuid)
	if err != nil {
		beego.Error("get reportset err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	reportset.WebPath = uri
	reportset.WebfileName = h.Filename
	err = models.UpdateReportSet(reportset, "web_path", "webfile_name")
	if err != nil {
		beego.Error("update reportset err :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) GetReportSetWebFile() {
	//	uuid := c.GetString("uuid")
	//
	//	reportset, err := models.GetReportSetByReportUuid(uuid)
	//	if err != nil {
	//		beego.Error("get report err :", err)
	//		c.ReplyErr(errcode.ErrDownloadFileFailed)
	//		return
	//	}
	//	filedata, err := ossfile.GetFile(reportset.WebPath)
	//	if err != nil {
	//		beego.Error("get file err :", err)
	//		c.ReplyErr(errcode.ErrDownloadFileFailed)
	//		return
	//	}
	//	c.ReplyFile("application/octet-stream", reportset.WebfileName, filedata)
	return
}

func (c *Controller) GetReportData() {
	var reqbody map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqbody)
	if err != nil {
		beego.Error("unmarshal json err", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	//	reportuuid, ok := reqbody["report_uuid"]
	reportsetuuid, ok1 := reqbody["reportset_uuid"]
	if !ok1 {
		beego.Error("miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	conditions := reqbody["conditions"]
	condition_interfaces := conditions.([]interface{})
	var conditionList []map[string]interface{}
	for _, v := range condition_interfaces {
		conditionList = append(conditionList, v.(map[string]interface{}))
	}
	//	conditionList := conditions.([]map[string]interface{})
	datas, err := reportset.GetData(reportsetuuid.(string), conditionList)
	if err != nil {
		beego.Error("get Data err", err)
		c.ReplyErr(errcode.ErrActionGetReportData)
		return
	}
	c.ReplySucc(datas)
}
