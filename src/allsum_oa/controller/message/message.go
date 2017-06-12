package message

import (
	"allsum_oa/controller/base"
	"allsum_oa/service"
	"common/lib/errcode"
	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) GetLatestMsg() {
	uid := int(c.UserID)
	maxId, err := c.GetInt("max_msgid")
	if err != nil {
		beego.Error("parameters wrong!")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msg, err := service.GetMsgBatch(maxId, uid)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(msg)
}

func (c *Controller) DelOneMsg() {
	msgId, err := c.GetInt("msgid")
	if err != nil {
		beego.Error("parameters wrong!")
		c.ReplyErr(errcode.ErrParams)
		return
	}

	err = service.DelMsg(msgId)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc("ok")
}
