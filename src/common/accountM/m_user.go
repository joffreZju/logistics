package accountM

import (
	"common/lib/keycrypt"
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

type User struct {
	Id        int       `gorm:"AUTO_INCREMENT;primary_key"` // 用户ID，表内自增
	Tel       string    `gorm:"unique_index;size:15;not null" json:",omitempty"`
	Password  string    `json:"-"`                         // 密码
	UserName  string    `gorm:"size:64" json:",omitempty"` // 用户名
	Icon      string    `gorm:"size:64" json:",omitempty"`
	Desc      string    `json:",omitempty"`
	Gender    int8      `gorm:"default:1" json:",omitempty"`
	Address   string    `gorm:"size:64" json:",omitempty"`
	LoginTime time.Time `gorm:"timestamp" json:",omitempty"`                 //登录时间
	Ctime     time.Time `gorm:"default:current_timestamp" json:",omitempty"` //
	Mail      string    `gorm:"size:64" json:",omitempty"`
	UserType  int       `gorm:"default:1" json:",omitempty"` //1 普通用户
	Status    int       `gorm:"not null" json:",omitempty"`
	Companys  []Company `gorm:"-" json:",omitempty"` // 用户的所在组织
}

func (User) TableName() string {
	return "allsum_user"
}

func CreateUser(u *User) (err error) {
	err = NewOrm().Create(u).Error
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

func GetUserById(id int) (u *User, err error) {
	u = &User{}
	err = NewOrm().First(u, id).Error
	if err != nil {
		return
	}
	sql := `select * from allsum_company as t1 inner join allsum_user_company as t2
		on t1.no = t2.cno
		where t2.user_id = ?`
	var cmp []Company
	err = NewOrm().Raw(sql, u.Id).Scan(&cmp).Error
	u.Companys = cmp
	return
}

func GetUserByTel(tel string) (u *User, err error) {
	u = new(User)
	err = NewOrm().Find(u, User{Tel: tel}).Error
	if err != nil {
		return
	}
	sql := `select * from allsum_company as t1 inner join allsum_user_company as t2
		on t1.no = t2.cno
		where t2.user_id = ?`
	var cmp []Company
	err = NewOrm().Raw(sql, u.Id).Scan(&cmp).Error
	u.Companys = cmp
	return
}

const (
	CompApproveWait = iota
	CompApproveAccessed
	CompApproveNotAccessed
	CompanyDeleted
)

type Company struct {
	Id          int    `gorm:"auto_increment;not null"`
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
	err = NewOrm().First(c, cno).Error
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

func InsertCompany(c *Company) (err error) {
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

func AuditCompany(cno string, approverId int, status int, msg string) (err error) {
	c := NewOrm().Model(&Company{}).Where("no=?", cno).
		Updates(&Company{Status: status, Approver: approverId, ApproveTime: time.Now(), ApproveMsg: msg}).RowsAffected
	if c != 1 {
		err = errors.New("approve compony failed")
	}
	return
}

type UserCompany struct {
	Id     int    `gorm:"auto_increment"`
	UserId int    `gorm:"not null"`
	Cno    string `gorm:"not null"`
}

func (UserCompany) TableName() string {
	return "allsum_user_company"
}

func DelCompanyUser(cno string, uid int) (err error) {
	uc := UserCompany{
		UserId: uid,
		Cno:    cno,
	}
	c := NewOrm().Delete(&uc, UserCompany{UserId: uid, Cno: cno}).RowsAffected
	if c != 1 {
		err = errors.New("delete user failed")
	}
	return
}

func CreateCompanyUser(cno string, utel string) (err error) {
	user := User{
		Tel:      utel,
		Password: keycrypt.Sha256Cal("123456"),
		UserType: UserTypeNormal,
		Status:   UserStatusOk,
	}
	err = NewOrm().FirstOrCreate(user, User{Tel: utel}).Error
	if err != nil {
		return
	}
	err = AddUserToCompany(cno, user.Id)
	return
}

func AddUserToCompany(cno string, uid int) (err error) {
	uc := UserCompany{
		UserId: uid,
		Cno:    cno,
	}
	err = NewOrm().Create(&uc).Error
	return
}
