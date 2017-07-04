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
	msgId, err := c.GetInt("msgId")
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
