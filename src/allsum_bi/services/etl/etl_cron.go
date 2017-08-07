package etl

import (
	"allsum_bi/models"
	"allsum_bi/util"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

//var CronEtls map[int]*cron.Cron
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
		err = StartEtl(etlrecord)
		if err != nil {
			etlrecord.Status = util.SYNC_ERROR
			models.UpdateSynchronous(etlrecord, "status")
		}
		//	AddCronWithFullScript(etlrecord.Id, etlrecord.Cron, etlrecord.Script)
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
		if p, ok := etltaskmap[id]; ok {
			delete(etltaskmap, id)
			defer func() {
				if r := recover(); r != nil {
					beego.Error("do etl crash ", r)
				}
			}()
			p.Stop()
		}
		script, err := MakeRunScript(fullscript)
		if err != nil {
			return
		}
		DoETL(id, []byte(script))
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
	err = CronEtls[id].AddFunc(cronstr, func() {
		if p, ok := etltaskmap[id]; ok {
			delete(etltaskmap, id)
			defer func() {
				if r := recover(); r != nil {
					beego.Error("do etl crash ", r)
				}
			}()
			p.Stop()
		}
		DoETL(id, []byte(script))
	})
	return
}

func StopCronBySyncUuid(uuid string) (err error) {
	sync, err := models.GetSynchronousByUuid(uuid)
	if err != nil {
		return
	}
	StopCronById(sync.Id)
	sync.Status = util.SYNC_STOP
	err = models.UpdateSynchronous(sync, "status")
	return
}

func StopCronById(id int) {
	CronEtls[id].Stop()
	p := etltaskmap[id]
	delete(etltaskmap, id)
	defer func() {
		if r := recover(); r != nil {
			beego.Error("do etl crash ", r)
		}
	}()
	p.Stop()
}

func StopAll() {
	defer func() {
		if r := recover(); r != nil {
			beego.Error("do etl crash ", r)
		}
	}()
	for id, v := range CronEtls {
		v.Stop()
		p := etltaskmap[id]
		delete(etltaskmap, id)
		p.Stop()
	}
}
