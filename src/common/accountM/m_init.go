package accountM

import (
	"common/lib/keycrypt"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

const ReadOnly = 1

type (
	DBPool struct {
		db  *gorm.DB
		rdb *gorm.DB
	}
)

var (
	hasReadOnly    = false
	readOnlyDBName = "alaccountro"
	ormer          DBPool
)

func NewOrm(readOnly ...int) *gorm.DB {
	if hasReadOnly && len(readOnly) > 0 && readOnly[0] == ReadOnly {
		return ormer.rdb
	}
	return ormer.db
}

func ModelInit(db *gorm.DB) (err error) {
	db.AutoMigrate(new(User), new(Company), new(UserCompany))
	ormer.db = db
	return nil
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
	ormer.db.AutoMigrate(new(User), new(Company), new(UserCompany))

	if beego.BConfig.RunMode == "prod" {
		ormer.db.LogMode(false)
	} else {
		ormer.db.LogMode(true)
		fmt.Printf("Set orm debug open--------------------\n")
	}
	//Ormer.db.SetLogger(beego.BeeLogger)

	return
}
