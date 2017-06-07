package user

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"common/lib/push"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
	"math/rand"
	"strconv"
	"strings"
	"time"
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
	pwd := c.GetString("password")
	firm_type := c.GetString("firm_type")
	firm_name := c.GetString("firm_name")

	pwdEncode := keycrypt.Sha256Cal(pwd)
	beego.Info("register tel:", tel)

	u := model.User{
		Tel:      tel,
		No:       uniqueNo("U"),
		Password: pwdEncode,
		UserType: model.UserTypeNormal,
		Status:   model.UserStatusOk,
		Ctime:    time.Now(),
	}
	err := model.CreateUser("public", &u)
	if err != nil {
		beego.Error(err)
		if strings.Contains(err.Error(), "duplicate key") {
			err = errcode.ErrUserAlreadyExisted
		} else {
			err = errcode.ErrUserCreateFailed
		}
		c.ReplyErr(err)
		return
	}
	comp := model.Company{
		No:       uniqueNo("C"),
		FirmName: firm_name,
		FirmType: firm_type,
		Creator:  u.Id,
		Status:   0,
	}
	err = model.InsertCompany(&comp)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmCreateFailed)
		return
	}
	model.AddUserToCompany(comp.No, u.Id)
	u.Companys = append(u.Companys, comp)
	//生成token失败的话也注册成功，客户端提示用户重新登录
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(u.Id), comp.No, "", c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, u.Tel)
	}
	beego.Info("login ok,token:%+v", token)
	c.ReplySucc(u)
}

func (c *Controller) GetUserInfo() {
	user, e := service.GetUserById("public", int(c.UserID))
	if e != nil {
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		beego.Error(e)
		return
	}
	c.ReplySucc(user)
	return
}

func (c *Controller) UserLogin() {
	tel := c.GetString("tel")
	passwd := c.GetString("password")
	user, err := service.GetUserByTel("public", tel)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}
	if !keycrypt.CheckSha256(passwd, user.Password) {
		c.ReplyErr(errcode.ErrUserPasswordError)
		return
	}
	c.loginAction(user)
}

func (c *Controller) UserLoginPhone() {
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
	user, err := service.GetUserByTel("public", tel)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrUserNotExisted)
		return
	}
	c.loginAction(user)
}

func (c *Controller) loginAction(user *model.User) {
	var comNo string
	if len(user.Companys) == 1 {
		comNo = user.Companys[0].No
	} else {
		comNo = ""
	}
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), comNo, "", c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, *user)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	} else {
		if len(comNo) != 0 {
			u, e := service.GetUserByTel(comNo, user.Tel)
			if e == nil {
				u.Companys = user.Companys
				c.ReplySucc(u)
				return
			}
		}
		user.LoginTime = time.Now()
		model.UpdateUser("public", user)
		c.ReplySucc(user)
		beego.Debug("login ok,token:%+v", token)
		return
	}
}

func (c *Controller) Test() {
	a, e := c.RedisClient.Hmset("group_1", map[string]interface{}{
		"roles":  "1_2_3",
		"groups": "4_5_6",
	})
	beego.Info(a, e)

	b, e := c.RedisClient.Hmget("group_1", []string{"roles", "groups"})
	beego.Info(b, e)

	d, e := c.RedisClient.HDel("group_1", "roles")
	beego.Info(d, e)

	f, e := c.RedisClient.Hmget("group_1", []string{"roles", "groups"})
	beego.Info(f, e, len(f["roles"]))

	c.ReplySucc(nil)
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
	uid := (int)(c.UserID)
	pwd := keycrypt.Sha256Cal(c.GetString("password"))
	owd := keycrypt.Sha256Cal(c.GetString("oldpassword"))
	prefix := "public"
	user, e := service.GetUserById(prefix, uid)
	if e != nil {
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		beego.Error(e)
		return
	}
	if len(user.Password) == 0 {
		user.Password = pwd
		e = model.UpdateUser(prefix, user)
		if e != nil {
			c.ReplyErr(e)
		} else {
			c.ReplySucc(nil)
		}
	} else {
		if owd != user.Password {
			e = errcode.ErrUserPasswordError
			c.ReplyErr(e)
		} else {
			user.Password = pwd
			e = model.UpdateUser(prefix, user)
			if e != nil {
				c.ReplyErr(e)
			} else {
				c.ReplySucc(nil)
			}
		}
	}
}

func uniqueNo(prefix string) string {
	str := strings.Replace(time.Now().Format("0102150405.000"), ".", "", 1)
	str = prefix + str
	return str
}

//注册之后增加公司资质信息
func (c *Controller) AddLicenseFile() {
	fileUrl := c.GetString("url")
	compNo := c.UserComp
	uid := int(c.UserID)
	comp := model.Company{
		No:          compNo,
		Creator:     uid,
		LicenseFile: fileUrl,
	}
	e := model.UpdateCompany(&comp)
	if e != nil {
		c.ReplyErr(errcode.ErrFirmUpdateFailed)
		beego.Error(e)
	}
	c.ReplySucc(nil)
}

func (c *Controller) AdminGetFirmInfo() {
	no := c.GetString("cno")
	f, err := model.GetCompany(no)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmNotExisted)
		return
	}
	c.ReplySucc(*f)
}

func (c *Controller) AdminGetFirmList() {
	list, err := model.GetCompanyList()
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(list)
}

func (c *Controller) AdminFirmAudit() {
	uid := int(c.UserID)
	cno := c.GetString("cno")
	status, err := c.GetInt("status")
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	msg := c.GetString("msg")
	err = service.AuditCompany(cno, uid, status, msg)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(nil)
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

	err := model.CreateCompanyUser(cno, tel)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc("ok")
}

//**********************************************************
//暂未开放的接口
func (c *Controller) EditProfile() {
	gender, _ := c.GetInt8("gender")
	username := c.GetString("username")
	descp := c.GetString("descp")
	address := c.GetString("address")
	id := int(c.UserID)
	user := model.User{
		Id:       id,
		Gender:   gender,
		Desc:     descp,
		UserName: username,
		Address:  address,
	}

	err := model.UpdateUser("public", &user)
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
	user, err := service.GetUserByTel("public", tel)
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
func (c *Controller) Retrievepwd() {
	tel := c.GetString("tel")
	user, err := service.GetUserByTel("public", tel)
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
		err = model.UpdateUser("public", user)
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
func (c *Controller) FirmRegister() {
	uid, _ := c.GetInt("uid")
	name := c.GetString("firm_name")
	desc := c.GetString("desc")
	phone := c.GetString("phone")
	lf := c.GetString("license_file")
	tp := c.GetString("firm_type")

	firm := model.Company{
		No:          uniqueNo("O"),
		Creator:     uid,
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
func (c *Controller) FirmModify() {
	no := c.GetString("no")
	name := c.GetString("name")
	desc := c.GetString("desc")
	phone := c.GetString("phone")
	lf := c.GetString("license_file")
	uid := int(c.UserID)
	firm := model.Company{
		No:          no,
		Creator:     uid,
		FirmName:    name,
		Desc:        desc,
		Phone:       phone,
		LicenseFile: lf,
	}
	err := model.UpdateCompany(&firm)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmUpdateFailed)
		return
	}
	c.ReplySucc("ok")
}

//登录之后切换当前公司
func (c *Controller) SwitchCurrentFirm() {
	cno := c.GetString("cno")
	uid := int(c.UserID)
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(uid), cno, "", c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	} else {
		c.ReplySucc(nil)
		beego.Debug("switch company success,token:%+v", token)
		return
	}
}
