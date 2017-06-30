package db

import (
	"allsum_bi/db/conn"
	"allsum_bi/util"
	"fmt"

	"github.com/astaxie/beego"
)

func GetCompanyTable(ownerid string, table string) (schema_table string) {
	tableName := fmt.Sprintf("%s%s.%s", util.BI_COMMENT_PREFIX, ownerid, table)
	return tableName
}

func GetCompanySchema(ownerid string) (schema string) {
	return util.BI_COMMENT_PREFIX + ownerid
}

func GetSystemTable(table string) (schema_table string) {
	tableName := util.BI_SCHEMA + "." + table
	return tableName
}

func GetManagerTable(table string) (schema_table string) {
	tableName := util.BI_MANAGER + "." + table
	return tableName
}

func CreateManagerSchema() (err error) {
	return CreateSchema(util.BI_MANAGER)
}

func CreateSystemSchema() (err error) {
	return CreateSchema(util.BI_SCHEMA)
}

func CreateSchema(schemaName string) (err error) {
	db, err := conn.GetConn(util.BASEDB_CONNID)
	if err != nil {
		return
	}
	var exist bool
	db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = ?)", schemaName).Row().Scan(&exist)
	beego.Debug("exist:", exist)
	if exist {
		return
	} else {
		sql := fmt.Sprintf("create schema %v", schemaName)
		db.Exec(sql)
	}
	return
}

func SchemaExist(schemaName string, connid string) (exist bool) {
	db, err := conn.GetConn(connid)
	if err != nil {
		return
	}
	db.Raw("SELECT EXISTS(SELECT 1 FROM information_schema.schemata WHERE schema_name = '?')", schemaName).Scan(&exist)
	return exist
}

func ListSchemaTable(dbid string, schema string) (tablenames []string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT table_name FROM information_schema.tables WHERE table_schema = ?", schema).Rows()
	if err != nil {
		beego.Error("1")
		return
	}
	defer rows.Close()

	for rows.Next() {
		var tableName string
		err = rows.Scan(&tableName)
		if err != nil {
			beego.Error("2")
			return tablenames, err
		}
		tablenames = append(tablenames, tableName)
	}
	return
}

func ListScheme(dbid string) (schemalist []string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Raw("SELECT schema_name FROM information_schema.schemata").Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var schema string
		err = rows.Scan(&schema)
		if err != nil {
			return schemalist, err
		}
		schemalist = append(schemalist, schema)
	}
	return
}
