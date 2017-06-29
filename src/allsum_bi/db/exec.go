package db

import (
	"allsum_bi/db/conn"

	"github.com/astaxie/beego"
)

func Exec(dbid string, Sql string, params ...interface{}) (err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	err = db.Exec(Sql, params...).Error
	return

}

func QueryDatas(dbid string, sql string, params ...interface{}) (datas [][]interface{}, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	beego.Debug("sql: ", sql, params)
	rows, err := db.Raw(sql, params...).Rows()
	if err != nil {
		beego.Error("get rows err: ", err)
		return
	}
	defer rows.Close()
	columns, err := rows.ColumnTypes()
	collen := len(columns)
	if err != nil {
		beego.Error("get rows err ", err)
		return
	}
	for rows.Next() {
		data := make([]interface{}, collen)
		dataaddr := make([]interface{}, collen)
		for i, _ := range dataaddr {
			dataaddr[i] = &data[i]
		}
		rows.Scan(dataaddr...)
		datas = append(datas, data)
	}
	beego.Debug("datas : ", datas)
	return

	//	for rows.Next() {
	//		var data []interface{}
	//		rows.Scan(&data)
	//		beego.Debug("testdata:", rows)
	//
	//		err = db.ScanRows(rows, &data)
	//		if err != nil {
	//			beego.Error("err scan", err)
	//			return
	//		}
	//		datas = append(datas, data)
	//	}
	//	return
}
