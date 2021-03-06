package service

import (
	"allsum_oa/model"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"math"
	"time"
)

const (
	Interval    = 60
	MaxDistance = 90
)

//每Interval扫描一次，如果设置的开始时间小于MaxDistance就执行操作
func Ticker() {
	tick := time.Tick(time.Minute * Interval)
	for {
		select {
		case <-tick:
			go ScanAllSchema(MaxDistance)
		}
	}
}

func ScanAllSchema(MaxDistance float64) {
	db := model.NewOrm()
	comps := []model.Company{}
	e := db.Table(model.Public+"."+model.Company{}.TableName()).
		Find(&comps, model.Company{Status: model.CompanyStatApproveAccessed}).Error
	if e != nil {
		return
	}
	for _, v := range comps {
		go handleGroupOperation(v.No, MaxDistance)
		go handleFormtpl(v.No, MaxDistance)
		go handleApprovaltpl(v.No, MaxDistance)
	}
}

//处理定时生效的表单模板
func handleFormtpl(prefix string, MaxDistance float64) {
	ftpls := []*model.Formtpl{}
	db := model.NewOrm().Table(prefix + "." + model.Formtpl{}.TableName())
	e := db.Find(&ftpls, "status=?", model.TplInit).Error
	if e != nil {
		beego.Error(e)
		return
	}
	for _, v := range ftpls {
		if math.Abs(v.BeginTime.Sub(time.Now()).Minutes()) <= MaxDistance {
			e = db.Model(v).Update("status", model.TplAbled).Error
			if e != nil {
				beego.Error(e)
			}
		}
	}
}

//处理定时生效的审批单模板
func handleApprovaltpl(prefix string, MaxDistance float64) {
	atpls := []*model.Approvaltpl{}
	db := model.NewOrm().Table(prefix + "." + model.Approvaltpl{}.TableName())
	e := db.Find(&atpls, "status=?", model.TplInit).Error
	if e != nil {
		beego.Error(e)
		return
	}
	for _, v := range atpls {
		if math.Abs(v.BeginTime.Sub(time.Now()).Minutes()) <= MaxDistance {
			e = db.Model(v).Update("status", model.TplAbled).Error
			if e != nil {
				beego.Error(e)
			}
		}
	}
}

//处理定时生效的组织树修改操作
func handleGroupOperation(prefix string, MaxDistance float64) {
	op := &model.GroupOperation{}
	db := model.NewOrm().Table(prefix + "." + op.TableName())
	e := db.Find(op, "status=?", model.GroupOpStatFuture).Error
	if e != nil {
		if e != gorm.ErrRecordNotFound {
			beego.Error(e)
		}
		return
	}
	if math.Abs(op.BeginTime.Sub(time.Now()).Minutes()) > MaxDistance {
		return
	}
	//从json解析新组织树
	newGroups := []*model.Group{}
	e = json.Unmarshal([]byte(op.Groups), &newGroups)
	if e != nil {
		beego.Error("执行定时任务失败:", e)
		return
	}
	//从group表中删除旧组织树
	tx := model.NewOrm().Table(prefix + "." + model.Group{}.TableName()).Begin()
	e = tx.Delete(&model.Group{}).Error
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	//向group中插入新组织树
	for _, v := range newGroups {
		e = tx.Create(v).Error
		if e != nil {
			beego.Error(e)
			tx.Rollback()
			return
		}
	}
	op.Status = model.GroupOpStatHistory
	e = tx.Table(prefix + "." + op.TableName()).Model(op).Updates(op).Error
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	e = tx.Commit().Error
	if e != nil {
		beego.Error(e)
	}
}
