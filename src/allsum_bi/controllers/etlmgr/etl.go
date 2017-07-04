package etlmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/etl"
	"allsum_bi/util"
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
		schemaTables = append(schemaTables, schema+"."+v)
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
		schema_table := schema + "." + table
		if _, ok := sync[schema_table]; !ok {
			tableMap = map[string]interface{}{
				"name":   table,
				"status": util.SYNC_NONE,
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
				"name":        table,
				"status":      sync[schema_table].Status,
				"syncuuid":    syncuuid,
				"owner":       sync[schema_table].Owner,
				"script":      sync[schema_table].Script,
				"sourcetable": sync[schema_table].SourceTable,
				"desttable":   sync[schema_table].DestTable,
				"cron":        sync[schema_table].Cron,
				"documents":   sync[schema_table].Documents,
				"errorlimit":  sync[schema_table].ErrorLimit,
				"errornum":    errornum,
				"lasttime":    sync[schema_table].Lasttime,
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
	script := c.GetString("script")
	cron := c.GetString("cron")
	documents := c.GetString("documents")
	errorlimit, err := c.GetInt("error_limit")
	if err != nil {
		errorlimit = 10
	}
	if syncUuid == "" || script == "" || cron == "" || documents == "" {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("set etl some params is null ")
	}
	setdata := map[string]interface{}{
		"sync_uuid":   syncUuid,
		"script":      script,
		"cron":        cron,
		"documents":   documents,
		"error_limit": errorlimit,
	}
	//no params check  ,so do once etl befor etl
	err = etl.SetAndDoEtl(setdata)
	if err != nil {
		//because this  action is set, so return params err
		c.ReplyErr(errcode.ErrServerError)
		beego.Error("set etl error : ", err)
		return
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) StopEtl() {
	uuid := c.GetString("uuid")
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
