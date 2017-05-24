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
	//prefix := c.UserComp
	//id, e := c.GetInt("Id")
	//if e != nil {
	//	c.ReplyErr(errcode.ErrParams)
	//	beego.Error(e)
	//	return
	//}
}
