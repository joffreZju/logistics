package service

import (
	"allsum_oa/model"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func addFuncsToRole(tx *gorm.DB, r *model.Role, fids []int) (e error) {
	for _, id := range fids {
		f := model.Func{}
		e = tx.First(&f, id).Error
		if e == gorm.ErrRecordNotFound || e != nil {
			//tx.Rollback()
			return fmt.Errorf("func %d is not found", id)
		}
		rf := model.RoleFunc{
			RoleId: r.Id,
			FuncId: id,
			Ctime:  time.Now(),
		}
		e = tx.Create(&rf).Error
		if e != nil {
			//tx.Rollback()
			return
		}
	}
	return nil
}

func delFuncsOfRole(tx *gorm.DB, r *model.Role) (e error) {
	e = tx.Where("role_id = ? ", r.Id).Delete(model.RoleFunc{}).Error
	return
}

func AddRole(prefix string, r *model.Role, fids []int) (e error) {
	tx := model.NewOrm().Table(prefix + r.TableName()).Begin()
	e = tx.Create(r).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = addFuncsToRole(tx, r, fids)
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func UpdateRole(prefix string, r *model.Role, fids []int) (e error) {
	tx := model.NewOrm().Table(prefix + r.TableName()).Begin()
	e = tx.Model(r).Updates(model.Role{Name: r.Name, Desc: r.Desc}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = addFuncsToRole(tx, r, fids)
	if e != nil {
		tx.Rollback()
		return
	}
	e = delFuncsOfRole(tx, r)
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}
