package models

import (
	"allsum_bi/db/conn"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type ReportSet struct {
	Id               int
	Uuid             string
	Reportid         int
	Script           string
	Params           string
	Resttype         int
	Conditions       string
	EnableEventTypes string
	WebPath          string
	WebfileName      string
	Status           string
}

func InsertReportSet(reportset ReportSet) (uuidstr string, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	reportset.Uuid = uuid.NewV4().String()
	//	_, err = GetReportSetByReportid(reportset.Reportid)
	//	if !strings.Contains(err.Error(), "not found") {
	//		return uuidstr, fmt.Errorf("exsit")
	//	}

	err = db.Table(GetReportSetTableName()).Create(&reportset).Error
	uuidstr = reportset.Uuid
	return
}

func GetReportSetByReportid(reportid int) (reportset ReportSet, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportSetTableName()).Where("reportid=?", reportid).First(&reportset).Error
	if err != nil {
		return
	}
	return
}

func GetReportSet(id int) (reportset ReportSet) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportSetTableName()).Where("id=?", id).First(&reportset).Error
	return
}

func GetReportSetByUuid(uuid string) (reportset ReportSet, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportSetTableName()).Where("uuid=?", uuid).First(&reportset).Error
	return
}

func GetReportSetsByReportUuid(uuid string) (reportsets []ReportSet, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	var report Report
	err = db.Table(GetReportTableName()).Where("uuid=?", uuid).First(&report).Error
	if err != nil {
		return
	}
	rows, err := db.Table(GetReportSetTableName()).Where("reportid=?", report.Id).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var reportset ReportSet
		err = db.ScanRows(rows, &reportset)
		if err != nil {
			return reportsets, err
		}
		reportsets = append(reportsets, reportset)
	}
	return
}

func ListReportSetByField(fields []string, values []interface{}, limit int, index int) (reportsets []ReportSet, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	rows, err := db.Table(GetReportSetTableName()).Where(condition).Limit(limit).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var reportset ReportSet
		err = db.ScanRows(rows, &reportset)
		if err != nil {
			return reportsets, err
		}
		reportsets = append(reportsets, reportset)
	}
	return
}

func UpdateReportSet(reportset ReportSet, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportSetTableName()).Where("id=?", reportset.Id).Updates(reportset).Update(fields).Error
	return
}
