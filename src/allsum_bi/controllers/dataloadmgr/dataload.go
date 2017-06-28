package dataloadmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/models"
	"allsum_bi/services/aggregation"
	"allsum_bi/services/dataload"
	"allsum_bi/util"
	"allsum_bi/util/errcode"
	"allsum_bi/util/ossfile"
	"io/ioutil"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListDataload() {
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 20
	}
	index, err := c.GetInt("index")
	if err != nil {
		index = 0
	}
	dataloads, err := models.ListDataLoadByField([]string{}, []interface{}{}, limit, index)
	if err != nil {
		beego.Error("list dataload err: ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	var dataloadres []map[string]interface{}
	for _, dataload := range dataloads {
		dataloadmap := map[string]interface{}{
			"index": dataload.Id,
			"uuid":  dataload.Uuid,
			"name":  dataload.Name,
		}
		dataloadres = append(dataloadres, dataloadmap)
	}
	c.ReplySucc(dataloadres)
}

func (c *Controller) GetDataLoad() {
	uuid := c.GetString("uuid")
	if uuid == "" {
		beego.Error("get dataload miss uuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	dataload, err := models.GetDataLoadByUuid(uuid)
	if err != nil {
		beego.Error("get dataload db err ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	//TODO get Aggregateid
	aggregateid := dataload.Aggregateid
	var flush_script, cron string
	if aggregateid == 0 {
		flush_script = ""
		cron = ""
	} else {
		Aggregate, err := models.GetAggregateOps(aggregateid)
		if err != nil {
			beego.Error("get Aggregate ops err", err)
			c.ReplyErr(errcode.ErrServerError)
		}
		flush_script = Aggregate.Script
		cron = Aggregate.Cron
	}
	res := map[string]interface{}{
		"uuid":          dataload.Uuid,
		"name":          dataload.Name,
		"owner":         dataload.Owner,
		"table_name":    dataload.Basetable,
		"create_script": dataload.CreateScript,
		"alter_script":  dataload.AlterScript,
		"flush_script":  flush_script,
		"cron":          cron,
		"web_path":      dataload.WebPath,
		"webfile_name":  dataload.WebfileName,
		"documents":     dataload.Documents,
	}

	return
}

func (c *Controller) SaveDataload() {
	dataloadName := c.GetString("name")
	dataloadOwner, err := c.GetInt("owner")
	if err != nil {
		beego.Error("list dataloadOwner")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	uuid := c.GetString("uuid")

	if uuid == "" {
		dataloadstruct := models.DataLoad{
			Name:   dataloadName,
			Status: util.DATALOAD_BUILDING,
			Owner:  dataloadOwner,
		}
		uuid, err := models.InsertDataLoad(dataloadstruct)
		if err != nil {
			beego.Error("insert dataload err :", err)
			c.ReplyErr(errcode.ErrServerError)
			return
		}
		res := map[string]string{
			"uuid": uuid,
		}
		c.ReplySucc(res)
		return
	}
	table_name := c.GetString("table_name")
	create_script := c.GetString("create_script")
	flush_script := c.GetString("flush_script")
	alter_script := c.GetString("alter_script")
	cron := c.GetString("cron")
	documents := c.GetString("documents")
	dataloadmap := map[string]string{
		"uuid":          uuid,
		"name":          dataloadName,
		"table_name":    table_name,
		"create_script": create_script,
		"alter_script":  alter_script,
		"flush_script":  flush_script,
		"cron":          cron,
		"documents":     documents,
	}
	//have check
	err = dataload.AddDataLoad(dataloadmap)
	if err != nil {
		beego.Error("list dataload err :", err)
		c.ReplyErr(errcode.ErrServerError)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) TestDataLoadCreateScript() {
	create_script := c.GetString("create_script")
	table_name := c.GetString("table_name")
	dataload_uuid := c.GetString("uuid")
	err := dataload.TestCreateScript(dataload_uuid, table_name, create_script)
	if err != nil {
		beego.Error("test create script : ", err)
		c.ReplyErr(err)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) TestDataLoadAlterScript() {
	dataload_uuid := c.GetString("uuid")
	alter_script := c.GetString("alter_script")
	table_name := c.GetString("table_name")
	err := dataload.TestAlterScript(dataload_uuid, table_name, alter_script)
	if err != nil {
		beego.Error("test alter script : ", err)
		c.ReplyErr(err)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) TestAggregate() {
	flush_script := c.GetString("flush_script")
	dataload_uuid := c.GetString("uuid")
	table_name := c.GetString("table_name")
	cron := c.GetString("cron")
	err := aggregation.TestAddCronWithFlushScript(cron, flush_script)
	if err != nil {
		beego.Error("test flush script err :", err)
		c.ReplyErr(err)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) UploadDataLoadWeb() {
	dataloaduuid := c.GetString("uuid")
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
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	filename := h.Filename + uuid.NewV4().String()
	uri, err := ossfile.PutFile("dataload", filename, data)
	if err != nil {
		beego.Error("put file to oss err : ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}

	dataload, err := models.GetDataLoadByUuid(dataloaduuid)
	if err != nil {
		beego.Error("get dataload err : ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	dataload.WebPath = uri
	dataload.WebfileName = filename
	err = models.UpdateDataLoad(dataload, "web_path", "webfile_name")
	if err != nil {
		beego.Error("update dataload err :", err)
		c.ReplyErr(errcode.ErrServerError)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)

}