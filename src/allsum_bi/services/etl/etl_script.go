package etl

import (
	"allsum_bi/db"
	"allsum_bi/db/conn"
	"allsum_bi/services/util"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"os"

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

func MakeSinkJs() (sinkjs string, err error) {
	conninfo, err := conn.GetConninfo(util.BASEDB_CONNID)
	if err != nil {
		return
	}
	sinkjs = fmt.Sprintf(util.JS_TEMPLATE, "sink", conninfo.Host, conninfo.Port, conninfo.Dbname, conninfo.DbUser, conninfo.Passwd)
	return
}

func MakeSinkJsWithID(dbid string) (sinkjs string, err error) {
	conninfo, err := conn.GetConninfo(dbid)
	if err != nil {
		return
	}
	sinkjs = fmt.Sprintf(util.JS_TEMPLATE, "sink", conninfo.Host, conninfo.Port, conninfo.Dbname, conninfo.DbUser, conninfo.Passwd)
	return
}

func MakeTransportJs(dbid string, sourceschema string, sourcetable string, destschema string, desttable string, script string, params ...interface{}) (transportjs string, err error) {
	transportstr := ""
	if script != "" {
		transportstr = fmt.Sprintf(script, params...)
	}
	beego.Debug("transportstr: ", transportstr)
	sourceNameSpace, err := db.EncodeTableSchema(dbid, sourceschema, sourcetable)
	if err != nil {
		return
	}
	sinkNameSpace, err := db.EncodeTableSchema(dbid, destschema, sourcetable)
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
