package form

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"time"
)

type Controller struct {
	base.Controller
}

const (
	TimeFormat = "2006-01-02 15:04:05"
	DateFormat = "2006-01-02"
	CommonErr  = 99999
)

//用json参数 todo
func (c *Controller) AddFormtpl() {
	prefix := c.UserComp
	name := c.GetString("Name")
	formtype := c.GetString("Type")
	desc := c.GetString("Desc")
	content := c.GetString("Content")
	//attachment := 文件名字符串数组 todo
	begintimeStr := c.GetString("BginTime")
	ftpl := model.Formtpl{
		No:      fmt.Sprintf("ftpl%d", time.Now().Unix()),
		Name:    name,
		Type:    formtype,
		Desc:    desc,
		Content: content,
		Ctime:   time.Now(),
		//Attachment:model.NewStrSlice(),
	}
	begintime, e := time.Parse(DateFormat, begintimeStr)
	if e != nil {
		t := time.Now()
		begintime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Local())
		ftpl.Status = model.Abled
	}
	ftpl.Status = model.Init

	ftpl.BeginTime = begintime
	e = service.AddFormtpl(prefix, &ftpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

func (c *Controller) UpdateFormtpl() {
	prefix := c.UserComp
	no := c.GetString("No")
	name := c.GetString("Name")
	formtype := c.GetString("Type")
	desc := c.GetString("Desc")
	content := c.GetString("Content")
	//attachment := 文件名字符串数组 todo
	begintimeStr := c.GetString("BginTime")
	ftpl := model.Formtpl{
		No:      no,
		Name:    name,
		Type:    formtype,
		Desc:    desc,
		Content: content,
		//Attachment:model.NewStrSlice(),
	}
	begintime, e := time.Parse(DateFormat, begintimeStr)
	if e != nil {
		t := time.Now()
		begintime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		ftpl.Status = model.Abled
	}
	ftpl.Status = model.Init
	ftpl.BeginTime = begintime

	e = service.UpdateFormtpl(prefix, &ftpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

func (c *Controller) ControlFormtpl() {
	prefix := c.UserComp
	no := c.GetString("No")
	status, e := c.GetInt("Status")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.ControlFormtpl(prefix, no, status)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

func (c *Controller) DelFormtpl() {
	prefix := c.UserComp
	no := c.GetString("No")
	e := service.DelFormtpl(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplyErr("success")
	}
}

func (c *Controller) AddApprovaltpl() {
	str := c.GetString("approvaltpl")
	atpl := model.Approvaltpl{}
	e := json.Unmarshal([]byte(str), &atpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	atpl.No = fmt.Sprintf("atpl%d", time.Now().Unix())
	atpl.Ctime = time.Now()
	if atpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		atpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		atpl.Status = model.Abled
	} else {
		atpl.Status = model.Init
	}

}
