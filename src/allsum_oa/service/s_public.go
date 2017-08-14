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
	sql := fmt.Sprintf(
		`SELECT *
		FROM "public"."%s" AS t1 INNER JOIN "public"."%s" AS t2
			ON t1.no = t2.cno
		WHERE t2.user_id = ?`, model.Company{}.TableName(), model.UserCompany{}.TableName())
	comps = []model.Company{}
	e = db.Raw(sql, uid).Scan(&comps).Error
	return
}

func GetRolesOfUser(prefix string, uid int) (roles []model.Role, e error) {
	db := model.NewOrm()
	roles = []model.Role{}
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."%s" AS t1 INNER JOIN "%s"."%s" AS t2
			ON t1.id = t2.role_id
		WHERE t2.user_id = %d`, prefix, model.Role{}.TableName(), prefix, model.UserRole{}.TableName(), uid)
	e = db.Raw(sql).Scan(&roles).Error
	if e != nil {
		return
	}
	return
}

func GetGroupsOfUser(prefix string, uid int) (groups []model.Group, e error) {
	db := model.NewOrm()
	groups = []model.Group{}
	sql := fmt.Sprintf(
		`SELECT *
		FROM "%s"."%s" AS t1 INNER JOIN "%s"."%s" AS t2
			ON t1.id = t2.group_id
		WHERE t2.user_id = %d`, prefix, model.Group{}.TableName(), prefix, model.UserGroup{}.TableName(), uid)
	e = db.Raw(sql).Scan(&groups).Error
	if e != nil {
		return
	}
	return
}

func GetFuncIdsOfUser(prefix string, uid int) (functions []model.Function, e error) {
	db := model.NewOrm()
	rids := []int{}
	sql := fmt.Sprintf(`SELECT DISTINCT (role_id) FROM "%s"."%s" WHERE user_id = %d`,
		prefix, model.UserRole{}.TableName(), uid)
	e = db.Raw(sql).Pluck("role_id", &rids).Error
	if e != nil {
		return
	}
	functions = []model.Function{}
	sql = fmt.Sprintf(`
		SELECT DISTINCT(t2.*)
		FROM "%s"."%s" AS t1 INNER JOIN "public"."%s" AS t2
			ON t1.func_id = t2."id"
		WHERE role_id IN (?)
		ORDER BY t2.pid`, prefix, model.RoleFunc{}.TableName(), model.Function{}.TableName())
	e = db.Raw(sql, rids).Scan(&functions).Error
	return
}

func GetUserListOfCompany(prefix string) (users []*model.User, e error) {
	users = []*model.User{}
	e = model.NewOrm().Table(prefix + "." + model.User{}.TableName()).
		Order("id").Find(&users).Error
	if e != nil {
		return
	}
	for _, v := range users {
		v.Roles, e = GetRolesOfUser(prefix, v.Id)
		if e != nil {
			return
		}
		v.Groups, e = GetGroupsOfUser(prefix, v.Id)
		if e != nil {
			return
		}
	}
	return
}

func SearchUsersByName(prefix, uname string) (users []*model.User, e error) {
	users = []*model.User{}
	e = model.NewOrm().Table(prefix+"."+model.User{}.TableName()).
		Where("user_name like ?", "%"+uname+"%").Find(&users).Error
	if e != nil {
		return
	}
	for _, v := range users {
		v.Roles, e = GetRolesOfUser(prefix, v.Id)
		if e != nil {
			return
		}
		v.Groups, e = GetGroupsOfUser(prefix, v.Id)
		if e != nil {
			return
		}
	}
	return
}

func LockUser(prefix string, user *model.User) (e error) {
	tx := model.NewOrm().Begin()
	count := tx.Table(prefix+"."+user.TableName()).Model(user).
		Update("status", model.UserStatusLocked).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("lock user failed")
	}
	e = tx.Table(prefix+"."+model.UserRole{}.TableName()).
		Delete(&model.UserRole{}, &model.UserRole{UserId: user.Id}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	e = tx.Table(prefix+"."+model.UserGroup{}.TableName()).
		Delete(&model.UserGroup{}, &model.UserGroup{UserId: user.Id}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func GetUserListByUids(prefix string, uids []int) (users []*model.User, e error) {
	users = []*model.User{}
	e = model.NewOrm().Table(prefix+"."+model.User{}.TableName()).
		Find(&users, "id in (?)", uids).Error
	return
}

func GetCompanyList() (interface{}, error) {
	type CompanyDetail struct {
		model.Company
		CreateUser model.User
	}
	list := []*CompanyDetail{}
	e := model.NewOrm().Table(model.Public + "." + model.Company{}.TableName()).
		Order("status, ctime desc").Find(&list).Error
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	for _, v := range list {
		e = model.NewOrm().Table(model.Public+"."+model.User{}.TableName()).
			First(&v.CreateUser, v.Creator).Error
		if e != nil {
			beego.Error(e)
			return nil, e
		}
	}
	return list, nil
}

func createSchema(schema string) (e error) {
	sql := fmt.Sprintf(`create schema "%s"`, schema)
	e = model.NewOrm().Exec(sql).Error
	if e != nil && (strings.Contains(e.Error(), model.DBErrStrAlreadyExists) ||
		strings.Contains(e.Error(), "已经存在")) {
		return nil
	}
	return
}

func createCreatorRole(prefix string) (e error) {
	tx := model.NewOrm().Begin()
	r := &model.Role{
		Name:   "管理员",
		Descrp: "公司初始注册者",
	}
	e = tx.Table(prefix + "." + r.TableName()).Create(r).Error
	if e != nil {
		tx.Rollback()
		return
	}
	funcs := []*model.Function{}
	//todo 菜单是否可以分配，是否要在这里要体现出来，只把开放给其他公司的菜单全部分配给创始人？
	e = tx.Table(model.Public+"."+model.Function{}.TableName()).
		Find(&funcs, "pid=? or sys_id=?", 0, "oa").Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, v := range funcs {
		//if len(strings.Split(v.Path, "-")) > 2 {
		rf := &model.RoleFunc{
			RoleId: r.Id,
			FuncId: v.Id,
		}
		e = tx.Table(prefix + "." + rf.TableName()).Create(rf).Error
		if e != nil {
			tx.Rollback()
			return
		}
		//}
	}
	comp := &model.Company{}
	e = tx.Table(model.Public+"."+comp.TableName()).Find(comp, "no=?", prefix).Error
	if e != nil {
		tx.Rollback()
		return
	}
	ur := &model.UserRole{
		UserId: comp.Creator,
		RoleId: r.Id,
	}
	e = tx.Table(prefix + "." + ur.TableName()).Create(ur).Error
	if e != nil {
		tx.Rollback()
		return
	}
	return tx.Commit().Error
}

func AuditCompany(cno string, approver *model.User, status int, msg string) (err error) {
	tx := model.NewOrm().Begin()
	c := tx.Table(model.Public+"."+model.Company{}.TableName()).Model(&model.Company{}).Where("no=?", cno).
		Updates(&model.Company{
			Approver:     approver.Id,
			ApproverName: approver.UserName,
			Status:       status,
			ApproveMsg:   msg,
			ApproveTime:  time.Now()}).RowsAffected
	if c != 1 {
		err = errors.New("approve compony failed")
		tx.Rollback()
		return
	}
	//tx.Model(&model.Company{}).Where("no=?", cno).Updates(&model.Company{ApproveTime: time.Now()})
	if status == model.CompanyStatApproveAccessed {
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
		sql := fmt.Sprintf(`select user_id from "public"."%s" where cno=?`, model.UserCompany{}.TableName())
		err = tx.Raw(sql, cno).Pluck("user_id", &uids).Error
		if err != nil {
			beego.Error(err)
			tx.Rollback()
			return
		}
		users, err := GetUserListByUids("public", uids)
		if err != nil {
			beego.Error(err)
			tx.Rollback()
			return err
		} else {
			for _, u := range users {
				err = model.FirstOrCreateUser(cno, u)
				if err != nil {
					beego.Error(err)
					tx.Rollback()
					return err
				}
			}
		}
		//创建创始人角色并将其functions赋予公司的创始人
		err = createCreatorRole(cno)
		if err != nil {
			beego.Error(err)
			tx.Rollback()
			return err
		}
	}
	return tx.Commit().Error
}

func AddFunction(f *model.Function) (e error) {
	db := model.NewOrm()
	ffather := &model.Function{}
	e = db.Table(model.Public+"."+model.Function{}.TableName()).Find(ffather, f.Pid).Error
	if e != nil || strings.Count(ffather.Path, "-") > 5 {
		//菜单最多五级
		return errors.New("父节点选取有误，菜单最多五级")
	}
	tx := db.Begin()
	e = tx.Table(model.Public + "." + model.Function{}.TableName()).Create(f).Error
	if e != nil {
		tx.Rollback()
		return
	}
	f.Path = fmt.Sprintf("%s-%d", ffather.Path, f.Id)
	count := tx.Table(model.Public + "." + model.Function{}.TableName()).Model(f).Updates(f).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("add failed, zero rows affected")
	}
	return tx.Commit().Error
}

func UpdateFunction(f *model.Function) (e error) {
	tx := model.NewOrm().Begin()
	count := tx.Table(model.Public + "." + model.Function{}.TableName()).Model(f).Updates(f).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("update failed, zero rows affected")
	}
	return tx.Commit().Error
}

func DelFunction(fid int) (e error) {
	db := model.NewOrm()
	f := &model.Function{}
	e = db.Table(model.Public+"."+model.Function{}.TableName()).First(f, fid).Error
	if e != nil {
		return
	}
	var count int64 = 0
	e = db.Table(model.Public+"."+model.Function{}.TableName()).Where("pid=?", fid).Count(&count).Error
	if e != nil || count != 0 {
		return errors.New("该功能节点不能被删除")
	}
	tx := db.Begin()
	count = tx.Table(model.Public + "." + model.Function{}.TableName()).Delete(f).RowsAffected
	if count != 1 {
		tx.Rollback()
		return errors.New("del failed, zero rows affected")
	}
	return tx.Commit().Error
}
