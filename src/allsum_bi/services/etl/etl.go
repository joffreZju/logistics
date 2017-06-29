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
	StartEtlCron()
}

//func TestETL() {
//
//	sync := models.Synchronous{
//		Owner:        "xxx",
//		CreateScript: "create table XX.XX",
//		Documents:    "文档",
//		Status:       "building",
//	}
//	syncid, err := models.InsertSynchronous(sync)
//	if err != nil {
//		beego.Error("err:", err)
//		return
//	}
//	sync.Id = syncid
//	beego.Debug("sync.id", sync.Id)
//	pipeline := NewPipeline()
//	pipeline.Id = sync.Id
//	err = pipeline.MakeSourceJs("earthsumuid")
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//	err = pipeline.MakeSinkJs()
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//
//	pipeline.MakeTransPortForm("", "", "")
//
//	err = pipeline.MakeTransportJs("public.route_base", "public.route_base")
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//
//	pipeline.SetCron("0 * * * * *")
//	pipeline.MakeFullJs()
//	beego.Debug("fulljs:", pipeline.FullJs)
//	runjs := pipeline.MakeRunJs()
//
//	err = Save(&pipeline)
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//	beego.Debug("runjs:", runjs)
//	err = AddCronWithScript(pipeline.Id, pipeline.Cron, runjs)
//	if err != nil {
//		beego.Error(err)
//		return
//	}
//}
//
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
	beego.Debug("create sql:", createsql)
	err = db.Exec(util.BASEDB_CONNID, createsql)
	if err != nil {
		return
	}
	var pipeline Pipeline
	calletl := make(chan string)
	defer close(calletl)
	go func() {
		pipeline, err = callEtl(dbid, schema, table, "", "", "")
		if err != nil {
			beego.Error("call etl", err)
			db.DeleteTable(util.BASEDB_CONNID, schema, table)
			return
		}
		calletl <- "end"
	}()
	select {
	case <-calletl:
		beego.Debug("etl done")
	case <-time.After(time.Minute * 3):
		beego.Error("timeout")
		return
	}
	createscript, err := pipeline.MakeFullCreateScript(createsql)
	if err != nil {
		db.DeleteTable(util.BASEDB_CONNID, schema, table)
		return
	}
	script, err := pipeline.MakeDefaultTransPortScript(util.BASEDB_CONNID)
	if err != nil {
		db.DeleteTable(util.BASEDB_CONNID, schema, table)
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
		Lasttime:     time.Now(),
	}
	_, err = models.InsertSynchronous(sync)
	if err != nil {
		db.DeleteTable(util.BASEDB_CONNID, schema, table)
		return
	}
	return
}

func DoEtlCalibration(dbid string, schema string, table string) (err error) {
	//TODO
	syncmap, err := models.ListSyncInSourceTables(dbid, []string{schema + "." + table})
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
	fieldinterfaces := []interface{}{}
	for _, field := range fields {
		fieldinterfaces = append(fieldinterfaces, field)
	}
	sqlstr := fmt.Sprintf(scriptMap["paramsql"], fieldinterfaces...)
	params, err := db.QueryDatas(util.BASEDB_CONNID, sqlstr)
	_, err = callEtl(dbid, sync.Owner, sync.SourceTable, scriptMap["transport"], scriptMap["transform"], params[0])
	if err != nil {
		return
	}
	sync.Lasttime = time.Now()
	err = models.UpdateSynchronous(sync, "lasttime")
	if err != nil {
		beego.Error("update Lasttime error", err)
	}
	return
}

func StartEtl(uuid string) (err error) {
	sync, err := models.GetSynchronousByUuid(uuid)
	if err != nil {
		return
	}
	scriptMap, err := DecodeScript(sync.Script)
	if err != nil {
		return fmt.Errorf("sync db script err: ", err)
	}
	fields := strings.Split(scriptMap["paramfields"], ",")
	fieldinterfaces := []interface{}{}
	for _, field := range fields {
		fieldinterfaces = append(fieldinterfaces, field)
	}
	sqlstr := fmt.Sprintf(scriptMap["paramsql"], fieldinterfaces...)
	params, err := db.QueryDatas(util.BASEDB_CONNID, sqlstr)
	pipeline, err := buildEtl(sync.SourceDbId, sync.Owner, sync.SourceTable, scriptMap["transport"], scriptMap["transform"], params[0])
	pipeline.SetCron(sync.Cron)
	runjs := pipeline.MakeRunJs()

	err = AddCronWithScript(sync.Id, pipeline.Cron, runjs)
	if err != nil {
		return
	}
	return
}

func callEtl(dbid string, schema string, table string, transportform string, transform string, params ...interface{}) (pipeline Pipeline, err error) {
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
	beego.Debug("runjs:", runjs)
	err = DoETL([]byte(runjs))
	if err != nil {
		return
	}
	return
}

func buildEtl(dbid string, schema string, table string, transportform string, transform string, params ...interface{}) (pipeline Pipeline, err error) {
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
	//	runjs := pipeline.MakeRunJs()
	//	beego.Debug("runjs:", runjs)
	//	err = DoETL([]byte(runjs))
	//	if err != nil {
	//		return
	//	}
	return
}

func SetAndDoEtl(setdata map[string]interface{}) (err error) {
	syncid := setdata["sync_uuid"]
	script := setdata["script"]
	cron := setdata["cron"]
	documents := setdata["documents"]
	error_limit := setdata["error_limit"]
	sync, err := models.GetSynchronousByUuid(syncid.(string))
	if err != nil {
		return
	}
	scriptMap, err := DecodeScript(script.(string))
	fields := strings.Split(scriptMap["paramfields"], ",")
	var fieldinterfaces []interface{}
	for _, v := range fields {
		fieldinterfaces = append(fieldinterfaces, v)
	}
	sqlstr := fmt.Sprintf(scriptMap["paramsql"], fieldinterfaces...)
	params, err := db.QueryDatas(util.BASEDB_CONNID, sqlstr)
	pipeline, err := callEtl(util.BASEDB_CONNID, sync.Owner, sync.SourceTable, scriptMap["transport"], scriptMap["transform"], params[0])
	if err != nil {
		return
	}

	pipeline.SetCron(cron.(string))
	runjs := pipeline.MakeRunJs()

	err = AddCronWithScript(sync.Id, pipeline.Cron, runjs)
	if err != nil {
		return
	}

	sync = models.Synchronous{
		Id:         sync.Id,
		Script:     script.(string),
		Cron:       cron.(string),
		Documents:  documents.(string) + "\n Update @ time: " + time.Now().Format("2006-01-02 15:04:05"),
		ErrorLimit: error_limit.(int),
		Lasttime:   time.Now(),
		Status:     util.SYNC_STARTED,
	}
	err = models.UpdateSynchronous(sync, "script", "cron", "documents", "error_limit", "status")
	return
}
