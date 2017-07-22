package kettlemgr

import (
	"allsum_bi/models"
	"allsum_bi/services/kettlesvs"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"encoding/json"
	"io/ioutil"
	"path"
	"strings"

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
			if !strings.Contains(kfname, ".ktr") {
				beego.Error("must ktr file :", kfname)
				c.ReplyErr(errcode.ErrUploadFileFailed)
				return
			}
			ktrdatamap[kfname] = ktrdata
		}
	}

	kettlejob, err := kettlesvs.AddJobKtrfile(Name, Cron, jobfilename, jobfiledata, ktrdatamap)
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

func (c *Controller) DownloadKJob() {
	uuid := c.GetString("uuid")
	kettlejob, err := models.GetKettleJobByUuid(uuid)
	if err != nil {
		return
	}
	//jobfile
	Kjbpath := kettlejob.Kjbpath
	var JobMap map[string][]byte
	var kettlejobmap map[string]string
	err = json.Unmarshal([]byte(Kjbpath), &kettlejobmap)
	if err != nil {
		beego.Error("unmarshall json err: ", err)
		c.ReplyErr(errcode.ErrActionGetJobInfo)
		return
	}
	urlpath := kettlejobmap["urlpath"]
	jobfilename := path.Base(urlpath)
	filedata, err := ossfile.GetFile(urlpath)
	if err != nil {
		beego.Error("get oss file err ", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	JobMap[jobfilename] = filedata
	//ktrfile
	var ktrmap map[string]string
	err = json.Unmarshal([]byte(kettlejob.Ktrpath), &ktrmap)
	if err != nil {
		beego.Error("json unmarshal err: ", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	for _, ktrurl := range ktrmap {
		ktrfilename := path.Base(ktrurl)
		ktrdata, err := ossfile.GetFile(ktrurl)
		if err != nil {
			beego.Error("this url cant load :", ktrurl, err)
			c.ReplyErr(errcode.ErrDownloadFileFailed)
			return
		}
		JobMap[ktrfilename] = ktrdata
	}
	zipdata, err := util.Zip(JobMap)
	if err != nil {
		beego.Error("zip err :", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	c.ReplyFile("application/octet-stream", "kettle.zip", zipdata)
	return

}

func (c *Controller) SetJobEnable() {
	kettleJoblimit, err := beego.AppConfig.Int("kettle::joblimit")
	if err != nil {
		kettleJoblimit = 5
	}
	var reqmap map[string]int
	reqbody := c.Ctx.Request.Body
	body, err := ioutil.ReadAll(reqbody)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &reqmap)
	if err != nil {
		beego.Error("unmarshal json err", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	setRunNum := 0
	for _, status := range reqmap {
		if status == util.KETTLEJOB_RIGHT {
			setRunNum += 1
		}
	}
	if setRunNum > kettleJoblimit {
		beego.Error("run num err", setRunNum, kettleJoblimit)
		c.ReplyErr(errcode.ErrActionSetJobNum)
		return
	}
	kettlejobs, err := models.ListKettleJobByField([]string{}, []interface{}{}, 0, 0)
	if err != nil {
		beego.Error("list kettle job err:", err)
		c.ReplyErr(errcode.ErrActionGetJobInfo)
		return
	}
	for _, job := range kettlejobs {
		status := reqmap[job.Uuid]
		job.Status = status
		if status == util.KETTLEJOB_RIGHT {
			var kjbmap map[string]string
			err = json.Unmarshal([]byte(job.Kjbpath), &kjbmap)
			if err != nil {
				beego.Error(" err: ", err)
				c.ReplyErr(errcode.ErrActionSetJobNum)
				return
			}
			kettleWorkPath := beego.AppConfig.String("kettle::workpath")
			kettlesvs.AddCron(job.Id, job.Cron, kettleWorkPath+path.Base(kjbmap["urlpath"]))
			err = models.UpdateKettleJob(job, "status")
			if err != nil {
				beego.Error(" err: ", err)
				c.ReplyErr(errcode.ErrActionSetJobNum)
				return
			}

		} else if status == util.KETTLEJOB_FAIL {
			kettlesvs.StopCron(job.Id)
			err = models.UpdateKettleJob(job, "status")
			if err != nil {
				beego.Error(" err: ", err)
				c.ReplyErr(errcode.ErrActionSetJobNum)
				return
			}
		}
	}
}

func (c *Controller) DeleteKJob() {
	uuid := c.GetString("uuid")
	err := models.DeleteKettleJobByUuid(uuid)
	if err != nil {
		beego.Error(" err: ", err)
		c.ReplyErr(errcode.ErrActionGetJobInfo)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}
