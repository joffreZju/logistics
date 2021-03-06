package api

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
	service.ApiService
}

const commonErr = 99999

//获取所有公司的schema
func (c *Controller) GetSchemaList() {
	schemas, e := c.ApiService.GetSchemaList()
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(schemas)
	}
}

//获取某公司的角色列表
func (c *Controller) GetRoleList() {
	company := c.GetString("companyNo")
	roles, e := service.GetRoleList(company)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(roles)
	}
}

//获取角色中的所有用户
func (c *Controller) GetUsersOfRole() {
	company := c.GetString("companyNo")
	rid, e := c.GetInt("roleId")
	if e != nil || rid == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	users, e := service.GetUsersOfRole(company, rid)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(users)
	}
}

//获取用户的基本信息
func (c *Controller) GetUserInfo() {
	uid, e := c.GetInt("userId")
	if e != nil || uid == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	user, e := service.GetUserById(model.Public, uid)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(user)
	}
}
