package models

import (
	"allsum_bi/db/conn"
	"fmt"
)

type KettleJobLog struct {
	Id          int
	KettleJobId int
	ErrorInfo   string
	Status      int
}

func InsertKettleJobLog(kettlejoblog KettleJobLog) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}

	err = db.Table(GetKettleJobLogTableName()).Create(kettlejoblog).Error
	return
}

func ListKettleJobLogByField(fields []string, values []interface{}, limit int, index int) (kettlejoblogs []KettleJobLog, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	rows, err := db.Table(GetKettleJobLogTableName()).Where(condition).Limit(limit).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var kettlejoblog KettleJobLog
		err = db.ScanRows(rows, &kettlejoblog)
		if err != nil {
			return kettlejoblogs, err
		}
		kettlejoblogs = append(kettlejoblogs, kettlejoblog)
	}
	return
}
