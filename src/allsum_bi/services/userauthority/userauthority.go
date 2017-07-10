package userauthority

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/util"
	"strings"
	"time"
)

func AddUserReportsetAuthority(companyid string, roleid int, rolename string, reportsetuuid string) (err error) {
	schema := db.GetCompanySchema(companyid)
	exist := db.CheckTableExist(util.BASEDB_CONNID, models.GetUserAuthorityTableName(companyid))
	if !exist {
		exist = db.SchemaExist(schema, util.BASEDB_CONNID)
		if !exist {
			err = db.CreateSchema(schema)
			if err != nil {
				return
			}
		}
		sql := strings.Replace(util.CREATE_USER_AUTHOIRTY, util.SCRIPT_SCHEMA, schema, -1)
		err = db.Exec(util.BASEDB_CONNID, sql)
		if err != nil {
			return
		}
	}
	reportset, err := models.GetReportSetByUuid(reportsetuuid)
	if err != nil {
		return
	}
	userauth, err := models.GetAuthorityByRoleReport(companyid, roleid, reportset.Reportid)
	if err != nil && strings.Contains(err.Error(), "not found") {
		userauth := models.UserAuthority{
			Roleid:       roleid,
			Rolename:     rolename,
			Reportid:     reportset.Reportid,
			Reportsetids: []int{reportset.Id},
			Createtime:   time.Now(),
			Limittime:    0,
		}
		err = models.InsertUserAuthority(companyid, userauth)
		return err
	}
	reportsetids := userauth.Reportsetids.([]int)
	for _, v := range reportsetids {
		if v == reportset.Id {
			return
		}
	}
	reportsetids = append(reportsetids, reportset.Id)
	userauth.Reportsetids = reportsetids

	err = models.UpdateUserAuthority(companyid, userauth, "reportsetids")
	return err
}

func RemoveReportSetAuthority(companyid string, roleid int, reportsetuuid string) (err error) {
	reportset, err := models.GetReportSetByUuid(reportsetuuid)
	if err != nil {
		return
	}
	userauth, err := models.GetAuthorityByRoleReport(companyid, roleid, reportset.Reportid)
	if err != nil {
		return
	}
	newreportsetids := []int{}
	for _, v := range userauth.Reportsetids.([]int) {
		newreportsetids = append(newreportsetids, v)
	}
	userauth.Reportsetids = newreportsetids

	err = models.UpdateUserAuthority(companyid, userauth, "reportsetids")
	return err
}

func RemoveReportAuthority(companyid string, roleid int, reportuuid string) (err error) {
	report, err := models.GetReportByUuid(reportuuid)
	if err != nil {
		return
	}
	userauth, err := models.GetAuthorityByRoleReport(companyid, roleid, report.Id)
	if err != nil {
		return
	}
	err = models.DeleteUserAuthority(companyid, userauth.Id)
	return err
}
