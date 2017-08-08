package etlmgr

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/etl"
	"allsum_bi/util"
	base "common/lib/baseController"
	"common/lib/errcode"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ShowSycnList() {
	dbid := c.GetString("dbid")
	schema := c.GetString("schema")
	beego.Debug("schema ", dbid, schema)
	tableNames, err := db.ListSchemaTable(dbid, schema)
	if err != nil {
		beego.Error("ListSchemaTable err ", err)
		c.ReplyErr(errcode.ErrActionGetSchemaTable)
		return
	}
	var schemaTables []string
	for _, v := range tableNames {
		schema_v, _ := db.EncodeTableSchema(dbid, schema, v)
		schemaTables = append(schemaTables, schema_v)
	}
	sync, err := models.ListSyncInSourceTables(dbid, schemaTables)
	if err != nil {
		beego.Error("GetSyncBy err ", err)
		c.ReplyErr(errcode.ErrActionGetSycn)
		return
	}
	var res []map[string]interface{}
	for _, table := range tableNames {
		var tableMap map[string]interface{}
		schema_table, _ := db.EncodeTableSchema(dbid, schema, table)
		if _, ok := sync[schema_table]; !ok {
			tableMap = map[string]interface{}{
				"name":          table,
				"create_script": "",
				"alter_script":  "",
				"param_script":  "",
				"script":        "",
				"cron":          "",
				"documents":     "",
				"errorlimit":    0,
				"errornum":      0,
				"lasttime":      "",
				"desttable":     "",
				"sourcetable":   "",
				"owner":         "",
				"sync_uuid":     "",
				"status":        util.SYNC_NONE,
			}
		} else {
			syncuuid := sync[schema_table].Uuid
			syncid := sync[schema_table].Id
			errornum, err := models.CountSynchronousLogsBySyncid(syncid, util.SYNC_ENABLE)
			if err != nil {
				beego.Error("count synclog num err: ", err)
				errornum = 0
			}
			tableMap = map[string]interface{}{
				"name":          table,
				"status":        sync[schema_table].Status,
				"sync_uuid":     syncuuid,
				"owner":         sync[schema_table].Owner,
				"create_script": sync[schema_table].CreateScript,
				"alter_script":  sync[schema_table].AlterScript,
				"param_script":  sync[schema_table].ParamScript,
				"script":        sync[schema_table].Script,
				"sourcetable":   sync[schema_table].SourceTable,
				"desttable":     sync[schema_table].DestTable,
				"cron":          sync[schema_table].Cron,
				"documents":     sync[schema_table].Documents,
				"errorlimit":    sync[schema_table].ErrorLimit,
				"errornum":      errornum,
				"lasttime":      sync[schema_table].Lasttime,
			}
		}
		res = append(res, tableMap)
	}
	c.ReplySucc(res)
}

func (c *Controller) DataCalibration() {
	dbid := c.GetString("dbid")
	schema := c.GetString("schema")
	table := c.GetString("table")
	schema_table := schema + "." + table
	//	syncmap, err := models.ListSyncInSourceTables([]string{schema_table})
	//	if err != nil {
	//		beego.Error("ListSyncInSourceTables fail err :", err)
	//		c.ReplyErr(errcode.ErrServerError)
	//		return
	//	}
	checkres := db.CheckTableExist(util.BASEDB_CONNID, schema_table)
	beego.Debug("checkres", checkres)
	var err error
	if checkres {
		err = etl.DoEtlCalibration(dbid, schema, table)
	} else {
		err = etl.DoEtlWithoutTable(dbid, schema, table)
	}
	if err != nil {
		c.ReplyErr(errcode.ErrServerError)
		beego.Error("do etl err: ", err)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) SetEtl() {
	syncUuid := c.GetString("sync_uuid")
	createScript := c.GetString("create_script")
	alterScript := c.GetString("alter_script")
	paramScript := c.GetString("param_script")
	script := c.GetString("script")
	cron := c.GetString("cron")
	documents := c.GetString("documents")
	is_all_schema, err := c.GetBool("is_all_schema")
	if err != nil {
		is_all_schema = false
	}
	errorlimit, err := c.GetInt("errorlimit")
	if err != nil {
		errorlimit = 10
	}
	if syncUuid == "" || cron == "" || documents == "" {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("set etl some params is null ", syncUuid, script, cron, documents)
		return
	}
	setdata := map[string]interface{}{
		"sync_uuid":     syncUuid,
		"create_script": createScript,
		"alter_script":  alterScript,
		"param_script":  paramScript,
		"script":        script,
		"cron":          cron,
		"documents":     documents,
		"is_all_schema": is_all_schema,
		"error_limit":   errorlimit,
	}
	//no params check  ,so do once etl befor etl
	err = etl.SetAndDoEtl(setdata)
	if err != nil {
		//because this  action is set, so return params err
		beego.Error("set etl error : ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) StopEtl() {
	uuid := c.GetString("sync_uuid")
	err := etl.StopCronBySyncUuid(uuid)
	if err != nil {
		beego.Error("stop etl err", err)
		c.ReplyErr(err)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}
