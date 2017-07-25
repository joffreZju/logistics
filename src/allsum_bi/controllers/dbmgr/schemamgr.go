package dbmgr

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"common/lib/errcode"

	"github.com/astaxie/beego"
)

func (c *Controller) ListSchema() {
	dbid := c.GetString("dbid")
	schemas, err := db.ListSchema(dbid)
	if err != nil {
		beego.Error("listdb err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string][]string{
		"schemas": schemas,
	}
	c.ReplySucc(res)
}

func (c *Controller) ListSchemaTable() {
	dbid := c.GetString("dbid")
	schema := c.GetString("schema")
	tableNames, err := db.ListSchemaTable(dbid, schema)
	if err != nil {
		beego.Error("listdb err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string][]string{
		"tableNames": tableNames,
	}
	c.ReplySucc(res)
}

func (c *Controller) ListAllDbSchema() {
	databases, err := models.ListDatabaseManager()
	if err != nil {
		beego.Error("list db err :", err)
		c.ReplyErr(errcode.ErrActionGetDbMgr)
		return
	}
	schemadata := []map[string]interface{}{}
	for _, database := range databases {
		schemas, err := db.ListSchema(database.Dbid)
		if err != nil {
			beego.Error("list schema err:", err)
			c.ReplyErr(errcode.ErrActionGetDbMgr)
			return
		}
		databaseschema := map[string]interface{}{
			"name":   database.Name,
			"dbid":   database.Dbid,
			"schema": schemas,
		}
		schemadata = append(schemadata, databaseschema)
	}
	c.ReplySucc(schemadata)
}
