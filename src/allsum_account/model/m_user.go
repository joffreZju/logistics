package model

import (
	//"errors"
	"fmt"
	"strings"
	"time"
)

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
)

type User struct {
	Id         int       `gorm:"AUTO_INCREMENT;primary_key" ` // 用户ID，表内自增
	Tel        string    `gorm:"unique_index;size:15;not null" json:",omitempty"`
	Password   string    `json:"-"`                         // 密码
	UserName   string    `gorm:"size:64" json:",omitempty"` // 用户名
	Icon       string    `gorm:"size:64" json:",omitempty"`
	Descp      string    `gorm:"" json:",omitempty"`
	Gender     int8      `gorm:"default:1" json:",omitempty"`
	Address    string    `gorm:"size:64" json:",omitempty"`
	LoginTime  time.Time `gorm:"timestamp" json:",omitempty"`                 //登录时间
	CreateTime time.Time `gorm:"default:current_timestamp" json:",omitempty"` //
	Mail       string    `gorm:"size:64" json:",omitempty"`
	UserType   int       `gorm:"default:1" json:",omitempty"` //1 普通用户
	Companys   []Company `orm:"-" json:",omitempty"`          // 用户的所在组织
}

func (User) TableName() string {
	return "allsum_user"
}

func CreateUser(u *User) (err error) {
	err = Ormer.db.Create(u).Error
	return
}

func CreateUserIfNotExist(u *User) (err error) {
	err = NewOrm().Where("Tel= ?", u.Tel).FirstOrCreate(u).Error
	return
}

func UpdateUser(u *User, fields ...string) (err error) {
	if len(fields) == 0 {
		fields = append(fields, "Id", "Icon",
			"Gender", "Descp", "Address", "LoginTime",
			"Tel", "UserName", "Password", "Mail")
	}
	sql := fmt.Sprintf("update allsum_user set PARAMS where id = ?")

	params, values := "", []interface{}{}
	for _, f := range fields {
		switch f {
		case "UserName":
			params += " user_name= ? ,"
			values = append(values, u.UserName)
		case "Descp":
			params += " descp= ? ,"
			values = append(values, u.Descp)
		case "Gender":
			params += " gender= ? ,"
			values = append(values, u.Gender)
		case "Address":
			params += " address= ? ,"
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
	err = NewOrm().Exec(sql, values...).Error
	if err != nil {
		return
	}
	//_, err = db.RowsAffected
	return
}

func GetUser(id int) (u *User, err error) {
	u = &User{}
	err = NewOrm(ReadOnly).First(u, id).Error
	return
}

func GetUserByTel(tel string) (u *User, err error) {
	u = new(User)
	err = NewOrm(ReadOnly).Where("Tel = ?", tel).Find(u).Error
	if err != nil {
		return
	}
	sql := `select * from allsum_company as t1 inner join allsum_user_company as t2
		on t1.no = t2.cno
		where t2.uid = ?`
	var cmp []Company
	//err = Ormer.db.Table("allsum_user_company").Where("uid=?", u.Id).Find(&cmp).Error
	err = Ormer.db.Raw(sql, u.Id).Scan(&cmp).Error
	u.Companys = cmp
	return
}

type Company struct {
	Id          int    `gorm:"auto_increment;not null"`
	No          string `gorm:"unique"`
	Creater     int    `gorm:"not null"`
	FirmName    string
	Desc        string
	Phone       string
	LicenseFile string    `gorm:"size:255;not null"`
	Status      int       //0:待审核;1:审核通过;2:审核不通过3:删除;
	FirmType    int       //0:普通公司，1:个体户
	Approver    int       //审核人
	ApproveTime time.Time //批复时间
	Msg         string    //审批意见
	CreateTime  time.Time `orm:"type(datetime)" json:",omitempty"` //申请时间
}

func (Company) TableName() string {
	return "allsum_company"
}

func GetCompany(no string) (c *Company, err error) {
	c = new(Company)
	err = NewOrm().Table("allsum_company").Where("no=?", no).First(c).Error
	return
}

func GetCompanies() (list []Company, err error) {
	err = NewOrm().Table("allsum_company").Find(&list).Error
	return
}

func DeleteCompany(cno string) (err error) {
	err = NewOrm().Table("allsum_company").Where("no=?", cno).Update("status", 3).Error
	return
}

func InsertCompany(c *Company) (err error) {
	err = NewOrm().Create(c).Error
	return
}

func UpdateCompany(c *Company) (err error) {
	err = NewOrm().Table("allsum_company").Where("no=? and status <> 1", c.No).Update(map[string]interface{}{
		"name":         c.FirmName,
		"desc":         c.Desc,
		"phone":        c.Phone,
		"license_file": c.LicenseFile,
	}).Error
	return
}

func AuditCompany(cno string, uid int, st int, msg string) (err error) {
	err = NewOrm().Model(&Company{}).Where("no=?", cno).Updates(Company{Status: st, Approver: uid, ApproveTime: time.Now(), Msg: msg}).Error
	return
}

type UserCompany struct {
	Id  int    `gorm:"auto_increment"`
	Uid int    `gorm:"not null"`
	Cno string `gorm:"not null"`
}

func (UserCompany) TableName() string {
	return "allsum_user_company"
}

func DelCompanyUser(cno string, uid int) (err error) {
	uc := UserCompany{
		Uid: uid,
		Cno: cno,
	}
	err = NewOrm().Delete(&uc, "uid=? and cno=?", uid, cno).Error
	return
}

func AddCompanyUser(cno string, tel string) (err error) {
	user := User{
		Tel: tel,
	}
	err = CreateUserIfNotExist(&user)
	if err != nil {
		return
	}
	uc := UserCompany{
		Uid: user.Id,
		Cno: cno,
	}
	err = NewOrm().Create(&uc).Error
	return
}

func AddUserToCompany(cno string, uid int) (err error) {
	uc := UserCompany{
		Uid: uid,
		Cno: cno,
	}
	err = NewOrm().Create(&uc).Error
	return
}

func DeleteCompanyUser(cno string, uid int) (err error) {
	err = NewOrm().Where("uid = ? and cno= ?", uid, cno).Delete(&UserCompany{}).Error
	return
}
