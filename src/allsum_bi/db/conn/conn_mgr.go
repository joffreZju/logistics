package conn

import (
	"allsum_bi/util"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	_ "github.com/lib/pq"
)

var conn map[string]Conn

func init() {
	conn = map[string]Conn{}
}

func CreateConn(conninfo Conn) (err error) {
	_, ok := conn[conninfo.Id]
	if ok {
		return fmt.Errorf("this db is conned")
	}
	if conninfo.Dbtype != util.PG_DB_TYPE {
		return fmt.Errorf("CreatConn: db type not import !")
	}
	connurl := get_pg_url(&conninfo)
	db, err := gorm.Open("postgres", connurl)
	if err != nil {
		errinfo := fmt.Sprintf("CreatConn: %v", err)
		beego.Error("%s", errinfo)
		return err
	}
	db.DB().SetMaxIdleConns(5)
	db.DB().SetMaxOpenConns(30)
	isdebug := beego.AppConfig.String("log::debug")
	if isdebug == "true" {
		db.LogMode(true)
	}
	conninfo.Db = db
	conninfo.Status = true
	conninfo.Lastusetime = time.Now()
	conn[conninfo.Id] = conninfo
	return
}

func RemoveConn(connid string) {
	conninfo, ok := conn[connid]
	if ok {
		conninfo.Db.Close()
		delete(conn, connid)
	}
}
func GetConninfo(connid string) (conninfo Conn, err error) {
	conninfo, ok := conn[connid]
	if ok {
		return
	} else {
		err = fmt.Errorf("ERROR ConnID: %s", connid)
		return
	}
}

func GetBIConn() (db *gorm.DB, err error) {
	return GetConn(util.BASEDB_CONNID)
}
func GetConn(connid string) (db *gorm.DB, err error) {
	conninfo, ok := conn[connid]
	if ok {
		db = conninfo.Db
		return
	} else {
		err = fmt.Errorf("ERROR ConnID: %s", connid)
		return
	}
}

func GetAllConn() (conninfos map[string]Conn) {
	return conn
}

func get_pg_url(conninfo *Conn) (connstr string) {
	url := fmt.Sprintf("host=%s user=%s password=%s port=%d dbname=%s sslmode=disable", conninfo.Host, conninfo.DbUser, conninfo.Passwd, conninfo.Port, conninfo.Dbname)
	beego.Debug("pgurl: ", url)
	return url
}
