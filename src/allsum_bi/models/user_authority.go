package models

import (
	"allsum_bi/db/conn"
	"database/sql"
	"fmt"
	"strconv"
	"strings"
	"time"
)

type UserAuthority struct {
	Id           int
	Roleid       int
	Rolename     string
	Reportid     int
	Reportsetids interface{}
	Createtime   time.Time
	Limittime    int
}

func InsertUserAuthority(companyid string, userauth UserAuthority) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	var reportsetstrs []string
	for _, v := range userauth.Reportsetids.([]int) {
		reportsetstrs = append(reportsetstrs, strconv.Itoa(v))
	}
	userauth.Reportsetids = "{" + strings.Join(reportsetstrs, ",") + "}"
	err = db.Table(GetUserAuthorityTableName(companyid)).Create(&userauth).Error
	return
}

func GetUserAuthorityByRoleid(companyid string, roleid int) (userauths []UserAuthority, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetUserAuthorityTableName(companyid)).Where("roleid=?", roleid).Rows()
	if err != nil {
		return
	}
	for rows.Next() {
		var userauth UserAuthority
		err = db.ScanRows(rows, &userauth)
		if err != nil {
			return userauths, err
		}
		userauth.Reportsetids = string(userauth.Reportsetids.([]byte))
		userauth.Reportsetids = strings.TrimRight(strings.TrimPrefix(userauth.Reportsetids.(string), "{"), "}")
		reportsets := []int{}
		for _, v := range strings.Split(userauth.Reportsetids.(string), ",") {
			Reportsetid, err := strconv.Atoi(v)
			if err != nil {
				return userauths, err
			}
			reportsets = append(reportsets, Reportsetid)
		}
		userauth.Reportsetids = reportsets
		userauths = append(userauths, userauth)
	}
	return
}

func GetAuthority(companyid string, id int) (userauth UserAuthority, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("id=?", id).First(&userauth).Error
	userauth.Reportsetids = string(userauth.Reportsetids.([]byte))
	userauth.Reportsetids = strings.TrimRight(strings.TrimPrefix(userauth.Reportsetids.(string), "{"), "}")
	reportsets := []int{}
	for _, v := range strings.Split(userauth.Reportsetids.(string), ",") {
		Reportsetid, err := strconv.Atoi(v)
		if err != nil {
			return userauth, err
		}
		reportsets = append(reportsets, Reportsetid)
	}
	userauth.Reportsetids = reportsets
	return
}

func GetAuthorityByRoleReport(companyid string, roleid int, reportid int) (userauth UserAuthority, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("roleid=? and reportid=?", roleid, reportid).First(&userauth).Error
	reportsets := []int{}
	for _, v := range strings.Split(userauth.Reportsetids.(string), ",") {
		Reportsetid, err := strconv.Atoi(v)
		if err != nil {
			return userauth, err
		}
		reportsets = append(reportsets, Reportsetid)
	}
	userauth.Reportsetids = reportsets
	return
}

func CheckAuthorityReport(companyid string, roleid int, reportid int) (checkres bool) {
	db, err := conn.GetBIConn()
	if err != nil {
		return false
	}
	var userauth UserAuthority
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("roleid=? and reportid=?", roleid, reportid).First(&userauth).Error
	if err == nil {
		return true
	}
	return false
}

func CheckAuthorityReportSet(companyid string, roleid int, reportsetid int) (checkres bool) {
	db, err := conn.GetBIConn()
	if err != nil {
		return false
	}
	Reportset, err := GetReportSet(reportsetid)
	if err != nil {
		return false
	}

	var userauth UserAuthority
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("roleid=? and reportid=?", roleid, Reportset.Reportid).First(&userauth).Error
	if err != nil {
		return false
	}
	userauth.Reportsetids = string(userauth.Reportsetids.([]byte))
	userauth.Reportsetids = strings.TrimRight(strings.TrimPrefix(userauth.Reportsetids.(string), "{"), "}")
	for _, v := range strings.Split(userauth.Reportsetids.(string), ",") {
		Reportsetid, err := strconv.Atoi(v)
		if err != nil {
			return false
		}
		if Reportsetid == reportsetid {
			return true
		}
	}
	return false
}

func ListUserAuthority(companyid string, fields []string, values []interface{}, limit int, index int) (userauthoritys []UserAuthority, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("companyid=%s and id>%d", companyid, index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	var rows *sql.Rows
	if limit == 0 {
		rows, err = db.Table(GetUserAuthorityTableName(companyid)).Where(condition).Rows()
	} else {
		rows, err = db.Table(GetUserAuthorityTableName(companyid)).Where(condition).Limit(limit).Rows()
	}
	if err != nil {
		return
	}
	defer rows.Close()

	for rows.Next() {
		var userauthority UserAuthority
		err = db.ScanRows(rows, &userauthority)
		if err != nil {
			return userauthoritys, err
		}
		userauthoritys = append(userauthoritys, userauthority)
	}
	return
}

func UpdateUserAuthority(companyid string, userauth UserAuthority, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	var reportsetstrs []string
	for _, v := range userauth.Reportsetids.([]int) {
		reportsetstrs = append(reportsetstrs, strconv.Itoa(v))
	}
	userauth.Reportsetids = "{" + strings.Join(reportsetstrs, ",") + "}"
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("id=?", userauth.Id).Updates(userauth).Update(fields).Error
	return
}

func DeleteUserAuthority(companyid string, id int) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetUserAuthorityTableName(companyid)).Where("id=?", id).Delete(UserAuthority{}).Error
	return
}
