package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func addFuncsToRole(prefix string, tx *gorm.DB, r *model.Role, fids []int) (e error) {
	for _, id := range fids {
		f := model.Func{}
		e = tx.Table(model.Public+model.Func{}.TableName()).First(&f, id).Error
		if e == gorm.ErrRecordNotFound || e != nil {
			return fmt.Errorf("func %d is not found", id)
		}
		rf := model.RoleFunc{
			RoleId: r.Id,
			FuncId: id,
			Ctime:  time.Now(),
		}
		e = tx.Table(prefix + rf.TableName()).Create(&rf).Error
		if e != nil {
			return
		}
	}
	return nil
}

func delFuncsOfRole(prefix string, tx *gorm.DB, r *model.Role) (e error) {
	e = tx.Table(prefix+model.RoleFunc{}.TableName()).
		Delete(model.RoleFunc{}, "role_id = ? ", r.Id).Error
	return
}

func AddRole(prefix string, r *model.Role, fids []int) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + r.TableName()).Create(r).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = addFuncsToRole(prefix, tx, r, fids)
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func UpdateRole(prefix string, r *model.Role, fids []int) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + r.TableName()).Model(r).
		Updates(model.Role{Name: r.Name, Desc: r.Desc}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = delFuncsOfRole(prefix, tx, r)
	if e != nil {
		tx.Rollback()
		return
	}
	e = addFuncsToRole(prefix, tx, r, fids)
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func DelRole(prefix string, rid int) (e error) {
	db := model.NewOrm()
	count := 0
	e = db.Table(prefix+model.UserRole{}.TableName()).
		Where("role_id = ?", rid).Count(&count).Error
	if e != nil || count != 0 {
		return fmt.Errorf("there some users in this role!%v", e)
	}
	e = db.Table(prefix + model.Role{}.TableName()).
		Delete(model.Role{Id: rid}).Error
	if e != nil {
		return
	}
	return nil
}

func AddUsersToRole(prefix string, rid int, uids []int) (e error) {
	db := model.NewOrm().Table(prefix + model.UserRole{}.TableName())
	ug := model.UserRole{}
	for _, uid := range uids {
		e = db.FirstOrCreate(&ug, &model.UserRole{UserId: uid, RoleId: rid}).Error
		if e != nil {
			return
		}
	}
	return nil
}

func DelUsersFromRole(prefix string, rid int, uids []int) (e error) {
	tx := model.NewOrm().Table(prefix + model.UserRole{}.TableName()).Begin()
	del := tx.Delete(&model.UserRole{}, "role_id = ? and user_id in (?)", rid, uids)
	if int(del.RowsAffected) != len(uids) {
		tx.Rollback()
		return errors.New("del failed,amount of users in this role is not match")
	} else if del.Error != nil {
		tx.Rollback()
		return del.Error
	} else {
		tx.Commit()
	}
	return nil
}
