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

//todo 检测用户权限，不符合直接返回
//在进行组织树管理时前端先校验用户权限？每一个controller里面还需要校验吗？

//更新和增加组织属性
func (c *Controller) UpsertAttr() {
	//uid := c.UserID
	prefix := c.UserComp
	update := c.GetString("Update")
	a := &model.Attribute{
		No:   c.GetString("No"),
		Name: c.GetString("Name"),
		Desc: c.GetString("Desc"),
	}
	var e error
	if update == "true" {
		//todo flag
		a.Utime = time.Now()
		e = service.UpdateAttr(prefix, a)
	} else if update == "false" {
		a.Ctime = time.Now()
		e = service.CreateAttr(prefix, a)
	} else {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}

//新增上下级
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
