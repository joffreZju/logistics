package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TimeFormatWithLocal = "2006-01-02T15:04:05+08:00"
	TimeFormat          = "2006-01-02 15:04:05"
	DateFormat          = "2006-01-02"
)

//用户类型
const (
	UserTypeNormal = iota + 1
)

//用户状态
const (
	UserStatusOk = iota + 1
	UserStatusLocked
)

//公司状态
const (
	CompanyStatApproveWait = iota + 1
	CompanyStatApproveAccessed
	CompanyStatApproveNotAccessed
	CompanyStatDeleted
)

type User struct {
	Id        int    `gorm:"primary_key" ` // 用户id,继承自public
	No        string `gorm:"unique;size:64"`
	Tel       string `gorm:"unique;size:15;not null"`
	Password  string `json:"-"` // 密码
	UserName  string `gorm:""`
	Icon      string `gorm:""`
	Descrp    string
	Gender    int       // 1:男 2:女
	Address   string    `gorm:"size:64"`
	Ctime     time.Time `gorm:"default:current_timestamp"`
	LoginTime time.Time `gorm:"timestamp"`
	Mail      string    `gorm:"size:64"`
	Status    int       `gorm:"not null"`
	UserType  int       `gorm:"default:1"` //1 普通用户
	Companys  []Company `gorm:"-"`
	Groups    []Group   `gorm:"-"`
	Roles     []Role    `gorm:"-"`
	Funcs     []int     `gorm:"-"`
}

func (User) TableName() string {
	return "allsum_user"
}

func CreateUser(prefix string, u *User) (err error) {
	err = NewOrm().Table(prefix + "." + u.TableName()).Create(u).Error
	return
}

func FirstOrCreateUser(prefix string, u *User) (err error) {
	err = NewOrm().Table(prefix+"."+u.TableName()).
		FirstOrCreate(u, User{Tel: u.Tel}).Error
	return
}

func UpdateUser(prefix string, u *User) (e error) {
	c := NewOrm().Table(prefix+"."+u.TableName()).Where("id=? or tel=?", u.Id, u.Tel).
		Updates(u).RowsAffected
	if c != 1 {
		e = errors.New("update user failed")
		return
	}
	return nil
}

func Update___User(prefix string, u *User, fields ...string) (err error) {
	if len(prefix) == 0 {
		prefix = "public"
	}
	if len(fields) == 0 {
		fields = append(fields, "Id", "Icon",
			"Gender", "Descp", "Address", "LoginTime",
			"Tel", "UserName", "Password", "Mail")
	}
	sql := fmt.Sprintf(`update "%s".user set PARAMS where id = ?`, prefix)

	params, values := "", []interface{}{}
	for _, f := range fields {
		switch f {
		case "UserName":
			params += " user_name= ? ,"
			values = append(values, u.UserName)
		case "Descp":
			params += " desc= ? ,"
			values = append(values, u.Descrp)
		case "Gender":
			params += " gender= ? ,"
			values = append(values, u.Gender)
		case "Address":
			params += " addr= ? ,"
			values = append(values, u.Address)
		case "LoginTime":
			params += " login_time= ? ,"
			values = append(values, time.Now().Format(TimeFormat))
		case "Mail":
			params += " mail= ? ,"
			values = append(values, u.Mail)
		case "Password":
			params += " password= ? ,"
			values = append(values, u.Password)
		}
	}
	values = append(values, u.Id)
	if len(params) > 1 {
		params = params[:len(params)-1]
	}
	sql = strings.Replace(sql, "PARAMS", params, 1)
	c := NewOrm().Exec(sql, values...).RowsAffected
	if c != 1 {
		err = errors.New("update failed")
		return
	}
	return nil
}

type Company struct {
	Id           int    `gorm:"auto_increment;primary_key"`
	No           string `gorm:"unique"`
	Creator      int    `gorm:"not null"`
	AdminId      int    `gorm:"not null"`
	FirmName     string
	FirmType     string
	Descrp       string
	Address      string
	Phone        string
	LicenseFile  StrSlice  `gorm:"type:text[]"`
	Status       int       //1:待审核;2:审核通过;3:审核不通过4:删除;
	Approver     int       //审核人
	ApproverName string    //审核人
	ApproveTime  time.Time //批复时间
	ApproveMsg   string    //审批意见
	Ctime        time.Time `gorm:"default:current_timestamp"` //申请时间
}

func (Company) TableName() string {
	return "allsum_company"
}

func GetCompany(cno string) (c *Company, err error) {
	c = new(Company)
	err = NewOrm().Where("no=?", cno).First(c).Error
	return
}

func DeleteCompany(cno string) (err error) {
	err = NewOrm().Table(Company{}.TableName()).Where("no=?", cno).Update("status", CompanyStatDeleted).Error
	return
}

func CreateCompany(c *Company) (err error) {
	err = NewOrm().Create(c).Error
	return
}

func UpdateCompany(c *Company) (err error) {
	cond := &Company{
		No:      c.No,
		AdminId: c.AdminId,
	}
	count := NewOrm().Table(c.TableName()).Where(cond).Updates(&c).RowsAffected
	if count != 1 {
		err = errors.New("update company info failed")
	}
	return
}

type UserCompany struct {
	Id     int    `gorm:"auto_increment,primary_key"`
	UserId int    `gorm:"not null"`
	Cno    string `gorm:"not null"`
}

func (UserCompany) TableName() string {
	return "allsum_user_company"
}

func AddUserToCompany(cno string, uid int) (err error) {
	uc := UserCompany{
		UserId: uid,
		Cno:    cno,
	}
	err = NewOrm().Create(&uc).Error
	return
}

type Function struct {
	Id       int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name     string `gorm:"not null"`
	Descrp   string
	Pid      int       `gorm:"not null"`
	Ctime    time.Time `gorm:"default:current_timestamp"`
	Path     string    `gorm:""`
	Icon     string    //菜单图标
	SysId    string    //各系统id
	Services StrSlice  `gorm:"type:text[]"` //此功能可访问的API接口集合
}

func (Function) TableName() string {
	return "function"
}

func GetFunctions(sysIds []string) (funcs []*Function, e error) {
	funcs = []*Function{}
	e = NewOrm().Where("sys_id in (?)", sysIds).Find(&funcs).Error
	return
}

type AppVersion struct {
	Id          int `gorm:"primary_key;AUTO_INCREMENT"`
	Version     string
	Environment int      //1:开发2:测试3:预发布4:生产
	DownloadUrl StrSlice `gorm:"type:text[]"` //多个下载地址
	UpgradeType int      //1:透明2:友好提示3:强制升级
	Descrp      string
	Ctime       time.Time `gorm:"default:current_timestamp"`
}

func (AppVersion) TableName() string {
	return "app_version"
}
