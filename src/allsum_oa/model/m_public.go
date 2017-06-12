package model

import (
	"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

const (
	UserTypeNormal = iota + 1
)

const (
	UserStatusOk = iota
	UserStatusLocked
)

const (
	CompApproveWait = iota
	CompApproveAccessed
	CompApproveNotAccessed
	CompanyDeleted
)

type User struct {
	Id        int       `gorm:"primary_key" ` // 用户id,继承自public
	No        string    `gorm:"unique;size:64"`
	Tel       string    `gorm:"size:15;not null" json:",omitempty"`
	Password  string    `json:"-"` // 密码
	UserName  string    `gorm:"size:64" json:",omitempty"`
	Icon      string    `gorm:"size:64" json:",omitempty"`
	Desc      string    `gorm:"" json:",omitempty"`
	Gender    int8      `gorm:"" json:",omitempty"`
	Address   string    `gorm:"size:64" json:",omitempty"`
	Ctime     time.Time `gorm:"default:current_timestamp" json:",omitempty"`
	LoginTime time.Time `gorm:"timestamp" json:",omitempty"`
	Mail      string    `gorm:"size:64" json:",omitempty"`
	Status    int       `gorm:"not null" json:",omitempty"`
	UserType  int       `gorm:"default:1" json:",omitempty"` //1 普通用户
	Companys  []Company `gorm:"-"`
	Groups    []Group   `gorm:"-"`
	Roles     []Role    `gorm:"-"`
	Funcs     []int     `gorm:"-"`
}

func (User) TableName() string {
	return "user"
}

func CreateUser(prefix string, u *User) (err error) {
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
			values = append(values, u.Desc)
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
	Id          int    `gorm:"auto_increment;primary_key"`
	No          string `gorm:"unique"`
	Creator     int    `gorm:"not null"`
	FirmName    string
	FirmType    string
	Desc        string
	Phone       string
	LicenseFile string    `gorm:"not null"`
	Status      int       //0:待审核;1:审核通过;2:审核不通过3:删除;
	Approver    int       //审核人
	ApproveTime time.Time //批复时间
	ApproveMsg  string    //审批意见
	Ctime       time.Time `gorm:"default:current_timestamp" json:",omitempty"` //申请时间
}

func (Company) TableName() string {
	return "allsum_company"
}

func GetCompany(cno string) (c *Company, err error) {
	c = new(Company)
	err = NewOrm().Where("no=?", cno).First(c).Error
	return
}

func GetCompanyList() (list []Company, err error) {
	list = []Company{}
	err = NewOrm().Find(&list).Error
	return
}

func DeleteCompany(cno string) (err error) {
	err = NewOrm().Table(Company{}.TableName()).Where("no=?", cno).Update("status", CompanyDeleted).Error
	return
}

func CreateCompany(c *Company) (err error) {
	err = NewOrm().Create(c).Error
	return
}

func UpdateCompany(c *Company) (err error) {
	count := NewOrm().Table(c.TableName()).Where("no=? and creator = ? and status <> ?", c.No, c.Creator, CompApproveAccessed).
		Updates(&c).RowsAffected
	if count != 1 {
		err = errors.New("update license file failed")
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
	Id    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"not null"`
	Desc  string
	Pid   int       `gorm:"not null"`
	Ctime time.Time `gorm:"default:current_timestamp"`
	Path  string    `gorm:""`
}

func (Function) TableName() string {
	return "function"
}

func GetFunctions() (funcs []*Function, e error) {
	funcs = []*Function{}
	e = NewOrm().Find(&funcs).Error
	return
}
