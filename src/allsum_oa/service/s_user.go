package service

import (
	"allsum_oa/model"
	"fmt"
)

func GetUserById(prefix string, uid int) (user *model.User, e error) {
	db := model.NewOrm()
	user = &model.User{}
	e = db.Table(prefix+"."+user.TableName()).First(user, uid).Error
	if e != nil {
		return
	}
	user.Roles, e = GetRolesOfUser(prefix, uid)
	if e != nil {
		return
	}
	user.Groups, e = GetGroupsOfUser(prefix, uid)
	if e != nil {
		return
	}
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
