package user

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"common/lib/push"
	"fmt"
	"math/rand"
	"strconv"
	"time"

	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
)

const commonErr = 99999

type Controller struct {
	base.Controller
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
		No:       model.UniqueNo("U"),
		Password: pwdEncode,
		UserType: model.UserTypeNormal,
		Status:   model.UserStatusOk,
		Ctime:    time.Now(),
	}
	err := model.CreateUser("public", &u)
	if err != nil {
		beego.Error(err)
		err = errcode.ErrUserCreateFailed
		c.ReplyErr(err)
		return
	}
	comp := model.Company{
		No:       model.UniqueNo("C"),
		FirmName: firm_name,
		FirmType: firm_type,
		Creator:  u.Id,
		Status:   model.CompanyStatApproveWait,
	}
	err = model.CreateCompany(&comp)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmCreateFailed)
		return
	}
	model.AddUserToCompany(comp.No, u.Id)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmCreateFailed)
		return
	}
	u.Companys = append(u.Companys, comp)
	//生成token失败的话也注册成功，客户端提示用户重新登录
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(u.Id), "", "", c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, u.Tel)
	}
	beego.Info("register and login ok,token:%+v", token)
	key := fmt.Sprintf("%d_%s", u.Id, token.Value)
	err = c.saveUserInfoToRedis(key, comp.No, &u)

	c.ReplySucc(u)
}

func (c *Controller) GetUserInfo() {
	user, e := service.GetUserById("public", c.UserID)
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
	user.LoginTime = time.Now()
	defer model.UpdateUser("public", user)
	token, err := o2o.Auth.NewSingleToken(strconv.Itoa(user.Id), "", "", c.Ctx.ResponseWriter)
	if err != nil {
		beego.Error("o2o.Auth.NewSingleToken error:", err, *user)
		c.ReplyErr(errcode.ErrAuthCreateFailed)
		return
	}
	for _, v := range user.Companys {
		if v.Status != model.CompanyStatApproveAccessed {
			continue
		}
		userInSchema, e := service.GetUserByTel(v.No, user.Tel)
		if e != nil || userInSchema.Status != model.UserStatusOk {
			continue
		}
		uKey := fmt.Sprintf("%d_%s", userInSchema.Id, token.Value)
		e = c.saveUserInfoToRedis(uKey, v.No, userInSchema)
		if e != nil {
			c.ReplyErr(errcode.New(commonErr, e.Error()))
			beego.Error(e)
			return
		} else {
			c.ReplySucc(userInSchema)
			beego.Info("login ok:%+v", token)
			return
		}
	}
	c.ReplySucc(user)
	beego.Info("login ok:%+v", token)
}

func (c *Controller) saveUserInfoToRedis(key, cno string, u *model.User) (e error) {
	roles, groups := "", ""
	for _, v := range u.Roles {
		roles += fmt.Sprintf("%d_", v.Id)
	}
	for _, v := range u.Groups {
		groups += fmt.Sprintf("%d_", v.Id)
	}
	_, e = c.RedisClient.Hmset(key, map[string]interface{}{
		"company": cno,
		"roles":   roles,
		"groups":  groups,
	})
	if e != nil {
		return
	}
	_, e = c.RedisClient.Expire(key, int64(tokenauth.TokenPeriod+10))
	return
}

//登录之后切换当前公司
func (c *Controller) SwitchCurrentFirm() {
	cno := c.GetString("cno")
	tokenStr := c.Ctx.Request.Header.Get("access_token")
	uid := c.UserID
	if len(cno) == 0 {
		c.ReplyErr(errcode.ErrParams)
		return
	}
	company, e := model.GetCompany(cno)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		return
	}
	user := &model.User{}
	if company.Status == model.CompanyStatApproveAccessed {
		user, e = service.GetUserById(cno, uid)
		if user.Status != model.UserStatusOk {
			c.ReplyErr(errcode.ErrUserLocked)
			return
		}
	} else {
		user, e = service.GetUserById("public", uid)
	}
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		return
	}
	var currentCompanyIndex int
	for k, v := range user.Companys {
		if v.No == cno {
			currentCompanyIndex = k
			break
		}
	}
	user.Companys = user.Companys[currentCompanyIndex : currentCompanyIndex+1]
	//将用户的company,groups和roles放入缓存
	key := fmt.Sprintf("%d_%s", user.Id, tokenStr)
	e = c.saveUserInfoToRedis(key, cno, user)

	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(user)
		beego.Info("switch company success with cno:%s", cno)
	}
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

	key := fmt.Sprintf("%s_%s", token.SingleID, token.Value)
	_, err = c.RedisClient.Expire(key, 2)
	if err != nil {
		beego.Error(err)
	}
	c.ReplySucc(nil)
}

func (c *Controller) Forgetpwd() {
	tel := c.GetString("tel")
	code := c.GetString("code")
	newPwd := keycrypt.Sha256Cal(c.GetString("password"))
	mycode := c.Cache.Get(tel)
	if mycode == nil {
		c.ReplyErr(errcode.ErrAuthCodeExpired)
		return
	} else if fmt.Sprintf("%s", mycode) != code {
		c.ReplyErr(errcode.ErrAuthCodeError)
		return
	}
	prefixPublic := model.Public
	user, e := service.GetUserByTel(prefixPublic, tel)
	if e != nil {
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		beego.Error(e)
		return
	}
	user.Password = newPwd
	e = model.UpdateUser(prefixPublic, user)
	if e != nil {
		c.ReplyErr(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) Resetpwd() {
	uid := c.UserID
	pwd := keycrypt.Sha256Cal(c.GetString("password"))
	owd := keycrypt.Sha256Cal(c.GetString("oldpassword"))
	prefix := model.Public
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

//注册之后增加公司资质信息
func (c *Controller) AddLicenseFile() {
	urlstr := c.GetString("url")
	urlslice := model.StrSlice{}
	e := json.Unmarshal([]byte(urlstr), &urlslice)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	compNo := c.UserComp
	uid := c.UserID
	comp := model.Company{
		No:          compNo,
		Creator:     uid,
		LicenseFile: urlslice,
	}
	e = model.UpdateCompany(&comp)
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
	uid := c.UserID
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

func (c *Controller) FirmGetUserList() {
	prefix := c.UserComp
	users, e := service.GetUserListOfCompany(prefix)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(users)
	}
}

func (c *Controller) FirmAddUser() {
	cno := c.GetString("cno")
	tel := c.GetString("tel")
	name := c.GetString("name")
	gender, e := c.GetInt("gender")
	if e != nil {
		gender = 0
	}
	groups, roles := []int{}, []int{}
	gstr, rstr := c.GetString("groups"), c.GetString("roles")
	e = json.Unmarshal([]byte(gstr), &groups)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	e = json.Unmarshal([]byte(rstr), &roles)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrParams)
		return
	}

	user := &model.User{
		Tel:      tel,
		No:       model.UniqueNo("U"),
		Password: keycrypt.Sha256Cal("123456"),
		UserName: name,
		Gender:   gender,
		UserType: model.UserTypeNormal,
		Status:   model.UserStatusOk,
	}
	e = model.FirstOrCreateUser("public", user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = model.FirstOrCreateUser(cno, user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = model.AddUserToCompany(cno, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = service.AddUserToGroups(cno, groups, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = service.AddUserToRoles(cno, roles, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(nil)
}

func (c *Controller) FirmControlUserStatus() {
	prefix := c.UserComp
	tel := c.GetString("tel")
	status, e := c.GetInt("status")
	if e != nil || (status != model.UserStatusOk && status != model.UserStatusLocked) {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = model.UpdateUser(prefix, &model.User{Tel: tel, Status: status})
	e = model.UpdateUser(prefix, &model.User{Tel: tel, Status: status})
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(nil)
}

func (c *Controller) UpdateUserProfile() {
	cno := c.GetString("cno")
	tel := c.GetString("tel")
	uname := c.GetString("username")
	mail := c.GetString("mail")
	rstr := c.GetString("roles")
	gstr := c.GetString("groups")
	beego.Info(rstr, gstr)
	user, e := service.GetUserByTel("public", tel)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	user.UserName = uname
	user.Mail = mail
	e = model.UpdateUser("public", user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = model.UpdateUser(cno, user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	if len(rstr) != 0 {
		rids := []int{}
		e = json.Unmarshal([]byte(rstr), &rids)
		if e != nil {
			c.ReplyErr(errcode.ErrParams)
			beego.Error(e)
			return
		}
		e = service.UpdateRolesOfUser(cno, rids, user.Id)
		if e != nil {
			c.ReplyErr(errcode.ErrServerError)
			beego.Error(e)
			return
		}
	}
	if len(gstr) != 0 {
		gids := []int{}
		e = json.Unmarshal([]byte(gstr), &gids)
		if e != nil {
			c.ReplyErr(errcode.ErrParams)
			beego.Error(e)
			return
		}
		e = service.UpdateGroupssOfUser(cno, gids, user.Id)
		if e != nil {
			c.ReplyErr(errcode.ErrServerError)
			beego.Error(e)
			return
		}
	}
	c.ReplySucc(nil)
}

func (c *Controller) GetFunctionsTree() {
	funcs, e := model.GetFunctions()
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(funcs)
}

//暂未开放的接口**********************************************************
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
	//lf := c.GetString("license_file")
	tp := c.GetString("firm_type")

	firm := model.Company{
		No:       model.UniqueNo("O"),
		Creator:  uid,
		FirmName: name,
		Desc:     desc,
		Phone:    phone,
		//LicenseFile: lf,
		FirmType: tp,
		Status:   0,
	}
	err := model.CreateCompany(&firm)
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
	//lf := c.GetString("license_file")
	uid := c.UserID
	firm := model.Company{
		No:       no,
		Creator:  uid,
		FirmName: name,
		Desc:     desc,
		Phone:    phone,
		//LicenseFile: lf,
	}
	err := model.UpdateCompany(&firm)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmUpdateFailed)
		return
	}
	c.ReplySucc("ok")
}
