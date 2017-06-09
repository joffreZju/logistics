package group

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"encoding/json"
	"github.com/astaxie/beego"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999

//更新和增加组织属性
func (c *Controller) GetAttrList() {
	prefix := c.UserComp
	al, e := service.GetAttrList(prefix)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(al)
}

func (c *Controller) AddAttr() {
	prefix := c.UserComp
	name := c.GetString("name")
	desc := c.GetString("desc")
	if len(name) == 0 || len(desc) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a := &model.Attribute{
		No:    model.UniqueNo("GA"),
		Name:  name,
		Desc:  desc,
		Ctime: time.Now(),
	}
	e := service.AddAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateAttr() {
	prefix := c.UserComp
	aid, e := c.GetInt("id")
	name := c.GetString("name")
	desc := c.GetString("desc")
	if e != nil || len(name) == 0 || len(desc) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a := &model.Attribute{
		Id:    aid,
		Name:  name,
		Desc:  desc,
		Utime: time.Now(),
	}
	e = service.UpdateAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
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
	ng.No = model.UniqueNo("G")
	ng.Ctime = time.Now()
	if ng.Pid == 0 && len(sons) == 0 {
		e = service.AddRootGroup(prefix, ng)
	} else {
		e = service.AddGroup(prefix, ng, sons)
	}
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
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
	ng.No = model.UniqueNo("G")
	ng.Ctime = time.Now()
	e = service.MergeGroups(prefix, ng, oldIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//转让升级
func (c *Controller) MoveGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("Id")
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
		c.ReplySucc(nil)
	}
}

//删除组织
func (c *Controller) DelGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("Id")
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
		c.ReplySucc(nil)
	}
}

//更新组织
func (c *Controller) UpdateGroup() {
	prefix := c.UserComp
	str := c.GetString("group")
	g := new(model.Group)
	e := json.Unmarshal([]byte(str), g)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	//更新属性
	e = service.UpdateGroup(prefix, g)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//获取所有组织节点
func (c *Controller) GetGroupList() {
	prefix := c.UserComp
	gs, e := service.GetGroupList(prefix)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(gs)
}

//todo next begin
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
		c.ReplySucc(nil)
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
		c.ReplySucc(nil)
	}
}
