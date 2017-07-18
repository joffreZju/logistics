package etl

import (
	"allsum_bi/db"
	"allsum_bi/db/conn"
	"allsum_bi/models"
	"allsum_bi/util"
	"encoding/json"
	"errors"
	"fmt"
	"io/ioutil"
	"os"
	"strings"

	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

type Pipeline struct {
	Id            int
	SourceTable   string
	DestTable     string
	SourceJs      string
	SinkJs        string
	TransformJs   string
	TransformPath string
	TransPortJs   string
	FullJs        string
	TransPortForm string
	Cron          string
}

func Save(pipeline *Pipeline) (err error) {

	sync := models.Synchronous{
		Id:     pipeline.Id,
		Script: pipeline.FullJs,
		Cron:   pipeline.Cron,
		Status: util.SYNC_STARTED,
	}
	err = models.UpdateSynchronous(sync, "script", "cron", "status")
	return
}

func createTransformPath() (err error) {
	if _, err := os.Stat(util.TRANSFORM_PATH); os.IsNotExist(err) {
		err = os.Mkdir(util.TRANSFORM_PATH, os.ModePerm)
		if err != nil {
			return err
		}
		//		if err := os.Chmod(TRANSFORM_PATH, 0644); err != nil {
		//			return err
		//		}
	}
	return
}

func MakeRunScript(fullscript string) (runjs string, err error) {
	//	allJsScript := strings.Split(fullscript, "#")
	var fulljsMap map[string]string
	err = json.Unmarshal([]byte(fullscript), &fullscript)
	if err != nil {
		return
	}
	runjs = fulljsMap["sourcejs"] + "\n" + fulljsMap["sinkjs"] + "\n" + fulljsMap["transportjs"] + "\n"
	return
}

func MakeSourceJs(id string) (sourcejs string, err error) {
	conninfo, err := conn.GetConninfo(id)
	if err != nil {
		return
	}
	//TODO  check DBTYPE
	sourcejs = fmt.Sprintf(util.JS_TEMPLATE, "source", conninfo.Host, conninfo.Port, conninfo.Dbname, conninfo.DbUser, conninfo.Passwd)
	return
}

func MakeRunJs(sourcejs string, sinkjs string, transportjs string) (runjs string) {
	//	allJsScript := strings.Split(p.FullJs, "#")
	runjs = sourcejs + "\n" + sinkjs + "\n" + transportjs + "\n"
	return
}

//
//func (p *Pipeline) SetCron(cron string) {
//	p.Cron = cron
//	return
//}
//
func MakeSinkJs() (sinkjs string, err error) {
	conninfo, err := conn.GetConninfo(util.BASEDB_CONNID)
	if err != nil {
		return
	}
	sinkjs = fmt.Sprintf(util.JS_TEMPLATE, "sink", conninfo.Host, conninfo.Port, conninfo.Dbname, conninfo.DbUser, conninfo.Passwd)
	return
}

//func (p *Pipeline) MakeTransPortForm(transportform string, transform string, params ...interface{}) {
//	if transform == "" {
//		p.TransformJs = ""
//		p.TransPortForm = transportform
//	} else {
//		p.MakeTransformJs([]byte(transform), params...)
//		p.TransPortForm = fmt.Sprintf(transportform, p.TransformPath)
//	}
//}
//
func MakeTransportJs(dbid string, sourceschema string, sourcetable string, destschema string, desttable string, script string, params ...interface{}) (transportjs string, err error) {
	transportstr := ""
	if script != "" {
		transportstr = fmt.Sprintf(script, params...)
	}
	sourceNameSpace, err := db.EncodeTableSchema(dbid, sourceschema, sourcetable)
	if err != nil {
		return
	}
	sinkNameSpace, err := db.EncodeTableSchema(dbid, sourceschema, sourcetable)
	if err != nil {
		return
	}

	transportjs = fmt.Sprintf(util.JS_TRANSPORT, sourceNameSpace, transportstr, sinkNameSpace)
	return
}

func (p *Pipeline) MakeTransformJs(js []byte, params ...interface{}) (err error) {
	p.TransformJs = fmt.Sprintf(string(js), params)
	uid := uuid.NewV4()

	p.TransformPath = util.TRANSFORM_PATH + uid.String()

	err = ioutil.WriteFile(p.TransformPath, js, os.ModePerm)
	if err != nil {
		return
	}
	return
}

func (p *Pipeline) MakeFullJs() (err error) {
	fulljsMap := map[string]string{
		"sourcejs":      p.SourceJs,
		"sinkjs":        p.SinkJs,
		"transformjs":   p.TransformJs,
		"transformpath": p.TransformPath,
		"transportjs":   p.TransPortJs,
	}
	fulljsJson, err := json.Marshal(fulljsMap)
	if err != nil {
		return
	}
	p.FullJs = string(fulljsJson)
	return
}

func (p *Pipeline) MakeFullCreateScript(createsql string) (createscript string, err error) {
	createScriptMap := map[string]string{
		"createsql": createsql,
		"etlscript": p.FullJs,
	}
	createscriptbytes, err := json.Marshal(createScriptMap)
	createscript = string(createscriptbytes)
	return
}

func (p *Pipeline) MakeDefaultTransPortScript(dbid string) (script string, err error) {
	schema_table := strings.Split(p.DestTable, ".")
	fields, err := db.GetTableFields(dbid, schema_table[0], schema_table[1])
	beego.Debug("fields: ", fields)
	sqlstr := fmt.Sprintf(util.DEFAULT_PARAMS_SQL, fields[0], "maxnum", p.DestTable)
	fieldvalues, err := db.QueryDatas(dbid, sqlstr)
	if err != nil {
		return
	}
	beego.Debug("fieldvalues", fieldvalues)
	scriptMap := map[string]string{
		"transport":   fmt.Sprintf(util.TRANSPORTFORM_GOJA, fmt.Sprintf(util.DEFAULT_TRANSPORT, fields[0], fieldvalues[0][0])),
		"transform":   "",
		"paramsql":    util.DEFAULT_PARAMS_SQL,
		"paramfields": strings.Join([]string{fields[0], fields[0], p.DestTable}, ","),
	}
	script, err = EncodeScript(scriptMap)
	if err != nil {
		return
	}
	return
}

func UnMakeCreateScript(createscript string) (createsql string, js_script string, err error) {
	var script map[string]string
	err = json.Unmarshal([]byte(createscript), &script)
	if err != nil {
		return
	}
	createsql = script["createsql"]
	js_script = script["etlscript"]
	return
}

func DecodeScript(script string) (scriptMap map[string]string, err error) {
	err = json.Unmarshal([]byte(script), &scriptMap)
	if err != nil {
		return
	}
	_, ok1 := scriptMap["transport"]
	_, ok2 := scriptMap["transform"]
	_, ok3 := scriptMap["paramsql"]
	_, ok4 := scriptMap["paramfields"]
	if !(ok1 && ok2 && ok3 && ok4) {
		beego.Error("1 2 3 4", ok1, ok2, ok3, ok4)
		err = errors.New("miss script")
		return
	}
	return
}

func EncodeScript(scriptMap map[string]string) (script string, err error) {
	scriptbytes, err := json.Marshal(scriptMap)
	if err != nil {
		return
	}
	script = string(scriptbytes)
	return
}
