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
	CommonErr = 99999
)

func (c *Controller) AddFormtpl() {
	prefix := c.UserComp
	str := c.GetString("formtpl")
	ftpl := model.Formtpl{}
	e := json.Unmarshal([]byte(str), &ftpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	ftpl.No = model.UniqueNo("Ftpl")
	ftpl.Ctime = time.Now()
	if ftpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		ftpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		ftpl.Status = model.TplAbled
	} else {
		ftpl.Status = model.TplInit
	}
	ftpl.Attachment = model.NewStrSlice()
	e = service.AddFormtpl(prefix, &ftpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateFormtpl() {
	prefix := c.UserComp
	str := c.GetString("formtpl")
	ftpl := model.Formtpl{}
	e := json.Unmarshal([]byte(str), &ftpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if ftpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		ftpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		ftpl.Status = model.TplAbled
	} else {
		ftpl.Status = model.TplInit
	}
	//ftpl.Attachment = model.NewStrSlice() todo

	e = service.UpdateFormtpl(prefix, &ftpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
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
		c.ReplySucc(nil)
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
		c.ReplySucc(nil)
	}
}

func (c *Controller) AddApprovaltpl() {
	prefix := c.UserComp
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
		atpl.Status = model.TplAbled
	} else {
		atpl.Status = model.TplInit
	}
	e = service.AddApprovaltpl(prefix, &atpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateApprovaltpl() {
	prefix := c.UserComp
	str := c.GetString("approvaltpl")
	atpl := model.Approvaltpl{}
	e := json.Unmarshal([]byte(str), &atpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if atpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		atpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		atpl.Status = model.TplAbled
	} else {
		atpl.Status = model.TplInit
	}
	e = service.UpdateApprovaltpl(prefix, &atpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) ControlApprovaltpl() {
	prefix := c.UserComp
	no := c.GetString("No")
	status, e := c.GetInt("Status")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	e = service.ControlApprovaltpl(prefix, no, status)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) DelApprovaltpl() {
	prefix := c.UserComp
	no := c.GetString("No")
	e := service.DelApprovaltpl(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) AddApproval() {
	prefix := c.UserComp
	astr := c.GetString("approval")
	fstr := c.GetString("form")
	a := model.Approval{}
	f := model.Form{}
	e1 := json.Unmarshal([]byte(astr), &a)
	e2 := json.Unmarshal([]byte(fstr), &f)
	if e1 != nil || e2 != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e1, e2)
		return
	}
	if a.Status != model.ApproveInit && a.Status != model.Approving {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("approval status is wrong")
		return
	}
	f.No = fmt.Sprintf("form%d", time.Now().Unix())
	f.Ctime = time.Now()
	//f.Attachment = model.NewStrSlice() todo
	a.No = fmt.Sprintf("aprvl%d", time.Now().Unix())
	a.Ctime = time.Now()
	a.FormNo = f.No
	e := service.AddApproval(prefix, &f, &a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateApproval() {
	prefix := c.UserComp
	astr := c.GetString("approval")
	fstr := c.GetString("form")
	a := model.Approval{}
	f := model.Form{}
	e1 := json.Unmarshal([]byte(astr), &a)
	e2 := json.Unmarshal([]byte(fstr), &f)
	if e1 != nil || e2 != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e1, e2)
		return
	}
	if a.Status != model.ApproveInit && a.Status != model.Approving {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("approval status is wrong")
		return
	}
	//f.Attachment = model.NewStrSlice() todo
	e := service.UpdateApproval(prefix, &f, &a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) CancelApproval() {
	prefix := c.UserComp
	no := c.GetString("No")
	e := service.CancelApproval(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

func (c *Controller) Approve() {
	prefix := c.UserComp
	str := c.GetString("approve")
	aflow := model.ApproveFlow{}
	e := json.Unmarshal([]byte(str), &aflow)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if aflow.Opinion != model.ApproveOpinionAgree && aflow.Opinion != model.ApproveOpinionRefuse {
		c.ReplyErr(errcode.New(CommonErr, "opinion is wrong"))
		return
	}
	e = service.Approve(prefix, &aflow)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}
