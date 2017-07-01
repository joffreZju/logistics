package role

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

func (c *Controller) GetRoleList() {
	prefix := c.UserComp
	roles, e := service.GetRoles(prefix)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(roles)
	}
}

func (c *Controller) AddRole() {
	prefix := c.UserComp
	name := c.GetString("name")
	desc := c.GetString("descrp")
	funcIdsStr := c.GetString("functions")
	var funcIds []int
	e := json.Unmarshal([]byte(funcIdsStr), &funcIds)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	r := model.Role{
		Name:   name,
		Descrp: desc,
		Ctime:  time.Now(),
	}
	e = service.AddRole(prefix, &r, funcIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateRole() {
	prefix := c.UserComp
	id, e := c.GetInt("id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	name := c.GetString("name")
	desc := c.GetString("descrp")
	funcIdsStr := c.GetString("functions")
	var funcIds []int
	e = json.Unmarshal([]byte(funcIdsStr), &funcIds)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	r := model.Role{
		Id:     id,
		Name:   name,
		Descrp: desc,
	}
	e = service.UpdateRole(prefix, &r, funcIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) DelRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("id")
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
		c.ReplySucc(nil)
	}
}

//获取角色中的所有用户
func (c *Controller) GetUsersOfRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("id")
	users, e := service.GetUsersOfRole(prefix, rid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(users)
	}
}

//为角色添加用户
func (c *Controller) AddUsersToRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("roleId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("users")
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
		c.ReplySucc(nil)
	}
}

//从角色中删除批量用户
func (c *Controller) DelUsersFromRole() {
	prefix := c.UserComp
	rid, e := c.GetInt("roleId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("users")
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
		c.ReplySucc(nil)
	}
}
