package testmgr

import (
	"allsum_bi/models"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	base "common/lib/baseController"
	"common/lib/errcode"
	"encoding/json"
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
	c.ReplySucc(testinfos)
	//	var documents []string
	//	var filepaths []string
	//	for _, testinfo := range testinfos {
	//		documents = append(documents, testinfo.Documents)
	//		for _, filep := range testinfo.Filepaths {
	//			filepaths = append(filepaths, filep)
	//		}
	//	}
	//
	//	res := map[string]interface{}{
	//		"documents": documents,
	//		"filepaths": filepaths,
	//	}
	//	c.ReplySucc(res)
}

func (c *Controller) AddTestFile() {
	f, h, err := c.GetFile("uploadfile")
	if err != nil {
		beego.Error("get report err ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	data, err := ioutil.ReadAll(f)
	if err != nil {
		beego.Error("read filehandler err :", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	f.Close()
	filename := uuid.NewV4().String() + "-" + h.Filename
	uripath, err := ossfile.PutFile("test_info", filename, data)
	if err != nil {
		beego.Error("put file to oss err : ", err)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	res := map[string]string{
		"uripath": uripath,
	}
	c.ReplySucc(res)
}

func (c *Controller) AddTest() {
	reqbody := c.Ctx.Input.RequestBody
	var reqmap map[string]interface{}
	err := json.Unmarshal(reqbody, &reqmap)
	if err != nil {
		beego.Error("unmarshal json :", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportuuid, ok := reqmap["reportuuid"]
	if !ok {
		beego.Error("miss reportuuid")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	title, ok := reqmap["title"]
	if !ok {
		beego.Error("miss title")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	document, ok := reqmap["document"]
	if !ok {
		beego.Error("miss documents")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	weburis, ok := reqmap["weburis"]
	if !ok {
		beego.Error("miss weburis")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	weburiinterfaces := weburis.([]interface{})
	webpaths := []string{}
	for _, weburiinterface := range weburiinterfaces {
		webpaths = append(webpaths, weburiinterface.(string))
	}
	report, err := models.GetReportByUuid(reportuuid.(string))
	if err != nil {
		beego.Error("get report err", err)
		c.ReplyErr(errcode.ErrActionGetReport)
		return
	}
	demand, err := models.GetDemand(report.Demandid)
	if err != nil {
		beego.Error("get demand err :", err)
		c.ReplyErr(errcode.ErrActionGetDemand)
		return
	}

	//	var uripaths []string
	//	fileform := c.Ctx.Request.MultipartForm
	//	filelist := fileform.File
	//	for uploadkey, _ := range filelist {
	//		f, h, err := c.GetFile(uploadkey)
	//		if err != nil {
	//			beego.Error("get file err :", err)
	//			//c.ReplyErr(errcode.ErrParams)
	//			continue
	//		}
	//		data, err := ioutil.ReadAll(f)
	//		if err != nil {
	//			beego.Error("read filehandle err : ", err)
	//			//c.ReplyErr(errcode.ErrUploadFileFailed)
	//			continue
	//		}
	//		f.Close()
	//		filename := uuid.NewV4().String() + "-" + h.Filename
	//		uripath, err := ossfile.PutFile("test_info", filename, data)
	//		if err != nil {
	//			beego.Error("put file to oss err : ", err)
	//			continue
	//		}
	//		uripaths = append(uripaths, uripath)
	//	}
	//	documents := c.GetString("document")
	testinfo := models.TestInfo{
		Reportid:    report.Id,
		Testerid:    c.UserID,
		Title:       title.(string),
		Handlerid:   demand.Handlerid,
		HandlerName: demand.HandlerName,
		Documents:   document.(string),
		Filepaths:   webpaths,
		Status:      util.IS_OPEN,
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

func (c *Controller) RepairBug() {
	testuuid := c.GetString("uuid")
	testinfo := models.TestInfo{
		Uuid:   testuuid,
		Status: util.IS_CLOSE,
	}
	err := models.UpdateTestInfoByUuid(testinfo, "uuid")
	if err != nil {
		beego.Error("update testinfo by uuid :", testuuid, err)
		c.ReplyErr(errcode.ErrActionPutTestData)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}
