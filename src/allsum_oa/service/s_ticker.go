package service

import (
	"allsum_oa/model"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"math"
	"time"
)

//十分钟扫描一次，如果设置的开始时间小于
func Ticker() {
	tick := time.Tick(time.Minute * 10)
	for {
		select {
		case <-tick:
			fmt.Println(time.Now())
			go ScanGroupOperation(15)
		}
	}
}

func ScanGroupOperation(interval float64) {
	db := model.NewOrm()
	comps := []model.Company{}
	e := db.Find(&comps, model.Company{Status: model.CompanyApproveAccessed}).Error
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
	txx := model.NewOrm().Table(prefix + "." + op.TableName()).Begin()
	e := txx.Find(op, "is_future=?", model.GroupTreeIsFuture).Error
	if e != nil {
		beego.Error(e)
		txx.Rollback()
		return
	}
	if math.Abs(op.BeginTime.Sub(time.Now()).Minutes()) > interval {
		txx.Rollback()
		return
	}
	//删掉任务
	if txx.Delete(op).RowsAffected != 1 {
		txx.Rollback()
	} else {
		txx.Commit()
	}
	tx := model.NewOrm().Begin()
	newGroups, oldGroups := []*model.Group{}, []*model.Group{}
	//从json解析新组织树
	e = json.Unmarshal([]byte(op.Groups), &newGroups)
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	//拿到旧组织树并转为json
	e = tx.Table(prefix + "." + model.Group{}.TableName()).Find(&oldGroups).Error
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	oldStr, e := json.Marshal(&oldGroups)
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	//保存旧组织树
	hisOp := &model.GroupOperation{
		Groups:  string(oldStr),
		EndTime: op.BeginTime,
	}
	e = tx.Table(prefix + "." + hisOp.TableName()).Create(hisOp).Error
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	//从group表中删除旧组织树
	e = tx.Table(prefix + "." + model.Group{}.TableName()).Delete(&model.Group{}).Error
	if e != nil {
		beego.Error(e)
		tx.Rollback()
		return
	}
	//向group中插入新组织树
	for _, v := range newGroups {
		e = tx.Table(prefix + "." + model.Group{}.TableName()).Create(v).Error
		if e != nil {
			beego.Error(e)
			tx.Rollback()
			return
		}
	}
	e = tx.Commit().Error
	if e != nil {
		beego.Error(e)
	}
}
