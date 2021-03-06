package kettlesvs

import (
	"allsum_bi/models"
	"allsum_bi/services/util"
	"allsum_bi/services/util/ossfile"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
	"strings"
	"syscall"
	"time"

	"github.com/astaxie/beego"
	uuid "github.com/satori/go.uuid"
)

func Start() {
	InitJob()
	StartJobCron()
	go ReloadJobPath()
}

func InitJob() {
	kettleHomePath := beego.AppConfig.String("kettle::homepath")
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	os.MkdirAll(kettleWorkPath, 0771)
	os.MkdirAll(kettleHomePath, 0771)
}

func ExecJob(jobfile string) (fmtstr string, err error) {
	kettleHomePath := beego.AppConfig.String("kettle::homepath")
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	kettlejobpath := kettleWorkPath + jobfile
	beego.Info("start--job: ", jobfile)
	cmd := exec.Command(kettleHomePath+"kitchen.sh", "-file="+kettlejobpath)
	cmd.SysProcAttr = &syscall.SysProcAttr{Setsid: true}
	stdout, err := cmd.StdoutPipe()
	err = cmd.Start()
	if err != nil {
		return
	}

	beego.Info("running--job: ", jobfile)
	signal := make(chan int)
	go func() {
		err = cmd.Wait()
		content, err := ioutil.ReadAll(stdout)
		if err != nil {
			fmt.Println(err)
		}
		fmtstr = string(content)
		signal <- 1
	}()
	select {
	case <-signal:
		beego.Info("stoped--job: ", jobfile)
	case <-time.After(2 * time.Minute):
		beego.Info("time out two minutes jobfile:", jobfile)
		err = syscall.Kill(-cmd.Process.Pid, syscall.SIGKILL)
		if err != nil {
			beego.Error("syscall.Kill err: ", err)
		}
		cmd.Process.Release()
		fmtstr = jobfile + "ERROR: TIMEOUT"
	}
	return
}

func AddJobKtrfile(userid int, name string, cron string, filename string, jobfiledata []byte, jobktrs []map[string]string) (kettlejob models.KettleJob, err error) {
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	jobfiledatastr := string(jobfiledata)
	for _, jobktrmap := range jobktrs {
		basektrpath := path.Base(jobktrmap["uri"])
		jobfiledatastr = strings.Replace(jobfiledatastr, jobktrmap["name"], basektrpath, -1)
	}
	kjobfile := fmt.Sprintf("%d_%v_%s", userid, time.Now(), filename)
	jobfiledata = []byte(jobfiledatastr)
	urlpath, err := ossfile.PutFile("kettle", kjobfile, jobfiledata)
	if err != nil {
		return
	}
	kjobfilemap := map[string]string{
		"filename": filename,
		"urlpath":  urlpath,
	}
	err = ioutil.WriteFile(kettleWorkPath+kjobfile, jobfiledata, 0664)
	if err != nil {
		return
	}
	kjobfilejson, err := json.Marshal(kjobfilemap)
	if err != nil {
		return
	}
	ktrfilejson, err := json.Marshal(jobktrs)
	if err != nil {
		return
	}

	kettlejob = models.KettleJob{
		Name:     name,
		Cron:     cron,
		Kjbpath:  string(kjobfilejson),
		Ktrpaths: string(ktrfilejson),
		Status:   util.KETTLEJOB_FAIL,
	}
	kettlejob, err = models.InsertKettleJob(kettlejob)
	return
}

func AddJobKtrfile_OLD(name string, cron string, filename string, jobfiledata []byte, ktrdatas map[string][]byte) (kettlejob models.KettleJob, err error) {

	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	jobfiledatastr := string(jobfiledata)
	ktrfilemap := map[string]string{}
	for k, v := range ktrdatas {
		uuidktrfile := uuid.NewV4().String() + k
		ktrurlpath, err := ossfile.PutFile("kettle", uuidktrfile, v)
		if err != nil {
			return kettlejob, err
		}
		err = ioutil.WriteFile(kettleWorkPath+uuidktrfile, v, 0664)
		if err != nil {
			return kettlejob, err
		}
		ktrfilemap[k] = ktrurlpath
		jobfiledatastr = strings.Replace(jobfiledatastr, k, uuidktrfile, -1)
	}

	kjobfile := uuid.NewV4().String() + "-" + filename
	jobfiledata = []byte(jobfiledatastr)
	urlpath, err := ossfile.PutFile("kettle", kjobfile, jobfiledata)
	if err != nil {
		return
	}
	kjobfilemap := map[string]string{
		"filename": filename,
		"urlpath":  urlpath,
	}
	err = ioutil.WriteFile(kettleWorkPath+kjobfile, jobfiledata, 0664)
	if err != nil {
		return
	}
	kjobfilejson, err := json.Marshal(kjobfilemap)
	if err != nil {
		return
	}
	ktrfilejson, err := json.Marshal(ktrfilemap)
	if err != nil {
		return
	}

	kettlejob = models.KettleJob{
		Name:     name,
		Cron:     cron,
		Kjbpath:  string(kjobfilejson),
		Ktrpaths: string(ktrfilejson),
		Status:   util.KETTLEJOB_FAIL,
	}
	kettlejob, err = models.InsertKettleJob(kettlejob)
	return
}

//func AddKtrfile(uuid string, filename string, filedata []byte) (err error) {
//	kettlejob, err := models.GetKettleJobByUuid(uuid)
//	if err != nil {
//		return
//	}
//	var kjbmap map[string]string
//	err = json.Unmarshal([]byte(kettlejob.Kjbpath), &kjbmap)
//	if err != nil {
//		return
//	}
//	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
//	kjbfiledata, err := ioutil.ReadFile(kettleWorkPath + path.Base(kjbmap["urlpath"]))
//	if err != nil {
//		return
//	}
//	if !strings.Contains(string(kjbfiledata), filename) {
//		return fmt.Error("jobfile have not this ktr")
//	}
//	uuidktrname := uuid.NewV4.String() + "_" + filename
//	urlpath, err := ossfile.PutFile("kettle", kjobfile, filedata)
//	if err != nil {
//		return
//	}
//
//	var ktrmap map[string]string
//	err = json.Unmarshal([]byte(kettlejob.Ktrpath), &ktrmap)
//	if err != nil {
//		return
//	}
//	ktrmap[filename] = urlpath
//
//	kjbfiledatastr := string(kjbfiledata)
//	for k, v := range ktrmap {
//		kjbfiledatastr := strings.Replace(kjbfiledatastr, filename, path.Base(v), -1)
//	}
//
//	return
//}

//need in go func
func ReloadJobPath() (err error) {
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	//clean path
	err = os.RemoveAll(kettleWorkPath + "*")
	if err != nil {
		return
	}
	kettlejobs, err := models.ListKettleJobByField([]string{"status"}, []interface{}{util.KETTLEJOB_RIGHT}, 0, 0)
	if err != nil {
		return
	}
	for _, kettlejob := range kettlejobs {
		filedata, err := ossfile.GetFile(kettlejob.Kjbpath)
		if err != nil {
			return err
		}
		err = ioutil.WriteFile(kettleWorkPath+kettlejob.Kjbpath, filedata, 0644)
		if err != nil {
			return err
		}
		ktrpaths := strings.Split(kettlejob.Ktrpaths, "_")
		for _, ktrpath := range ktrpaths {
			filedata, err := ossfile.GetFile(ktrpath)
			if err != nil {
				return err
			}
			err = ioutil.WriteFile(kettleWorkPath+ktrpath, filedata, 0644)
			if err != nil {
				return err
			}
		}
	}
	return
}
