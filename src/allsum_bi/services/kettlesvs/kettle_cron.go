package kettlesvs

import (
	"allsum_bi/models"
	"allsum_bi/services/util"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"sync"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	"github.com/robfig/cron"
)

//var CronJobs map[int]*cron.Cron
var CronJobs map[int]map[string]interface{}
var maplock sync.Mutex

func init() {
	//	CronJobs = map[int]*cron.Cron{}
	CronJobs = map[int]map[string]interface{}{}
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
			jobmap := map[string]interface{}{
				"id":     job.Id,
				"status": job.Status,
			}
			models.UpdateKettleJob(jobmap, "status")
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
		jobc["cron"].(*cron.Cron).Stop()
		delete(CronJobs, jobid)

	}
	CronJobs[jobid] = map[string]interface{}{
		"cron": cron.New(),
	}
	CronJobs[jobid]["cron"].(*cron.Cron).Start()
	err = CronJobs[jobid]["cron"].(*cron.Cron).AddFunc(cronstr, func() {
		maplock.Lock()
		if cmdi, ok := CronJobs[jobid]["cmd"]; ok { //结束上一轮没有跑完的任务
			beego.Info("stop cmd for jobid:", jobid)
			content := make([]byte, 5000)
			num, err := CronJobs[jobid]["stdout"].(*os.File).Read(content)
			beego.Info("process out:", num)
			//content, err := ioutil.ReadAll(CronJobs[jobid]["stdout"].(*os.File))
			if err != nil {
				beego.Error("stop read stdout err", err)
			}
			fmtstr := string(content) + "ERROR TIMEOUT STOP BY NEXT JOB"
			pid := cmdi.(*exec.Cmd).Process.Pid
			beego.Info("pid:", pid)
			savelog(jobid, fmtstr)
			err = syscall.Kill(-pid, syscall.SIGKILL)
			if err != nil {
				beego.Error("syscall.Kill err :", err)
			}
			cmdi.(*exec.Cmd).Process.Release()
			beego.Info("pid:", cmdi.(*exec.Cmd).Process.Pid)
			delete(CronJobs[jobid], "cmd")
		}
		kettleHomePath := beego.AppConfig.String("kettle::homepath")
		kettleWorkPath := beego.AppConfig.String("kettle::workpath")
		kettlejobpath := kettleWorkPath + jobfilepath
		beego.Info("start--job: ", jobfilepath)
		cmd := exec.Command(kettleHomePath+"kitchen.sh", "-file="+kettlejobpath)
		cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
		CronJobs[jobid]["cmd"] = cmd
		stdout, err := cmd.StdoutPipe()
		CronJobs[jobid]["stdout"] = stdout
		maplock.Unlock()
		err = cmd.Start()
		if err != nil {
			cmd.Wait()
			return
		}
		beego.Info("running--job: ", jobfilepath)
		cmd.Wait()
		maplock.Lock()
		delete(CronJobs[jobid], "cmd")
		delete(CronJobs[jobid], "stdout")
		maplock.Unlock()
		beego.Info("stoped--job: ", jobfilepath)
		content, err := ioutil.ReadAll(stdout)
		if err != nil {
			return
		}
		fmtstr := string(content)
		savelog(jobid, fmtstr)

		beego.Info("kettle job jobid....: ", jobid)
	})
	return
}

func savelog(jobid int, fmtstr string) {
	//	fmt.Println("fmtstr :", fmtstr)
	if strings.Contains(fmtstr, "ERROR") || strings.Contains(fmtstr, "error") || strings.Contains(fmtstr, "Error") {
		beego.Error("ExecJob fail :", fmtstr)
		kettlelog := models.KettleJobLog{
			KettleJobId: jobid,
			ErrorInfo:   fmtstr,
			Timestamp:   time.Now(),
			Status:      util.KETTLEJOB_RIGHT,
		}
		models.InsertKettleJobLog(kettlelog)
		return
	}
}

func StopCron(jobid int) (err error) {
	if jobc, ok := CronJobs[jobid]; ok {
		if cmd, ok := jobc["cmd"]; ok {
			err = cmd.(*exec.Cmd).Process.Kill()
			if err != nil {
				beego.Error("add cron cmd: ", err)
			}
		}
		jobc["cron"].(*cron.Cron).Stop()
		delete(CronJobs, jobid)
	}
	return
}
