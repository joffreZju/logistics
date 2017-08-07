package demandsvs

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/util"
	"encoding/json"
	"fmt"
)

func ChangeStatus(demanduuid string, demandstatus int, reportstatus int, HandleCompany string) (err error) {
	demand, err := models.GetDemandByUuid(demanduuid)
	if err != nil {
		return
	}
	report, err := models.GetReport(demand.Reportid)
	if err != nil {
		return
	}
	demand.Status = demandstatus
	report.Status = reportstatus
	err = models.UpdateDemand(demand, "status")
	if err != nil {
		return
	}
	err = models.UpdateReport(report, "status")
	if err != nil {
		return
	}

	var authoritymaps []map[string]string
	err = json.Unmarshal([]byte(demand.AssignerAuthority), &authoritymaps)
	if err != nil {
		return
	}
	err = RevokeAuthority(authoritymaps, HandleCompany, demand.Handlerid)
	return
}

func AddAuthority(authoritymaps []map[string]string, UserComp string, Userid int) (err error) {
	userdbname := fmt.Sprintf("%s_%d", UserComp, Userid)
	for _, authoritymap := range authoritymaps {
		dbid := util.BASEDB_CONNID
		schema := authoritymap["schema"]
		err = db.CreateUser(dbid, userdbname)
		if err != nil {
			return
		}
		err = db.AddAuthority(dbid, userdbname, schema, "SELECT")
		if err != nil {
			return
		}
	}
	return
}

func RevokeAuthority(authoritymaps []map[string]string, UserComp string, Handleid int) (err error) {
	userdbname := fmt.Sprintf("%s_%d", UserComp, Handleid)
	for _, authoritymap := range authoritymaps {
		dbid := util.BASEDB_CONNID
		schema := authoritymap["schema"]
		err = db.RevokeAuthority(dbid, userdbname, schema, "SELECT")
		if err != nil {
			return
		}
	}
	return
}
