package etl

import (
	"allsum_bi/models"
	"allsum_bi/util"

	"github.com/robfig/cron"
)

var CronEtls map[int]*cron.Cron

func init() {
	CronEtls = map[int]*cron.Cron{}
}

func StartEtlCron() (err error) {
	etlRecords, err := models.ListSynchronous()
	if err != nil {
		return err
	}
	for _, etlrecord := range etlRecords {
		if etlrecord.Status != util.SYNC_STARTED {
			continue
		}
		AddCronWithFullScript(etlrecord.Id, etlrecord.Cron, etlrecord.Script)
	}
	return
}

func AddCronWithFullScript(id int, cronstr string, fullscript string) (err error) {
	if etlc, ok := CronEtls[id]; ok {
		etlc.Stop()
	}
	CronEtls[id] = cron.New()
	CronEtls[id].Start()
	err = CronEtls[id].AddFunc(cronstr, func() {
		script, err := MakeRunScript(fullscript)
		if err != nil {
			return
		}
		DoETL([]byte(script))
	})
	if err != nil {
		return
	}
	return
}

func AddCronWithScript(id int, cronstr string, script string) (err error) {
	if etlc, ok := CronEtls[id]; ok {
		etlc.Stop()
	}
	CronEtls[id] = cron.New()
	CronEtls[id].Start()
	err = CronEtls[id].AddFunc(cronstr, func() { DoETL([]byte(script)) })
	return
}

func StopCronById(id int) {
	CronEtls[id].Stop()
}

func StopAll() {
	for _, v := range CronEtls {
		v.Stop()
	}
}
