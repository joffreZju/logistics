package msg

import (
	"allsum_oa/controller/base"
	"allsum_oa/service"
	"common/lib/errcode"
	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) GetHistoryMsg() {
	uid := c.UserID
	company := c.UserComp
	minId, err := c.GetInt("minId")
	if err != nil {
		beego.Error("parameters wrong:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msgs, err := service.GetHistoryMsg(company, uid, minId)
	if err != nil {
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(msgs)
}

func (c *Controller) GetLatestMsg() {
	uid := c.UserID
	company := c.UserComp
	maxId, err := c.GetInt("maxId")
	if err != nil {
		beego.Error("parameters wrong:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msgs, err := service.GetLatestMsg(company, uid, maxId)
	if err != nil {
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(msgs)
}

func (c *Controller) DelMsgById() {
	msgId, err := c.GetInt("id")
	if err != nil {
		beego.Error("parameters wrong:", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	err = service.DelMsgById(msgId)
	if err != nil {
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(nil)
}
