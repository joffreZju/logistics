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
	e = model.NewOrm().Table(prefix + "." + model.Formtpl{}.TableName()).Find(&ftpls).Error
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
func GetApprocvaltplList(prefix string) (atpls []*model.Approvaltpl, e error) {
	db := model.NewOrm()
	atpls = []*model.Approvaltpl{}
	e = db.Table(prefix + "." + model.Approvaltpl{}.TableName()).Find(&atpls).Error
	if e != nil {
		return
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
				e = db.Table(prefix+"."+a.TableName()).Model(a).Update("status", model.ApprovalStatAccessed).Error
				if e != nil {
					beego.Error(e)
				}
				//todo 审批通过，新建消息，并发消息给审批发起人
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
		users, e := getApproverByRole(prefix, a, rid)
		if e == nil {
			//找到符合条件的审批人
			r := &model.Role{}
			e = db.Table(prefix+"."+r.TableName()).First(r, "id=?", rid).Error
			if e != nil {
				beego.Error(e)
			}
			var matchUsers string
			for _, v := range users {
				matchUsers += fmt.Sprintf("%d-", v.Id)
			}
			af := &model.ApproveFlow{
				ApprovalNo: a.No,
				MatchUsers: matchUsers,
				RoleId:     r.Id,
				RoleName:   r.Name,
				Status:     model.ApprovalStatWaiting,
			}
			//创建一步流程
			e = db.Table(prefix + "." + af.TableName()).Create(af).Error
			if e != nil {
				stop = true
				break
			}
			//更新审批单当前信息
			a.CurrentRole = rid
			e = db.Table(prefix + "." + a.TableName()).Model(a).Update(a).Error
			if e != nil {
				stop = true
				break
			}
			//todo 给所有users发需要审批的消息，带上af.Id
			_ = users
			return
		} else if e == gorm.ErrRecordNotFound && a.SkipBlankRole == model.SkipBlankRoleYes {
			//没有符合条件的审批人，跳过
			continue
		} else {
			//审批无法流转下去
			stop = true
			break
		}
	}
	if stop {
		beego.Error("审批单无法继续流转:", e)
		e = db.Table(prefix+"."+a.TableName()).Model(a).Update("status", model.ApprovalStatStop).Error
		if e != nil {
			beego.Error("尝试停止审批单:", e)
		}
	} else {
		//后面角色全部跳过，审批单完全通过
		e = db.Table(prefix+"."+a.TableName()).Model(a).Update("status", model.ApprovalStatAccessed).Error
		if e != nil {
			beego.Error(e)
		}
	}
	//todo 审批状态有更新，新建消息，并发消息给审批发起人
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
	for _, v := range strings.Split(g.Path, "-") {
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
	e = tx.Table(prefix+"."+af.TableName()).Where("id=? and status=?", af.Id, af.Status).Updates(af).Error
	if e != nil {
		tx.Rollback()
		return
	}
	if af.Status == model.ApprovalStatNotAccessed {
		a.Status = model.ApprovalStatNotAccessed
		e = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).Error
		if e != nil {
			tx.Rollback()
			return
		}
		//todo 审批状态有更新，新建消息，并发消息给审批发起人，审批单未通过
	} else {
		go nextStepOfApproval(prefix, a)
	}
	return tx.Commit().Error
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
	e = model.NewOrm().Table(prefix+"."+af.TableName()).Limit(1).Order("id desc").
		First(af, "approval_no=?", approvalNo).Error
	return
}

func GetApprovalsFromMe(prefix string, uid int) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	e = db.Table(prefix+"."+model.Approval{}.TableName()).
		Find(&alist, "user_id=?", uid).Error
	return
}

//func GetTodoApprovalsToMe(prefix string, uid int) (alist []*model.Approval, e error) {
//	db := model.NewOrm()
//	alist = []*model.Approval{}
//	e = db.Table(prefix + "." + model.Approval{}.TableName()).
//		Where("status=? and currentuser=?", model.ApprovalStatWaiting, uid).Find(&alist).Error
//	return
//}

func GetFinishedApprovalsToMe(prefix string, uid int) (alist []*model.Approval, e error) {
	db := model.NewOrm()
	alist = []*model.Approval{}
	sql := fmt.Sprintf(`select * from "%s".approval as t1 inner join "%s".approve_flow as t2
		on t1.no = t2.approval_no where t2.user_id=%d`, prefix, prefix, uid)
	e = db.Raw(sql).Scan(&alist).Error
	return
}
