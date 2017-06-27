package model

import (
	"common/lib/keycrypt"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/postgres"
	"strings"
	"time"
)

const (
	ReadOnly = 1
	Public   = "public"
)

type DBPool struct {
	db  *gorm.DB
	rdb *gorm.DB
}

var (
	hasReadOnly = false
	ormer       DBPool
)

func NewOrm(readOnly ...int) *gorm.DB {
	if hasReadOnly && len(readOnly) > 0 && readOnly[0] == ReadOnly {
		return ormer.rdb
	}
	return ormer.db
}

func InitPgSQL(key string) (err error) {
	username := beego.AppConfig.String("pgsql::username")
	password := beego.AppConfig.String("pgsql::password")
	addr := beego.AppConfig.String("pgsql::addr")
	port := beego.AppConfig.String("pgsql::port")
	addr_ro := beego.AppConfig.String("pgsql::addr_ro")
	dbname := beego.AppConfig.String("pgsql::dbname")

	if len(key) > 0 {
		password, err = keycrypt.Decode(key, password)
		if err != nil {
			return
		}
	}
	beego.Debug(username, password, addr, dbname)
	ormer.db, err = gorm.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			username, password, addr, port, dbname))
	if err != nil {
		return
	}
	if len(addr_ro) > 0 {
		ormer.rdb, err = gorm.Open("postgres",
			fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
				username, password, addr_ro, port, dbname))
		if err != nil {
			return
		}
		hasReadOnly = true
	}
	ormer.db.DB().SetConnMaxLifetime(time.Minute * 5)
	ormer.db.DB().SetMaxIdleConns(5)
	ormer.db.DB().SetMaxOpenConns(25)

	err = initModel()
	if err != nil {
		return
	}

	if beego.BConfig.RunMode == "prod" {
		ormer.db.LogMode(false)
	} else {
		ormer.db.LogMode(true)
	}
	//ormer.db.SetLogger(beego.BeeLogger)
	return
}

func initModel() (err error) {
	//init public model
	db := ormer.db
	if !db.HasTable(User{}.TableName()) {
		prefix := Public + "."
		db.Table(prefix + User{}.TableName()).AutoMigrate(new(User))
		db.Table(prefix + Company{}.TableName()).AutoMigrate(new(Company))
		db.Table(prefix + UserCompany{}.TableName()).AutoMigrate(new(UserCompany))
		db.Table(prefix + Function{}.TableName()).AutoMigrate(new(Function))
		db.Table(prefix + Message{}.TableName()).AutoMigrate(new(Message))
	}
	//init schema model
	comps := []Company{}
	err = db.Find(&comps, Company{Status: CompanyApproveAccessed}).Error
	if err != nil {
		return
	}
	for _, v := range comps {
		err = InitSchemaModel(v.No)
		if err != nil {
			return
		}
	}
	return
}

func InitSchemaModel(prefix string) (e error) {
	db := ormer.db
	prefix += "."
	e = db.Table(prefix + Group{}.TableName()).AutoMigrate(new(Group)).Error
	if e != nil && (strings.Contains(e.Error(), "already exists") ||
		strings.Contains(e.Error(), "已经存在")) {
		beego.Info(e)
		return nil
	}
	db.Table(prefix + User{}.TableName()).AutoMigrate(new(User))
	db.Table(prefix + UserGroup{}.TableName()).AutoMigrate(new(UserGroup))
	db.Table(prefix + Attribute{}.TableName()).AutoMigrate(new(Attribute))
	db.Table(prefix + GroupOperation{}.TableName()).AutoMigrate(new(GroupOperation))
	db.Table(prefix + Role{}.TableName()).AutoMigrate(new(Role))
	db.Table(prefix + RoleFunc{}.TableName()).AutoMigrate(new(RoleFunc))
	db.Table(prefix + UserRole{}.TableName()).AutoMigrate(new(UserRole))
	db.Table(prefix + Formtpl{}.TableName()).AutoMigrate(new(Formtpl))
	db.Table(prefix + Form{}.TableName()).AutoMigrate(new(Form))
	db.Table(prefix + Approvaltpl{}.TableName()).AutoMigrate(new(Approvaltpl))
	db.Table(prefix + Approval{}.TableName()).AutoMigrate(new(Approval))
	db.Table(prefix + ApproveFlow{}.TableName()).AutoMigrate(new(ApproveFlow))
	return nil
}

func UniqueNo(prefix string) string {
	str := strings.Replace(time.Now().Format("0102150405.000"), ".", "", 1)
	str = prefix + str
	return str
}
