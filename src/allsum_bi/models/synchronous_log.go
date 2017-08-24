package models

import (
	"allsum_bi/db/conn"
	"fmt"
	"time"
)

type SynchronousLog struct {
	Id        int
	Syncid    int
	Errormsg  string
	Res       string
	Timestamp time.Time
	Status    int
}

func InsertSynchronousLogs(syncLogs SynchronousLog) (id int, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	exist := db.NewRecord(syncLogs)
	if exist {
		return id, fmt.Errorf("exist")
	}
	err = db.Table(GetSynchronousLogTableName()).Create(&syncLogs).Error
	return syncLogs.Id, err
}

func ListSynchronousLogs() (synclogs []SynchronousLog, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetSynchronousLogTableName()).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var log SynchronousLog
		err = db.ScanRows(rows, &log)
		if err != nil {
			return synclogs, err
		}
		synclogs = append(synclogs, log)
	}
	return
}

func ListSynchronousLogsBySyncid(syncid int, Status int) (synclogs []SynchronousLog, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	rows, err := db.Table(GetSynchronousLogTableName()).Where("syncid=? AND status=?", syncid, Status).Rows()
	if err != nil {
		return
	}
	defer rows.Close()
	for rows.Next() {
		var log SynchronousLog
		err = db.ScanRows(rows, &log)
		if err != nil {
			return synclogs, err
		}
		synclogs = append(synclogs, log)
	}
	return
}

func CountSynchronousLogsBySyncid(syncid int, status int) (count int, err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetSynchronousLogTableName()).Where("syncid=? And status=?", syncid, status).Count(&count).Error
	return
}

func UpdateSynchronousLogBySyncId(synclog SynchronousLog, params ...string) (err error) {
	db, err := conn.GetBIConn()
	if err != nil {
		return
	}
	err = db.Table(GetSynchronousLogTableName()).Where("syncid=?", synclog.Syncid).Updates(synclog).Update(params).Error
	return
}
