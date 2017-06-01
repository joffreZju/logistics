package user

import (
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"common/lib/push"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
	"math/rand"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999
const host = "http://allsum.com:8093"

func (c *Controller) GetCode() {
	tel := c.GetString("tel")

	// 测试环境用123456，不发短信
	if beego.BConfig.RunMode != "prod" {
		err := c.Cache.Put(tel, 123456, time.Duration(300*time.Second))
		if err != nil {
			beego.Error("GetCode set redis error", err)
			c.ReplyErr(err)
			return
		}
		c.ReplySucc("success")
		return
	}

	// 先查看用户短信是否已经, 如果短信已经发送，60秒后重试
	if c.Cache.IsExist(tel) {
		c.ReplyErr(errcode.ErrUserCodeHasAlreadyExited)
		return
	}

	// 正试环境发短信, 60秒后过期
	code := fmt.Sprintf("%d", rand.Intn(9000)+1000)
	err := c.Cache.Put(tel, code, time.Duration(60*time.Second))
	if err != nil {
		beego.Error("GetCode set redis error", err)
		c.ReplyErr(err)
		return
	}

	if err = push.SendSmsCodeToMobile(tel, code); err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrSendSMSMsgError)
	} else {
		c.ReplySucc("success")
	}
}

func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	code := c.GetString("smscode")
	vcode := c.Cache.Get(tel)
	if vcode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if fmt.Sprintf("%s", vcode) != code {
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	}

	url := host + "/exempt/user/register"
	m := make(map[string]string)
	m["tel"] = tel
	m["password"] = c.GetString("password")
	m["addr"] = c.GetString("addr")
	m["desc"] = c.GetString("desc")
	m["gender"] = c.GetString("gender")
	req := httplib.Post(url)
	req.JSONBody(m)
	resp, e := req.Response()
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}

	bodystr, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}

	body := make(map[string]interface{})
	e = json.Unmarshal(bodystr, &body)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(body["data"])
}

func (c *Controller) UserLogin() {
	url := host + "/exempt/user/login_auth"
	m := make(map[string]string)
	m["tel"] = c.GetString("tel")
	m["password"] = c.GetString("password")
	req := httplib.Post(url)
	req.JSONBody(m)
	resp, e := req.Response()
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}

	bodystr, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}

	body := make(map[string]interface{})
	e = json.Unmarshal(bodystr, &body)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(body["data"])
}
