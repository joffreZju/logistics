package models

import (
	"allsum_bi/db/conn"
	"fmt"
)

type DatabaseManager struct {
	Dbid     string
	Dbname   string
	Dbtype   string
	Host     string
	Port     int
	Dbuser   string
	Password string
	Params   string
	Name     string
}

func InsertDatabaseManager(conninfo DatabaseManager) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	exist := db.NewRecord(conninfo)
	if !exist {
		return fmt.Errorf("exist")
	}
	err = db.Table(GetDatabaseManagerTableName()).Create(&conninfo).Error
	return
}

func GetDatabaseManager(connid string) (conninfo DatabaseManager, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDatabaseManagerTableName()).Where("dbid=?", connid).First(&conninfo).Error
	return
}

func ListDatabaseManager() (conninfos []DatabaseManager, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetDatabaseManagerTableName()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var conninfo DatabaseManager
		err = db.ScanRows(rows, &conninfo)
		if err != nil {
			return conninfos, err
		}
		conninfos = append(conninfos, conninfo)
	}
	return
}

func UpdateDatabaseManager(conninfo DatabaseManager, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDatabaseManagerTableName()).Where("dbid=?", conninfo.Dbid).Updates(conninfo).Update(fields).Error
	return
}
