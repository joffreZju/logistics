package group

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999

//更新和增加组织属性
func (c *Controller) AddAttr() {
	prefix := c.UserComp
	a := &model.Attribute{
		No:    c.GetString("No"),
		Name:  c.GetString("Name"),
		Desc:  c.GetString("Desc"),
		Ctime: time.Now(),
	}
	e := service.AddAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}
func (c *Controller) UpdateAttr() {
	prefix := c.UserComp
	a := &model.Attribute{
		No:    c.GetString("No"),
		Name:  c.GetString("Name"),
		Desc:  c.GetString("Desc"),
		Utime: time.Now(),
	}
	e := service.UpdateAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}

//新增组织树上下级
func (c *Controller) AddGroup() {
	uid := c.UserID
	prefix := c.UserComp
	newGroupStr := c.GetString("NewGroup")
	sonsStr := c.GetString("Sons")
	ng := &model.Group{}
	sons := make([]int, 0)
	e := json.Unmarshal([]byte(newGroupStr), ng)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if len(sonsStr) != 0 {
		e = json.Unmarshal([]byte(sonsStr), &sons)
		if e != nil {
			c.ReplyErr(errcode.ErrParams)
			beego.Error(e)
			return
		}
	}
	ng.AdminId = uid
	ng.CreatorId = uid
	ng.No = fmt.Sprintf("%d", time.Now().UnixNano()) //todo
	ng.Ctime = time.Now()
	e = service.AddGroup(prefix, ng, sons)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//合并
func (c *Controller) MergeGroups() {
	uid := c.UserID
	prefix := c.UserComp
	oldIdsStr := c.GetString("OldGroups")
	newGroupStr := c.GetString("NewGroup")
	ng := &model.Group{}
	oldIds := make([]int, 0)
	e := json.Unmarshal([]byte(newGroupStr), ng)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if len(oldIdsStr) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	} else {
		e = json.Unmarshal([]byte(oldIdsStr), &oldIds)
		if e != nil {
			c.ReplyErr(errcode.ErrParams)
			beego.Error(e)
			return
		}
	}
	ng.AdminId = uid
	ng.CreatorId = uid
	ng.No = fmt.Sprintf("%d", time.Now().UnixNano()) //todo
	ng.Ctime = time.Now()
	e = service.MergeGroups(prefix, ng, oldIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//转让升级
func (c *Controller) MoveGroup() {
	//uid := c.UserID
	prefix := c.UserComp
	gid, e := c.GetInt("GroupId")
	newPid, e2 := c.GetInt("NewPid")
	if e != nil || e2 != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e, e2)
		return
	}
	e = service.MoveGroup(prefix, gid, newPid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//删除组织
func (c *Controller) DelGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("GroupId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.DelGroup(prefix, gid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//编辑组织
func (c *Controller) EditGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("GroupId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	newName := c.GetString("NewName")
	e = service.EditGroup(prefix, newName, gid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//为组织添加用户
func (c *Controller) AddUsersToGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("GroupId")
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
	e = service.AddUsersToGroup(prefix, gid, uids)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

//从组织删除批量用户
func (c *Controller) DelUsersFromGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("GroupId")
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
	e = service.DelUsersFromGroup(prefix, gid, uids)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}
