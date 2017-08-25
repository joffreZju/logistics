package firm

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
)

const commonErr = 99999

type Controller struct {
	base.Controller
}

//注册之后修改公司信息
func (c *Controller) UpdateFirmInfo() {
	compNo := c.UserComp
	uid := c.UserID
	comp := &model.Company{
		No:      compNo,
		AdminId: uid,
	}
	urlstr := c.GetString("url")
	addrStr := c.GetString("address")
	name := c.GetString("name")
	phone := c.GetString("phone")
	if len(urlstr) != 0 {
		urllist := model.StrSlice{}
		e := json.Unmarshal([]byte(urlstr), &urllist)
		if e != nil {
			c.ReplyErr(errcode.ErrParams)
			beego.Error(e)
			return
		}
		if len(urllist) != 0 {
			comp.LicenseFile = urllist
		}
	}
	if len(addrStr) != 0 {
		comp.Address = addrStr
	}
	if len(phone) != 0 {
		comp.Phone = phone
	}
	if len(name) != 0 {
		comp.FirmName = name
	}
	e := model.UpdateCompany(comp)
	if e != nil {
		c.ReplyErr(errcode.ErrFirmUpdateFailed)
		beego.Error(e)
	}
	c.ReplySucc(nil)
}

//公司管理员获取用户列表
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

//用名字搜索用户
func (c *Controller) FirmSearchUsersByName() {
	prefix := c.UserComp
	uname := c.GetString("username")
	users, e := service.SearchUsersByName(prefix, uname)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.New(commonErr, e.Error()))
	} else {
		c.ReplySucc(users)
	}
}

//公司管理员添加用户
// 1 先在public下创建用户,并将用户和公司关联
// 2 然后在schema下创建用户
// 3 如果用组织和角色信息,那么同时将组织和角色进行关联
func (c *Controller) FirmAddUser() {
	//cno := c.GetString("cno")
	prefix := c.UserComp
	tel := c.GetString("tel")
	name := c.GetString("name")
	gender, e := c.GetInt("gender")
	if e != nil {
		gender = 1
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
	e = model.CreateUser(prefix, user)
	if e != nil {
		beego.Error(e)
		if strings.Contains(e.Error(), model.DBErrStrDuplicateKey) {
			c.ReplyErr(errcode.ErrUserAlreadyExisted)
		} else {
			c.ReplyErr(errcode.ErrServerError)
		}
		return
	}
	e = model.AddUserToCompany(prefix, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = service.AddUserToGroups(prefix, groups, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = service.AddUserToRoles(prefix, roles, user.Id)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(nil)
}

//锁定用户，并删除其角色和组织信息,同时删除该用户在redis存在的token。锁定之后不能解锁
func (c *Controller) FirmControlUserStatus() {
	prefix := c.UserComp
	currentUid := c.UserID
	tel := c.GetString("tel")
	status, e := c.GetInt("status")
	if e != nil || status != model.UserStatusLocked {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	user, e := service.GetUserByTel(model.Public, tel)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	if currentUid == user.Id {
		c.ReplyErr(errcode.ErrLockUserFailed)
		return
	}
	e = service.LockUser(prefix, user)
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(nil)
	infoKeys, e := c.RedisClient.Keys(fmt.Sprintf("%d-*", user.Id))
	if e != nil {
		beego.Error(e)
	}
	for _, v := range infoKeys {
		index := strings.Index(v, "-")
		if index == -1 {
			beego.Error("拼接符错误，未能删除用户token")
		}
		c.RedisClient.Del(v)
		c.RedisClient.Del(v[index+1:])
		beego.Info("锁定用户成功,id:", user.Id)
	}
}

//公司管理员更新用户基本信息
func (c *Controller) FirmUpdateUserProfile() {
	prefix := c.UserComp
	tel := c.GetString("tel")
	uname := c.GetString("username")
	mail := c.GetString("mail")
	gender, e := c.GetInt("gender")
	if e != nil {
		gender = 1
	}
	user, e := service.GetUserByTel("public", tel)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		return
	}
	user.UserName = uname
	user.Mail = mail
	user.Gender = gender
	e = model.UpdateUser("public", user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	e = model.UpdateUser(prefix, user)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	c.ReplySucc(nil)
}

//公司管理员更新用户组织和角色信息
func (c *Controller) FirmUpdateUserRoleAndGroup() {
	prefix := c.UserComp
	tel := c.GetString("tel")
	rstr := c.GetString("roles")
	gstr := c.GetString("groups")
	user, e := service.GetUserByTel("public", tel)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrParams)
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
		e = service.UpdateRolesOfUser(prefix, rids, user.Id)
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
		e = service.UpdateGroupsOfUser(prefix, gids, user.Id)
		if e != nil {
			c.ReplyErr(errcode.ErrServerError)
			beego.Error(e)
			return
		}
	}
	c.ReplySucc(nil)
}
