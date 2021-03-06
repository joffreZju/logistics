package user

import (
	"allsum_account/controller/base"
	"allsum_account/model"
	"allsum_account/service"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"common/lib/push"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
)

type Controller struct {
	base.Controller
}

var AllsumUserList = []string{"15158134537", "15558085697", "18867543358", "18694582678", "18667907711", "13735544671", "15301107268"}

func getGroup(tel string) string {
	for _, v := range AllsumUserList {
		if tel == v {
			return "aaa"
		}
	}
	return "bbb"
}
func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	desc := c.GetString("desc")
	addr := c.GetString("address")
	mail := c.GetString("mail")

	var err error

	passwdc := keycrypt.Sha256Cal(passwd)
	beego.Debug("tel:", tel)

	u := model.User{
		Tel:      tel,
		Password: passwdc,
		Descp:    desc,
		Address:  addr,
		Mail:     mail,
		UserType: 1,
		//CreateTime: time.Now(),
	}
	err = service.UserCreate(&u)
	if err != nil {
		beego.Error("user register failed", err)
		c.ReplyErr(err)
		return
	}
	c.ReplySucc(u)

}

/*
func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	desc := c.GetString("desc")
	addr := c.GetString("address")
	mail := c.GetString("mail")

	var err error

	passwdc := keycrypt.Sha256Cal(passwd)
	beego.Debug("tel:", tel)
	var code int
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
			Address:  addr,
			Mail:     mail,
			UserType: 1,
			//CreateTime: time.Now(),
		}
		err := service.UserCreate(&u)
		if err != nil {
			beego.Error("user register failed", err)
			c.ReplyErr(err)
			return
		}

		jsonstr := make(map[string]interface{})
		jsonstr["tel"] = tel
		jsonstr["password"] = passwd
		c.ReplySucc(jsonstr)
	}
}*/

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

func (c *Controller) UserLoginAuth() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	user, err := service.GetUserByTel(tel)
	if err != nil {
		beego.Error(errcode.ErrUserNotExisted)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}
	if !keycrypt.CheckSha256(passwd, user.Password) {
		c.ReplyErr(errcode.ErrUserPasswordError)
		return
	}
	c.ReplySucc(user)

	return
}

//获取用户信息和其公司信息
func (c *Controller) GetUserCompanys() {
	tel := c.GetString("tel")
	user, err := service.GetUserByTel(tel)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}
	c.ReplySucc(user)
}

//用户登陆
func (c *Controller) UserLogin() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	user, err := service.GetUserByTel(tel)
	if err != nil {
		beego.Error(errcode.ErrUserNotExisted)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}
	if len(passwd) < 6 {
		beego.Error(errcode.ErrUserNeedInit)
	}
	if !keycrypt.CheckSha256(passwd, user.Password) {
		c.ReplyErr(errcode.ErrUserPasswordError)
		return
	}
	var comNo string
	if len(user.Companys) == 1 {
		comNo = user.Companys[0].No
	} else {
		comNo = "none"
	}

	var bizInfo string = ""
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), comNo, bizInfo, c.Ctx.ResponseWriter)
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

func (c *Controller) Retrievepwd() {
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
		pwd := c.GetString("password")
		pwdc := keycrypt.Sha256Cal(pwd)
		user.Password = pwdc
		err = service.UserUpdate(user, "Password")
		if err != nil {
			c.ReplyErr(err)
			return
		}
		retstr := make(map[string]interface{})
		retstr["tel"] = tel
		retstr["password"] = pwd
		c.ReplySucc(retstr)
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
		var comNo string
		if len(user.Companys) == 1 {
			comNo = user.Companys[0].No
		} else {
			comNo = "none"
		}

		token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), comNo, "", c.Ctx.ResponseWriter)
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

	// 测试环境用1234，不发短信
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
	if err = push.SendSmsCodeToMobile(tel, code); err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrSendSMSMsgError)
	} else {
		c.ReplySucc("发送短信成功")
	}
}

func uniqueNo(prefix string) string {
	str := strings.Replace(time.Now().Format("0102150405.000"), ".", "", 1)
	str = prefix + str
	return str
}

func (c *Controller) FirmRegister() {
	uid, _ := c.GetInt("uid")
	name := c.GetString("firm_name")
	desc := c.GetString("desc")
	phone := c.GetString("phone")
	lf := c.GetString("license_file")
	tp, _ := c.GetInt("firm_type")

	firm := model.Company{
		No:          uniqueNo("O"),
		Creater:     uid,
		FirmName:    name,
		Desc:        desc,
		Phone:       phone,
		LicenseFile: lf,
		FirmType:    tp,
		Status:      0,
	}
	err := model.InsertCompany(&firm)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmCreateFailed)
		return
	}
	model.AddUserToCompany(firm.No, uid)
	c.ReplySucc("ok")
}

func (c *Controller) GetFirmInfo() {
	no := c.GetString("no")
	f, err := model.GetCompany(no)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmNotExisted)
		return
	}
	c.ReplySucc(*f)
}

func (c *Controller) GetFirmList() {
	list, err := model.GetCompanies()
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(list)

}

func (c *Controller) FirmModify() {
	no := c.GetString("no")
	name := c.GetString("name")
	desc := c.GetString("desc")
	phone := c.GetString("phone")
	lf := c.GetString("license_file")
	firm := model.Company{
		No:          no,
		FirmName:    name,
		Desc:        desc,
		Phone:       phone,
		LicenseFile: lf,
	}
	err := model.UpdateCompany(&firm)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmUpdate)
		return
	}
	c.ReplySucc("ok")
}
func (c *Controller) FirmAudit() {
	uid, _ := c.GetInt("uid")
	cno := c.GetString("cno")
	status, err := c.GetInt("status")
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msg := c.GetString("msg")
	err = model.AuditCompany(cno, uid, status, msg)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc("ok")
}

func (c *Controller) FirmDelUser() {
	uno, _ := c.GetInt("uid")
	cno := c.GetString("cno")
	err := model.DelCompanyUser(cno, uno)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc("ok")
}

func (c *Controller) FirmAddUser() {
	//uid, _ := c.GetInt("uid")
	cno := c.GetString("cno")
	tel := c.GetString("tel")

	err := model.AddCompanyUser(cno, tel)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc("ok")
}
