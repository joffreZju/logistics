package db

import (
	"allsum_bi/util"
	"strings"

	"github.com/astaxie/beego"
)

func AddAuthority(dbid string, username string, schema string, option string) (err error) {
	username = strings.ToLower(username)
	sqlstr1 := "GRANT " + option + " ON ALL TABLES IN SCHEMA " + schema + " TO " + username
	err = Exec(dbid, sqlstr1)
	if err != nil {
		return
	}
	sqlstr2 := "ALTER DEFAULT PRIVILEGES IN SCHEMA " + schema + " GRANT " + option + " ON TABLES TO " + username
	err = Exec(dbid, sqlstr2)
	return
}

func RevokeAuthority(dbid string, username string, schema string, option string) (err error) {
	username = strings.ToLower(username)
	sqlstr := "revoke " + option + " on ALL tables in schema " + schema + " FROM " + username
	err = Exec(dbid, sqlstr)
	return
}

func CreateUser(dbid string, username string) (err error) {
	username = strings.ToLower(username)
	if IsUserExist(dbid, username) {
		return
	}
	passwd := util.DEFAULT_PASSWD
	sqlstr := "CREATE USER " + username + " WITH PASSWORD '" + passwd + "'"
	//	params := []interface{}{passwd}
	return Exec(dbid, sqlstr)
}

func IsUserExist(dbid, username string) (exist bool) {
	username = strings.ToLower(username)
	sqlstr := "select usename from pg_user "
	datas, err := QueryDatas(dbid, sqlstr)
	if datas == nil || err != nil || len(datas) == 0 {
		return false
	}
	beego.Debug("data:", username)
	for _, data := range datas {
		datastr := string(data[0].([]byte))
		beego.Debug("data:", datastr, datastr == username, datastr, username)
		if strings.EqualFold(datastr, username) {
			return true
		}
	}
	return false
}

func ALterPassWD(dbid string, username string, passwd string) (err error) {
	username = strings.ToLower(username)
	sqlstr := "alter user " + username + " with password '" + passwd + "'"
	return Exec(dbid, sqlstr)
}
