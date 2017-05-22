package service

import "allsum_oa/model"

func CreateAttr(pref string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(pref + a.TableName()).Create(a).Error
	return
}

func UpdateAttr(pref string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(pref+a.TableName()).
		Where("no = ?", a.No).Update("desc", "name", "utime").Error
	return
}
