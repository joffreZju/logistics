package dbmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/db/models"
	"allsum_bi/util/errcode"
	"fmt"

	"github.com/astaxie/beego"
	"github.com/satori/go.uuid"
)

type Controller struct {
	base.Controller
}

func (c *Controller) ListDb() {
	var res []map[string]string
	databases, err := models.ListDatabaseManager()
	if err != nil {
		beego.Error("listdb err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	for _, v := range databases {
		var resbase map[string]string
		resbase = map[string]string{
			"id":   v.Dbid,
			"name": v.Name,
		}
		res = append(res, resbase)
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) ListDbDetail() {
	var res []map[string]string
	databases, err := models.ListDatabaseManager()
	if err != nil {
		beego.Error("listdb err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	for _, v := range databases {
		var resbase map[string]string
		resbase = map[string]string{
			"id":     v.Dbid,
			"name":   v.Name,
			"dbname": v.Dbname,
			"host":   v.Host,
			"port":   fmt.Sprintf("%d", v.Port),
			//		"username": v.Dbuser,
			//		"password": v.Password,
		}
		res = append(res, resbase)
	}
	c.ReplySucc(res)
	return
}

func (c *Controller) AddDb() {
	Name := c.GetString("name")
	Host := c.GetString("host")
	Port, err := c.GetInt("port")
	if err != nil {
		beego.Error("beego err ", err)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	DbType := c.GetString("dbtype")
	DbName := c.GetString("dbname")
	DbUser := c.GetString("dbuser")
	DbPasswd := c.GetString("dbpasswd")

	databasemgr := models.DatabaseManager{
		Dbid:     uuid.NewV4().String(),
		Dbname:   DbName,
		Dbtype:   DbType,
		Host:     Host,
		Port:     Port,
		Dbuser:   DbUser,
		Password: DbPasswd,
		Name:     Name,
	}
	err = models.InsertDatabaseManager(databasemgr)
	if err != nil {
		beego.Error("listdb err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string]string{
		"res": "ok",
	}
	c.ReplySucc(res)
	return
}
