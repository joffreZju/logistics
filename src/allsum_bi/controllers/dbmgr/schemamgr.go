package dbmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/db"
	"stowage/common/lib/errcode"

	"github.com/astaxie/beego"
)

type Controller struct {
	base.Controller
}

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
