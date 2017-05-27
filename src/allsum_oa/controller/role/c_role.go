package group

import (
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"encoding/json"
	"github.com/astaxie/beego"

	"allsum_oa/model"
	"allsum_oa/service"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999

func (c *Controller) AddRole() {
	prefix := c.UserComp
	name := c.GetString("Name")
	desc := c.GetString("Desc")
	funcIdsStr := c.GetString("FuncIds")
	var funcIds []int
	e := json.Unmarshal([]byte(funcIdsStr), &funcIds)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	r := model.Role{
		Name:  name,
		Desc:  desc,
		Ctime: time.Now(),
	}
	e = service.AddRole(prefix, &r, funcIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}

func (c *Controller) UpdateRole() {
	prefix := c.UserComp
	id, e := c.GetInt("Id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	name := c.GetString("Name")
	desc := c.GetString("Desc")
	funcIdsStr := c.GetString("FuncIds")
	var funcIds []int
	e = json.Unmarshal([]byte(funcIdsStr), &funcIds)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	r := model.Role{
		Id:   id,
		Name: name,
		Desc: desc,
	}
	e = service.UpdateRole(prefix, &r, funcIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}

func (c *Controller) DelRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("Id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.DelRole(prefix, rid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}

//为组织添加用户
func (c *Controller) AddUsersToRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("RoleId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("Users")
	uids := make([]int, 0)
	e = json.Unmarshal([]byte(usersStr), &uids)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.AddUsersToRole(prefix, rid, uids)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//从组织删除批量用户
func (c *Controller) DelUsersFromRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("RoleId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("Users")
	uids := make([]int, 0)
	e = json.Unmarshal([]byte(usersStr), &uids)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.DelUsersFromRole(prefix, rid, uids)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}
