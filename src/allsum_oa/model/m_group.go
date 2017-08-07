package model

import (
	"time"
)

const (
	GroupOpStatHistory = iota + 1
	GroupOpStatFuture
)

type UserGroup struct {
	Id      int       `gorm:"AUTO_INCREMENT;primary_key"`
	UserId  int       `gorm:"not null"`
	GroupId int       `gorm:"not null"`
	Ctime   time.Time `gorm:"default:current_timestamp"`
}

func (UserGroup) TableName() string {
	return "oa_user_group"
}

type Group struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	No        string `gorm:"unique;not null"`
	AdminId   int    `gorm:"not null"`
	CreatorId int    `gorm:"not null"`
	Descrp    string
	AttrId    int       `gorm:""` //属性id
	Name      string    `gorm:"not null"`
	Pid       int       `gorm:"not null"` //父节点id
	Ctime     time.Time `gorm:"default:current_timestamp"`
	Utime     time.Time
	Path      string //`gorm:"not null"` 需要先插入记录再更新path 1-2-3
}

func (Group) TableName() string {
	return "oa_group"
}

type Attribute struct {
	Id     int    `gorm:"primary_key;AUTO_INCREMENT"`
	No     string `gorm:"unique;not null"`
	Descrp string
	Name   string    `gorm:"not null"`
	Ctime  time.Time `gorm:"default:current_timestamp"`
	Utime  time.Time
}

func (Attribute) TableName() string {
	return "oa_attribute"
}

type GroupOperation struct {
	Id        int `gorm:"primary_key;AUTO_INCREMENT"`
	Descrp    string
	Groups    string `gorm:"type:jsonb;not null" json:",omitempty"`
	BeginTime time.Time
	Status    int `gorm:"default:1"` // 1:历史组织树，2:还未生效的组织树
}

func (GroupOperation) TableName() string {
	return "oa_group_operation"
}
