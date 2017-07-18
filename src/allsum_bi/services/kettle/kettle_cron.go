package kettle

import (
	"allsum_bi/models"
	"allsum_bi/util"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

var CronJobs map[int]*cron.Cron

func init() {
	CronJobs = map[int]*cron.Cron{}
}
func StartJobCron() (err error) {
	jobs, err := models.ListKettleJobByField([]string{"status"}, []interface{}{util.KETTLEJOB_RIGHT}, 0, 0)
	if err != nil {
		return
	}
	joblimit, err := beego.AppConfig.Int("kettle::joblimit")
	if err != nil {
		joblimit = 5
	}
	jobnum := 0
	for _, job := range jobs {
		jobnum += 1
		if jobnum > joblimit {
			return nil
		}
		err := AddCron(job.Id, job.Cron, job.Kjbpath)
		if err != nil {
			return err
		}
	}
	return
}

func AddCron(jobid int, cronstr string, jobfilepath string) (err error) {
	if jobc, ok := CronJobs[jobid]; ok {
		jobc.Stop()
	}
	CronJobs[jobid] = cron.New()
	CronJobs[jobid].Start()
	err = CronJobs[jobid].AddFunc(cronstr, func() {
		fmtstr, err := ExecJob(jobfilepath)
		if err != nil {
			beego.Error("ExecJob fail :", fmtstr)
			kettlelog := models.KettleJobLog{
				KettleJobId: jobid,
				ErrorInfo:   fmtstr,
				Status:      util.KETTLEJOB_RIGHT,
			}
			models.InsertKettleJobLog(kettlelog)
			return
		}
	})
	return
}
