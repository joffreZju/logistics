package message

import (
	"allsum_oa/controller/base"
	"github.com/astaxie/beego"
	"common/lib/errcode"
	"allsum_oa/service"
)

type Controller struct{
	base.Controller
}

func (c *Controller)GetLatestMsg(){
	uid:=int(c.UserId)
	maxId,err := c.GetInt("max_msgid")
	if err!=nil{
		beego.Error("parameters wrong!")
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msg,err:=service.GetMsgBatch(maxId,uid)
	if err!=nil{
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(msg)
}

func (c *Controller)DelOneMsg(){
	msgId,err := c.GetInt("msgid")
	if err!=nil{
		beego.Error("parameters wrong!")
		c.ReplyErr(errcode.ErrParams)
		return
	}

	err=service.DelMsg(msgId)
	if err!=nil{
		c.ReplyErr(error)
		return
	}
	c.ReplySucc("ok")
}


