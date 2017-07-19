package kettle

import (
	"allsum_bi/models"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	"io/ioutil"
	"os"
	"os/exec"
	"strings"

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
	cmd := exec.Command(kettleHomePath+"kitchen.sh", "-file="+kettlejobpath)
	fmtbytes, err := cmd.Output()
	if err != nil {
		fmtstr = string(fmtbytes)
		beego.Error("error format:", fmtstr)
		return fmtstr, err
	}
	return
}

func AddJobfile(name string, cron string, filename string, filedata []byte) (kettlejob models.KettleJob, err error) {
	kjobfile := uuid.NewV4().String() + "-" + filename
	urlpath, err := ossfile.PutFile("kettle", kjobfile, filedata)
	if err != nil {
		return
	}
	kjobfile_path := filename + "," + urlpath

	kettlejob = models.KettleJob{
		Name:    name,
		Cron:    cron,
		Kjbpath: kjobfile_path,
		Status:  util.KETTLEJOB_FAIL,
	}
	kettlejob, err = models.InsertKettleJob(kettlejob)
	return
}

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
		ktrpaths := strings.Split(kettlejob.Ktrpath, "_")
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
