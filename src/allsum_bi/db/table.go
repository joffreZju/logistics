package db

import (
	"allsum_bi/db/conn"
	"allsum_bi/services/util"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

func GetTableDescFromSource(dbid string, schema string, table string, destschema string, desttable string) (createsql string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT column_name, udt_name, character_maximum_length, is_nullable FROM information_schema.columns WHERE table_name=? and table_schema=?", table, schema).Rows()
	if err != nil {
		beego.Error("get create sql fail :", err)
		return
	}
	defer rows.Close()
	formatsql := "create table " + destschema + "." + desttable + " (%s)"
	fieldstr := ""
	for rows.Next() {
		var column_name, data_type, is_nullable string
		var character_maximum_length interface{}
		err = rows.Scan(&column_name, &data_type, &character_maximum_length, &is_nullable)
		if err != nil {
			return createsql, err
		}
		if character_maximum_length != nil {
			data_type = data_type + fmt.Sprintf("(%v)", character_maximum_length)
		}
		if is_nullable == "YES" {
			is_nullable = "null"
		} else {
			is_nullable = "not null"
		}
		if fieldstr == "" {
			fieldstr = fieldstr + "xminstr varchar(32) not null, " + column_name + " " + data_type + " " + is_nullable
		} else {
			fieldstr = fieldstr + "," + column_name + " " + data_type + " " + is_nullable
		}
	}
	createsql = fmt.Sprintf(formatsql, fieldstr)
	return
}

func GetTableDesc(dbid string, schema string, table string, destschema string, desttable string) (createsql string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT column_name, udt_name, character_maximum_length, is_nullable FROM information_schema.columns WHERE table_name=? and table_schema=?", table, schema).Rows()
	if err != nil {
		beego.Error("get create sql fail :", err)
		return
	}
	defer rows.Close()
	formatsql := "create table " + destschema + "." + desttable + " (%s)"
	fieldstr := ""
	for rows.Next() {
		var column_name, data_type, is_nullable string
		var character_maximum_length interface{}
		err = rows.Scan(&column_name, &data_type, &character_maximum_length, &is_nullable)
		if err != nil {
			return createsql, err
		}
		if character_maximum_length != nil {
			data_type = data_type + fmt.Sprintf("(%v)", character_maximum_length)
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
	rows, err := db.Raw("SELECT column_name, data_type from information_schema.columns where table_name=? and table_schema=?", table, schema).Rows()
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
		"where pg_class.relname = ? "+
		"and pg_constraint.contype='p' "+
		"and pg_namespace.nspname= ?;", table, schema).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	pks = make(map[string]bool)
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
	schema_table := schema + "." + table
	sql := "DROP TABLE " + schema_table
	err = Exec(dbid, sql)
	return
}

func DeleteSchemaTable(dbid string, table string) (err error) {
	sql := "DROP TABLE " + table
	err = Exec(dbid, sql)
	return
}

func GetTableFields(dbid string, schema string, table_name string) (fields []string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("select column_name from information_schema.columns where table_name = ? and table_schema = ?", table_name, schema).Rows()
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
		beego.Error("get conn err", err)
		return false
	}
	schema_table := strings.Split(table_name, ".")
	var count int
	err = db.Raw("select count(*) from information_schema.tables where table_schema=? and table_type='BASE TABLE' and table_name=?", schema_table[0], schema_table[1]).Row().Scan(&count)
	if err != nil {
		beego.Error("CheckTableExist err", err)
		return false
	}
	beego.Debug("table count", count, table_name)
	if count == 0 {
		return false
	}
	return true
}

func ListTableData(dbid string, table_name string, conditions []map[string]interface{}, limit int) (datas []map[string]interface{}, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	var condition_strs []string
	for _, condition := range conditions {
		str := fmt.Sprintf("%s %s %v", condition["key"], condition["opt"], condition["value"])
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

func EncodeTableSchema(dbid string, schema string, table string) (schema_table string, err error) {
	conninfo, err := conn.GetConninfo(dbid)
	if err != nil {
		return
	}
	if conninfo.Dbtype == util.PG_DB_TYPE {
		schema_table = schema + "." + table
	} else {
		schema_table = table
	}
	return
}

func DecodeTableSchema(dbid string, schema_table string) (schema string, table string, err error) {
	conninfo, err := conn.GetConninfo(dbid)
	if err != nil {
		return
	}
	if conninfo.Dbtype == util.PG_DB_TYPE {
		schema_table_slice := strings.Split(schema_table, ".")
		if len(schema_table_slice) < 2 {
			return schema, table, fmt.Errorf("error")
		}
		schema = schema_table_slice[0]
		table = schema_table_slice[1]
	} else {
		schema_table_slice := strings.Split(schema_table, ".")
		schema = ""
		table = schema_table_slice[0]
	}
	return
}
