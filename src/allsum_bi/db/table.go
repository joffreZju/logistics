package db

import (
	"allsum_bi/db/conn"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

func GetTableDesc(dbid string, schema string, table string, destschema string, desttable string) (createsql string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT column_name, data_type, character_maximum_length, is_nullable FROM information_schema.columns WHERE table_name='?' and table_schema='?'").Rows()
	if err != nil {
		beego.Error("get create sql fail :", err)
		return
	}
	defer rows.Close()
	formatsql := "create table " + destschema + "." + desttable + " (%s)"
	fieldstr := ""
	for rows.Next() {
		var column_name, data_type, is_nullable string
		var character_maximum_length int
		err = rows.Scan(&column_name, &data_type, &character_maximum_length, &is_nullable)
		if err != nil {
			return createsql, err
		}
		beego.Debug("character_maximum_length :", character_maximum_length)
		if character_maximum_length != 0 {
			data_type = data_type + fmt.Sprintf("(%d)", character_maximum_length)
		}
		if is_nullable == "YES" {
			is_nullable = "null"
		} else {
			is_nullable = "not null"
		}
		if fieldstr == "" {
			fieldstr = fieldstr + column_name + " " + data_type + " " + is_nullable
		} else {
			fieldstr = fieldstr + "," + column_name + " " + data_type + " " + is_nullable
		}
	}
	createsql = fmt.Sprintf(formatsql, fieldstr)
	return
}

func GetTableColumes(dbid string, table string, schema string) (columns []map[string]string, err error) {
	pks, err := GetTablePk(dbid, schema, table)
	if err != nil {
		return
	}

	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT column_name, data_type from information_schema.columns where table_name='?' and table_schema='?'").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var column_name, data_type string
		err = rows.Scan(&column_name, &data_type)
		if err != nil {
			return columns, err
		}
		is_pk := "false"
		if _, ok := pks[column_name]; ok {
			is_pk = "true"
		}
		column := map[string]string{
			"column_name": column_name,
			"data_type":   data_type,
			"is_pk":       is_pk,
		}
		columns = append(columns, column)
	}
	return
}

func GetTablePk(dbid string, schema string, table string) (pks map[string]bool, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("select pg_attribute.attname as colname from pg_constraint "+
		"inner join pg_class on pg_constraint.conrelid = pg_class.oid "+
		"inner join pg_attribute on pg_attribute.attrelid = pg_class.oid and  pg_attribute.attnum = ANY(pg_constraint.conkey) "+
		"inner join pg_type on pg_type.oid = pg_attribute.atttypid "+
		"inner join pg_namespace on pg_constraint.connamespace = pg_namespace.oid "+
		"where pg_class.relname = '?' "+
		"and pg_constraint.contype='p' "+
		"and pg_namespace.nspname='?';", table, schema).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var colname string
		err = rows.Scan(&colname)
		if err != nil {
			return pks, err
		}
		pks[colname] = true
	}
	return
}

func DeleteTable(dbid string, schema string, table string) (err error) {
	err = Exec(dbid, "DROP TABLE ?.?", schema, table)
	return
}

func DeleteSchemaTable(dbid string, table string) (err error) {
	err = Exec(dbid, "DROP TABLE ?", table)
	return
}

func GetTableFields(dbid string, schema string, table_name string) (fields []string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("select column_name from information_schema.columns where table_name = '?' and table_schema = '?'", table_name, schema).Rows()
	if err != nil {
		beego.Error("get table fields err : ", err)
		return
	}
	defer rows.Close()
	for rows.Next() {
		var column_name string
		err = rows.Scan(&column_name)
		if err != nil {
			return fields, err
		}
		fields = append(fields, column_name)
	}
	return
}

func CheckTableExist(dbid string, table_name string) (isexist bool) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	isexist = db.HasTable(table_name)
	return
}

func ListTableData(dbid string, table_name string, conditions []map[string]interface{}, limit int) (datas []map[string]interface{}, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	var condition_strs []string
	for _, conditionmap := range conditions {
		str := fmt.Sprintf("%s %s %v", conditionmap["key"], conditionmap["opt"], conditionmap["value"])
		condition_strs = append(condition_strs, str)
	}
	condition_str := strings.Join(condition_strs, " AND ")
	//schema_table := strings.Split(table_name, ".")
	//columns, err := GetTableColumes(dbid, schema_table[0], schema_table[1])
	if err != nil {
		return
	}
	rows, err := db.Table(table_name).Where(condition_str).Limit(limit).Rows()
	if err != nil {
		return
	}
	rows.Close()
	for rows.Next() {
		var data map[string]interface{}
		err = db.ScanRows(rows, &data)
		if err != nil {
			return datas, err
		}
		datas = append(datas, data)
	}
	return
}
