package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"time"
)

func GetRoleList(prefix string) (roles []*model.Role, e error) {
	roles = []*model.Role{}
	e = model.NewOrm().Table(prefix + "." + model.Role{}.TableName()).Order("id").Find(&roles).Error
	return
}

func GetRolesDetail(prefix string) (roles []*model.Role, e error) {
	db := model.NewOrm()
	roles = []*model.Role{}
	e = db.Table(prefix + "." + model.Role{}.TableName()).Order("id").Find(&roles).Error
	if e != nil {
		return
	}
	sql := fmt.Sprintf(`SELECT * from "public"."%s" as t1 inner join "%s"."%s" as t2
		on t1."id" = t2.func_id
		where t2.role_id = ?`, model.Function{}.TableName(), prefix, model.RoleFunc{}.TableName())
	for _, v := range roles {
		e = db.Raw(sql, v.Id).Scan(&v.Funcs).Error
		if e != nil {
			return
		}
	}
	return
}

func addFuncsToRole(prefix string, tx *gorm.DB, r *model.Role, fids []int) (e error) {
	for _, id := range fids {
		f := model.Function{}
		e = tx.Table(model.Public+"."+model.Function{}.TableName()).First(&f, id).Error
		if e != nil {
			return fmt.Errorf("func %d is not found", id)
		}
		rf := model.RoleFunc{
			RoleId: r.Id,
			FuncId: id,
			Ctime:  time.Now(),
		}
		e = tx.Table(prefix + "." + rf.TableName()).Create(&rf).Error
		if e != nil {
			return
		}
	}
	return nil
}

func delAllFuncsOfRole(prefix string, tx *gorm.DB, r *model.Role) (e error) {
	e = tx.Table(prefix+"."+model.RoleFunc{}.TableName()).
		Delete(model.RoleFunc{}, "role_id=?", r.Id).Error
	return
}

func AddRole(prefix string, r *model.Role, fids []int) (e error) {
	tx := model.NewOrm().Begin()
	e = tx.Table(prefix + "." + r.TableName()).Create(r).Error
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
	e = tx.Table(prefix + "." + r.TableName()).Model(r).
		Updates(r).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = delAllFuncsOfRole(prefix, tx, r)
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
	tx := model.NewOrm().Begin()
	count := 0
	e = model.NewOrm().Table(prefix+"."+model.UserRole{}.TableName()).
		Where("role_id = ?", rid).Count(&count).Error
	if e != nil || count != 0 {
		return fmt.Errorf("there some users in this role!%v", e)
	}
	e = tx.Table(prefix + "." + model.Role{}.TableName()).
		Delete(&model.Role{Id: rid}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = tx.Table(prefix+"."+model.RoleFunc{}.TableName()).
		Delete(&model.RoleFunc{}, &model.RoleFunc{RoleId: rid}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func GetUsersOfRole(prefix string, rid int) (users []*model.User, e error) {
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."%s" AS t1 INNER JOIN "%s"."%s" AS t2
			ON t1.id = t2.user_id
		WHERE t2.role_id = %d
		ORDER BY t1.id`, prefix, model.User{}.TableName(), prefix, model.UserRole{}.TableName(), rid)
	users = []*model.User{}
	e = model.NewOrm().Raw(sql).Scan(&users).Error
	return
}

func AddUsersToRole(prefix string, rid int, uids []int) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.UserRole{}.TableName())
	for _, uid := range uids {
		ur := &model.UserRole{UserId: uid, RoleId: rid}
		e = db.FirstOrCreate(ur, ur).Error
		if e != nil {
			return
		}
	}
	return nil
}

func AddUserToRoles(prefix string, rids []int, uid int) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.UserRole{}.TableName())
	for _, rid := range rids {
		ur := &model.UserRole{UserId: uid, RoleId: rid}
		e = db.FirstOrCreate(ur, ur).Error
		if e != nil {
			return
		}
	}
	return nil
}

func DelUsersFromRole(prefix string, rid int, uids []int) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.UserRole{}.TableName()).Begin()
	del := tx.Delete(&model.UserRole{}, "role_id=? and user_id in (?)", rid, uids)
	if int(del.RowsAffected) != len(uids) {
		tx.Rollback()
		return errors.New("del failed,amount of users in this role is not match")
	} else if del.Error != nil {
		tx.Rollback()
		return del.Error
	} else {
		return tx.Commit().Error
	}
}

func UpdateRolesOfUser(prefix string, newRids []int, uid int) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.UserRole{}.TableName()).Begin()
	e = tx.Where("user_id=?", uid).Delete(&model.UserRole{}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, rid := range newRids {
		ur := &model.UserRole{UserId: uid, RoleId: rid}
		e = tx.FirstOrCreate(ur, ur).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}
