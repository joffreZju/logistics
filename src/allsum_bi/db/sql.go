package db

import (
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

func MakeInsetSql(tablename string, fields []string, data []map[string]interface{}) (sql string, params []interface{}, err error) {
	fieldsstr := strings.Join(fields, ",")
	prefixsql := "insert into " + tablename + " (%s) values %s ;"
	var values []string
	for _, v := range data {
		var value []string
		for _, field := range fields {
			value = append(value, "?")
			params = append(params, v[field])
		}
		valuestr := strings.Join(value, ",")
		values = append(values, "("+valuestr+")")
	}
	values_str := strings.Join(values, ",")
	sql = fmt.Sprintf(prefixsql, fieldsstr, values_str)
	beego.Debug("insert sql: ", sql, params)
	return
}

func MakeUpdateSql(tablename string, condition map[string]interface{}, fields []string, data map[string]interface{}) (sql string, params []interface{}, err error) {
	prefixsql := "update " + tablename + " set %s where %s;"

	var updates []string
	for _, field := range fields {
		updatestr := fmt.Sprintf("%s=?", field)
		updates = append(updates, updatestr)
		params = append(params, data[field])
	}
	var conditions []string
	for k, v := range condition {
		conditionstr := fmt.Sprintf("%s=?", k)
		conditions = append(conditions, conditionstr)
		params = append(params, v)
	}
	condition_str := strings.Join(conditions, " AND ")
	update_str := strings.Join(updates, ",")
	sql = fmt.Sprintf(prefixsql, update_str, condition_str)
	return
}
