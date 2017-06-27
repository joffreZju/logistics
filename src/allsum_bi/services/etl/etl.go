package etl

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/util"
	"fmt"
	_ "io/ioutil"
	"strings"
	"time"

	"github.com/astaxie/beego"
	_ "github.com/compose/transporter/adaptor/all"
	_ "github.com/compose/transporter/function/all"
)

func Start() {
	err := createTransformPath()
	if err != nil {
		beego.Error("create transporter error: ", err)
		return
	}
}

func TestETL() {

	sync := models.Synchronous{
		Owner:        "xxx",
		CreateScript: "create table XX.XX",
		Documents:    "文档",
		Status:       "building",
	}
	syncid, err := models.InsertSynchronous(sync)
	if err != nil {
		beego.Error("err:", err)
		return
	}
	sync.Id = syncid
	beego.Debug("sync.id", sync.Id)
	pipeline := NewPipeline()
	pipeline.Id = sync.Id
	err = pipeline.MakeSourceJs("earthsumuid")
	if err != nil {
		beego.Error(err)
		return
	}
	err = pipeline.MakeSinkJs()
	if err != nil {
		beego.Error(err)
		return
	}

	pipeline.MakeTransPortForm("", "", "")

	err = pipeline.MakeTransportJs("public.route_base", "public.route_base")
	if err != nil {
		beego.Error(err)
		return
	}

	pipeline.SetCron("0 * * * * *")
	pipeline.MakeFullJs()
	beego.Debug("fulljs:", pipeline.FullJs)
	runjs := pipeline.MakeRunJs()

	err = Save(&pipeline)
	if err != nil {
		beego.Error(err)
		return
	}
	beego.Debug("runjs:", runjs)
	err = AddCronWithScript(pipeline.Id, pipeline.Cron, runjs)
	if err != nil {
		beego.Error(err)
		return
	}
}

func DoETL(scriptbuff []byte) (err error) {
	transporter, err := newBuilder(scriptbuff)
	if err != nil {
		return
	}
	err = transporter.run()
	return
}

func DoEtlWithoutTable(dbid string, schema string, table string) (err error) {
	createsql, err := db.GetTableDesc(dbid, schema, table, schema, table)
	if err != nil {
		return
	}
	err = db.Exec(dbid, createsql)
	if err != nil {
		return
	}

	pipeline, err := callEtl(dbid, schema, table, "", "", "")
	if err != nil {
		db.DeleteTable(dbid, schema, table)
		return
	}
	createscript, err := pipeline.MakeFullCreateScript(createsql)
	if err != nil {
		return
	}
	script, err := pipeline.MakeDefaultTransPortScript(dbid)
	if err != nil {
		return
	}
	sync := models.Synchronous{
		Owner:        schema,
		CreateScript: createscript,
		SourceDbId:   dbid,
		SourceTable:  schema + "." + table,
		DestDbId:     util.BASEDB_CONNID,
		DestTable:    schema + "." + table,
		Script:       script,
		Status:       util.SYNC_BUILDING,
	}
	_, err = models.InsertSynchronous(sync)
	if err != nil {
		return
	}
	return
}

func DoEtlCalibration(dbid string, schema string, table string) (err error) {
	//TODO
	syncmap, err := models.ListSyncInSourceTables([]string{schema + "." + table})
	if err != nil {
		return err
	}
	sync, ok := syncmap[schema+"."+table]
	if !ok {
		return fmt.Errorf("no find sync task")
	}
	scriptMap, err := DecodeScript(sync.Script)
	if err != nil {
		return fmt.Errorf("sync db script err: ", err)
	}
	fields := strings.Split(scriptMap["paramfields"], ",")
	params, err := db.QueryToFields(util.BASEDB_CONNID, scriptMap["paramsql"], fields)
	_, err = callEtl(util.BASEDB_CONNID, sync.Owner, sync.SourceTable, scriptMap["transport"], scriptMap["transform"], params...)
	if err != nil {
		return
	}
	return
}

func callEtl(dbid string, schema string, table string, transportform string, transform string, params ...string) (pipeline Pipeline, err error) {
	pipeline = NewPipeline()
	err = pipeline.MakeSourceJs(dbid)
	if err != nil {
		beego.Error("make source js err : ", err)
		return
	}
	err = pipeline.MakeSinkJs()
	if err != nil {
		beego.Error("make sink js err : ", err)
		return
	}
	pipeline.MakeTransPortForm(transportform, transform, params...)

	err = pipeline.MakeTransportJs(schema+"."+table, schema+"."+table)
	if err != nil {
		beego.Error("make transporter js err ", err)
		return
	}
	pipeline.MakeFullJs()
	beego.Debug("fulljs: ", pipeline.FullJs)
	runjs := pipeline.MakeRunJs()

	err = DoETL([]byte(runjs))
	if err != nil {
		return
	}
	return
}

func SetAndDoEtl(setdata map[string]interface{}) (err error) {
	syncid := setdata["syncid"]
	script := setdata["script"]
	cron := setdata["cron"]
	documents := setdata["documents"]
	error_limit := setdata["error_limit"]
	sync, err := models.GetSynchronous(syncid.(int))
	if err != nil {
		return
	}
	scriptMap, err := DecodeScript(script.(string))
	fields := strings.Split(scriptMap["paramfields"], ",")
	params, err := db.QueryToFields(util.BASEDB_CONNID, scriptMap["paramsql"], fields)
	pipeline, err := callEtl(util.BASEDB_CONNID, sync.Owner, sync.SourceTable, scriptMap["transport"], scriptMap["transform"], params...)
	if err != nil {
		return
	}

	pipeline.SetCron(cron.(string))
	runjs := pipeline.MakeRunJs()

	err = AddCronWithScript(syncid.(int), pipeline.Cron, runjs)
	if err != nil {
		return
	}

	sync = models.Synchronous{
		Id:         syncid.(int),
		Script:     script.(string),
		Cron:       cron.(string),
		Documents:  documents.(string) + "\n Update @ time: " + time.Now().Format("2006-01-02 15:04:05"),
		ErrorLimit: error_limit.(int),
		Status:     util.SYNC_STARTED,
	}
	err = models.UpdateSynchronous(sync, "script", "cron", "documents", "error_limit", "status")
	return
}
