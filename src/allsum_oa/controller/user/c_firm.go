package user

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"common/lib/keycrypt"
	"encoding/json"
	"github.com/astaxie/beego"
)

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

func (c *Controller) FirmAddUser() {
	cno := c.GetString("cno")
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
	if e != nil {
		c.ReplyErr(errcode.New(commonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(nil)
}

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
		e = service.UpdateGroupssOfUser(prefix, gids, user.Id)
		if e != nil {
			c.ReplyErr(errcode.ErrServerError)
			beego.Error(e)
			return
		}
	}
	c.ReplySucc(nil)
}
