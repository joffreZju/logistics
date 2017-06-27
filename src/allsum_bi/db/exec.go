package db

import (
	"allsum_bi/db/conn"
	"fmt"
)

func Exec(dbid string, Sql string, params ...interface{}) (err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	err = db.Exec(Sql, params...).Error
	return

}

func QueryToFields(dbid string, sql string, fields []string, params ...interface{}) (res []string, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	row := db.Raw(sql, params...)
	for _, field := range fields {
		v, ok := row.Get(field)
		if !ok {
			err = fmt.Errorf("miss data")
			return res, err
		}
		res = append(res, v.(string))
	}
	return
}

func QueryDatas(dbid string, sql string, params ...interface{}) (datas []map[string]interface{}, err error) {
	db, err := conn.GetConn(dbid)
	if err != nil {
		return
	}
	rows, err := db.Exec(sql).Rows()

	if err != nil {
		return
	}
	for rows.Next() {
		var data map[string]interface{}
		err = db.ScanRows(rows, &data)
		if err != nil {
			return
		}
		datas = append(datas, data)
	}
	return
}
