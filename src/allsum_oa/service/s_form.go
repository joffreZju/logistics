package service

import (
	"allsum_oa/model"
	"errors"
)

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
	tx := model.NewOrm().Begin()
	count := 0
	if status == model.TplDisabled {
		e = tx.Table(prefix+"."+model.Approvaltpl{}.TableName()).
			Where("formtpl_no=?", no).Count(&count).Error
		if e != nil {
			return
		} else if count != 0 {
			return errors.New("some approvaltpl are using this formtpl")
		}
	}
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
	tx := model.NewOrm().Begin()
	count := 0
	e = tx.Table(prefix+"."+model.Approvaltpl{}.TableName()).
		Where("formtpl_no=?", no).Count(&count).Error
	if e != nil {
		return
	} else if count != 0 {
		return errors.New("some approvaltpl are using this formtpl")
	}
	c := tx.Table(prefix + "." + model.Formtpl{}.TableName()).
		Delete(&model.Formtpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	return tx.Commit().Error
}

func GetApprocvaltplList(prefix string) (atpls []*model.Approvaltpl, e error) {
	db := model.NewOrm()
	atpls = []*model.Approvaltpl{}
	e = db.Table(prefix + "." + model.Approvaltpl{}.TableName()).Find(&atpls).Error
	if e != nil {
		return
	}
	for _, v := range atpls {
		e = db.Table(prefix+"."+model.Formtpl{}.TableName()).
			Find(&v.FormtplContent, "no=?", v.FormtplNo).Error
		if e != nil {
			return
		}
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

func AddApproval(prefix string, a *model.Approval) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + "." + a.FormContent.TableName()).Create(&(a.FormContent)).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = tx.Table(prefix + "." + a.TableName()).Create(a).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func UpdateApproval(prefix string, a *model.Approval) (e error) {
	tx := model.NewOrm().Begin()
	aprvl := model.Approval{}
	e = tx.Table(prefix+"."+aprvl.TableName()).First(&aprvl, "no=?", a.No).Error
	if e != nil {
		return
	}
	if aprvl.Status != model.ApproveDraft {
		e = errors.New("approval is already commited")
		return
	}
	c := tx.Table(prefix + "." + a.FormContent.TableName()).
		Model(&(a.FormContent)).Updates(&(a.FormContent)).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong form no")
		return
	}
	c = tx.Table(prefix + "." + a.TableName()).Model(a).Updates(a).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approval no")
		return
	}
	return tx.Commit().Error
}

func CancelApproval(prefix, no string) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.Approval{}.TableName()).Begin()
	a := model.Approval{}
	e = tx.First(&a, "no=?", no).Error
	if e != nil {
		return
	}
	if a.Status == model.ApproveAccessed || a.Status == model.ApproveNotAccessed {
		e = errors.New("approval has been finished")
		return
	}
	c := tx.Model(&a).Update("status", model.ApproveCanceled).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("approval no is wrong")
		return
	}
	return tx.Commit().Error
}

func Approve(prefix string, aflow *model.ApproveFlow) (e error) {
	tx := model.NewOrm().Begin()
	a := model.Approval{}
	//检查该审批单当前状态
	e = tx.Table(prefix+"."+a.TableName()).First(&a, "no=?", aflow.ApprovalNo).Error
	if e != nil {
		return
	}
	if a.Status != model.Approving {
		return errors.New("approval has been finished")
	}
	//检查是否审批过
	count := 0
	e = tx.Table(prefix+"."+aflow.TableName()).
		Where("approval_no=? and user_id=?", aflow.ApprovalNo, aflow.UserId).Count(&count).Error
	if e != nil {
		return
	} else if count != 0 {
		return errors.New("您已经审批过了")
	}
	//审批，修改审批单状态
	newStatus := -1
	for k, u := range a.UserFlow {
		if u == aflow.UserId {
			e = tx.Table(prefix + "." + aflow.TableName()).Create(aflow).Error
			if e != nil {
				tx.Rollback()
				return
			}
			if aflow.Opinion == model.ApproveOpinionRefuse {
				//不同意
				newStatus = model.ApproveNotAccessed
			} else if aflow.Opinion == model.ApproveOpinionAgree && k == len(a.UserFlow)-1 {
				//最后一位审批人同意
				newStatus = model.ApproveAccessed
			} else {
				//中间审批人同意
				//todo 给下一位审批人发推送
				//userNext := a.UserFlow[k + 1]
			}
			break
		}
	}
	if newStatus != -1 {
		//todo 给发起审批的人推送
		c := tx.Table(prefix+"."+a.TableName()).
			Model(&a).Update("status", newStatus).RowsAffected
		if c != 1 {
			tx.Rollback()
			e = errors.New("update approval status failed")
			return
		}
	}
	return tx.Commit().Error
}
