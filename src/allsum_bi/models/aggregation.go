package models

import (
	"allsum_bi/db/conn"
	"database/sql"
	"fmt"

	uuid "github.com/satori/go.uuid"
)

type AggregateOps struct {
	Id           int
	Uuid         string
	Reportid     int
	Name         string
	CreateScript string
	AlterScript  string
	Script       string
	ScriptType   string
	DestTable    string
	Cron         string
	Documents    string
	Status       string
}

func InsertAggregateOps(aggregate_ops AggregateOps) (uuidstr string, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	aggregate_ops.Uuid = uuid.NewV4().String()
	exsit := db.NewRecord(aggregate_ops)
	if !exsit {
		return uuidstr, fmt.Errorf("exsit")
	}
	err = db.Table(GetAggregateOpsTableName()).Create(&aggregate_ops).Error
	uuidstr = aggregate_ops.Uuid
	return
}

func InsertAggregateReturnAggregate(aggregate_ops *AggregateOps) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	aggregate_ops.Uuid = uuid.NewV4().String()

	err = db.Table(GetAggregateOpsTableName()).Create(&aggregate_ops).Error
	return
}

func GetAggregateOps(id int) (aggregate_ops AggregateOps, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetAggregateOpsTableName()).Where("id=?", id).First(&aggregate_ops).Error
	return
}

func GetAggregateOpsByUuid(uuid string) (aggregate_ops AggregateOps, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetAggregateOpsTableName()).Where("uuid=?", uuid).First(&aggregate_ops).Error
	return
}

func ListAggregateOpsByField(fields []string, values []interface{}, limit int, index int) (aggregate_opses []AggregateOps, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	var rows *sql.Rows
	if limit == 0 {
		rows, err = db.Table(GetAggregateOpsTableName()).Where(condition).Rows()
	} else {
		rows, err = db.Table(GetAggregateOpsTableName()).Where(condition).Limit(limit).Rows()
	}
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aggregate_ops AggregateOps
		err = db.ScanRows(rows, &aggregate_ops)
		if err != nil {
			return aggregate_opses, err
		}
		aggregate_opses = append(aggregate_opses, aggregate_ops)
	}
	return
}

func ListAllAggregateOps() (aggregate_opses []AggregateOps, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetAggregateOpsTableName()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var aggregate_ops AggregateOps
		err = db.ScanRows(rows, &aggregate_ops)
		if err != nil {
			return aggregate_opses, err
		}
		aggregate_opses = append(aggregate_opses, aggregate_ops)
	}
	return
}

func UpdateAggregate(aggregate_ops map[string]interface{}, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetAggregateOpsTableName()).Where("id=?", aggregate_ops["id"]).Updates(aggregate_ops).Update(fields).Error
	return
}
