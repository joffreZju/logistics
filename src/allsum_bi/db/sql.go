package db

import (
	"fmt"
	"strings"
)

func MakeInsetSql(tablename string, fields []string, data []map[string]interface{}) (sql string, err error) {
	fieldsstr := strings.Join(fields, ",")
	prefixsql := "insert into " + tablename + " (%s) values %s ;"
	var values []string
	for _, v := range data {
		var value []string
		for _, field := range fields {
			value = append(value, v[field].(string))
		}
		valuestr := strings.Join(value, ",")
		values = append(values, "("+valuestr+")")
	}
	values_str := strings.Join(values, ",")
	sql = fmt.Sprintf(prefixsql, fieldsstr, values_str)
	return
}

func MakeUpdateSql(tablename string, condition map[string]interface{}, fields []string, data map[string]interface{}) (sql string, err error) {
	prefixsql := "update " + tablename + " set %s where %s;"
	var conditions []string
	for k, v := range condition {
		conditionstr := fmt.Sprintf("%s=%v", k, v)
		conditions = append(conditions, conditionstr)
	}
	condition_str := strings.Join(conditions, " AND ")
	var updates []string
	for _, field := range fields {
		updatestr := fmt.Sprintf("%s=%v", field, data[field])
		updates = append(updates, updatestr)
	}
	update_str := strings.Join(updates, ",")
	sql = fmt.Sprintf(prefixsql, update_str, condition_str)
	return
}
