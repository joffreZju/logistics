package models

import (
	"allsum_bi/db/conn"
	"fmt"
	"time"

	"github.com/satori/go.uuid"
)

type Demand struct {
	Id                int
	Uuid              string
	Owner             int
	OwnerName         string
	Reportid          int
	Description       string
	Contactid         int
	Handleid          int
	HandleName        string
	Assignerid        int
	AssignerName      string
	Price             float32
	Deadline          time.Time
	Resultcode        string
	Inittime          time.Time
	Assignetime       time.Time
	Complettime       time.Time
	AssignerAuthority string
	DocUrl            string
	DocName           string
	Status            int
}

func InsertDemand(demand Demand) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	demand.Uuid = uuid.NewV4().String()
	exist := db.NewRecord(demand)
	if exist {
		return fmt.Errorf("exist")
	}
	err = db.Table(GetDemandTableName()).Create(&demand).Error
	return
}

func GetDemandByUuid(uuid string) (demand Demand, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDemandTableName()).Where("uuid=?", uuid).First(&demand).Error
	return
}

func ListDemandByField(fields []string, values []interface{}, limit int, index int) (demands []Demand, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	condition := fmt.Sprintf("id>%d", index)
	for i, v := range fields {
		condition = condition + fmt.Sprintf(" and %s=%v", v, values[i])
	}
	rows, err := db.Table(GetDemandTableName()).Where(condition).Limit(limit).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var demand Demand
		err = db.ScanRows(rows, &demand)
		if err != nil {
			return demands, err
		}
		demands = append(demands, demand)
	}
	return
}

func UpdateDemand(demand Demand, fields ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetDemandTableName()).Where("id=?", demand.Id).Updates(demand).Update(fields).Error
	return
}
