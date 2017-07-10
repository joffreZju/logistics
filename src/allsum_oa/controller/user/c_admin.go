package user

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"encoding/json"
	"github.com/astaxie/beego"
)

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
	uid := c.UserID
	cno := c.GetString("cno")
	status, err := c.GetInt("status")
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msg := c.GetString("msg")
	err = service.AuditCompany(cno, uid, status, msg)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
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
	if e != nil || f.Pid != 0 {
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

func (c *Controller) AdminDelFunction() {
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

func (c *Controller) GetLatestAppVersion() {
	app, e := service.GetLatestAppVersion()
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(app)
	}
}
