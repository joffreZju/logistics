package aggregatemgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/models"
	"allsum_bi/services/aggregation"
	"common/lib/errcode"
	"encoding/json"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListAggregate() {
	reportuuid := c.GetString("reportuuid")
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		beego.Error("get reportid err", err)
		c.ReplyErr(errcode.ErrActionGetReport)
	}
	limit, _ := c.GetInt("limit", 10)

	index, _ := c.GetInt("index", 0)
	aggregates, err := models.ListAggregateOpsByField([]string{"reportid"}, []interface{}{report.Id}, limit, index)
	if err != nil {
		beego.Error("list aggregate err", err)
		c.ReplyErr(errcode.ErrActionGetAggregate)
		return
	}
	var res []map[string]interface{}
	for _, v := range aggregates {
		resmap := map[string]interface{}{
			"index": v.Id,
			"uuid":  v.Uuid,
			"name":  v.Name,
		}
		res = append(res, resmap)
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) GetAggregate() {
	uuid := c.GetString("uuid")
	aggregateops, err := models.GetAggregateOpsByUuid(uuid)
	if err != nil {
		beego.Error("get aggregate err", err)
		c.ReplyErr(errcode.ErrActionGetAggregate)
		return
	}
	res := map[string]interface{}{
		"uuid":          aggregateops.Uuid,
		"reportid":      aggregateops.Reportid,
		"name":          aggregateops.Name,
		"create_script": aggregateops.CreateScript,
		"alter_script":  aggregateops.AlterScript,
		"script":        aggregateops.Script,
		"dest_table":    aggregateops.DestTable,
		"cron":          aggregateops.Cron,
		"documents":     aggregateops.Documents,
		"status":        aggregateops.Status,
	}
	c.ReplySucc(res)

}

func (c *Controller) SaveAggregate() {
	uuid := c.GetString("uuid")
	name := c.GetString("name")
	reportuuid := c.GetString("reportuuid")
	beego.Debug("uuid :", reportuuid, name)
	if reportuuid == "" {
		beego.Error("aggregate reprotuuid miss ")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		beego.Error("get report err:", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	if uuid == "" {
		aggregate := models.AggregateOps{
			Name:     name,
			Reportid: report.Id,
		}
		uuid, err := models.InsertAggregateOps(aggregate)
		if err != nil {
			beego.Error("insert aggregate err :", err)
			c.ReplyErr(errcode.ErrActionPutAggregate)
			return
		}
		res := map[string]string{
			"uuid": uuid,
		}
		c.ReplySucc(res)
		return
	}
	var reqbody map[string]interface{}
	err = json.Unmarshal(c.Ctx.Input.RequestBody, &reqbody)
	if err != nil {
		beego.Error("umarshal fail err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	create_table_name, ok1 := reqbody["create_table_name"]

	create_script := reqbody["create_script"]
	alter_script := reqbody["alter_script"]
	flush_script, ok4 := reqbody["flush_script"]
	cron, ok6 := reqbody["cron"]
	documents, ok7 := reqbody["documents"]
	if !(ok1 && ok4 && ok6 && ok7) {
		beego.Error("miss params :", ok1, ok4, ok6, ok7)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	//	aggregateMap := map[string]interface{}{
	//		"create_table_name": create_table_name,
	//		"create_script":     create_script,
	//		"alter_script":      alter_script,
	//		"flush_script":      flush_script,
	//		"cron":              cron,
	//		"ducuments":         documents,
	//	}
	err = aggregation.AddAggregate(uuid, create_table_name.(string), create_script.(string), alter_script.(string), flush_script.(string), cron.(string), documents.(string))
	if err != nil {
		beego.Error("add aggregate err :", err)
		c.ReplyErr(errcode.ErrActionPutAggregate)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) TestAggregateCreateScript() {
	create_script := c.GetString("create_script")
	table_name := c.GetString("table_name")
	uuid := c.GetString("uuid")
	err := aggregation.TestCreateScript(uuid, table_name, create_script)
	if err != nil {
		beego.Error("test create script :", err)
		c.ReplyErr(&errcode.CodeError{500, err.Error()})
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) TestAggregateAlterScript() {
	uuid := c.GetString("uuid")
	alter_script := c.GetString("alter_script")
	table_name := c.GetString("table_name")
	err := aggregation.TestAlterScript(uuid, table_name, alter_script)
	if err != nil {
		beego.Error("test alter script : ", err)
		c.ReplyErr(&errcode.CodeError{500, err.Error()})
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) TestAggregateFlushScript() {
	flush_script := c.GetString("flush_script")
	uuid := c.GetString("uuid")
	table_name := c.GetString("table_name")
	cron := c.GetString("cron")
	err := aggregation.TestFlushScript(uuid, table_name, flush_script, cron)
	if err != nil {
		beego.Error("test flush script err :", err)
		c.ReplyErr(&errcode.CodeError{500, err.Error()})
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}
