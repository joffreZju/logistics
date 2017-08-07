package service

import "allsum_oa/model"

type ApiService struct {
}

func (*ApiService) GetSchemaList() (schemas []string, e error) {
	e = model.NewOrm().Table(model.Public+"."+model.Company{}.TableName()).
		Where(&model.Company{Status: model.CompanyStatApproveAccessed}).Pluck("no", &schemas).Error
	return
}
