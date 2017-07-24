package kettlesvs

import (
	"allsum_bi/models"
	"allsum_bi/util"
	"encoding/json"
	"path"
	"strings"

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
			job.Status = util.KETTLEJOB_FAIL
			models.UpdateKettleJob(job, "status")
			continue
		}
		var kjbmap map[string]string
		err = json.Unmarshal([]byte(job.Kjbpath), &kjbmap)
		if err != nil {
			beego.Error("err: ", err)
			return
		}
		err := AddCron(job.Id, job.Cron, path.Base(kjbmap["urlpath"]))
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
		fmtstr, _ := ExecJob(jobfilepath)
		if strings.Contains(fmtstr, "ERROR") {
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

func StopCron(jobid int) (err error) {
	if jobc, ok := CronJobs[jobid]; ok {
		jobc.Stop()
	}
	return
}
