package aggregation

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/util"
	"strings"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var CronAggregate map[int]*cron.Cron

func init() {
	CronAggregate = map[int]*cron.Cron{}
}

func StartAggregateCron() (err error) {
	aggregates, err := models.ListAllAggregateOps()
	if err != nil {
		return err
	}
	for _, aggregate := range aggregates {
		if aggregate.Status != util.AGGREGATE_STARTED {
			continue
		}
		AddCronWithFlushScript(aggregate.Id, aggregate.Cron, aggregate.Script)
	}
	return
}

func AddCronWithFlushScript(id int, cronstr string, flushscript string) (err error) {
	if c, ok := CronAggregate[id]; ok {
		c.Stop()
	}
	CronAggregate[id] = cron.New()
	CronAggregate[id].Start()
	err = CronAggregate[id].AddFunc(cronstr, func() {
		err = DoAggregate(id, flushscript)
		if err != nil {
			return
		}
	})
	return
}

func DoAggregate(id int, flushsqlscript string) (err error) {
	aggregate, err := models.GetAggregateOps(id)
	if err != nil {
		return
	}
	demand, err := models.GetReportDemand(aggregate.Reportid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(demand.Owner)
	flush_script_real := strings.Replace(aggregate.Script, util.SCRIPT_TABLE, aggregate.DestTable, util.SCRIPT_LIMIT)
	flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, util.SCRIPT_LIMIT)

	err = db.Exec(util.BASEDB_CONNID, flushsqlscript)
	if err != nil {
		return
	}
	return
}

func TestAddCronWithFlushScript(cronstr string, flushscript string) (err error) {
	tc := cron.New()
	tc.Start()
	err = tc.AddFunc(cronstr, func() {
		beego.Debug("cron str is right")
		//err = DoAggregate(flushscript)
		//if err != nil {
		//	return
		//}
	})
	tc.Stop()
	return
}
