package demandsvs

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/util"
	"common/lib/service_client/oaclient"
	"encoding/json"
	"fmt"
	"strings"
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
		db.CreateSchema(schema)
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

func GetHandlerDbUser(compid string, handlerid int) (dbusername string) {
	return fmt.Sprintf("%s_%d", compid, handlerid)
}

func GetHandlerUserFromOA() (users []map[string]interface{}, err error) {
	roles, err := oaclient.GetAllRoleByCompany(util.COMPANY_NO)
	if err != nil {
		return
	}
	var biroles []map[string]interface{}
	for _, rolebi := range roles {
		name := rolebi["Name"].(string)
		if strings.Contains(name, util.REG_ROLE_BI) {
			biroles = append(biroles, rolebi)
		}
	}
	if len(biroles) == 0 {
		return
	}
	var developroles []map[string]interface{}
	for _, rolebi := range biroles {
		name := rolebi["Name"].(string)
		if strings.Contains(name, util.REG_ROLE_DEVELOP) {
			developroles = append(developroles, rolebi)
		}
	}
	if len(developroles) == 0 {
		developroles = biroles
	}
	for _, developrole := range developroles {
		var roleid int
		roleid = int(developrole["Id"].(float64))
		userapis, err := oaclient.GetAllUserByRole(util.COMPANY_NO, roleid)
		if err != nil {
			continue
		}
		users = append(users, userapis...)
	}
	return
}
