package user

import (
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"common/lib/push"
	"fmt"
	"github.com/astaxie/beego"
	"math/rand"
	"time"
)

type Controller struct {
	base.Controller
}

const CommonErr = 99999

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

	url := "/exempt/user/register"
	m := make(map[string]string)
	m["tel"] = tel
	m["password"] = c.GetString("password")
	m["addr"] = c.GetString("addr")
	m["desc"] = c.GetString("desc")
	m["gender"] = c.GetString("gender")
	resp, ecode := c.post_account(url, m)
	if ecode != nil {
		c.ReplyErr(ecode)
		//beego.Error(ecode)
		fmt.Println(ecode)

		return
	}
	beego.Debug(resp["data"])
	c.ReplySucc(resp["data"])
}

func (c *Controller) UserLogin() {
	url := "/exempt/user/login_auth"
	m := make(map[string]string)
	m["tel"] = c.GetString("tel")
	m["password"] = c.GetString("password")
	resp, ecode := c.post_account(url, m)
	if ecode != nil {
		c.ReplyErr(ecode)
		beego.Error(ecode.Error())
		return
	}
	c.ReplySucc(resp["data"])
}

func (c *Controller) Resetpwd() {

}

func (c *Controller) GetUserInfo() {

}

func (c *Controller) EditProfile() {

}

//用户注册公司相关操作
func (c *Controller) FirmRegister() {

}
func (c *Controller) FirmModify() {

}
func (c *Controller) FirmAddUser() {

}
func (c *Controller) FirmDelUser() {

}

//allsum管理员接口,审核公司
func (c *Controller) GetFirmList() {

}
func (c *Controller) FirmAudit() {

}
func (c *Controller) GetFirmInfo() {

}
