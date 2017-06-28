package form

import (
	"allsum_oa/controller/base"
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/errcode"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

type Controller struct {
	base.Controller
}

const (
	CommonErr = 99999
)

//表单模板增删改查*************************
func (c *Controller) GetFormtplList() {
	prefix := c.UserComp
	ftpls, e := service.GetFormtplList(prefix)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(ftpls)
	}
}

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
	no := c.GetString("no")
	status, e := c.GetInt("status")
	if e != nil || len(no) == 0 {
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
	no := c.GetString("no")
	e := service.DelFormtpl(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//审批单模板增删改查*************************
func (c *Controller) GetApprovaltplList() {
	prefix := c.UserComp
	atpls, e := service.GetApprocvaltplList(prefix)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(atpls)
	}
}

func (c *Controller) GetApprovaltplDetail() {
	prefix := c.UserComp
	no := c.GetString("no")
	atpl, e := service.GetApprovaltplDetail(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(atpl)
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
	atpl.No = model.UniqueNo("Atpl")
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
	no := c.GetString("no")
	status, e := c.GetInt("status")
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
	no := c.GetString("no")
	e := service.DelApprovaltpl(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//审批流相关接口***************************
//获取符合条件的审批人
func (c *Controller) GetApproverList() {
	prefix := c.UserComp
	currentGroup, e := c.GetInt("currentGroup")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	atplNo := c.GetString("approvaltplNo")
	if strings.Contains(c.UserGroups, fmt.Sprintf("%d", currentGroup)) == false {
		c.ReplyErr(errcode.ErrGroupOfUser)
		return
	}
	rolemap, e := service.GetApproverList(prefix, atplNo, currentGroup)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(rolemap)
	}
}

func (c *Controller) AddApproval() {
	prefix := c.UserComp
	str := c.GetString("approval")
	a := model.Approval{}
	e := json.Unmarshal([]byte(str), &a)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if (a.Status != model.ApproveDraft && a.Status != model.Approving) || len(a.UserFlow) == 0 {
		c.ReplyErr(errcode.New(CommonErr, "审批单状态错误"))
		beego.Error("approval status is wrong")
		return
	}
	a.FormContent.No = model.UniqueNo("F")
	a.FormContent.Ctime = time.Now()
	a.No = model.UniqueNo("A")
	a.Ctime = time.Now()
	a.FormNo = a.FormContent.No
	a.Currentuser = a.UserFlow[0]
	e = service.AddApproval(prefix, &a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		if a.Status == model.Approving {
			//todo 向第一个审批人推送消息,修改状态为approving
		}
		c.ReplySucc(nil)
	}
}

func (c *Controller) UpdateApproval() {
	prefix := c.UserComp
	str := c.GetString("approval")
	a := model.Approval{}
	e := json.Unmarshal([]byte(str), &a)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if (a.Status != model.ApproveDraft && a.Status != model.Approving) || len(a.UserFlow) == 0 {
		c.ReplyErr(errcode.New(CommonErr, "审批单状态错误"))
		beego.Error("approval status is wrong")
		return
	}
	if len(a.No) == 0 || len(a.FormNo) == 0 || a.FormNo != a.FormContent.No {
		c.ReplyErr(errcode.New(CommonErr, "审批单编号有误"))
		return
	}
	a.Currentuser = a.UserFlow[0]
	e = service.UpdateApproval(prefix, &a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		if a.Status == model.Approving {
			//todo 向第一个审批人推送消息,修改状态为approving
		}
		c.ReplySucc(nil)
	}
}

func (c *Controller) CancelApproval() {
	prefix := c.UserComp
	no := c.GetString("no")
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
	aflow.Ctime = time.Now()
	aflow.UserId = c.UserID
	if aflow.UserId <= 0 || aflow.UserId != c.UserID {
		c.ReplyErr(errcode.New(CommonErr, "user id is wrong"))
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

//获取我发起的所有审批单
func (c *Controller) GetApprovalsFromMe() {
	prefix := c.UserComp
	uid := c.UserID
	alist, e := service.GetApprovalsFromMe(prefix, uid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}

//获取需要我审批的审批单
func (c *Controller) GetTodoApprovalsToMe() {
	prefix := c.UserComp
	uid := c.UserID
	alist, e := service.GetTodoApprovalsToMe(prefix, uid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}

//获取我审批过的审批单
func (c *Controller) GetFinishedApprovalsToMe() {
	prefix := c.UserComp
	uid := c.UserID
	alist, e := service.GetFinishedApprovalsToMe(prefix, uid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}

//获取我审批过的审批单
func (c *Controller) GetApprovalDetail() {
	prefix := c.UserComp
	no := c.GetString("no")
	alist, e := service.GetApprovalDetail(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}
