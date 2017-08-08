package admin

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"encoding/json"
	"github.com/astaxie/beego"
)

const commonErr = 99999

type Controller struct {
	base.Controller
}

func (c *Controller) AdminGetFirmInfo() {
	no := c.GetString("cno")
	f, err := model.GetCompany(no)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmNotExisted)
		return
	}
	c.ReplySucc(*f)
}

func (c *Controller) AdminGetFirmList() {
	companylist, err := service.GetCompanyList()
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(companylist)
}

func (c *Controller) AdminFirmAudit() {
	cno := c.GetString("cno")
	msg := c.GetString("msg")
	status, err := c.GetInt("status")
	if err != nil || (status != model.CompanyStatApproveAccessed && status != model.CompanyStatApproveNotAccessed) {
		beego.Error(err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	approver, e := service.GetUserById(c.UserComp, c.UserID)
	if e != nil {
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		beego.Error(e)
		return
	}
	err = service.AuditCompany(cno, approver, status, msg)
	if err != nil {
		c.ReplyErr(errcode.ErrServerError)
		beego.Error(err)
		return
	}
	c.ReplySucc(nil)
}

func (c *Controller) AdminAddFunction() {
	fstr := c.GetString("function")
	f := &model.Function{}
	e := json.Unmarshal([]byte(fstr), f)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.AddFunction(f)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) AdminUpdateFunction() {
	fstr := c.GetString("function")
	f := &model.Function{}
	e := json.Unmarshal([]byte(fstr), f)
	if e != nil || f.Pid != 0 || len(f.SysId) != 0 {
		//不能修改Pid,SysId
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.UpdateFunction(f)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) AdminDelFunction() { //todo 如果这个功能还有公司在用的话，那么不能删除
	fid, e := c.GetInt("id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.DelFunction(fid)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) GetAppVersionList() {
	vlist, e := service.GetAppVersionList()
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(vlist)
	}
}

func (c *Controller) AddAppVersion() {
	appstr := c.GetString("appversion")
	app := &model.AppVersion{}
	e := json.Unmarshal([]byte(appstr), app)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.AddAppVersion(app)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}
