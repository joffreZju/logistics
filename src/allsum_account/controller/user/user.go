package user

import (
	"common/lib/errcode"
	"common/lib/keycrypt"
	"common/lib/push"
	"common/lib/util"
	"controller/base"
	"fmt"
	"math/rand"
	"model"
	"service"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
)

type Controller struct {
	base.Controller
}

var AllsumUserList = []string{"15158134537"}

func getGroup(tel string) int {
	for _, v := range AllsumUserList {
		if tel == v {
			return 1
		}
	}
	return 0
}

func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	desc := c.GetString("desc")
	gender, _ := c.GetInt8("gender")
	addr := c.GetString("address")
	mail := c.GetString("mail")
	referer := c.GetString("referer")

	passwdc := keycrypt.Sha256Cal(passwd)
	beego.Debug("tel:", tel)
	var code int
	var err error
	if code, err = strconv.Atoi(c.GetString("code")); err != nil {
		beego.Error("user.regist error:", err)
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	}
	vcode := c.Cache.Get(tel)
	if vcode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if v, _ := strconv.Atoi(fmt.Sprintf("%s", vcode)); v != code {
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	} else {
		u := model.User{
			Tel:      tel,
			Password: passwdc,
			Descp:    desc,
			Gender:   gender,
			Address:  addr,
			Mail:     mail,
			Referer:  referer,
			UserType: 1,
			//CreateTime: time.Now(),
		}
		err := service.UserCreate(&u)
		if err != nil {
			beego.Error("user register failed", err)
			c.ReplyErr(errcode.ErrUserCreateFailed)
			return
		}

		//同时创建个人账户
		a := model.Account{
			AccountNo: util.RandomByte16(),
			Userid:    u.Id,
			UserType:  1,
			Status:    1,
		}
		err = service.AccountCreate(&a)
		if err != nil {
			beego.Error("create user account failed:", err)
		}

		c.ReplySucc("success")
	}
}

func (c *Controller) GetUserInfo() {
	user, err := service.GetUserInfo(int(c.UserID))
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(user)
	return
}

func (c *Controller) EditProfile() {
	gender, _ := c.GetInt8("gender")
	username := c.GetString("username")
	descp := c.GetString("descp")
	address := c.GetString("address")
	id := int(c.UserID)
	user := model.User{
		Id:       id,
		Gender:   gender,
		Descp:    descp,
		UserName: username,
		Address:  address,
	}

	err := service.UserUpdate(&user, "Gender", "UserName", "Descp", "Address")
	if err != nil {
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(user)
	return
}

//用户登陆
func (c *Controller) UserLogin() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	user, err := service.GetUserByTel(tel)
	if err != nil {
		beego.Error(errcode.ErrUserNotExisted)
		c.ReplyErr(errcode.ErrUserPasswordError)
		return
	}
	if !keycrypt.CheckSha256(passwd, user.Password) {
		c.ReplyErr(errcode.ErrUserPasswordError)
		return
	}

	g := getGroup(tel)
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), g, c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, *user)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	} else {
		user.LoginTime = time.Now()
		service.UserUpdate(user, "LoginTime")

		jsonstr := make(map[string]interface{})
		jsonstr["Token"] = token.Value
		jsonstr["User"] = user
		c.ReplySucc(jsonstr)
		beego.Debug("login ok,token:%+v", token)
		return
	}

}

func (c *Controller) UserLoginPhoneCode() {
	tel := c.GetString("tel")
	user, err := service.GetUserByTel(tel)
	if err != nil {
		beego.Error(errcode.ErrUserNotExisted)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}

	var code int
	if code, err = strconv.Atoi(c.GetString("code")); err != nil {
		beego.Error("user.login error:", err)
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	}
	vcode := c.Cache.Get(tel)
	if vcode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if v, _ := strconv.Atoi(fmt.Sprintf("%s", vcode)); v != code {
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	} else {
		token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), getGroup(tel), c.Ctx.ResponseWriter)
		if err != nil {
			beego.Error("o2o.Auth.NewSingleToken error:", err, *user)
			c.ReplyErr(errcode.ErrAuthCreateFailed)
			return
		} else {
			user.LoginTime = time.Now()
			service.UserUpdate(user, "LoginTime")

			jsonstr := make(map[string]interface{})
			jsonstr["Token"] = token.Value
			jsonstr["User"] = user
			c.ReplySucc(jsonstr)
			beego.Debug("login ok,token:%+v", token)
			return
		}

	}

}

func (c *Controller) LoginOut() {
	token, err := o2o.Auth.CheckToken(c.Ctx.Request)
	if err != nil {
		beego.Error("o2o.Auth.CheckToken error:", err)
		c.ReplySucc("OK")
		return
	}
	err = tokenauth.Store.DeleteToken(token.Value)
	if err != nil {
		beego.Error("tokenauth.Store.DeleteToken:", err)
	}
	c.ReplySucc("OK")
}
func (c *Controller) Resetpwd() {
	uid := (int)(c.UserID)
	if uid == 0 {
		uid, _ = c.GetInt("id")
	}
	user, err := service.GetUserInfo(uid)
	if err != nil {
		c.ReplyErr(err)
		return
	}
	if len(user.Password) == 0 {
		pwd := c.GetString("password")
		pwd = keycrypt.Sha256Cal(pwd)
		user.Password = pwd
		err = service.UserUpdate(user, "Password")
		if err != nil {
			c.ReplyErr(err)
			return
		} else {
			c.ReplySucc("OK")
		}
	} else {
		owd := c.GetString("oldpassword")
		owd = keycrypt.Sha256Cal(owd)
		if owd != user.Password {
			err = errcode.ErrUserPasswordError
			c.ReplyErr(err)
			return
		}
		pwd := c.GetString("password")
		pwd = keycrypt.Sha256Cal(pwd)
		user.Password = pwd
		err = service.UserUpdate(user, "Password")
		if err != nil {
			c.ReplyErr(err)
			return
		} else {
			c.ReplySucc("OK")
		}
	}
}

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

	//msg := fmt.Sprintf("您好，感谢您使用算配载服务，您的登录验证码是%v，验证码有效期为1分钟。", code)
	if push.SendMsgWithDayuToMobile(tel, code, "壹算科技") {
		c.ReplySucc("发送短信成功")
	} else {
		c.ReplyErr(errcode.ErrSendSMSMsgError)
	}
}
