package dataloadmgr

import (
	"allsum_bi/services/dataload"
	"allsum_bi/services/util"
	"common/lib/errcode"
	"encoding/json"

	"github.com/astaxie/beego"
)

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
	var conditionMap []map[string]interface{}
	condition, ok := reqdata["condition"]
	if !ok {
		conditionMap = []map[string]interface{}{}
	} else {
		for _, v := range condition.([]interface{}) {
			conditionMap = append(conditionMap, v.(map[string]interface{}))
		}
	}
	limit, ok := reqdata["limit"]
	if !ok {
		limit = 10
	}

	columns, datas, err := dataload.GetData(uuid.(string), conditionMap, int(limit.(float64)))
	if err != nil {
		beego.Error("get data err ", err)
		c.ReplyErr(errcode.ErrActionGetDataload)
		return
	}
	res := map[string]interface{}{
		"columns": columns,
		"datas":   datas,
	}
	c.ReplySucc(res)
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
	fieldinterface, ok := reqdata["fields"]
	if !ok {
		beego.Error("miss fields")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	var fields []string
	for _, v := range fieldinterface.([]interface{}) {
		fields = append(fields, v.(string))
	}
	datainterface, ok := reqdata["data"]
	if !ok {
		beego.Error("miss data")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	var data []map[string]interface{}
	for _, v := range datainterface.([]interface{}) {
		data = append(data, v.(map[string]interface{}))
	}
	if inputType == util.IS_UPDATE {
		err = dataload.UpdateData(uuid.(string), fields, data)
		if err != nil {
			beego.Error("update data fail err : ", err)
			c.ReplyErr(errcode.ErrActionInputData)
			return
		}
	} else if inputType == util.IS_INSERT {
		err = dataload.InsertNewData(uuid.(string), fields, data)
		if err != nil {
			beego.Error("insert data fail err : ", err)
			c.ReplyErr(errcode.ErrActionInputData)
			return
		}
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}
