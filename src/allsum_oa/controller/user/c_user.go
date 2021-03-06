package user

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"fmt"
	"strconv"
	"time"

	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
	"strings"
)

const commonErr = 99999

type Controller struct {
	base.Controller
}

//用户注册
//1 创建用户
//2 创建公司
//3 将用户和公司关联
//4 生成token,登陆成功
func (c *Controller) UserRegister() {
	tel := c.GetString("tel")
	username := c.GetString("username")
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
		UserName: username,
		Gender:   1,
		Password: pwdEncode,
		UserType: model.UserTypeNormal,
		Status:   model.UserStatusOk,
		Ctime:    time.Now(),
	}
	err := model.CreateUser("public", &u)
	if err != nil {
		beego.Error(err)
		if strings.Contains(err.Error(), model.DBErrStrDuplicateKey) {
			c.ReplyErr(errcode.ErrUserAlreadyExisted)
		} else {
			c.ReplyErr(errcode.ErrUserCreateFailed)
		}
		return
	}
	comp := model.Company{
		No:       model.UniqueNo("C"),
		FirmName: firm_name,
		FirmType: firm_type,
		Creator:  u.Id,
		AdminId:  u.Id,
		Status:   model.CompanyStatApproveWait,
	}
	err = model.CreateCompany(&comp)
	if err != nil {
		beego.Error(err)
		c.ReplyErr(errcode.ErrFirmCreateFailed)
		return
	}
	err = model.AddUserToCompany(comp.No, u.Id)
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
	key := fmt.Sprintf("%d-%s", u.Id, token.Value)
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

func (c *Controller) UpdateUserInfo() {
	prefix := c.UserComp
	uid := c.UserID
	uname := c.GetString("username")
	icon := c.GetString("icon")
	mail := c.GetString("mail")
	address := c.GetString("address")
	gender, e := c.GetInt("gender")
	if e != nil {
		gender = 1
	}
	user, e := service.GetUserById("public", uid)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		return
	}
	if len(uname) != 0 {
		user.UserName = uname
	}
	if len(icon) != 0 {
		user.Icon = icon
	}
	if gender != 0 {
		user.Gender = gender
	}
	if len(mail) != 0 {
		user.Mail = mail
	}
	if len(address) != 0 {
		user.Address = address
	}
	e = model.UpdateUser("public", user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	company, e := model.GetCompany(prefix)
	if company.Status == model.CompanyStatApproveAccessed {
		e = model.UpdateUser(prefix, user)
		if e != nil {
			beego.Error(e)
			c.ReplyErr(errcode.ErrServerError)
			return
		}
	}
	c.ReplySucc(nil)
}

//用户登录
//1 获取public下的用户
//2 检测用户是否存在有效的公司(schema),
// 如果没有,那么登录成功(没有功能集)
// 如果有,那么选择一家有效公司登录成功
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
		uKey := fmt.Sprintf("%d-%s", userInSchema.Id, token.Value)
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

//登录时将用户信息拼接存储在redis中,拼接使用 "-"
func (c *Controller) saveUserInfoToRedis(key, cno string, u *model.User) (e error) {
	roles, groups, funcs := "-", "-", "-"
	for _, v := range u.Roles {
		roles += fmt.Sprintf("%d-", v.Id)
	}
	for _, v := range u.Groups {
		groups += fmt.Sprintf("%d-", v.Id)
	}
	for _, v := range u.Funcs {
		for _, url := range v.Services {
			funcs += fmt.Sprintf("%s-", url)
		}
	}
	_, e = c.RedisClient.Hmset(key, map[string]interface{}{
		"company":   cno,
		"roles":     roles,
		"groups":    groups,
		"functions": funcs,
	})
	if e != nil {
		return
	}
	_, e = c.RedisClient.Expire(key, int64(tokenauth.TokenPeriod))
	return
}

//登录之后切换当前公司,用新公司的用户信息update redis中的信息
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
	key := fmt.Sprintf("%d-%s", user.Id, tokenStr)
	e = c.saveUserInfoToRedis(key, cno, user)

	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(user)
		beego.Info("switch company success with cno:%s", cno)
	}
}

//退出登录,删除token已经redis缓存
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

	key := fmt.Sprintf("%s-%s", token.SingleID, token.Value)
	_, err = c.RedisClient.Del(key)
	if err != nil {
		beego.Error(err)
	}
	c.ReplySucc(nil)
}

//忘记密码,用验证码修改
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

//重置密码
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
