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
		Updates(*ftpl).RowsAffected
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
		e = errors.New("formtpl no does not exists")
		return
	}
	e = tx.Table(prefix+model.Approvaltpl{}.TableName()).
		Model(&model.Approvaltpl{FormtplNo: no}).Update("status", status).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return nil
}

func DelFormtpl(prefix, no string) (e error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(prefix + model.Formtpl{}.TableName()).
		Delete(&model.Formtpl{No: no}).RowsAffected
	if c != 1 {
		tx.Rollback()
		e = errors.New("formtpl no does not exists")
		return
	}
	e = tx.Table(prefix+model.Approvaltpl{}.TableName()).
		Model(&model.Approvaltpl{FormtplNo: no}).Update("status", model.Disabled).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return nil
}
