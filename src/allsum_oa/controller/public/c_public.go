package public

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"common/lib/push"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"time"
)

const commonErr = 99999

type Controller struct {
	base.Controller
}

func (c *Controller) Test() {
	a, e := c.RedisClient.Hmset("group_1", map[string]interface{}{
		"roles":  "1_2_3",
		"groups": "4_5_6",
	})
	beego.Info(a, e)

	a, e = c.RedisClient.Hmset("group_1", map[string]interface{}{
		"roles":  "1_2",
		"groups": "4_5",
	})
	beego.Info(a, e)

	b, e := c.RedisClient.Hmget("group_1", []string{"roles", "groups"})
	beego.Info(b, e)

	d, e := c.RedisClient.HDel("group_1", "roles")
	beego.Info(d, e)

	f, e := c.RedisClient.Hmget("group_1", []string{"roles", "groups"})
	beego.Info(f, e, len(f["roles"]))

	g, e := c.RedisClient.Expire("group_1", 120)
	beego.Info(g, e)

	g, e = c.RedisClient.Expire("group_1", 30)
	beego.Info(g, e)

	c.ReplySucc(nil)
}

func (c *Controller) GetCode() {
	tel := c.GetString("tel")

	// 测试环境用123456，不发短信
	if beego.BConfig.RunMode != "prod" {
		err := c.Cache.Put(tel, 123456, time.Duration(600*time.Second))
		if err != nil {
			beego.Error("GetCode set redis error", err)
			c.ReplyErr(err)
			return
		}
		c.ReplySucc(nil)
		return
	}

	// 先查看用户短信是否已经, 如果短信已经发送，60秒后重试
	if c.Cache.IsExist(tel) {
		c.ReplyErr(errcode.ErrUserCodeHasAlreadyExited)
		return
	}

	// 正试环境发短信, 60秒后过期
	code := fmt.Sprintf("%d", rand.Intn(9000)+1000)
	err := c.Cache.Put(tel, code, time.Duration(600*time.Second))
	if err != nil {
		beego.Error("GetCode set redis error", err)
		c.ReplyErr(err)
		return
	}

	if err = push.SendSmsCodeToMobile(tel, code); err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrSendSMSMsgError)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) GetFunctionsTree() {
	idstr := c.GetString("sysIds")
	sysIds := []string{}
	e := json.Unmarshal([]byte(idstr), &sysIds)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	funcs, e := model.GetFunctions(sysIds)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(funcs)
}

func (c *Controller) GetLatestAppVersion() {
	app, e := service.GetLatestAppVersion()
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(app)
	}
}
