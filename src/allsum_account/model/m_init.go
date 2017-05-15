package model

import (
	"common/lib/keycrypt"
	"fmt"

	"github.com/astaxie/beego"
	orm "github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const ReadOnly = 1

type (
	DBPool struct {
		db  *orm.DB
		rdb *orm.DB
	}
)

var (
	hasReadOnly    = false
	readOnlyDBName = "alaccountro"
	Ormer          DBPool
)

func NewOrm(readOnly ...int) *orm.DB {
	if hasReadOnly && len(readOnly) > 0 && readOnly[0] == ReadOnly {
		return Ormer.rdb
	}
	return Ormer.db
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
	Ormer.db, err = orm.Open("postgres",
		fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
			username, password, addr, port, dbname))
	if err != nil {
		return
	}
	if len(addr_ro) > 0 {
		Ormer.rdb, err = orm.Open("postgres",
			fmt.Sprintf("user=%s password=%s host=%s port=%s dbname=%s sslmode=disable",
				username, password, addr_ro, port, dbname))
		if err != nil {
			return
		}
		hasReadOnly = true
	}
	Ormer.db.AutoMigrate(new(User), new(File), new(Document))

	if beego.BConfig.RunMode == "prod" {
		Ormer.db.LogMode(false)
	} else {
		Ormer.db.LogMode(true)
		fmt.Printf("Set orm debug open--------------------\n")
	}
	//Ormer.db.SetLogger(beego.BeeLogger)

	return
}
