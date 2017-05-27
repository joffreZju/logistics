package model

import (
	"common/lib/keycrypt"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
)

const (
	ReadOnly = 1
	Public   = "public."
)

type DBPool struct {
	db  *gorm.DB
	rdb *gorm.DB
}

//配置或数据表
var schemas = []string{"group."}

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
	for _, v := range schemas {
		ormer.db.Table(v + Group{}.TableName()).AutoMigrate(new(Group))
		ormer.db.Table(v + UserGroup{}.TableName()).AutoMigrate(new(UserGroup))
		ormer.db.Table(v + Attribute{}.TableName()).AutoMigrate(new(Attribute))
		ormer.db.Table(v + Operation{}.TableName()).AutoMigrate(new(Operation))
		ormer.db.Table(v + HistoryGroup{}.TableName()).AutoMigrate(new(HistoryGroup))

		ormer.db.Table(v + Role{}.TableName()).AutoMigrate(new(Role))
		ormer.db.Table(v + RoleFunc{}.TableName()).AutoMigrate(new(RoleFunc))
		ormer.db.Table(v + Func{}.TableName()).AutoMigrate(new(Func))
		ormer.db.Table(v + UserRole{}.TableName()).AutoMigrate(new(UserRole))

		ormer.db.Table(v + Formtpl{}.TableName()).AutoMigrate(new(Formtpl))
		ormer.db.Table(v + Form{}.TableName()).AutoMigrate(new(Form))
		ormer.db.Table(v + Approvaltpl{}.TableName()).AutoMigrate(new(Approvaltpl))
		ormer.db.Table(v + Approval{}.TableName()).AutoMigrate(new(Approval))
		ormer.db.Table(v + ApproveFlow{}.TableName()).AutoMigrate(new(ApproveFlow))
	}

	if beego.BConfig.RunMode == "prod" {
		ormer.db.LogMode(false)
	} else {
		ormer.db.LogMode(true)
	}
	//Ormer.db.SetLogger(beego.BeeLogger)

	return
}
