package models

import (
	"allsum_bi/db/conn"
	"database/sql"
	"fmt"
	"strings"

	uuid "github.com/satori/go.uuid"
)

type KettleJob struct {
	Id       int
	Uuid     string
	Name     string
	Cron     string
	Kjbpath  string
	Ktrpaths string
	Lock     string
	Status   int
}

func InsertKettleJob(kettlejob KettleJob) (kettlejobres KettleJob, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	kettlejob.Uuid = uuid.NewV4().String()

	err = db.Table(GetKettleJobTableName()).Create(&kettlejob).Error
	kettlejobres = kettlejob
	return
}

func ListKettleJobByField(fields []string, values []interface{}, limit int, index int) (kettlejobs []KettleJob, err error) {
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
		rows, err = db.Table(GetKettleJobTableName()).Where(condition).Rows()

	} else {
		rows, err = db.Table(GetKettleJobTableName()).Where(condition).Limit(limit).Rows()

	}
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var kettlejob KettleJob
		err = db.ScanRows(rows, &kettlejob)
		if err != nil {
			return kettlejobs, err
		}
		kettlejobs = append(kettlejobs, kettlejob)
	}
	return
}

func CountKettleJobByField(fields []string, values []interface{}) (count int, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	var conditions []string
	for _, v := range fields {
		conditions = append(conditions, fmt.Sprintf("%s=?", v))
	}
	conditionstr := strings.Join(conditions, " and ")
	err = db.Table(GetKettleJobTableName()).Where(conditionstr, values...).Count(&count).Error
	return
}

func GetKettleJobByUuid(uuid string) (kettle KettleJob, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetKettleJobTableName()).Where("uuid=?", uuid).First(&kettle).Error
	return
}

func UpdateKettleJob(kettle map[string]interface{}, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetKettleJobTableName()).Where("id=?", kettle["id"]).Updates(kettle).Update(fields).Error
	return
}

func DeleteKettleJobByUuid(uuid string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Exec("delete from "+GetKettleJobTableName()+" where uuid = ?", uuid).Error
	return
}
