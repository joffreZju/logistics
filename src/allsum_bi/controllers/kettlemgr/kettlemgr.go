package kettlemgr

import (
	"allsum_bi/models"
	"allsum_bi/services/kettle"
	"allsum_bi/util"
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"encoding/json"
	"io/ioutil"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) AddKJob() {
	Name := c.GetString("name")
	Cron := c.GetString("cron")
	fileform := c.Ctx.Request.MultipartForm
	filelist := fileform.File
	var jobfilename string
	var jobfiledata []byte
	var ktrdatamap map[string][]byte
	for k, _ := range filelist {
		if k == "kettlejob" {
			f, h, err := c.GetFile("kettlejob")
			if err != nil {
				beego.Error("get file err:", err)
				c.ReplyErr(errcode.ErrUploadFileFailed)
				return
			}
			jobfiledata, err = ioutil.ReadAll(f)
			if err != nil {
				beego.Error("ioread file err ", err)
				c.ReplyErr(errcode.ErrUploadFileFailed)
				return
			}
			jobfilename = h.Filename
		} else {
			f, h, err := c.GetFile(k)
			if err != nil {
				beego.Error("get ktr file err :", err)
				c.ReplyErr(errcode.ErrUploadFileFailed)
				return
			}
			ktrdata, err := ioutil.ReadAll(f)
			if err != nil {
				beego.Error("ioread file err ", err)
				c.ReplyErr(errcode.ErrUploadFileFailed)
				return
			}
			kfname := h.Filename
			ktrdatamap[kfname] = ktrdata
		}
	}

	kettlejob, err := kettle.AddJobKtrfile(Name, Cron, jobfilename, jobfiledata, ktrdatamap)
	if err != nil {
		beego.Error("add job file :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	var Kjbpath map[string]string
	err = json.Unmarshal([]byte(kettlejob.Kjbpath), &Kjbpath)
	if err != nil {
		Kjbpath = map[string]string{}
	}
	var Ktrpath map[string]string
	err = json.Unmarshal([]byte(kettlejob.Ktrpath), &Ktrpath)
	if err != nil {
		Ktrpath = map[string]string{}
	}
	res := map[string]interface{}{
		"uuid":    kettlejob.Uuid,
		"name":    Name,
		"cron":    Cron,
		"lock":    kettlejob.Lock,
		"kjbpath": Kjbpath,
		"ktrpath": Ktrpath,
		"status":  kettlejob.Status,
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) ListKJob() {
	kettlejobs, err := models.ListKettleJobByField([]string{"status"}, []interface{}{util.KETTLEJOB_RIGHT}, 0, 0)
	if err != nil {
		beego.Error("list kettle job err", err)
		c.ReplyErr(errcode.ErrActionGetJobInfo)
		return
	}
	var res []map[string]interface{}
	for _, kettlejob := range kettlejobs {
		var kjbmap map[string]string
		err = json.Unmarshal([]byte(kettlejob.Kjbpath), &kjbmap)
		if err != nil {
			kjbmap = map[string]string{}
		}
		var ktrmap map[string]interface{}
		err = json.Unmarshal([]byte(kettlejob.Ktrpath), &ktrmap)
		if err != nil {
			ktrmap = map[string]interface{}{}
		}
		subres := map[string]interface{}{
			"uuid":   kettlejob.Uuid,
			"name":   kettlejob.Name,
			"cron":   kettlejob.Cron,
			"kjb":    kjbmap,
			"ktr":    ktrmap,
			"lock":   kettlejob.Lock,
			"status": kettlejob.Status,
		}
		res = append(res, subres)
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) SetJobEnable() {

}

func (c *Controller) DeleteKJob() {

}
