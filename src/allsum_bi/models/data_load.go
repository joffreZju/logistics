package models

import (
	"allsum_bi/db/conn"
	"database/sql"
	"fmt"
	_ "time"

	"github.com/satori/go.uuid"
)

type DataLoad struct {
	Id           int
	Uuid         string
	Name         string
	Owner        string
	Columns      string
	CreateScript string
	AlterScript  string
	Basetable    string
	Documents    string
	WebPath      string
	WebfileName  string
	Aggregateid  int
	Status       string
}

func InsertDataLoad(dataload DataLoad) (uuidstr string, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	dataload.Uuid = uuid.NewV4().String()
	exist := db.NewRecord(dataload)
	if !exist {
		return uuidstr, fmt.Errorf("exist")
	}
	err = db.Table(GetDataLoadTableName()).Create(&dataload).Error
	uuidstr = dataload.Uuid
	return
}

func GetDataLoad(id int) (dataload DataLoad, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDataLoadTableName()).Where("id=?", id).First(&dataload).Error
	return
}

func GetDataLoadByUuid(uuid string) (dataload DataLoad, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDataLoadTableName()).Where("uuid=?", uuid).First(&dataload).Error
	return
}

func ListDataLoadByField(fields []string, values []interface{}, limit int, index int) (dataloads []DataLoad, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	var rows *sql.Rows
	if limit == 0 {
		rows, err = db.Table(GetDataLoadTableName()).Where(condition).Rows()
	} else {
		rows, err = db.Table(GetDataLoadTableName()).Where(condition).Limit(limit).Rows()

	}
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var dataload DataLoad
		err = db.ScanRows(rows, &dataload)
		if err != nil {
			return dataloads, err
		}
		dataloads = append(dataloads, dataload)
	}
	return
}

func UpdateDataLoad(dataload map[string]interface{}, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDataLoadTableName()).Where("id=?", dataload["id"]).Updates(dataload).Update(fields).Error
	return
}
