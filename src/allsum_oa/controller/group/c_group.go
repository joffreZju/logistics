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
	desc := c.GetString("descrp")
	if len(name) == 0 || len(desc) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a := &model.Attribute{
		No:     model.UniqueNo("GA"),
		Name:   name,
		Descrp: desc,
		Ctime:  time.Now(),
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
	desc := c.GetString("descrp")
	if e != nil || len(name) == 0 || len(desc) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a := &model.Attribute{
		Id:     aid,
		Name:   name,
		Descrp: desc,
		Utime:  time.Now(),
	}
	e = service.UpdateAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) DelAttr() {
	prefix := c.UserComp
	aid, e := c.GetInt("id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a := &model.Attribute{
		Id: aid,
	}
	e = service.DelAttr(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) getBeginTimeOfOperation() (t time.Time, e error) {
	timeStr := c.GetString("beginTime")
	t = time.Now()
	if len(timeStr) != 0 {
		t, e = time.ParseInLocation(model.TimeFormatWithLocal, timeStr, time.Local)
		if e != nil {
			return
		}
	}
	t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	//todo beginTime可以设置一天中任何时间，方便测试
	return t, nil
}

//新增组织树上下级
func (c *Controller) AddGroup() {
	//检测是否有未提交的修改
	e := service.CheckFutureGroupOperation(c.UserComp)
	if e != nil {
		c.ReplyErr(errcode.ErrUpdateGroupTree)
		return
	}
	uid := c.UserID
	prefix := c.UserComp
	newGroupStr := c.GetString("newGroup")
	sonsStr := c.GetString("sons")
	desc := c.GetString("descrp")
	beginTime, e := c.getBeginTimeOfOperation()
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	ng := &model.Group{}
	sons := make([]int, 0)
	e = json.Unmarshal([]byte(newGroupStr), ng)
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
		e = service.AddRootGroup(prefix, desc, beginTime, ng)
	} else {
		e = service.AddGroup(prefix, desc, beginTime, ng, sons)
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
	//检测是否有未提交的修改
	e := service.CheckFutureGroupOperation(c.UserComp)
	if e != nil {
		c.ReplyErr(errcode.ErrUpdateGroupTree)
		return
	}
	uid := c.UserID
	prefix := c.UserComp
	oldIdsStr := c.GetString("oldGroups")
	newGroupStr := c.GetString("newGroup")
	desc := c.GetString("descrp")
	beginTime, e := c.getBeginTimeOfOperation()
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	ng := &model.Group{}
	oldIds := make([]int, 0)
	e = json.Unmarshal([]byte(newGroupStr), ng)
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
	e = service.MergeGroups(prefix, desc, beginTime, ng, oldIds)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//转让升级
func (c *Controller) MoveGroup() {
	//检测是否有未提交的修改
	e := service.CheckFutureGroupOperation(c.UserComp)
	if e != nil {
		c.ReplyErr(errcode.ErrUpdateGroupTree)
		return
	}
	prefix := c.UserComp
	gid, e := c.GetInt("id")
	newPid, e2 := c.GetInt("newPid")
	desc := c.GetString("descrp")
	if e != nil || e2 != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e, e2)
		return
	}
	beginTime, e := c.getBeginTimeOfOperation()
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.MoveGroup(prefix, desc, beginTime, gid, newPid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//删除组织
func (c *Controller) DelGroup() {
	//检测是否有未提交的修改
	e := service.CheckFutureGroupOperation(c.UserComp)
	if e != nil {
		c.ReplyErr(errcode.ErrUpdateGroupTree)
		return
	}
	prefix := c.UserComp
	desc := c.GetString("descrp")
	gid, e := c.GetInt("id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	beginTime, e := c.getBeginTimeOfOperation()
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.DelGroup(prefix, desc, beginTime, gid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//更新组织
func (c *Controller) UpdateGroup() {
	//检测是否有未提交的修改
	e := service.CheckFutureGroupOperation(c.UserComp)
	if e != nil {
		c.ReplyErr(errcode.ErrUpdateGroupTree)
		return
	}
	prefix := c.UserComp
	desc := c.GetString("descrp")
	str := c.GetString("group")
	g := new(model.Group)
	e = json.Unmarshal([]byte(str), g)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	beginTime, e := c.getBeginTimeOfOperation()
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	//更新属性
	e = service.UpdateGroup(prefix, desc, beginTime, g)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//获取历史组织操作记录
func (c *Controller) GetGroupOpList() {
	prefix := c.UserComp
	limit, e := c.GetInt("limit")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if limit == 0 || limit < -1 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	//limit = -1 所有记录
	ops, e := service.GetGroupOpList(prefix, limit)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(ops)
	}
}

//搜索历史组织操作记录
func (c *Controller) SearchGroupOpsByTime() {
	prefix := c.UserComp
	begin, e := time.ParseInLocation(model.TimeFormatWithLocal, c.GetString("begin"), time.Local)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	end, e := time.ParseInLocation(model.TimeFormatWithLocal, c.GetString("end"), time.Local)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	ops, e := service.SearchGroupOpsByTime(prefix, begin, end)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(ops)
	}
}

//获取历史组织树详情
func (c *Controller) GetGroupOpDetail() {
	prefix := c.UserComp
	opId, e := c.GetInt("opId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	groups, e := service.GetGroupOpDetail(prefix, opId)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(groups)
	}
}

//取消未生效的组织树操作记录
func (c *Controller) CancelGroupOp() {
	prefix := c.UserComp
	//t, e := time.ParseInLocation(model.TimeFormatWithLocal, c.GetString("beginTime"), time.Local)
	//if e != nil {
	//	c.ReplyErr(errcode.ErrParams)
	//	beego.Error(e)
	//	return
	//}
	//t = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
	opId, e := c.GetInt("opId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.CancelGroupOp(prefix, opId)
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

//获取组织节点下的所有用户
func (c *Controller) GetUsersOfGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("id")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	users, e := service.GetUsersOfGroup(prefix, gid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(users)
}

//为组织添加用户
func (c *Controller) AddUsersToGroup() {
	prefix := c.UserComp
	gid, e := c.GetInt("groupId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("users")
	uids := []int{}
	e = json.Unmarshal([]byte(usersStr), &uids)
	if e != nil || len(uids) == 0 {
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
	gid, e := c.GetInt("groupId")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	usersStr := c.GetString("users")
	uids := make([]int, 0)
	e = json.Unmarshal([]byte(usersStr), &uids)
	if e != nil || len(uids) == 0 {
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
