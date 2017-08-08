package msg

import (
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

const (
	commonErr = 99999
)

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

func (c *Controller) GetMsgsByPage() {
	page, e := c.GetInt("page")
	limit, e1 := c.GetInt("limit")
	if e != nil || e1 != nil {
		c.ReplyErr(e)
		beego.Error(e, e1)
		return
	}
	sum, msgs, e := service.GetMsgsByPage(c.UserComp, c.UserID, page, limit)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(map[string]interface{}{
			"Sum":  sum,
			"Msgs": msgs,
		})
	}
}
