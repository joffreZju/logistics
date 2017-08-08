package api

import (
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
	service.ApiService
}

const commonErr = 99999

func (c *Controller) GetSchemaList() {
	schemas, e := c.ApiService.GetSchemaList()
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(schemas)
	}
}
