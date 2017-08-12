package etl

import (
	"allsum_bi/models"
	"allsum_bi/services/util"
	"time"

	"github.com/astaxie/beego"
)

func SetEtlError(syncid int, msg string) (err error) {
	synclog := models.SynchronousLog{
		Syncid:    syncid,
		Errormsg:  msg,
		Timestamp: time.Now(),
		Status:    util.SYNC_ENABLE,
	}
	_, err = models.InsertSynchronousLogs(synclog)
	if err != nil {
		return
	}
	count, err := models.CountSynchronousLogsBySyncid(syncid, util.SYNC_ENABLE)
	sync, err := models.GetSynchronous(syncid)
	if err != nil {
		return
	}
	if count >= sync.ErrorLimit {
		sync.Status = util.SYNC_ERROR
		err = models.UpdateSynchronous(sync, "status")
		if err != nil {
			beego.Error("update sync err :", err)
		}
		StopCronById(syncid)
	}
	return
}

func CleanEtlError(syncid int) (err error) {
	synclog := models.SynchronousLog{
		Syncid: syncid,
		Status: util.SYNC_DISABLE,
	}
	err = models.UpdateSynchronousLogBySyncId(synclog, "status")
	return

}
