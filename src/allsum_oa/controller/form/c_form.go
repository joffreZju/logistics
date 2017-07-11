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
	ftpl := &model.Formtpl{}
	e := json.Unmarshal([]byte(str), ftpl)
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
	e = service.AddFormtpl(prefix, ftpl)
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
	ftpl := &model.Formtpl{}
	e := json.Unmarshal([]byte(str), ftpl)
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
	e = service.UpdateFormtpl(prefix, ftpl)
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
	if status != model.TplAbled && status != model.TplDisabled {
		c.ReplyErr(errcode.ErrParams)
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
	name := c.GetString("name")
	var atpls []*model.Approvaltpl
	var e error
	if len(name) != 0 {
		atpls, e = service.GetApprocvaltplList(prefix, name)
	} else {
		atpls, e = service.GetApprocvaltplList(prefix)
	}
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
	atpl := &model.Approvaltpl{}
	e := json.Unmarshal([]byte(str), atpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if (atpl.SkipBlankRole != model.SkipBlankRoleNo && atpl.SkipBlankRole != model.SkipBlankRoleYes) ||
		(atpl.TreeFlowUp != model.TreeFlowUpNo && atpl.TreeFlowUp != model.TreeFlowUpYes) ||
		len(atpl.RoleFlow) == 0 {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("审批单模板设置错误")
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
	e = service.AddApprovaltpl(prefix, atpl)
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
	atpl := &model.Approvaltpl{}
	e := json.Unmarshal([]byte(str), atpl)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if (atpl.SkipBlankRole != model.SkipBlankRoleNo && atpl.SkipBlankRole != model.SkipBlankRoleYes) ||
		(atpl.TreeFlowUp != model.TreeFlowUpNo && atpl.TreeFlowUp != model.TreeFlowUpYes) ||
		len(atpl.RoleFlow) == 0 {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("审批单模板设置错误")
		return
	}
	if atpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		atpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		atpl.Status = model.TplAbled
	} else {
		atpl.Status = model.TplInit
	}
	e = service.UpdateApprovaltpl(prefix, atpl)
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
	if status != model.TplAbled && status != model.TplDisabled {
		c.ReplyErr(errcode.ErrParams)
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
func (c *Controller) AddApproval() {
	prefix := c.UserComp
	atplNo := c.GetString("approvaltplNo")
	atpl, e := service.GetApprovaltpl(prefix, atplNo)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	astr := c.GetString("approval")
	a := &model.Approval{}
	e = json.Unmarshal([]byte(astr), a)
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if a.Status != model.ApprovalStatWaiting {
		c.ReplyErr(errcode.New(CommonErr, "审批单设置错误"))
		beego.Error("审批单设置错误")
		return
	}
	//处理表单内容
	a.FormContent.No = model.UniqueNo("F")
	a.FormContent.Ctime = time.Now()
	//生成编号
	a.No = model.UniqueNo("A")
	a.Ctime = time.Now()
	a.FormNo = a.FormContent.No
	//从模板中抽取审批流设定的条件
	a.TreeFlowUp = atpl.TreeFlowUp
	a.SkipBlankRole = atpl.SkipBlankRole
	a.RoleFlow = atpl.RoleFlow
	e = service.AddApproval(prefix, a)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
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
	uid := c.UserID
	ano := c.GetString("approvalNo")
	comment := c.GetString("comment")
	status, e := c.GetInt("status")
	if e != nil || (status != model.ApprovalStatAccessed && status != model.ApprovalStatNotAccessed) {
		beego.Error(e)
		c.ReplyErr(errcode.ErrParams)
		return
	}
	a, e := service.GetApproval(prefix, ano)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	af, e := service.GetLatestFlowOfApproval(prefix, a.No)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	if a.Status != model.ApprovalStatWaiting || af.Status != model.ApprovalStatWaiting {
		c.ReplyErr(errcode.ErrStatOfApproval)
		return
	}
	if !strings.Contains(af.MatchUsers, fmt.Sprintf("_%d_", uid)) {
		c.ReplyErr(errcode.ErrRoleOfUser)
		return
	}
	user, e := service.GetUserById("public", uid)
	if e != nil {
		c.ReplyErr(errcode.ErrGetUserInfoFailed)
		return
	}
	af.Status = status
	af.Comment = comment
	af.UserId = uid
	af.UserName = user.UserName
	e = service.Approve(prefix, a, af)
	if e != nil {
		c.ReplyErr(errcode.ErrStatOfApproval)
		beego.Error(e)
	} else {
		c.ReplySucc(a)
	}
}

//获取我发起的所有审批单
func (c *Controller) GetApprovalsFromMe() {
	prefix := c.UserComp
	uid := c.UserID
	beginTime := c.GetString("beginTime")
	condition := c.GetString("statusCondition")
	alist, e := service.GetApprovalsFromMe(prefix, uid, beginTime, condition)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}

//获取我收到的审批单
func (c *Controller) GetApprovalsToMe() {
	prefix := c.UserComp
	uid := c.UserID
	beginTime := c.GetString("beginTime")
	condition := c.GetString("statusCondition")
	var alist []*model.Approval
	var e error
	if condition == model.GetApprovalApproving {
		alist, e = service.GetTodoApprovalsToMe(prefix, uid, beginTime)
	} else if condition == model.GetApprovalFinished {
		alist, e = service.GetFinishedApprovalsToMe(prefix, uid, beginTime)
	} else {
		alist, e = service.GetTodoApprovalsToMe(prefix, uid, beginTime)
		if e != nil {
			beego.Error(e)
		}
		var alistFinished []*model.Approval
		alistFinished, e = service.GetFinishedApprovalsToMe(prefix, uid, beginTime)
		alist = append(alist, alistFinished...)
	}
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(alist)
}

//获取审批单详情
func (c *Controller) GetApprovalDetail() {
	prefix := c.UserComp
	no := c.GetString("no")
	a, e := service.GetApprovalDetail(prefix, no)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	c.ReplySucc(a)
}
