package reportset

import (
	"allsum_bi/models"
	"allsum_bi/services/reportset"
	"allsum_bi/util"
	"allsum_bi/util/errcode"
	"allsum_bi/util/ossfile"
	"encoding/json"
	"stowage/common/controller/base"
	"strconv"

	"github.com/astaxie/beego"
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
	reportsets, err := models.ListReportSetByField([]string{}, []interface{}{}, limit, index)
	if err != nil {
		beego.Error("list dataload err: ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	var reportsetres []map[string]interface{}
	for _, reportset := range reportsets {
		reportmap := map[string]interface{}{
			"index": reportset.Id,
			"uuid":  reportset.Uuid,
			"name":  "reprot-" + strconv.Itoa(reportset.Reportid),
		}
		reportsetres = append(reportsetres, reportmap)
	}
	c.ReplySucc(reportsetres)
}

type condition struct {
	field   string      `json:"field"`
	Type    string      `json:"type"`
	Greater interface{} `json:"omitempty"`
	Smaller interface{} `json:"omitempty"`
	Enum    []string    `json:"omitempty"`
}

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
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	var conditions []condition
	err = json.Unmarshal([]byte(reportset.Conditions), &conditions)

	if err != nil {
		beego.Error("unmarshal condition err ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string]interface{}{
		"uuid":         reportset.Uuid,
		"name":         "report-" + strconv.Itoa(reportset.Reportid),
		"script":       reportset.Script,
		"conditions":   conditions,
		"web_path":     reportset.WebPath,
		"webfile_name": reportset.WebfileName,
		"status":       reportset.Status,
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
		reportid, ok := reqbody["reportid"]
		if !ok {
			beego.Error("miss Reportid")
			c.ReplyErr(errcode.ErrParams)
			return
		}
		reportset := models.ReportSet{
			Reportid: reportid.(int),
			Status:   util.REPORTSET_BUILDING,
		}
		uuidstr, err := models.InsertReportSet(reportset)
		if err != nil {
			beego.Error("insert report set err: ", err)
			c.ReplyErr(errcode.ErrServerError)
			return
		}
		res := map[string]string{
			"uuid": uuidstr,
		}
		c.ReplySucc(res)
		return
	}
	reportset, err := models.GetReportSetByUuid(uuid.(string))
	if err != nil {
		beego.Error("get report set err:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	getsql, ok := reqbody["get_script"]
	if !ok {
		beego.Error("miss get_script")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportset.Script = getsql.(string)
	conditionbytes := reqbody["conditions"].([]byte)
	var conditions condition
	err = json.Unmarshal(conditionbytes, &conditions)
	if err != nil {
		beego.Error("unmarshal json err", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	err = models.UpdateReportSet(reportset, "script", "conditions")
	if err != nil {
		beego.Error("update report set ")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) GetReportSetWebFile() {
	uuid := c.GetString("uuid")

	reportset, err := models.GetReportSetByReportUuid(uuid)
	if err != nil {
		beego.Error("get report err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	filedata, err := ossfile.GetFile(reportset.WebPath)
	if err != nil {
		beego.Error("get file err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplyFile("application/octet-stream", reportset.WebfileName, filedata)
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
	reportuuid, ok := reqbody["reportUuid"]
	if !ok {
		beego.Error("miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	conditions := reqbody["conditions"]
	conditionList := conditions.([]map[string]interface{})
	datas, err := reportset.GetData(reportuuid.(string), conditionList)
	if err != nil {
		beego.Error("get Data err", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	c.ReplySucc(datas)
}
