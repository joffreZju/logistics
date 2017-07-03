package reportset

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/util"
	"fmt"
	"strings"
)

func GetData(uuid string, conditions []map[string]interface{}) (datas [][]interface{}, err error) {
	reportset, err := models.GetReportSetByReportUuid(uuid)
	if err != nil {
		return
	}
	checkres := checkCondition(reportset.Conditions, conditions)
	if !checkres {
		return datas, fmt.Errorf("paramater error conditions not right")
	}
	sqlstr := reportset.Script
	for _, conditionMap := range conditions {
		field := conditionMap["field"].(string)
		value := conditionMap["value"]
		sqlstr = strings.Replace(sqlstr, "{"+field+"}", value.(string), 1)
	}
	datas, err = db.QueryDatas(util.BASEDB_CONNID, sqlstr)
	return

}

func checkCondition(conditionDbFormat string, Conditions []map[string]interface{}) (checkres bool) {
	//TODO
	return true
}

func CheckConditionFormat(format string) (checkres bool) {
	//TODO
	return true
}
