package models

import (
	"allsum_bi/db/conn"
	"fmt"
	"time"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

type Synchronous struct {
	Id           int
	Uuid         string
	Owner        string
	CreateScript string
	Script       string
	SourceDbId   string
	SourceTable  string
	DestDbId     string
	DestTable    string
	Cron         string
	Documents    string
	ErrorLimit   int
	Status       string
	Lasttime     time.Time
}

func InsertSynchronous(sync Synchronous) (id int, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	sync.Uuid = uuid.NewV4().String()
	exist := db.NewRecord(sync)
	if !exist {
		return id, fmt.Errorf("exist")
	}
	err = db.Table(GetSynchronousTableName()).Create(&sync).Error
	return sync.Id, err
}

func GetSynchronous(id int) (sync Synchronous, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetSynchronousTableName()).Where("id=?", id).First(&sync).Error
	return
}

func GetSynchronousByUuid(uuid string) (sync Synchronous, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetSynchronousTableName()).Where("uuid=?", uuid).First(&sync).Error
	return
}

func GetSynchronousByOwner(owner string) (syncs []Synchronous, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetSynchronousTableName()).Where("owner=?", owner).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sync Synchronous
		db.ScanRows(rows, &sync)
		syncs = append(syncs, sync)
	}
	return
}

func ListSynchronous() (syncs []Synchronous, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetSynchronousTableName()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var sync Synchronous
		db.ScanRows(rows, &sync)
		syncs = append(syncs, sync)
	}
	return
}

func UpdateSynchronous(sync Synchronous, params ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetSynchronousTableName()).Where("id=?", sync.Id).Updates(sync).Update(params).Error
	return
}

func ListSyncInSourceTables(dbid string, tableNames []string) (syncs map[string]Synchronous, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetSynchronousTableName()).Where("source_db_id = ? and source_table in (?)", dbid, tableNames).Rows()
	syncs = make(map[string]Synchronous)
	for rows.Next() {
		var sync Synchronous
		err = db.ScanRows(rows, &sync)
		if err != nil {
			return syncs, err
		}
		beego.Debug("sourcetable", sync.SourceTable)
		syncs[sync.SourceTable] = sync
	}
	return
}
