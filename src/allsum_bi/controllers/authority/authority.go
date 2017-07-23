package authority

import (
	"allsum_bi/services/userauthority"
	"allsum_oa/controller/base"
	"common/lib/errcode"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) AddUserReportsetAuthority() {
	companyid := c.GetString("companyid")
	roleid, err := c.GetInt("roleid")
	if err != nil {
		beego.Error("get role id fail")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	rolename := c.GetString("rolename")
	reportsetuuid := c.GetString("reportsetuuid")
	err = userauthority.AddUserReportsetAuthority(companyid, roleid, rolename, reportsetuuid)
	if err != nil {
		beego.Error("add user reportset authority err :", err)
		c.ReplyErr(errcode.ErrActionAddUserAuthority)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) RemoveReportSetAuthority() {
	companyid := c.GetString("companyid")
	roleid, err := c.GetInt("roleid")
	if err != nil {
		beego.Error("get role id fail")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportsetuuid := c.GetString("reportsetuuid")
	err = userauthority.RemoveReportSetAuthority(companyid, roleid, reportsetuuid)
	if err != nil {
		beego.Error("remove report set authority err : ", err)
		c.ReplyErr(errcode.ErrActionDeleteUserAuthirity)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) RemoveReportAuthority() {
	companyid := c.GetString("companyid")
	roleid, err := c.GetInt("roleid")
	if err != nil {
		beego.Error("get role id fail")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	reportuuid := c.GetString("reportuuid")
	err = userauthority.RemoveReportSetAuthority(companyid, roleid, reportuuid)
	if err != nil {
		beego.Error("remove report set authority err : ", err)
		c.ReplyErr(errcode.ErrActionDeleteUserAuthirity)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
	return
}
