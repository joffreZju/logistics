package user

import (
	"allsum_oa/controller/base"
	"allsum_oa/service"
	"common/lib/errcode"
	"common/lib/push"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/bitly/go-simplejson"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
	"math/rand"
	"strconv"
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

func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	code := c.GetString("smscode")
	mycode := c.Cache.Get(tel)
	if mycode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if fmt.Sprintf("%s", mycode) != code {
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
	resp, e := c.post_account(url, m)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.Data["json"] = resp
	c.ServeJSON()
}

//解析allsum_account的返回data字段，获取uid和companys的no
func (c *Controller) getUidAndCmps(data interface{}) (uid int, cmpsNo []string, cmps []interface{}, e error) {
	b, e := json.Marshal(data)
	if e != nil {
		return 0, nil, nil, e
	}
	js, e := simplejson.NewJson(b)
	if e != nil {
		return 0, nil, nil, e
	}
	uid, e = js.Get("Id").Int()
	if e != nil {
		return 0, nil, nil, e
	}
	cmpsNo = []string{}
	for i := 0; i >= 0; i++ {
		no, e := js.Get("Companys").GetIndex(i).Get("No").String()
		if e != nil {
			break
		} else {
			cmpsNo = append(cmpsNo, no)
		}
	}
	cmps, e = js.Get("Companys").Array()
	if e != nil {
		beego.Debug(e)
		cmps = nil
	}
	return uid, cmpsNo, cmps, nil
}

func (c *Controller) GetUserCompanys() {
	url := "/exempt/user/getcompanys"
	m := make(map[string]string)
	m["tel"] = c.GetString("tel")
	resp, e := c.get_account(url, m)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	if resp.Code != 0 {
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	_, _, cmps, e := c.getUidAndCmps(resp.Data)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(cmps)
}

func (c *Controller) UserLogin() {
	company := c.GetString("company")
	url := "/exempt/user/login_auth"
	m := make(map[string]string)
	m["tel"] = c.GetString("tel")
	m["password"] = c.GetString("password")
	resp, e := c.post_account(url, m)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	if resp.Code != 0 {
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	uid, cmps, _, e := c.getUidAndCmps(resp.Data)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	flag := true
	for _, v := range cmps {
		if company == v {
			flag = false
		}
	}
	if flag {
		company = ""
	}
	var bizInfo string = ""
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(uid), company, bizInfo, c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, uid)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	} else {
		if len(company) != 0 {
			user, e := service.GetUserById(company, uid)
			if e == nil {
				c.ReplySucc(user)
				return
			}
		}
		c.Data["json"] = resp
		c.ServeJSON()
		beego.Info("login ok,token:%+v", token)
		return
	}
}

func (c *Controller) UserLoginPhone() {
	company := c.GetString("company")
	tel := c.GetString("tel")
	code := c.GetString("code")
	mycode := c.Cache.Get(tel)
	if mycode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if fmt.Sprintf("%s", mycode) != code {
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	}
	//请求用户信息
	url := "/exempt/user/getcompanys"
	m := make(map[string]string)
	m["tel"] = tel
	resp, e := c.get_account(url, m)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	if resp.Code != 0 {
		c.Data["json"] = resp
		c.ServeJSON()
		return
	}
	//准备token信息
	uid, cmps, _, e := c.getUidAndCmps(resp.Data)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	flag := true
	for _, v := range cmps {
		if company == v {
			flag = false
		}
	}
	if flag {
		company = ""
	}
	var bizInfo string = ""
	//登录
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(uid), company, bizInfo, c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, uid)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	} else {
		if len(company) != 0 {
			user, e := service.GetUserById(company, uid)
			if e == nil {
				c.ReplySucc(user)
				return
			}
		}
		c.Data["json"] = resp
		c.ServeJSON()
		beego.Info("login ok,token:%+v", token)
		return
	}
}

func (c *Controller) LoginOut() {
	token, err := o2o.Auth.CheckToken(c.Ctx.Request)
	if err != nil {
		beego.Error("o2o.Auth.CheckToken error:", err)
		c.ReplySucc(nil)
		return
	}
	err = tokenauth.Store.DeleteToken(token.Value)
	if err != nil {
		beego.Error("tokenauth.Store.DeleteToken:", err)
	}
	beego.Info("login_out success token:", token.Value)
	c.ReplySucc(nil)
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
