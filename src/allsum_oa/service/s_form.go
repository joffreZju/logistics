package service

import (
	"allsum_oa/model"
	"errors"
)

func AddFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	e = model.NewOrm().Table(prefix + ftpl.TableName()).Create(ftpl).Error
	return
}

func UpdateFormtpl(prefix string, ftpl *model.Formtpl) (e error) {
	c := model.NewOrm().Table(prefix + ftpl.TableName()).Model(&model.Formtpl{No: ftpl.No}).
		Updates(ftpl).RowsAffected
	if c != 1 {
		e = errors.New("wrong formtpl no")
		return
	}
	return nil
}

func ControlFormtpl(prefix, no string, status int) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix+model.Formtpl{}.TableName()).
		Model(&model.Formtpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	if status == model.TplDisabled {
		e = tx.Table(prefix+model.Approvaltpl{}.TableName()).
			Where("formtpl_no = ?", no).Update("status", status).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}

func DelFormtpl(prefix, no string) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + model.Formtpl{}.TableName()).
		Delete(&model.Formtpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong formtpl no")
		return
	}
	e = tx.Table(prefix+model.Approvaltpl{}.TableName()).
		Where("formtpl_no = ?", no).Update("status", model.TplDisabled).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func AddApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	e = model.NewOrm().Table(prefix + atpl.TableName()).Create(atpl).Error
	return
}

func UpdateApprovaltpl(prefix string, atpl *model.Approvaltpl) (e error) {
	c := model.NewOrm().Table(prefix + atpl.TableName()).Model(&model.Approvaltpl{No: atpl.No}).
		Updates(atpl).RowsAffected
	if c != 1 {
		e = errors.New("wrong approvaltpl no")
		return
	}
	return nil
}

func ControlApprovaltpl(prefix, no string, status int) (e error) {
	c := model.NewOrm().Table(prefix+model.Approvaltpl{}.TableName()).
		Model(&model.Approvaltpl{No: no}).Update("status", status).RowsAffected
	if c != 1 {
		e = errors.New("wrong approvaltpl no")
		return
	}
	return nil
}

func DelApprovaltpl(prefix, no string) (e error) {
	c := model.NewOrm().Begin().Table(prefix + model.Approvaltpl{}.TableName()).
		Delete(&model.Approvaltpl{No: no}).RowsAffected
	if c != 1 {
		e = errors.New("wrong approvaltpl no")
		return
	}
	return nil
}

func AddApproval(prefix string, f *model.Form, a *model.Approval) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + f.TableName()).Create(f).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = tx.Table(prefix + a.TableName()).Create(a).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func UpdateApproval(prefix string, f *model.Form, a *model.Approval) (e error) {
	tx := model.NewOrm().Begin()
	aprvl := model.Approval{}
	e = tx.Table(prefix+aprvl.TableName()).First(&aprvl, a.No).Error
	if e != nil {
		return
	}
	if aprvl.Status != model.ApproveInit {
		e = errors.New("approval is already commited")
		return
	}
	c := tx.Table(prefix+f.TableName()).Where("no = ?", f.No).Updates(f).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong form no")
		return
	}
	c = tx.Table(prefix+a.TableName()).Where("no =?", a.No).Updates(a).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("wrong approval no")
		return
	}
	return tx.Commit().Error
}

func CancelApproval(prefix, no string) (e error) {
	tx := model.NewOrm().Table(prefix + model.Approval{}.TableName()).Begin()
	a := model.Approval{}
	e = tx.First(&a, no).Error
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
	e = tx.First(&a, aflow.ApprovalNo).Error
	if e != nil {
		return
	}
	status := -1
	for k, u := range a.UserFlow {
		if u == aflow.UserId {
			e = tx.Table(prefix + aflow.TableName()).Create(aflow).Error
			if e != nil {
				tx.Rollback()
				return
			}
			if aflow.Opinion == model.ApproveFlowAgree && k == len(a.UserFlow)-1 {
				//最后一位审批人同意
				status = model.ApproveAccessed
			} else if aflow.Opinion == model.ApproveFlowRefuse {
				//不同意
				status = model.ApproveNotAccessed
			}
			break
		}
	}
	if status != -1 {
		c := tx.Table(prefix+a.TableName()).
			Model(&a).Update("status", status).RowsAffected
		if c != 1 {
			tx.Rollback()
			e = errors.New("update approval status failed")
			return
		}
	}
	return tx.Commit().Error
}
