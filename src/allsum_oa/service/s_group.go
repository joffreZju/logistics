package service

import (
	"allsum_oa/model"
	"errors"
	"fmt"
	"github.com/jinzhu/gorm"
	"strings"
	"time"
)

func CreateAttr(prefix string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(prefix + a.TableName()).Create(a).Error
	return
}

func UpdateAttr(prefix string, a *model.Attribute) (e error) {
	e = model.NewOrm().Table(prefix+a.TableName()).
		Where("no = ?", a.No).Update("desc", "name", "utime").Error
	return
}

func GetGroup(prefix string, id int) (g *model.Group, e error) {
	g = new(model.Group)
	e = model.NewOrm().Table(prefix+g.TableName()).First(g, id).Error
	return
}

func AddGroup(prefix string, ng *model.Group, sonIds []int) (e error) {
	tx := model.NewOrm().Table(prefix + ng.TableName()).Begin()
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
	ng.Path = father.Path + "-" + fmt.Sprintf("%d", ng.Id)
	e = tx.Model(ng).Update("path").Error
	if e != nil {
		tx.Rollback()
		return
	}
	if len(sonIds) != 0 {
		sons := []*model.Group{}
		e = tx.Where("id in (?)", sonIds).Find(&sons).Error
		if e != nil {
			return
		}
		for _, v := range sons {
			if v.Pid != ng.Pid {
				return errors.New("所选节点不在同一父节点下")
			}
		}
		for _, v := range sons {
			newPath := ng.Path + "-" + fmt.Sprintf("%d", v.Id)
			//找到v的所有子孙节点，修改其path
			children := []*model.Group{}
			e = tx.Where("path like ?", v.Path+"-%").Find(&children).Error
			if e == gorm.ErrRecordNotFound {
				continue
			} else {
				tx.Rollback()
				return
			}
			for _, ch := range children {
				ch.Path = strings.Replace(ch.Path, v.Path, newPath, 1)
				ch.Utime = time.Now()
				e = tx.Model(ch).Update("path", "utime")
				if e != nil {
					tx.Rollback()
					return
				}
			}
			v.Pid = ng.Id
			v.Path = newPath
			v.Utime = time.Now()
			e = tx.Model(v).Update("pid", "path", "utime").Error
			if e != nil {
				tx.Rollback()
				return
			}
		}
	}
	return tx.Commit().Error
}

func MergeGroups(prefix string, ng *model.Group, oldIds []int) (e error) {
	tx := model.NewOrm().Table(prefix + ng.TableName()).Begin()
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
	ng.Path = father.Path + "-" + fmt.Sprintf("%d", ng.Id)
	e = tx.Model(ng).Update("path").Error
	if e != nil {
		tx.Rollback()
		return
	}
	olds := []*model.Group{}
	e = tx.Where("id in (?)", oldIds).Find(&olds).Error
	if e != nil {
		return
	}
	for _, v := range olds {
		if v.Pid != ng.Pid {
			return errors.New("所选节点不在同一父节点下")
		}
	}
	for _, v := range olds {
		//找到v的所有子孙节点，修改其path
		children := []*model.Group{}
		e = tx.Where("path like ?", v.Path+"-%").Find(&children).Error
		if e == gorm.ErrRecordNotFound {
			continue
		} else {
			tx.Rollback()
			return
		}
		for _, ch := range children {
			ch.Path = strings.Replace(ch.Path, v.Path, ng.Path, 1)
			ch.Pid = ng.Id
			ch.Utime = time.Now()
			e = tx.Model(ch).Update("path", "utime")
			if e != nil {
				tx.Rollback()
				return
			}
		}
		//删除旧节点
		e = tx.Delete(v).Error
		if e != nil {
			tx.Rollback()
			return
		}
	}
	return tx.Commit().Error
}

func MoveGroup(prefix string, gid, newPid int) (e error) {
	g := new(model.Group)
	tx := model.NewOrm().Table(prefix + g.TableName()).Begin()
	e = tx.First(g, gid).Error
	if e != nil {
		return
	}
	gNewFather := new(model.Group)
	e = tx.First(gNewFather, newPid).Error
	if e != nil {
		return
	}
	gNewPath := gNewFather.Path + "-" + fmt.Sprintf("%d", gid)
	//找到v的所有子孙节点，修改其path
	children := []*model.Group{}
	e = tx.Where("path like ?", g.Path+"-%").Find(&children).Error
	if e != nil && e != gorm.ErrRecordNotFound {
		return
	}
	for _, ch := range children {
		ch.Path = strings.Replace(ch.Path, g.Path, gNewPath, 1)
		ch.Utime = time.Now()
		e = tx.Model(ch).Update("path", "utime")
		if e != nil {
			tx.Rollback()
			return
		}
	}
	g.Path = gNewPath
	g.Utime = time.Now()
	tx.Model(g).Update("path", "utime")
	return tx.Commit().Error
}

func DelGroup(prefix string, gid int) (e error) {
	g := new(model.Group)
	tx := model.NewOrm().Table(prefix + g.TableName()).Begin()
	e = tx.First(g, gid).Error
	if e != nil {
		return
	}
	gFather := new(model.Group)
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
		e = tx.Model(ch).Update("path", "utime")
		if e != nil {
			tx.Rollback()
			return
		}
	}
	e = tx.Delete(&g).Error
	if e != nil {
		tx.Rollback()
	}
	return tx.Commit().Error
}

func EditGroup(prefix, newName string, gid int) (e error) {
	tx := model.NewOrm().Table(prefix + model.Group{}.TableName())
	r := tx.Model(&model.Group{}).Where("id = ?", gid).Update("name", newName).RowsAffected
	if r == 0 {
		e = errors.New("change group name failed")
		return
	}
	return nil
}
