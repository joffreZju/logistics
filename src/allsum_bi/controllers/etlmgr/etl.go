package etlmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/etl"
	"allsum_bi/util"
	"allsum_bi/util/errcode"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ShowSycnList() {
	dbid := c.GetString("dbid")
	schema := c.GetString("schema")
	tableNames, err := db.ListSchemaTable(dbid, schema)
	if err != nil {
		beego.Error("ListSchemaTable err ", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	sync, err := models.ListSyncInSourceTables(tableNames)
	if err != nil {
		beego.Error("GetSyncBy err ", err)
		c.ReplyErr(errcode.ErrServerError)
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
			syncid := sync[schema_table].Id
			errornum, err := models.CountSynchronousLogsBySyncid(syncid, util.SYNC_ENABLE)
			if err != nil {
				beego.Error("count synclog num err: ", err)
				errornum = 0
			}
			tableMap = map[string]interface{}{
				"name":        table,
				"status":      sync[schema_table].Status,
				"syncid":      syncid,
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
	syncmap, err := models.ListSyncInSourceTables([]string{schema_table})
	if err != nil {
		beego.Error("ListSyncInSourceTables fail err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}

	if _, ok := syncmap[schema_table]; ok {
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
	syncid := c.GetString("syncid")
	script := c.GetString("script")
	cron := c.GetString("cron")
	documents := c.GetString("documents")
	errorlimit, err := c.GetInt("error_limit")
	if err != nil {
		errorlimit = 10
	}
	if syncid == "" || script == "" || cron == "" || documents == "" {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("set etl some params is null ")
	}
	setdata := map[string]interface{}{
		"syncid":      syncid,
		"script":      script,
		"cron":        cron,
		"documents":   documents,
		"error_limit": errorlimit,
	}
	//no params check  ,so do once etl befor etl
	err = etl.SetAndDoEtl(setdata)
	if err != nil {
		//because this  action is set, so return params err
		c.ReplyErr(errcode.ErrParams)
		beego.Error("set etl error : ", err)
	}
	res := map[string]string{
		"res": "success",
	}
	c.ReplySucc(res)
}

func (c *Controller) StartEtl() {

}
