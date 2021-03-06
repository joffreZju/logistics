package form

import (
	"allsum_oa/model"
	"allsum_oa/service"
	"common/lib/baseController"
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
	//判断生效时间，并赋予对应的状态值
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

//禁用或启用表单模板
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
		atpls, e = service.GetApprovaltplList(prefix, name)
	} else {
		atpls, e = service.GetApprovaltplList(prefix)
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

//获取某角色可以匹配到的组织
func (c *Controller) GetMatchGroupsOfRole() {
	rid, e := c.GetInt("rid")
	if e != nil {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	groups, e := service.GetMatchGroupsOfRole(c.UserComp, rid)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(groups)
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
	if atpl.EmailMsg != model.EmailMsgNo && atpl.EmailMsg != model.EmailMsgYes {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("审批单模板设置错误")
		return
	}
	if len(atpl.AllowRoles) == 0 || len(atpl.FlowContent) == 0 {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("审批单模板设置错误")
		return
	}
	atpl.No = model.UniqueNo("Atpl")
	//判断生效时间并置状态值
	atpl.Ctime = time.Now()
	if atpl.BeginTime.Sub(time.Now()).Hours() < 0 {
		t := time.Now()
		atpl.BeginTime = time.Date(t.Year(), t.Month(), t.Day(), 0, 0, 0, 0, t.Location())
		atpl.Status = model.TplAbled
	} else {
		atpl.Status = model.TplInit
	}
	//检测流程设置是否合法
	for _, v := range atpl.FlowContent {
		if (v.Necessary != model.FlowNecessaryNo && v.Necessary != model.FlowNecessaryYes) ||
			v.RoleId == 0 {
			c.ReplyErr(errcode.ErrParams)
			beego.Error("审批单模板设置错误")
			return
		}
		v.ApprovaltplNo = atpl.No
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
	if e != nil || len(atpl.No) == 0 {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	if atpl.EmailMsg != model.EmailMsgNo && atpl.EmailMsg != model.EmailMsgYes {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("审批单模板设置错误")
		return
	}
	if len(atpl.AllowRoles) == 0 || len(atpl.FlowContent) == 0 {
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
	//检测流程设置是否合法
	for _, v := range atpl.FlowContent {
		if (v.Necessary != model.FlowNecessaryNo && v.Necessary != model.FlowNecessaryYes) ||
			v.RoleId == 0 {
			c.ReplyErr(errcode.ErrParams)
			beego.Error("审批单模板设置错误")
			return
		}
		v.ApprovaltplNo = atpl.No
	}
	e = service.UpdateApprovaltpl(prefix, atpl)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//禁用或启用审批单模板
func (c *Controller) ControlApprovaltpl() {
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
	if len(no) == 0 {
		c.ReplyErr(errcode.ErrParams)
		beego.Error("no长度为零")
		return
	}
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
	astr := c.GetString("approval")
	a := &model.Approval{}
	e := json.Unmarshal([]byte(astr), a)
	if e != nil || a.UserId != c.UserID {
		c.ReplyErr(errcode.ErrParams)
		beego.Error(e)
		return
	}
	//处理表单内容
	a.FormContent.No = model.UniqueNo("F")
	a.FormContent.Ctime = time.Now()
	//生成编号
	a.No = model.UniqueNo("A")
	a.Ctime = time.Now()
	a.Status = model.ApprovalStatWaiting
	a.FormNo = a.FormContent.No
	e = service.AddApproval(prefix, a, atplNo)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
	} else {
		c.ReplySucc(nil)
	}
}

//发起人撤销未完成的审批流
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

//审批人进行审批
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
	af, e := service.GetApproveFlowById(prefix, a.CurrentFlow)
	if e != nil {
		c.ReplyErr(errcode.New(CommonErr, e.Error()))
		beego.Error(e)
		return
	}
	//检测整个审批单的状态和当前一步的状态，如果不是等待审批，那么退出
	if a.Status != model.ApprovalStatWaiting || af.Status != model.ApprovalStatWaiting {
		c.ReplyErr(errcode.ErrStatOfApproval)
		return
	}
	//检测当前一步的审批人列表，是否包含当前用户
	if !strings.Contains(af.MatchUsers, fmt.Sprintf("-%d-", uid)) {
		c.ReplyErr(errcode.ErrInfoOfUser)
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

//获取我发起的所有审批单，可以根据beginTime和审批单状态进行过滤
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

//获取我收到的审批单，可以根据beginTime和审批单状态进行过滤
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
