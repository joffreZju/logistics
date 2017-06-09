package service

import (
	"allsum_oa/model"
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
		Updates(model.Attribute{Name: a.Name, Desc: a.Desc, Utime: a.Utime}).Error
	return
}

func GetGroup(prefix string, id int) (g *model.Group, e error) {
	g = new(model.Group)
	e = model.NewOrm().Table(prefix+g.TableName()).First(g, id).Error
	return
}

func AddRootGroup(prefix string, ng *model.Group) (e error) {
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
	return tx.Commit().Error
}

func AddGroup(prefix string, ng *model.Group, sonIds []int) (e error) {
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
	ng.Path = father.Path + fmt.Sprintf("-%d", ng.Id)
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
			newPath := ng.Path + fmt.Sprintf("-%d", v.Id)
			//找到v的所有子孙节点
			children := []*model.Group{}
			e = tx.Find(&children, "path like ?", v.Path+"-%").Error
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
	return tx.Commit().Error
}

func MergeGroups(prefix string, ng *model.Group, oldIds []int) (e error) {
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
	ng.Path = father.Path + fmt.Sprintf("-%d", ng.Id)
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
		e = tx.Table(groupTb).Where("path like ?", old.Path+"-%").Find(&children).Error
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
	return tx.Commit().Error
}

func MoveGroup(prefix string, gid, newPid int) (e error) {
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
	gNewPath := gNewFather.Path + fmt.Sprintf("-%d", gid)
	//找到gid的所有子孙节点，修改其path
	children := []*model.Group{}
	e = tx.Find(&children, "path like ?", g.Path+"-%").Error
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
	return tx.Commit().Error
}

func DelGroup(prefix string, gid int) (e error) {
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
	e = tx.Where("path like ?", g.Path+"-%").Find(&children).Error
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
	return tx.Commit().Error
}

func UpdateGroup(prefix string, g *model.Group) (e error) {
	tx := model.NewOrm().Table(prefix + "." + model.Group{}.TableName())
	e = tx.Where("id = ?", g.Id).Updates(g).Error
	return
}

func GetGroupList(prefix string) (gs []*model.Group, e error) {
	gs = []*model.Group{}
	e = model.NewOrm().Table(prefix + "." + model.Group{}.TableName()).Find(&gs).Error
	return
}

func AddUsersToGroup(prefix string, gid int, uids []int) (e error) {
	db := model.NewOrm().Table(prefix + model.UserGroup{}.TableName())
	ug := model.UserGroup{}
	for _, uid := range uids {
		e = db.FirstOrCreate(&ug, &model.UserGroup{UserId: uid, GroupId: gid}).Error
		if e != nil {
			return
		}
	}
	return nil
}

func DelUsersFromGroup(prefix string, gid int, uids []int) (e error) {
	tx := model.NewOrm().Table(prefix + model.UserGroup{}.TableName()).Begin()
	del := tx.Delete(&model.UserGroup{}, "group_id = ? and user_id in (?)", gid, uids)
	if int(del.RowsAffected) != len(uids) {
		tx.Rollback()
		return errors.New("del failed, amount of users in this group is not match")
	} else if del.Error != nil {
		tx.Rollback()
		return del.Error
	} else {
		tx.Commit()
	}
	return nil
}
