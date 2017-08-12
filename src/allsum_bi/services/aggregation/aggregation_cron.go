package aggregation

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/util"
	"common/lib/service_client/oaclient"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var CronAggregate map[int]*cron.Cron
var AggregateLock map[int]bool

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
	if lock, ok := AggregateLock[id]; ok && lock {
		beego.Info("aggregate locked wait to Next round")
		return
	}
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
	if report.Reporttype == util.REPORT_TYPE_PRIVATE {
		schemas = []string{db.GetCompanySchema(demand.Owner)}
	} else {
		schemas, err = oaclient.GetAllCompanySchema()
		if err != nil {
			return
		}
	}
	AggregateLock[id] = true
	for _, schema := range schemas {
		desttable, _ := db.EncodeTableSchema(util.BASEDB_CONNID, schema, aggregate.DestTable)
		flush_script_real := strings.Replace(aggregate.Script, util.SCRIPT_TABLE, desttable, util.SCRIPT_LIMIT)
		flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, util.SCRIPT_LIMIT)

		err = db.Exec(util.BASEDB_CONNID, flushsqlscript)

		if err != nil {
			aggregates := models.AggregateLog{
				Aggregateid: id,
				Reportid:    aggregate.Reportid,
				Error:       err.Error(),
				Res:         "",
				Timestamp:   time.Now(),
				Status:      util.IS_OPEN,
			}
			models.InsertAggregateLog(aggregates)
			return
		}
	}
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
