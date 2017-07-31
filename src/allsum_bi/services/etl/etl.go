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
	"github.com/compose/transporter/pipeline"
)

var etltaskmap map[int]*pipeline.Pipeline

func init() {
	etltaskmap = map[int]*pipeline.Pipeline{}
}

func Start() {
	err := createTransformPath()
	if err != nil {
		beego.Error("create transporter error: ", err)
		return
	}
	StartEtlCron()
}

func DoETL(syncid int, scriptbuff []byte) (err error) {
	if p, ok := etltaskmap[syncid]; ok {
		p.Stop()
	}
	transporter, err := newBuilder(scriptbuff)
	if err != nil {
		return
	}
	defer func() {
		if r := recover(); r != nil {
			beego.Error("do etl crash ", r)
		}
	}()
	fmtstr, err := transporter.run(syncid)

	beego.Debug("etl fmt: ", fmtstr)
	if err != nil || strings.Contains(fmtstr, "ERROR") || strings.Contains(fmtstr, "Error") || strings.Contains(fmtstr, "error") {
		SetEtlError(syncid, err.Error()+"|"+fmtstr)
	} else {
		CleanEtlError(syncid)
	}
	return
}

func DoEtlWithoutTable(dbid string, schema string, table string) (err error) {
	sourceTable, err := db.EncodeTableSchema(dbid, schema, table)
	if err != nil {
		return
	}
	destTable, err := db.EncodeTableSchema(util.BASEDB_CONNID, schema, table)
	if err != nil {
		return
	}
	sync, err := models.GetSynchronousByTableName(dbid, sourceTable)
	var isinsert bool
	var createsql string
	if err == nil {
		isinsert = false
		createsql = strings.Replace(sync.CreateScript, util.SCRIPT_TABLE, sync.DestTable, -1)
		createsql = strings.Replace(createsql, util.SCRIPT_SCHEMA, schema, -1)
	} else {
		isinsert = true
		createsql, err = db.GetTableDesc(dbid, schema, table, schema, table)
		if err != nil {
			return err
		}
	}
	beego.Debug("create sql:", createsql)
	err = db.Exec(util.BASEDB_CONNID, createsql)
	if err != nil {
		return
	}

	sync_res := models.Synchronous{
		Owner:        schema,
		CreateScript: createsql,
		AlterScript:  "",
		ParamScript:  "",
		Script:       fmt.Sprintf(util.TRANSPORTFORM_GOJA, util.DEFAULT_TRANSPORT),
		SourceDbId:   dbid,
		SourceTable:  sourceTable,
		DestDbId:     util.BASEDB_CONNID,
		DestTable:    destTable,
		Status:       util.SYNC_BUILDING,
		Lasttime:     time.Now(),
	}
	fmt.Println("insert res: ", sync_res)
	if isinsert {
		syncid, err := models.InsertSynchronous(sync_res)
		if err != nil {
			db.DeleteTable(util.BASEDB_CONNID, schema, table)
			return err
		}
		sync_res.Id = syncid
	} else {
		sync_res.Id = sync.Id
		err = models.UpdateSynchronous(sync_res, "owner", "create_script", "alter_script", "script", "source_db_id", "source_table", "dest_db_id", "dest_table", "script", "status", "lasttime")
		if err != nil {
			return
		}
	}

	go func() {
		err = callEtl(sync_res.Id, dbid, sourceTable, destTable, "", "")
		if err != nil {
			beego.Error("call etl", err)
			db.DeleteTable(util.BASEDB_CONNID, schema, table)
			return
		}
	}()

	return
}

func DoEtlCalibration(dbid string, schema string, table string) (err error) {
	sourceTable, err := db.EncodeTableSchema(dbid, schema, table)
	if err != nil {
		return
	}
	sync, err := models.GetSynchronousByTableName(dbid, sourceTable)
	params := [][]interface{}{nil}
	if sync.ParamScript == "" {
		sync.Script = ""
	} else {
		sqlstr := strings.Replace(sync.ParamScript, util.SCRIPT_TABLE, sync.DestTable, -1)
		sqlstr = strings.Replace(sqlstr, util.SCRIPT_SCHEMA, schema, -1)
		params, err = db.QueryDatas(util.BASEDB_CONNID, sqlstr)
		if err != nil {
			return err
		}
	}
	go func() {
		err = callEtl(sync.Id, dbid, sourceTable, sync.DestTable, sync.Script, params[0])
		if err != nil {
			return
		}
	}()
	sync.Lasttime = time.Now()
	err = models.UpdateSynchronous(sync, "lasttime")
	if err != nil {
		beego.Error("update Lasttime error", err)
	}
	return
}

func StartEtl(sync models.Synchronous) (err error) {
	params := [][]interface{}{nil}
	if sync.ParamScript == "" {
		sync.Script = ""
	} else {
		sqlstr := strings.Replace(sync.ParamScript, util.SCRIPT_TABLE, sync.SourceTable, -1)
		sqlstr = strings.Replace(sqlstr, util.SCRIPT_SCHEMA, sync.Owner, -1)
		params, err = db.QueryDatas(util.BASEDB_CONNID, sqlstr)
		if err != nil {
			return err
		}
	}

	//
	//TODO makerunjs
	runjs, err := buildEtl(sync.SourceDbId, sync.SourceTable, sync.DestTable, sync.Script, params[0]...)
	if err != nil {
		return
	}
	err = AddCronWithScript(sync.Id, sync.Cron, runjs)
	if err != nil {
		return
	}
	return
}

func callEtl(syncid int, dbid string, SourceTable string, DestTable string, script string, params ...interface{}) (err error) {
	//	pipeline = NewPipeline()
	sourcejs, err := MakeSourceJs(dbid)
	if err != nil {
		beego.Error("make source js err : ", err)
		return
	}
	sinkjs, err := MakeSinkJs()
	if err != nil {
		beego.Error("make sink js err : ", err)
		return
	}
	//	pipeline.MakeTransPortForm(transportform, transform, params...)
	//TODO
	sourceSchema, sourceTable, err := db.DecodeTableSchema(dbid, SourceTable)
	if err != nil {
		return
	}
	destSchema, destTable, err := db.DecodeTableSchema(util.BASEDB_CONNID, DestTable)
	if err != nil {
		return
	}
	transportjs, err := MakeTransportJs(dbid, sourceSchema, sourceTable, destSchema, destTable, script, params...)
	if err != nil {
		beego.Error("make transporter js err ", err)
		return
	}
	//	pipeline.MakeFullJs()
	//	beego.Debug("fulljs: ", pipeline.FullJs)
	runjs := MakeRunJs(sourcejs, sinkjs, transportjs)
	beego.Debug("runjs:", runjs)
	err = DoETL(syncid, []byte(runjs))
	if err != nil {
		return
	}
	return
}

func buildEtl(dbid string, SourceTable string, DestTable string, script string, params ...interface{}) (runjs string, err error) {
	//	pipeline = NewPipeline()
	sourcejs, err := MakeSourceJs(dbid)
	if err != nil {
		beego.Error("make source js err : ", err)
		return
	}
	sinkjs, err := MakeSinkJs()
	if err != nil {
		beego.Error("make sink js err : ", err)
		return
	}
	//	pipeline.MakeTransPortForm(transportform, transform, params...)
	sourceSchema, sourceTable, err := db.DecodeTableSchema(dbid, SourceTable)
	if err != nil {
		return
	}
	destSchema, destTable, err := db.DecodeTableSchema(util.BASEDB_CONNID, DestTable)
	if err != nil {
		return
	}
	transportjs, err := MakeTransportJs(dbid, sourceSchema, sourceTable, destSchema, destTable, script, params...)
	if err != nil {
		beego.Error("make transporter js err ", err)
		return
	}
	runjs = MakeRunJs(sourcejs, sinkjs, transportjs)
	return
}

func SetAndDoEtl(setdata map[string]interface{}) (err error) {
	syncid := setdata["sync_uuid"].(string)
	script := setdata["script"].(string)
	alter_script := setdata["alter_script"].(string)
	cron := setdata["cron"].(string)
	documents := setdata["documents"].(string)
	error_limit := setdata["error_limit"]
	param_script := setdata["param_script"].(string)

	sync, err := models.GetSynchronousByUuid(syncid)
	if err != nil {
		return
	}

	sourceSchema, sourceTable, err := db.DecodeTableSchema(sync.SourceDbId, sync.SourceTable)
	if err != nil {
		return
	}
	destSchema, destTable, err := db.DecodeTableSchema(util.BASEDB_CONNID, sync.DestTable)
	if err != nil {
		return
	}

	//alter table
	if alter_script != "" {
		alter_script = strings.Replace(alter_script, util.SCRIPT_TABLE, sync.DestTable, -1)
		alter_script = strings.Replace(alter_script, util.SCRIPT_SCHEMA, destSchema, -1)

		err = db.Exec(util.BASEDB_CONNID, alter_script)
		if err != nil {
			beego.Error("alter table err :", err)
			return fmt.Errorf("alter table error")
		}
		sync.CreateScript, err = db.GetTableDesc(util.BASEDB_CONNID, sourceSchema, sourceTable, destSchema, destTable)
		if err != nil {
			beego.Error("get new create sql err:", err)
			return fmt.Errorf("get new create sql error")
		}
		sync.CreateScript = strings.Replace(sync.CreateScript, sync.DestTable, util.SCRIPT_TABLE, -1)
		sync.AlterScript = ""

	}
	params := [][]interface{}{nil}
	if param_script == "" {
		script = ""
	} else {
		sqlstr := strings.Replace(param_script, util.SCRIPT_TABLE, sync.DestTable, -1)
		sqlstr = strings.Replace(sqlstr, util.SCRIPT_SCHEMA, sync.Owner, -1)
		params, err = db.QueryDatas(util.BASEDB_CONNID, sqlstr)
		if err != nil {
			return err
		}
	}
	err = callEtl(sync.Id, sync.SourceDbId, sync.SourceTable, sync.DestTable, script, params[0])
	if err != nil {
		return
	}

	//	pipeline.SetCron(cron.(string))
	//TODO runjs
	runjs, err := buildEtl(sync.SourceDbId, sync.SourceTable, sync.DestTable, script, params[0])
	if err != nil {
		return
	}

	err = AddCronWithScript(sync.Id, cron, runjs)
	if err != nil {
		return
	}
	sync.Script = script
	sync.Cron = cron
	sync.Documents = documents + "\n Update @ time: " + time.Now().Format("2006-01-02 15:04:05")
	sync.ErrorLimit = error_limit.(int)
	sync.Lasttime = time.Now()
	sync.Status = util.SYNC_STARTED

	err = models.UpdateSynchronous(sync, "create_script", "alter_script", "script", "cron", "documents", "error_limit", "status")
	return
}
