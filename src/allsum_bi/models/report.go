package models

import (
	"allsum_bi/db/conn"

	"github.com/satori/go.uuid"
)

type Report struct {
	Id          int
	Uuid        string
	Demandid    int
	Owner       string
	Name        string
	Reporttype  int
	Description string
	Grouppath   string
	Level       int
	Status      int
}

func InsertReport(report Report) (res_report Report, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	report.Uuid = uuid.NewV4().String()
	//	_, err = GetReportByDemandid(report.Demandid)
	//	if err == nil || !strings.Contains(err.Error(), "not found") {
	//		return res_report, fmt.Errorf("exist", err)
	//	}
	err = db.Table(GetReportTableName()).Create(&report).Error
	res_report = report
	return
}

func GetReportByDemandid(demandid int) (report Report, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportTableName()).Where("demandid=?", demandid).First(&report).Error
	return
}

func GetReportByUuid(uuid string) (report Report, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportTableName()).Where("uuid=?", uuid).First(&report).Error
	return
}

func GetReport(id int) (report Report, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportTableName()).Where("id=?", id).First(&report).Error
	return
}

func GetReportDemand(id int) (demand Demand, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	var report Report
	err = db.Table(GetReportTableName()).Where("id=?", id).First(&report).Error
	if err != nil {
		return
	}
	err = db.Table(GetDemandTableName()).Where("id=?", report.Demandid).First(&demand).Error
	if err != nil {
		return
	}
	return
}

func UpdateReport(report map[string]interface{}, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetReportTableName()).Where("id=?", report["id"]).Updates(report).Update(fields).Error
	return
}
