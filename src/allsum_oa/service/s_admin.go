package service

import "allsum_oa/model"

func GetAppVersionList() (vlist []*model.AppVersion, e error) {
	vlist = []*model.AppVersion{}
	e = model.NewOrm().Order("ctime desc").Find(&vlist).Error
	return
}

func AddAppVersion(app *model.AppVersion) (e error) {
	e = model.NewOrm().Create(app).Error
	return
}

func GetLatestAppVersion() (app *model.AppVersion, e error) {
	app = &model.AppVersion{}
	e = model.NewOrm().Order("id desc").Limit(1).Find(app).Error
	return
}
