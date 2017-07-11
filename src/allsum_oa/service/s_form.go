package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"strconv"
	"strings"
)

//表单模板相关
func GetFormtplList(prefix string) (ftpls []*model.Formtpl, e error) {
	ftpls = []*model.Formtpl{}
	e = model.NewOrm().Table(prefix + "." + model.Formtpl{}.TableName()).Order("no desc").Find(&ftpls).Error
	return
}

func AddFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	e = model.NewOrm().Table(prefix + "." + ftpl.TableName()).Create(ftpl).Error
	return
}

func UpdateFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + ftpl.TableName()).
		Model(ftpl).Updates(ftpl).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

func ControlFormtpl(prefix, no string, status int) (e error) {
	count := 0
	if status == model.TplDisabled {
		e = model.NewOrm().Table(prefix+"."+model.Approvaltpl{}.TableName()).
			Where("formtpl_no=?", no).Count(&count).Error
		if e != nil {
			return
		} else if count != 0 {
			return errors.New("some approvaltpl are using this formtpl")
		}
	}
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix+"."+model.Formtpl{}.TableName()).
		Model(&model.Formtpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

func DelFormtpl(prefix, no string) (e error) {
	count := 0
	e = model.NewOrm().Table(prefix+"."+model.Approvaltpl{}.TableName()).
		Where("formtpl_no=?", no).Count(&count).Error
	if e != nil {
		return
	} else if count != 0 {
		return errors.New("some approvaltpl are using this formtpl")
	}
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + model.Formtpl{}.TableName()).
		Delete(&model.Formtpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

//审批单模板相关
func GetApprocvaltplList(prefix string, params ...string) (atpls []*model.Approvaltpl, e error) {
	db := model.NewOrm().Table(prefix + "." + model.Approvaltpl{}.TableName()).Order("no desc")
	atpls = []*model.Approvaltpl{}
	if len(params) != 0 {
		e = db.Where("name like ?", "%"+params[0]+"%").Find(&atpls).Error
	} else {
		e = db.Find(&atpls).Error
	}
	return
}

func GetApprovaltpl(prefix, atplno string) (atpl *model.Approvaltpl, e error) {
	atpl = &model.Approvaltpl{}
	e = model.NewOrm().Table(prefix+"."+model.Approvaltpl{}.TableName()).
		First(atpl, "no=?", atplno).Error
	return
}

func GetApprovaltplDetail(prefix, atplno string) (atpl *model.Approvaltpl, e error) {
	db := model.NewOrm()
	atpl = &model.Approvaltpl{}
	e = db.Table(prefix+"."+model.Approvaltpl{}.TableName()).
		First(atpl, "no=?", atplno).Error
	if e != nil {
		return
	}
	atpl.FormtplContent = new(model.Formtpl)
	e = db.Table(prefix+"."+model.Formtpl{}.TableName()).
		First(atpl.FormtplContent, "no=?", atpl.FormtplNo).Error
	if e != nil {
		return
	}
	return
}

func AddApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	e = model.NewOrm().Table(prefix + "." + atpl.TableName()).Create(atpl).Error
	return
}

func UpdateApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + atpl.TableName()).
		Model(atpl).Updates(atpl).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	return tx.Commit().Error
}

func ControlApprovaltpl(prefix, no string, status int) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix+"."+model.Approvaltpl{}.TableName()).
		Model(&model.Approvaltpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	return tx.Commit().Error
}

func DelApprovaltpl(prefix, no string) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + "." + model.Approvaltpl{}.TableName()).
		Delete(&model.Approvaltpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approvaltpl no")
		return
	}
	return tx.Commit().Error
}

//func UpdateApproval(prefix string, a *model.Approval) (e error) {
//	aprvl := model.Approval{}
//	e = model.NewOrm().Table(prefix + "." + aprvl.TableName()).First(&aprvl, "no=?", a.No).Error
//	if e != nil {
//		return
//	}
//	if aprvl.Status != model.ApprovalStatDraft {
//		e = errors.New("approval is already commited")
//		return
//	}
//	tx := model.NewOrm().Begin()
//	c := tx.Table(prefix + "." + a.FormContent.TableName()).
//		Model(a.FormContent).Updates(a.FormContent).RowsAffected
//	if c != 1 {
//		tx.Rollback()
//		e = errors.New("wrong form no")
//		return
//	}
//	c = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).RowsAffected
//	if c != 1 {
//		tx.Rollback()
//		e = errors.New("wrong approval no")
//		return
//	}
//	return tx.Commit().Error
//}

func CancelApproval(prefix, no string) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.Approval{}.TableName())
	a := model.Approval{}
	e = db.First(&a, "no=?", no).Error
	if e != nil {
		return
	}
	if a.Status == model.ApprovalStatAccessed || a.Status == model.ApprovalStatNotAccessed {
		e = errors.New("审批单已经完成")
		return
	}
	tx := db.Begin()
	c := tx.Model(&a).Update("status", model.ApprovalStatCanceled).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("审批单编号错误")
		return
	}
	return tx.Commit().Error
}

func AddApproval(prefix string, a *model.Approval) (e error) {
	db := model.NewOrm()
	user := &model.User{}
	group := &model.Group{}
	role := &model.Role{}
	e = db.Table(prefix+"."+user.TableName()).First(user, a.UserId).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+group.TableName()).First(group, a.GroupId).Error
	if e != nil {
		return
	}
	e = db.Table(prefix+"."+role.TableName()).First(role, a.RoleId).Error
	if e != nil {
		return
	}
	a.UserName = user.UserName
	a.GroupName = group.Name
	a.RoleName = role.Name
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + "." + a.TableName()).Create(a).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = tx.Table(prefix + "." + a.FormContent.TableName()).Create(a.FormContent).Error
	if e != nil {
		tx.Rollback()
		return
	}
	go nextStepOfApproval(prefix, a)
	return tx.Commit().Error
}

//当前角色审批完成后调用next，会判断整个审批单是否完成，如果没有那么继续下一步
func nextStepOfApproval(prefix string, a *model.Approval) {
	var e error
	db := model.NewOrm()
	var nextLoc int
	if a.CurrentRole == 0 {
		nextLoc = 0
	} else {
		for k, v := range a.RoleFlow {
			if v != a.CurrentRole {
				continue
			}
			if k == len(a.RoleFlow)-1 {
				//审批最后一步已完成
				a.Status = model.ApprovalStatAccessed
				e = db.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
				if e != nil {
					beego.Error(e)
				}
				go newMsgToCreator(prefix, a)
				return
			} else {
				nextLoc = k + 1
				break
			}
		}
	}
	stop := false
	for i := nextLoc; i < len(a.RoleFlow); i++ {
		rid := a.RoleFlow[i]
		var users []*model.User
		users, e = getApproverByRole(prefix, a, rid)
		if e == nil {
			//找到若干符合条件的审批人
			//开始创建一步流程
			var matchUsers string //拼接userId
			for _, v := range users {
				matchUsers += fmt.Sprintf("%d_", v.Id)
			}
			//skip oneself
			//if strings.Contains(matchUsers, fmt.Sprintf("%d_", a.UserId)){
			//	continue
			//}
			r := &model.Role{}
			e = db.Table(prefix+"."+r.TableName()).First(r, "id=?", rid).Error
			if e != nil {
				beego.Error(e)
			}
			af := &model.ApproveFlow{
				ApprovalNo: a.No,
				MatchUsers: matchUsers,
				RoleId:     rid,
				RoleName:   r.Name,
				Status:     model.ApprovalStatWaiting,
			}
			e = db.Table(prefix + "." + af.TableName()).Create(af).Error
			if e != nil {
				stop = true
				break
			}
			//更新审批单当前角色信息
			a.CurrentRole = rid
			e = db.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
			if e != nil {
				stop = true
				break
			}
			go newMsgToApprovers(prefix, users, a)
			return
		} else if e == gorm.ErrRecordNotFound && a.SkipBlankRole == model.SkipBlankRoleYes {
			//没有符合条件的审批人，跳过
			continue
		} else {
			//审批无法流转下去，没有审批人而且不允许跳过
			stop = true
			break
		}
	}
	if stop {
		beego.Error("审批单无法继续流转:", e)
		a.Status = model.ApprovalStatStop
		e = db.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
		if e != nil {
			beego.Error("尝试停止审批单:", e)
		}
		go newMsgToCreator(prefix, a)
	} else {
		//后面角色全部跳过，审批单完全通过
		a.Status = model.ApprovalStatAccessed
		e = db.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
		if e != nil {
			beego.Error(e)
		}
		go newMsgToCreator(prefix, a)
	}
	return
}

//找到审批角色下的用户
func getApproverByRole(prefix string, a *model.Approval, rid int) (users []*model.User, e error) {
	//完全按照角色流动
	if a.TreeFlowUp == model.TreeFlowUpNo {
		users, e = GetUsersOfRole(prefix, rid)
		return
	}
	//在组织树路径上寻找符合条件的角色
	db := model.NewOrm()
	g := &model.Group{}
	e = db.Table(prefix+"."+g.TableName()).First(g, "id=?", a.GroupId).Error
	if e != nil {
		return
	}
	pids := []int{}
	for _, v := range strings.Split(g.Path, "_") {
		pid, e := strconv.Atoi(v)
		if e != nil {
			return nil, e
		}
		pids = append(pids, pid)
	}
	//找到同时在组织路径上，在审批角色流里面的所有用户
	sql := fmt.Sprintf(`SELECT * from "%s".allsum_user WHERE id in
		(SELECT t1.user_id FROM "%s".user_group as t1 INNER JOIN "%s".user_role as t2
		on t1.user_id = t2.user_id
		where t1.group_id in (?) and t2.role_id=? )`, prefix, prefix, prefix)
	users = []*model.User{}
	e = db.Raw(sql, pids, rid).Scan(&users).Error
	return
}

func Approve(prefix string, a *model.Approval, af *model.ApproveFlow) (e error) {
	tx := model.NewOrm().Begin()
	//审批，修改一步审批流程状态
	count := tx.Table(prefix+"."+af.TableName()).Where("id=? and status=?", af.Id, model.ApprovalStatWaiting).Updates(af).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("审批失败")
	}
	if af.Status == model.ApprovalStatNotAccessed {
		a.Status = model.ApprovalStatNotAccessed
		count = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).RowsAffected
		if count != 1 {
			tx.Rollback()
			return errors.New("审批失败")
		}
		go newMsgToCreator(prefix, a)
	} else {
		go nextStepOfApproval(prefix, a)
	}
	return tx.Commit().Error
}

func newMsgToCreator(company string, a *model.Approval) {
	msg := &model.Message{
		CompanyNo: company,
		UserId:    a.UserId,
		MsgType:   model.MsgTypeApprove,
		Content: model.JsonMap{
			"ApprovalNo": a.No,
		},
	}
	if a.Status == model.ApprovalStatAccessed {
		msg.Title = "你的" + a.Name + "审批单已通过"
	} else if a.Status == model.ApprovalStatNotAccessed {
		msg.Title = "你的" + a.Name + "审批单被拒绝"
	} else if a.Status == model.ApprovalStatStop {
		msg.Title = "你的" + a.Name + "审批单无法继续流转，请咨询管理员或客服"
	}
	e := SaveAndSendMsg(msg)
	if e != nil {
		beego.Error(e)
	}
}

func newMsgToApprovers(company string, users []*model.User, a *model.Approval) {
	title := "来自$" + a.UserName + "$的审批消息"
	for _, v := range users {
		msg := &model.Message{
			CompanyNo: company,
			UserId:    v.Id,
			MsgType:   model.MsgTypeApprove,
			Title:     title,
			Content: model.JsonMap{
				"ApprovalNo": a.No,
			},
		}
		go SaveAndSendMsg(msg)
	}
}

func GetApprovalDetail(prefix, no string) (a *model.Approval, e error) {
	a = new(model.Approval)
	db := model.NewOrm()
	e = db.Table(prefix+"."+a.TableName()).First(a, "no=?", no).Error
	if e != nil {
		return
	}
	a.FormContent = new(model.Form)
	e = db.Table(prefix+"."+model.Form{}.TableName()).
		First(a.FormContent, "no=?", a.FormNo).Error
	if e != nil {
		return
	}
	e = db.Table(prefix + "." + model.ApproveFlow{}.TableName()).Order("ctime").
		Where(&model.ApproveFlow{ApprovalNo: a.No}).Find(&a.ApproveFLows).Error
	return
}

func GetApproval(prefix, no string) (a *model.Approval, e error) {
	a = new(model.Approval)
	e = model.NewOrm().Table(prefix+"."+a.TableName()).First(a, "no=?", no).Error
	return
}

func GetLatestFlowOfApproval(prefix, approvalNo string) (af *model.ApproveFlow, e error) {
	af = new(model.ApproveFlow)
	e = model.NewOrm().Table(prefix+"."+af.TableName()).Order("id desc").Limit(1).
		First(af, "approval_no=?", approvalNo).Error
	return
}

func GetApprovalsFromMe(prefix string, uid int, beginTime, condition string) (alist []*model.Approval, e error) {
	alist = []*model.Approval{}
	db := model.NewOrm().Table(prefix+"."+model.Approval{}.TableName()).
		Order("status, ctime desc").Where("user_id=?", uid)
	if len(beginTime) != 0 {
		db = db.Where("ctime>=?", beginTime)
	}
	if condition == model.GetApprovalApproving {
		db = db.Where("status=?", model.ApprovalStatWaiting)
	} else if condition == model.GetApprovalFinished {
		db = db.Where("status<>?", model.ApprovalStatWaiting)
	}
	e = db.Find(&alist).Error
	return
}

func GetTodoApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	sql := fmt.Sprintf(`select * from "%s".approval as t1 inner join "%s".approve_flow as t2
		on t1.no = t2.approval_no
		where t1.status=%d and t2.status=%d and t2.match_users like '%%%d_%%' `,
		prefix, prefix, model.ApprovalStatWaiting, model.ApprovalStatWaiting, uid)

	if len(params) != 0 && len(params[0]) != 0 {
		sql += fmt.Sprintf(`and t2.ctime>='%s' `, params[0])
	}
	sql += `order by t2.ctime desc`
	e = db.Raw(sql).Scan(&alist).Error
	return
}

func GetFinishedApprovalsToMe(prefix string, uid int, params ...string) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	//approve_flow 中user_id有值表示已经审批过
	sql := fmt.Sprintf(`select * from "%s".approval as t1 inner join "%s".approve_flow as t2
		on t1.no = t2.approval_no where t2.user_id=%d `, prefix, prefix, uid)
	if len(params) != 0 && len(params[0]) != 0 {
		sql += fmt.Sprintf(`and t2.ctime>='%s' `, params[0])
	}
	sql += `order by t2.ctime desc`
	e = db.Raw(sql).Scan(&alist).Error
	return
}
