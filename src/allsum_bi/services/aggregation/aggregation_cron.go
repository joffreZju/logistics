package aggregation

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/util"
	"common/lib/service_client/oaclient"
	"strings"
	"sync"
	"time"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var CronAggregate map[int]*cron.Cron
var AggregateLock map[int]bool

var maplock sync.Mutex

func init() {
	CronAggregate = map[int]*cron.Cron{}
	AggregateLock = map[int]bool{}
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

func StopAggregate(id int) (err error) {
	if c, ok := CronAggregate[id]; ok {
		c.Stop()
	}
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
	//TODO is easy
	report, err := models.GetReport(aggregate.Reportid)
	if err != nil {
		return
	}
	var schemas []string
	beego.Debug("reporttype: ", report.Reporttype)
	if report.Reporttype == util.REPORT_TYPE_COMMON {
		schemas, err = oaclient.GetAllCompanySchema()
		if err != nil {
			return
		}
	} else {
		schemas = []string{db.GetCompanySchema(demand.Owner)}
	}
	maplock.Lock()
	if lock, ok := AggregateLock[id]; ok && lock {
		beego.Info("aggregate locked wait to Next round")
		maplock.Unlock()
		return
	}
	AggregateLock[id] = true
	maplock.Unlock()
	for _, schema := range schemas {
		desttable, _ := db.EncodeTableSchema(util.BASEDB_CONNID, schema, aggregate.DestTable)
		flush_script_real := strings.Replace(aggregate.Script, util.SCRIPT_TABLE, desttable, util.SCRIPT_LIMIT)
		flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, util.SCRIPT_LIMIT)

		err = db.Exec(util.BASEDB_CONNID, flushsqlscript)

		if err != nil {
			beego.Error("aggregates err:", err)
			beego.Error("aggregates errsql:", flushsqlscript)
			aggregates := models.AggregateLog{
				Aggregateid: id,
				Reportid:    aggregate.Reportid,
				Error:       err.Error(),
				Res:         "",
				Timestamp:   time.Now(),
				Status:      util.IS_OPEN,
			}
			models.InsertAggregateLog(aggregates)
			maplock.Lock()
			defer maplock.Unlock()
			AggregateLock[id] = false
			return
		}
	}
	maplock.Lock()
	defer maplock.Unlock()
	AggregateLock[id] = false
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
