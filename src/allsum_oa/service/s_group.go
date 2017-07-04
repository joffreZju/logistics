package service

import (
	"allsum_oa/model"
	"encoding/json"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

func GetAttrList(prefix string) (al []model.Attribute, e error) {
	al = []model.Attribute{}
	e = model.NewOrm().Table(prefix + "." + model.Attribute{}.TableName()).Find(&al).Error
	return
}

func AddAttr(prefix string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(prefix + "." + a.TableName()).Create(a).Error
	return
}

func UpdateAttr(prefix string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(prefix+"."+a.TableName()).
		Where("id=?", a.Id).
		Updates(model.Attribute{Name: a.Name, Descrp: a.Descrp, Utime: a.Utime}).Error
	return
}

func DelAttr(prefix string, a *model.Attribute) (e error) {
	db := model.NewOrm()
	count := 0
	e = db.Table(prefix+"."+model.Group{}.TableName()).Where("attr_id=?", a.Id).Count(&count).Error
	if e != nil || count != 0 {
		return errors.New("仍有组织绑定此属性")
	}
	e = db.Table(prefix + "." + a.TableName()).Delete(a).Error
	return
}

func GetGroup(prefix string, id int) (g *model.Group, e error) {
	g = new(model.Group)
	e = model.NewOrm().Table(prefix+"."+g.TableName()).First(g, id).Error
	return
}

func CheckFutureGroupOperation(prefix string) (e error) {
	count := 0
	e = model.NewOrm().Table(prefix+"."+model.GroupOperation{}.TableName()).
		Where("status=? and begin_time>?", model.GroupOpStatFuture, time.Now()).Count(&count).Error
	if e != nil || count != 0 {
		return errors.New("当前有未生效的修改")
	}
	return nil
}

func handleTX(prefix, desc string, beginTime time.Time, tx *gorm.DB) (e error) {
	//确定是立即生效还是未来定时生效
	isFuture := true
	if beginTime.Sub(time.Now()).Nanoseconds() <= 0 {
		isFuture = false
	}
	//创建组织树操作记录
	groups := []*model.Group{}
	e = tx.Table(prefix + "." + model.Group{}.TableName()).Find(&groups).Error
	if e != nil {
		tx.Rollback()
		return
	}
	b, e := json.Marshal(groups)
	if e != nil {
		tx.Rollback()
		return e
	}
	op := &model.GroupOperation{
		Descrp:    desc,
		BeginTime: beginTime,
		Groups:    string(b),
	}
	if isFuture {
		op.Status = model.GroupOpStatFuture
	}
	e = model.NewOrm().Table(prefix + "." + op.TableName()).Create(op).Error //这里不能用tx
	if e != nil {
		tx.Rollback()
		return
	}
	//提交或回滚事务
	if isFuture {
		return tx.Rollback().Error
	} else {
		return tx.Commit().Error
	}
}

func AddRootGroup(prefix, desc string, beginTime time.Time, ng *model.Group) (e error) {
	tx := model.NewOrm().Table(prefix + "." + ng.TableName()).Begin()
	e = tx.Create(ng).Error
	if e != nil {
		tx.Rollback()
		return
	}
	ng.Path = fmt.Sprintf("%d", ng.Id)
	e = tx.Where("id=?", ng.Id).Updates(&ng).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//return tx.Commit().Error
	return handleTX(prefix, desc, beginTime, tx)
}

func AddGroup(prefix, desc string, beginTime time.Time, ng *model.Group, sonIds []int) (e error) {
	tx := model.NewOrm().Table(prefix + "." + ng.TableName()).Begin()
	father := &model.Group{}
	e = tx.First(&father, ng.Pid).Error
	if e != nil {
		return
	}
	e = tx.Create(ng).Error
	if e != nil {
		tx.Rollback()
		return
	}
	ng.Path = father.Path + fmt.Sprintf("_%d", ng.Id)
	e = tx.Where("id=?", ng.Id).Update("path", ng.Path).Error
	if e != nil {
		tx.Rollback()
		return
	}
	if len(sonIds) != 0 {
		sons := []*model.Group{}
		e = tx.Find(&sons, "id in (?)", sonIds).Error
		if e != nil {
			return
		}
		for _, v := range sons {
			if v.Pid != ng.Pid {
				return errors.New("所选节点不在同一父节点下")
			}
		}
		for _, v := range sons {
			newPath := ng.Path + fmt.Sprintf("_%d", v.Id)
			//找到v的所有子孙节点
			children := []*model.Group{}
			e = tx.Find(&children, "path like ?", v.Path+"_%").Error
			if e != nil && e != gorm.ErrRecordNotFound {
				tx.Rollback()
				return
			}
			//修改v的子孙节点的path
			for _, ch := range children {
				ch.Path = strings.Replace(ch.Path, v.Path, newPath, 1)
				ch.Utime = time.Now()
				e = tx.Where("id=?", ch.Id).Updates(model.Group{Path: ch.Path, Utime: ch.Utime}).Error
				if e != nil {
					tx.Rollback()
					return
				}
			}
			//修改v自己的父亲和path
			v.Pid = ng.Id
			v.Path = newPath
			v.Utime = time.Now()
			e = tx.Where("id=?", v.Id).Updates(model.Group{Pid: v.Pid, Path: v.Path, Utime: v.Utime}).Error
			if e != nil {
				tx.Rollback()
				return
			}
		}
	}
	//return tx.Commit().Error
	return handleTX(prefix, desc, beginTime, tx)
}

func MergeGroups(prefix, desc string, beginTime time.Time, ng *model.Group, oldIds []int) (e error) {
	groupTb := prefix + "." + ng.TableName()
	userGroupTb := prefix + "." + model.UserGroup{}.TableName()
	tx := model.NewOrm().Begin()
	father := &model.Group{}
	e = tx.Table(groupTb).First(&father, ng.Pid).Error
	if e != nil {
		return
	}
	//插入新节点
	e = tx.Table(groupTb).Create(ng).Error
	if e != nil {
		tx.Rollback()
		return
	}
	ng.Path = father.Path + fmt.Sprintf("_%d", ng.Id)
	e = tx.Table(groupTb).Where("id=?", ng.Id).Updates(ng).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//找到所有被合并的节点
	olds := []*model.Group{}
	e = tx.Table(groupTb).Where("id in (?)", oldIds).Find(&olds).Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, v := range olds {
		if v.Pid != ng.Pid {
			tx.Rollback()
			return errors.New("所选节点不在同一父节点下")
		}
	}
	for _, old := range olds {
		//找到每个old的所有子孙节点
		children := []*model.Group{}
		e = tx.Table(groupTb).Where("path like ?", old.Path+"_%").Find(&children).Error
		if e != nil && e != gorm.ErrRecordNotFound {
			tx.Rollback()
			return
		}
		//修改old所有子孙节点的Pid和path
		for _, ch := range children {
			ch.Path = strings.Replace(ch.Path, old.Path, ng.Path, 1)
			ch.Pid = ng.Id
			ch.Utime = time.Now()
			e = tx.Table(groupTb).Model(ch).Updates(ch).Error
			if e != nil {
				tx.Rollback()
				return
			}
		}
		//删除旧节点
		e = tx.Table(groupTb).Delete(old).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	//修改被合并节点的user关系
	e = tx.Table(userGroupTb).Where("group_id in (?)", oldIds).Update("group_id", ng.Id).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//return tx.Commit().Error
	return handleTX(prefix, desc, beginTime, tx)
}

func MoveGroup(prefix, desc string, beginTime time.Time, gid, newPid int) (e error) {
	g, gNewFather := new(model.Group), new(model.Group)
	tx := model.NewOrm().Table(prefix + "." + g.TableName()).Begin()
	e = tx.First(g, gid).Error
	if e != nil {
		return
	}
	e = tx.First(gNewFather, newPid).Error
	if e != nil {
		return
	}
	gNewPath := gNewFather.Path + fmt.Sprintf("_%d", gid)
	//找到gid的所有子孙节点，修改其path
	children := []*model.Group{}
	e = tx.Find(&children, "path like ?", g.Path+"_%").Error
	if e != nil && e != gorm.ErrRecordNotFound {
		return
	}
	for _, ch := range children {
		ch.Path = strings.Replace(ch.Path, g.Path, gNewPath, 1)
		ch.Utime = time.Now()
		e = tx.Model(ch).Updates(model.Group{Path: ch.Path, Utime: ch.Utime}).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	//修改g自己的path
	g.Path = gNewPath
	g.Pid = gNewFather.Id
	g.Utime = time.Now()
	e = tx.Model(g).Updates(g).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//return tx.Commit().Error
	return handleTX(prefix, desc, beginTime, tx)
}

func DelGroup(prefix, desc string, beginTime time.Time, gid int) (e error) {
	db := model.NewOrm()
	count := 0
	e = db.Table(prefix+"."+model.UserGroup{}.TableName()).
		Where("group_id=?", gid).Count(&count).Error
	if e != nil || count != 0 {
		return fmt.Errorf("there are some users in this group! %v", e)
	}
	g, gFather := new(model.Group), new(model.Group)
	tx := db.Table(prefix + "." + g.TableName()).Begin()
	e = tx.First(g, gid).Error
	if e != nil {
		return
	}
	e = tx.First(gFather, g.Pid).Error
	if e != nil {
		return
	}
	//找到g的所有子孙节点，修改其path
	children := []*model.Group{}
	e = tx.Where("path like ?", g.Path+"_%").Find(&children).Error
	if e != nil && e != gorm.ErrRecordNotFound {
		return
	}
	for _, ch := range children {
		ch.Path = strings.Replace(ch.Path, g.Path, gFather.Path, 1)
		if ch.Pid == g.Id {
			//g的直接子节点
			ch.Pid = gFather.Id
		}
		ch.Utime = time.Now()
		e = tx.Model(ch).Updates(ch).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	e = tx.Delete(&g).Error
	if e != nil {
		tx.Rollback()
		return
	}
	//return tx.Commit().Error
	return handleTX(prefix, desc, beginTime, tx)
}

func UpdateGroup(prefix, desc string, beginTime time.Time, g *model.Group) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.Group{}.TableName())
	e = tx.Where("id = ?", g.Id).Updates(g).Error
	return handleTX(prefix, desc, beginTime, tx)
}

func GetGroupOpList(prefix string, limit int) (ops []*model.GroupOperation, e error) {
	//limit = -1 所有记录
	ops = []*model.GroupOperation{}
	e = model.NewOrm().Table(prefix + "." + model.GroupOperation{}.TableName()).
		Select([]string{"id", "begin_time", "descrp", "status"}).Limit(limit).Find(&ops).Error
	return
}

func SearchGroupOpsByTime(prefix string, begin, end time.Time) (ops []*model.GroupOperation, e error) {
	ops = []*model.GroupOperation{}
	e = model.NewOrm().Table(prefix+"."+model.GroupOperation{}.TableName()).
		Select([]string{"id", "begin_time", "descrp", "status"}).
		Where("begin_time between ? and ?", begin, end).Find(&ops).Error
	return
}

func GetGroupOpDetail(prefix string, id int) (groups []*model.Group, e error) {
	groups = []*model.Group{}
	op := &model.GroupOperation{}
	e = model.NewOrm().Table(prefix+"."+model.GroupOperation{}.TableName()).
		Find(op, "id=?", id).Error
	if e != nil {
		return
	}
	e = json.Unmarshal([]byte(op.Groups), &groups)
	return
}

func CancelGroupOp(prefix string, beginTime time.Time) (e error) {
	op := &model.GroupOperation{
		BeginTime: beginTime,
		Status:    model.GroupOpStatFuture,
	}
	e = model.NewOrm().Table(prefix + "." + op.TableName()).
		Where(op).Delete(op).Error
	return
}

func GetGroupList(prefix string) (gs []*model.Group, e error) {
	gs = []*model.Group{}
	e = model.NewOrm().Table(prefix + "." + model.Group{}.TableName()).Find(&gs).Error
	return
}

func GetUsersOfGroup(prefix string, gid int) (users []*model.User, e error) {
	users = []*model.User{}
	sql := fmt.Sprintf(`select * from "%s"."allsum_user" as t1 inner join "%s".user_group as t2
		on t1.id=t2.user_id where t2.group_id=%d`, prefix, prefix, gid)
	e = model.NewOrm().Raw(sql).Scan(&users).Error
	return
}

func AddUsersToGroup(prefix string, gid int, uids []int) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.UserGroup{}.TableName())
	for _, uid := range uids {
		ug := &model.UserGroup{UserId: uid, GroupId: gid}
		e = db.FirstOrCreate(ug, ug).Error
		if e != nil {
			return
		}
	}
	return nil
}

func AddUserToGroups(prefix string, gids []int, uid int) (e error) {
	db := model.NewOrm().Table(prefix + "." + model.UserGroup{}.TableName())
	for _, gid := range gids {
		ug := &model.UserGroup{UserId: uid, GroupId: gid}
		e = db.FirstOrCreate(ug, ug).Error
		if e != nil {
			return
		}
	}
	return nil
}

func DelUsersFromGroup(prefix string, gid int, uids []int) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.UserGroup{}.TableName()).Begin()
	del := tx.Delete(&model.UserGroup{}, "group_id = ? and user_id in (?)", gid, uids)
	if int(del.RowsAffected) != len(uids) {
		tx.Rollback()
		return errors.New("del failed, amount of users in this group is not match")
	} else if del.Error != nil {
		tx.Rollback()
		return del.Error
	}
	return tx.Commit().Error
}

func UpdateGroupssOfUser(prefix string, newGids []int, uid int) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.UserGroup{}.TableName()).Begin()
	e = tx.Where("user_id=?", uid).Delete(&model.UserGroup{}).Error
	if e != nil {
		tx.Rollback()
		return
	}
	for _, gid := range newGids {
		ug := &model.UserGroup{UserId: uid, GroupId: gid}
		e = tx.FirstOrCreate(ug, ug).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}
