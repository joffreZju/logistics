package kettlemgr

import (
	"allsum_bi/models"
	"allsum_bi/services/kettlesvs"
	"allsum_bi/services/util"
	"allsum_bi/services/util/ossfile"
	base "common/lib/baseController"
	"common/lib/errcode"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"path"
	"strings"
	"time"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) AddKJob() {
	reqbody := c.Ctx.Input.RequestBody
	var reqmap map[string]interface{}
	err := json.Unmarshal(reqbody, &reqmap)
	Name, ok := reqmap["name"]
	if !ok {
		beego.Error("miss name")
		c.ReplyErr(errcode.ErrParams)
		return

	}
	Cron, ok := reqmap["cron"]
	if !ok {
		beego.Error("miss cron")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	var jobfilemap map[string]interface{}
	jobfile, ok := reqmap["jobfile"]
	if !ok {
		beego.Error("miss jobfile")
		c.ReplyErr(errcode.ErrParams)
		return
	}

	jobfilemap = jobfile.(map[string]interface{})
	jobfilename, ok := jobfilemap["name"]
	if !ok {
		beego.Error("miss jobfilename")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	jobfileuri, ok := jobfilemap["uri"]
	if !ok {
		beego.Error("miss joburi")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	jobdata, err := ossfile.GetFile(jobfileuri.(string))
	if err != nil {
		beego.Error("get job file err: ", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}

	var jobktrs []map[string]string
	jobktrinterfaces, ok := reqmap["jobktrs"]
	if !ok {
		beego.Error("miss jobktrs")
		c.ReplyErr(errcode.ErrParams)
		return
	}

	for _, jobktr := range jobktrinterfaces.([]interface{}) {
		jobktrmap := jobktr.(map[string]interface{})
		jobktrstrmap := map[string]string{}
		for k, v := range jobktrmap {
			jobktrstrmap[k] = v.(string)
		}

		jobktrs = append(jobktrs, jobktrstrmap)
	}
	kettlejob, err := kettlesvs.AddJobKtrfile(c.UserID, Name.(string), Cron.(string), jobfilename.(string), jobdata, jobktrs)
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
	var Ktrpath []map[string]string
	err = json.Unmarshal([]byte(kettlejob.Ktrpaths), &Ktrpath)
	if err != nil {
		Ktrpath = []map[string]string{}
	}
	res := map[string]interface{}{
		"uuid":    kettlejob.Uuid,
		"name":    Name,
		"cron":    Cron,
		"kjbpath": Kjbpath,
		"ktrpath": Ktrpath,
		"status":  kettlejob.Status,
	}
	c.ReplySucc(res)
	return
}
func (c *Controller) AddJobFile() {
	filetype := c.Ctx.Request.Header.Get("filetype")
	f, h, err := c.GetFile("uploadfile")
	if err != nil {
		beego.Error("get file err :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	jobfiledata, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Error("read file err :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	jobfilename := h.Filename
	var uri string
	if filetype == "job" {
		uri, err = ossfile.PutFile("kettle", uuid.NewV4().String()+jobfilename, jobfiledata)
		if err != nil {
			beego.Error("upload file err ", err)
			c.ReplyErr(errcode.ErrUploadFileFailed)
			return
		}
	} else {
		ktruriname := fmt.Sprintf("%d_%v_%s", c.UserID, time.Now(), jobfilename)
		kettleWorkPath := beego.AppConfig.String("kettle::workpath")
		err = ioutil.WriteFile(kettleWorkPath+ktruriname, jobfiledata, 0664)
		if err != nil {
			beego.Error("upload ktr file err :", err)
			c.ReplyErr(errcode.ErrUploadFileFailed)
			return
		}
		uri, err = ossfile.PutFile("kettle", ktruriname, jobfiledata)
		if err != nil {
			beego.Error("upload ktr file err :", err)
			c.ReplyErr(errcode.ErrUploadFileFailed)
			return
		}
	}
	res := map[string]string{
		"name": jobfilename,
		"uri":  uri,
	}
	c.ReplySucc(res)
}

//已经废弃。前端当时阶段做不了多文件上传
func (c *Controller) AddKJobOLD() {
	Name := c.GetString("name")
	Cron := c.GetString("cron")
	fileform := c.Ctx.Request.MultipartForm
	filelist := fileform.File
	var jobfilename string
	var jobfiledata []byte
	ktrdatamap := map[string][]byte{}
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

	kettlejob, err := kettlesvs.AddJobKtrfile_OLD(Name, Cron, jobfilename, jobfiledata, ktrdatamap)
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
	err = json.Unmarshal([]byte(kettlejob.Ktrpaths), &Ktrpath)
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
	kettlejobs, err := models.ListKettleJobByField([]string{}, []interface{}{}, 0, 0)
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
		var ktrmap []map[string]interface{}
		err = json.Unmarshal([]byte(kettlejob.Ktrpaths), &ktrmap)
		if err != nil {
			ktrmap = []map[string]interface{}{}
		}
		subres := map[string]interface{}{
			"uuid": kettlejob.Uuid,
			"name": kettlejob.Name,
			"cron": kettlejob.Cron,
			"kjb":  kjbmap,
			"ktr":  ktrmap,
			//		"lock":   kettlejob.Lock,
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
	JobMap := map[string][]byte{}
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
	var ktrmaps []map[string]string
	err = json.Unmarshal([]byte(kettlejob.Ktrpaths), &ktrmaps)
	if err != nil {
		beego.Error("json unmarshal err: ", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	for _, ktrmap := range ktrmaps {
		ktrfilename := path.Base(ktrmap["uri"])
		ktrdata, err := ossfile.GetFile(ktrmap["uri"])
		if err != nil {
			beego.Error("this url cant load :", ktrmap["uri"], err)
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
	reqbody := c.Ctx.Input.RequestBody
	err = json.Unmarshal(reqbody, &reqmap)
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
			kettlesvs.AddCron(job.Id, job.Cron, path.Base(kjbmap["urlpath"]))
			jobMap := map[string]interface{}{
				"id":     job.Id,
				"status": job.Status,
			}
			err = models.UpdateKettleJob(jobMap, "status")
			if err != nil {
				beego.Error(" err: ", err)
				c.ReplyErr(errcode.ErrActionSetJobNum)
				return
			}

		} else if status == util.KETTLEJOB_FAIL {
			kettlesvs.StopCron(job.Id)
			jobMap := map[string]interface{}{
				"id":     job.Id,
				"status": job.Status,
			}
			err = models.UpdateKettleJob(jobMap, "status")
			if err != nil {
				beego.Error(" err: ", err)
				c.ReplyErr(errcode.ErrActionSetJobNum)
				return
			}
		}
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
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
