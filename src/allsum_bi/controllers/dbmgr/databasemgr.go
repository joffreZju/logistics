package dbmgr

import (
	"allsum_bi/controllers/base"
	"allsum_bi/db/conn"
	"allsum_bi/models"
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

	conninfo := conn.Conn{
		Id:     databasemgr.Dbid,
		Dbtype: DbType,
		Name:   Name,
		Host:   Host,
		Port:   Port,
		DbUser: DbUser,
		Passwd: DbPasswd,
		Dbname: DbName,
	}
	err = conn.CreateConn(conninfo)
	if err != nil {
		return
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

func (c *Controller) UpdateDb() {
	Name := c.GetString("name")
	Host := c.GetString("host")
	dbid := c.GetString("dbid")
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

	conninfo := conn.Conn{
		Id:     dbid,
		Dbtype: DbType,
		Name:   Name,
		Host:   Host,
		Port:   Port,
		DbUser: DbUser,
		Passwd: DbPasswd,
		Dbname: DbName,
	}
	conn.RemoveConn(dbid)

	err = conn.CreateConn(conninfo)
	if err != nil {
		return
	}

	databasemgr := models.DatabaseManager{
		Dbid:     dbid,
		Dbname:   DbName,
		Dbtype:   DbType,
		Host:     Host,
		Port:     Port,
		Dbuser:   DbUser,
		Password: DbPasswd,
		Name:     Name,
	}
	err = models.UpdateDatabaseManager(databasemgr, "dbname", "dbtype", "host", "port", "dbuser", "db_passwd", "name")
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

func (c *Controller) DelDb() {
	dbid := c.GetString("dbid")
	conn.RemoveConn(dbid)
	err := models.DeleteDatabaseManager(dbid)
	if err != nil {
		beego.Error("delete db err :", err)
		c.ReplyErr(errcode.ErrServerError)
		return
	}
	res := map[string]string{
		"res": "ok",
	}
	c.ReplySucc(res)
}
