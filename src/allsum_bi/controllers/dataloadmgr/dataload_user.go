package dataloadmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/services/dataload"
	"allsum_bi/util"
	"allsum_bi/util/errcode"
	"encoding/json"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListData() {
	var reqdata map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqdata)
	if err != nil {
		beego.Error("umarshal fail err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	uuid, ok := reqdata["uuid"]
	if !ok {
		beego.Error("list data miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	condition, ok := reqdata["condition"]
	if !ok {
		condition = []map[string]interface{}{}
	}
	limit, ok := reqdata["limit"]
	if !ok {
		limit = 10
	}
	columns, datas, err := dataload.GetData(uuid, condition, limit)
	if err != nil {
		beego.Error("get data err ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	return
}

func (c *Controller) InputData() {
	var reqdata map[string]interface{}
	err := json.Unmarshal(c.Ctx.Input.RequestBody, &reqdata)
	if err != nil {
		beego.Error("umarshal fail err :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	uuid, ok := reqdata["uuid"]
	if !ok {
		beego.Error("miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	inputType, ok := reqdata["input_type"]
	if !ok {
		beego.Error("miss input_type")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	fields, ok := reqdata["fields"]
	if !ok {
		beego.Error("miss fields")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	data, ok := reqdata["data"]
	if !ok {
		beego.Error("miss data")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	if inputType == util.IS_UPDATE {
		err = dataload.UpdateData(uuid, fields, data)
		if err != nil {
			beego.Error("update data fail err : ", err)
			c.ReplyErr(errcode.ErrServerError)
			return
		}
	} else if inputType == util.IS_INSERT {
		err = dataload.InsertNewData(uuid, fields, data)
		if err != nil {
			beego.Error("insert data fail err : ", err)
			c.ReplyErr(errcode.ErrServerError)
			return
		}
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}
