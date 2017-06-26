package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/astaxie/beego"
	"strings"
	"time"
)

func GetUserByTel(prefix, tel string) (user *model.User, e error) {
	db := model.NewOrm()
	user = new(model.User)
	e = db.Table(prefix+"."+user.TableName()).First(user, model.User{Tel: tel}).Error
	if e != nil {
		return
	}
	uid := user.Id
	user.Companys, e = GetCompanysOfUser(uid)
	if e != nil {
		return
	}
	if prefix == "public" || len(prefix) == 0 {
		return user, nil
	}
	//获取schema下面的用户信息
	user.Roles, e = GetRolesOfUser(prefix, uid)
	if e != nil {
		return
	}
	user.Groups, e = GetGroupsOfUser(prefix, uid)
	if e != nil {
		return
	}
	user.Funcs, e = GetFuncIdsOfUser(prefix, uid)
	if e != nil {
		return
	}
	return
}

func GetUserById(prefix string, uid int) (user *model.User, e error) {
	db := model.NewOrm()
	user = &model.User{}
	e = db.Table(prefix+"."+user.TableName()).First(user, uid).Error
	if e != nil {
		return
	}

	user.Companys, e = GetCompanysOfUser(uid)
	if e != nil {
		return
	}
	if prefix == "public" || len(prefix) == 0 {
		return user, nil
	}
	user.Roles, e = GetRolesOfUser(prefix, uid)
	if e != nil {
		return
	}
	user.Groups, e = GetGroupsOfUser(prefix, uid)
	if e != nil {
		return
	}
	user.Funcs, e = GetFuncIdsOfUser(prefix, uid)
	if e != nil {
		return
	}
	return
}

func GetCompanysOfUser(uid int) (comps []model.Company, e error) {
	db := model.NewOrm()
	sql := `select * from allsum_company as t1 inner join allsum_user_company as t2
		on t1.no = t2.cno
		where t2.user_id = ?`
	comps = []model.Company{}
	e = db.Raw(sql, uid).Scan(&comps).Error
	return
}

func GetRolesOfUser(prefix string, uid int) (roles []model.Role, e error) {
	db := model.NewOrm()
	roles = []model.Role{}
	sql := fmt.Sprintf(`select * from "%s".role as t1 inner join "%s".user_role as t2
		on t1.id = t2.role_id where t2.user_id = %d`, prefix, prefix, uid)
	e = db.Raw(sql).Scan(&roles).Error
	if e != nil {
		return
	}
	return
}

func GetGroupsOfUser(prefix string, uid int) (groups []model.Group, e error) {
	db := model.NewOrm()
	groups = []model.Group{}
	sql := fmt.Sprintf(`select * from "%s".group as t1 inner join "%s".user_group as t2
		on t1.id = t2.group_id where t2.user_id = %d`, prefix, prefix, uid)
	e = db.Raw(sql).Scan(&groups).Error
	if e != nil {
		return
	}
	return
}

func GetFuncIdsOfUser(prefix string, uid int) (fids []int, e error) {
	db := model.NewOrm()
	rids := []int{}
	sql := fmt.Sprintf(`select distinct(role_id) from "%s".user_role where user_id = %d`, prefix, uid)
	e = db.Raw(sql).Pluck("role_id", &rids).Error
	//e = db.Table(prefix + "." + model.UserRole{}.TableName()).Where("user_id=?", uid).Pluck("role_id", &rids).Error
	if e != nil {
		return
	}
	fids = []int{}
	sql = fmt.Sprintf(`select t1.func_id from "%s".role_func as t1 INNER JOIN "public"."function" as t2
		on t1.func_id = t2."id"
		where role_id in (?)
		ORDER BY t2.pid`, prefix)
	//sql = fmt.Sprintf(`select distinct(func_id) from "%s".role_func where role_id in (?)`, prefix)
	e = db.Raw(sql, rids).Pluck("func_id", &fids).Error
	if e != nil {
		return
	}
	return
}

func createSchema(schema string) (e error) {
	sql := fmt.Sprintf(`create schema "%s"`, schema)
	e = model.NewOrm().Exec(sql).Error
	if e != nil && (strings.Contains(e.Error(), "already exists") ||
		strings.Contains(e.Error(), "已经存在")) {
		return nil
	}
	return
}

func GetUserList(prefix string, uids []int) (users []model.User, e error) {
	users = []model.User{}
	e = model.NewOrm().Table(prefix+"."+model.User{}.TableName()).
		Find(&users, "id in (?)", uids).Error
	return
}

func AuditCompany(cno string, approverId int, status int, msg string) (err error) {
	tx := model.NewOrm().Begin()
	c := tx.Model(&model.Company{}).Where("no=?", cno).
		Updates(&model.Company{
			Approver:   approverId,
			Status:     status,
			ApproveMsg: msg}).RowsAffected
	if c != 1 {
		err = errors.New("approve compony failed")
		tx.Rollback()
		return
	}
	tx.Model(&model.Company{}).Where("no=?", cno).Updates(&model.Company{ApproveTime: time.Now()})
	if status == model.CompanyApproveAccessed {
		//创建schema，直接提交
		err = createSchema(cno)
		if err != nil {
			tx.Rollback()
			return
		}
		//自动建表
		err = model.InitSchemaModel(cno)
		if err != nil {
			tx.Rollback()
			return
		}
		//迁移用户
		uids := []int{}
		sql := fmt.Sprint(`select user_id from "public".allsum_user_company where cno=?`)
		err = tx.Raw(sql, cno).Pluck("user_id", &uids).Error
		if err != nil {
			beego.Error(err)
			tx.Rollback()
			return
		}
		beego.Info(uids)
		users, err := GetUserList("public", uids)
		if err != nil {
			beego.Error(err)
			tx.Rollback()
			return err
		} else {
			beego.Info(users)
			for _, u := range users {
				err = model.FirstOrCreateUser(cno, &u)
				if err != nil {
					beego.Error(err)
					tx.Rollback()
					return err
				}
			}
		}
	}
	return tx.Commit().Error
}
