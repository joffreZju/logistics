package models

import (
	"allsum_bi/db/conn"
	"database/sql"
	"fmt"
	"time"
)

type AggregateLog struct {
	Id          int
	Aggregateid int
	Reportid    int
	Error       string
	Res         string
	Timestamp   time.Time
	Status      int
}

func InsertAggregateLog(aggregate_log AggregateLog) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetAggregateLogTableName()).Create(&aggregate_log).Error
	return
}

func ListAggregateLog(fields []string, values []interface{}, limit int, index int) (aggregate_logs []AggregateLog, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%v", index)
	for _, field := range fields {
		condition = condition + " and " + field + "=?"
	}
	var rows *sql.Rows
	if limit == 0 {
		rows, err = db.Table(GetAggregateLogTableName()).Where(condition, values...).Rows()
	} else {
		rows, err = db.Table(GetAggregateLogTableName()).Rows()
	}
	for rows.Next() {
		var aggregate_log AggregateLog
		err = db.ScanRows(rows, &aggregate_log)
		if err != nil {
			return aggregate_logs, err
		}
		aggregate_logs = append(aggregate_logs, aggregate_log)
	}
	return
}

func UpdateAggregateLog(aggregate_log map[string]interface{}, fields ...interface{}) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetAggregateLogTableName()).Updates(aggregate_log).Update(fields...).Error
	return
}
