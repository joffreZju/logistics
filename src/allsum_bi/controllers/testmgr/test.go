package testmgr

import (
	"allsum_bi/models"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	base "common/lib/baseController"
	"common/lib/errcode"
	"io/ioutil"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) GetTestInfo() {
	reportuuid := c.GetString("reportuuid")
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		beego.Error("get report err", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	testinfos, err := models.GetTestInfoByReportid(report.Id)
	if err != nil {
		beego.Error("get testinfo err: ", err)
		c.ReplyErr(errcode.ErrActionGetTestInfo)
		return
	}
	var documents []string
	var filepaths [][]string
	for _, testinfo := range testinfos {
		documents = append(documents, testinfo.Documents)
		filepaths = append(filepaths, testinfo.Filepaths.([]string))
	}

	res := map[string]interface{}{
		"documents": documents,
		"filepaths": filepaths,
	}
	c.ReplySucc(res)
}

func (c *Controller) AddTest() {
	reportuuid := c.GetString("reportuuid")
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		beego.Error("get report err", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	var uripaths []string
	fileform := c.Ctx.Request.MultipartForm
	filelist := fileform.File
	for uploadkey, _ := range filelist {
		f, h, err := c.GetFile(uploadkey)
		if err != nil {
			beego.Error("get file err :", err)
			//c.ReplyErr(errcode.ErrParams)
			continue
		}
		data, err := ioutil.ReadAll(f)
		if err != nil {
			beego.Error("read filehandle err : ", err)
			//c.ReplyErr(errcode.ErrUploadFileFailed)
			continue
		}
		f.Close()
		filename := uuid.NewV4().String() + "-" + h.Filename
		uripath, err := ossfile.PutFile("test_info", filename, data)
		if err != nil {
			beego.Error("put file to oss err : ", err)
			continue
		}
		uripaths = append(uripaths, uripath)
	}
	documents := c.GetString("document")
	testinfo := models.TestInfo{
		Reportid:  report.Id,
		Documents: documents,
		Filepaths: uripaths,
		Status:    util.IS_OPEN,
	}
	_, err = models.InsertTestInfo(testinfo)
	if err != nil {
		beego.Error("add test info err :", err)
		c.ReplyErr(errcode.ErrActionPutTestData)
		return
	}
	res := map[string]string{
		"res": "ok",
	}
	c.ReplySucc(res)
}

func (c *Controller) GetTestFile() {
	path := c.GetString("path")
	filedata, err := ossfile.GetFile(path)
	if err != nil {
		beego.Error("get file err : ", err)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	c.ReplyFile("application/octet-stream", path, filedata)
	return
}
