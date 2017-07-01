package service

import (
	"allsum_oa/model"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"math"
	"time"
)

//十分钟扫描一次，如果设置的开始时间小于
func Ticker() {
	tick := time.Tick(time.Minute * 10)
	for {
		select {
		case <-tick:
			go ScanAllSchema(15)
		}
	}
}

func ScanAllSchema(interval float64) {
	db := model.NewOrm()
	comps := []model.Company{}
	e := db.Find(&comps, model.Company{Status: model.CompanyStatApproveAccessed}).Error
	if e != nil {
		return
	}
	for _, v := range comps {
		go handleGroupOperation(v.No, interval)
		go handleFormtpl(v.No, interval)
		go handleApprovaltpl(v.No, interval)
	}
}

func handleFormtpl(prefix string, interval float64) {
	ftpls := []*model.Formtpl{}
	db := model.NewOrm().Table(prefix + "." + model.Formtpl{}.TableName())
	e := db.Find(&ftpls, "status=?", model.TplInit).Error
	if e != nil {
		beego.Error(e)
		return
	}
	for _, v := range ftpls {
		if math.Abs(v.BeginTime.Sub(time.Now()).Minutes()) < interval {
			e = db.Model(v).Update("status", model.TplAbled).Error
			if e != nil {
				beego.Error(e)
			}
		}
	}
}
func handleApprovaltpl(prefix string, interval float64) {
	atpls := []*model.Approvaltpl{}
	db := model.NewOrm().Table(prefix + "." + model.Approvaltpl{}.TableName())
	e := db.Find(&atpls, "status=?", model.TplInit).Error
	if e != nil {
		beego.Error(e)
		return
	}
	for _, v := range atpls {
		if math.Abs(v.BeginTime.Sub(time.Now()).Minutes()) < interval {
			e = db.Model(v).Update("status", model.TplAbled).Error
			if e != nil {
				beego.Error(e)
			}
		}
	}
}

func handleGroupOperation(prefix string, interval float64) {
	op := &model.GroupOperation{}
	db := model.NewOrm().Table(prefix + "." + op.TableName())
	e := db.Find(op, "is_future=?", model.GroupOpStatFuture).Error
	if e != nil {
		if e != gorm.ErrRecordNotFound {
			beego.Error(e)
		}
		return
	}
	if math.Abs(op.BeginTime.Sub(time.Now()).Minutes()) > interval {
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
	e = tx.Table(prefix + "." + op.TableName()).Model(op).Updates(op)
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
