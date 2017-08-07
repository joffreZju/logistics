package service

import (
	"allsum_oa/model"
	"github.com/astaxie/beego"
	"strings"
)

func GetAppVersionList() (vlist []*model.AppVersion, e error) {
	vlist = []*model.AppVersion{}
	e = model.NewOrm().Table(model.Public + "." + model.AppVersion{}.TableName()).Order("ctime desc").Find(&vlist).Error
	return
}

func AddAppVersion(app *model.AppVersion) (e error) {
	e = model.NewOrm().Table(model.Public + "." + app.TableName()).Create(app).Error
	return
}

func GetLatestAppVersion() (app *model.AppVersion, e error) {
	app = &model.AppVersion{}
	e = model.NewOrm().Table(model.Public + "." + app.TableName()).Order("id desc").Limit(1).Find(app).Error
	return
}

func UpdateDB() {
	//更新数据库中的_ 为 -
	tx := model.NewOrm().Begin()
	funcs := []*model.Function{}
	e := tx.Find(&funcs).Error
	defer beego.Info("error is:", e)
	if e != nil {
		return
	}
	for _, v := range funcs {
		v.Path = strings.Join(strings.Split(v.Path, "_"), "-")
		e = tx.Model(v).Updates(v).Error
		if e != nil {
			return
		}
	}

	comps := []model.Company{}
	e = tx.Find(&comps, &model.Company{Status: model.CompanyStatApproveAccessed}).Error
	if e != nil {
		return
	}
	for _, v := range comps {
		groups := []*model.Group{}
		e = tx.Table(v.No + "." + model.Group{}.TableName()).Find(&groups).Error
		if e != nil {
			return
		}
		for _, g := range groups {
			g.Path = strings.Join(strings.Split(g.Path, "_"), "-")
			e = tx.Table(v.No + "." + model.Group{}.TableName()).Model(g).Updates(g).Error
			if e != nil {
				return
			}
		}
	}
	tx.Commit()
}
