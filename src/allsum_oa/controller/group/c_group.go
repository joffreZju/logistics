package group

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"github.com/astaxie/beego"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999

//更新和增加组织属性
func (c *Controller) CreateAttr() {
	//uid := c.UserID
	//todo 检测用户权限，不符合直接返回
	ucomp := c.UserComp
	update := c.GetString("Update")
	a := &model.Attribute{
		No:   c.GetString("No"),
		Name: c.GetString("Name"),
		Desc: c.GetString("Desc"),
	}
	var e error
	if update == "true" {
		a.Utime = time.Now()
		e = service.UpdateAttr(ucomp, a)
	} else if update == "false" {
		a.Ctime = time.Now()
		e = service.CreateAttr(ucomp, a)
	} else {
		c.ReplyErr(errcode.ErrParams)
	}
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc("success")
	}
}
