package kettle

import (
	"allsum_bi/models"
	"allsum_bi/util"
	"allsum_bi/util/ossfile"
	"encoding/json"
	"io/ioutil"
	"os"
	"os/exec"
	"path"
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
	kjobfilemap := map[string]string{
		"filename": filename,
		"urlpath":  urlpath,
	}
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	err = ioutil.WriteFile(kettleWorkPath+path.Base(urlpath), filedata, 0664)
	if err != nil {
		return
	}
	kjobfilejson, err := json.Marshal(kjobfilemap)
	if err != nil {
		return
	}
	kettlejob = models.KettleJob{
		Name:    name,
		Cron:    cron,
		Kjbpath: string(kjobfilejson),
		Status:  util.KETTLEJOB_FAIL,
	}
	kettlejob, err = models.InsertKettleJob(kettlejob)
	return
}

func AddKtrfile(uuid string, filename string, filedata []byte) (err error) {
	kettlejob, err := models.GetKettleJobByUuid(uuid)
	if err != nil {
		return
	}
	var kjbmap map[string]string
	err = json.Unmarshal([]byte(kettlejob.Kjbpath), &kjbmap)
	if err != nil {
		return
	}
	kettleWorkPath := beego.AppConfig.String("kettle::workpath")
	kjbfiledata, err := ioutil.ReadFile(kettleWorkPath + path.Base(kjbmap["urlpath"]))
	if err != nil {
		return
	}
	if !strings.Contains(string(kjbfiledata), filename) {
		return
	}
	//	json = x2j.XmlToJson(kjbfiledata)
	//	uuidktrname := uuid.NewV4.String() + "_" + filename
	//	var ktrmap map[string]string
	//	err = json.Unmarshal([]byte(kettlejob.Ktrpath), &ktrmap)
	//	if err != nil {
	//		return
	//	}
	//	ktrmap[filename]
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
