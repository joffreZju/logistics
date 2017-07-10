package models

import (
	"allsum_bi/db/conn"
	"fmt"
	"strings"
)

type TestInfo struct {
	Id        int
	Uuid      string
	Reportid  int
	Documents string
	Filepaths interface{}
	Status    int
}

func InsertTestInfo(testinfo TestInfo) (uuidstr string, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	exsit := db.NewRecord(testinfo)
	if !exsit {
		return uuidstr, err
	}
	filepaths := testinfo.Filepaths.([]string)
	filepatharray := "{" + strings.Join(filepaths, ",") + "}"
	testinfo.Filepaths = filepatharray
	err = db.Table(GetTestInfoTableName()).Create(&testinfo).Error
	uuidstr = testinfo.Uuid
	return
}

func GetTestInfo(id int) (testinfo TestInfo, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetTestInfoTableName()).Where("id=?", id).First(&testinfo).Error
	testinfo.Filepaths = string(testinfo.Filepaths.([]byte))
	testinfo.Filepaths = strings.TrimRight(strings.TrimPrefix(testinfo.Filepaths.(string), "{"), "}")
	testinfo.Filepaths = strings.Split(testinfo.Filepaths.(string), ",")
	return
}

func ListTestInfos(fields []string, values []interface{}, limit int, index int) (testinfos []TestInfo, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	rows, err := db.Table(GetTestInfoTableName()).Where(condition).Limit(limit).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var testinfo TestInfo
		err = db.ScanRows(rows, &testinfo)
		if err != nil {
			return testinfos, err
		}
		testinfo.Filepaths = string(testinfo.Filepaths.([]byte))
		testinfo.Filepaths = strings.TrimRight(strings.TrimPrefix(testinfo.Filepaths.(string), "{"), "}")
		testinfo.Filepaths = strings.Split(testinfo.Filepaths.(string), ",")
		testinfos = append(testinfos, testinfo)
	}
	return
}

func GetTestInfoByReportid(reportid int) (testinfos []TestInfo, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetTestInfoTableName()).Where("reportid=?", reportid).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var testinfo TestInfo
		err = db.ScanRows(rows, &testinfo)
		if err != nil {
			return testinfos, err
		}
		testinfo.Filepaths = string(testinfo.Filepaths.([]byte))
		testinfo.Filepaths = strings.TrimRight(strings.TrimPrefix(testinfo.Filepaths.(string), "{"), "}")
		testinfo.Filepaths = strings.Split(testinfo.Filepaths.(string), ",")
		testinfos = append(testinfos, testinfo)
	}
	return
}

func ListAllTestInfos() (testinfos []TestInfo, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetTestInfoTableName()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var testinfo TestInfo
		err = db.ScanRows(rows, &testinfo)
		if err != nil {
			return testinfos, err
		}
		testinfo.Filepaths = string(testinfo.Filepaths.([]byte))
		testinfo.Filepaths = strings.TrimRight(strings.TrimPrefix(testinfo.Filepaths.(string), "{"), "}")
		testinfo.Filepaths = strings.Split(testinfo.Filepaths.(string), ",")
		testinfos = append(testinfos, testinfo)
	}
	return
}
