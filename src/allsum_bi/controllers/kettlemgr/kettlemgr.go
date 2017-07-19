package kettlemgr

import (
	"allsum_bi/models"
	"allsum_bi/services/kettle"
	"allsum_bi/util"
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"io/ioutil"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) AddKJob() {
	Name := c.GetString("name")
	Cron := c.GetString("cron")
	f, h, err := c.GetFile("kettlejob")
	if err != nil {
		beego.Error("get file err:", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Error("ioread file err ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	kettlejob, err := kettle.AddJobfile(Name, Cron, h.Filename, data)
	if err != nil {
		beego.Error("add job file :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}

	res := map[string]interface{}{
		"uuid":   kettlejob.Uuid,
		"name":   Name,
		"cron":   Cron,
		"lock":   kettlejob.Lock,
		"status": kettlejob.Status,
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) ListKJob() {
	limit, err := c.GetInt("limit")
	if err != nil {
		limit = 10
	}
	index, err := c.GetInt("index")
	if err != nil {
		index = 0
	}
	kettlejobs, err := models.ListKettleJobByField([]string{"status"}, []interface{}{util.KETTLEJOB_RIGHT}, limit, index)
	if err != nil {
		beego.Error("list kettle job err", err)
		c.ReplyErr(errcode.ErrActionGetJobInfo)
		return
	}
	var res []map[string]interface{}
	for _, kettlejob := range kettlejobs {
		subres := map[string]interface{}{
			"uuid":   kettlejob.Uuid,
			"name":   kettlejob.Name,
			"cron":   kettlejob.Cron,
			"lock":   kettlejob.Lock,
			"status": kettlejob.Status,
		}
		res = append(res, subres)
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) AddKtr() {

}

func (c *Controller) SetJobEnable() {

}

func (c *Controller) DeleteKJob() {

}
