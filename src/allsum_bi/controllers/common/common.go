package common

import (
	"allsum_bi/services/oa"
	base "common/lib/baseController"
	"common/lib/errcode"

	"github.com/astaxie/beego"
	_ "github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListDeveloper() {
	develops, err := oa.GetBiDevelops()
	if err != nil {
		beego.Error("oa get bi develops")
		c.ReplyErr(errcode.ErrServerError)
	}
	c.ReplySucc(develops)
}
