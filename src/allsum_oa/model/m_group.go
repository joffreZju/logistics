package model

import (
	"time"
)

const (
	GroupTreeIsFuture = 1
)

type UserGroup struct {
	Id      int       `gorm:"AUTO_INCREMENT;primary_key"`
	UserId  int       `gorm:"not null"`
	GroupId int       `gorm:"not null"`
	Ctime   time.Time `gorm:"default:current_timestamp"`
}

func (UserGroup) TableName() string {
	return "user_group"
}

type Group struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	No        string `gorm:"unique;not null"`
	AdminId   int    `gorm:"not null"`
	CreatorId int    `gorm:"not null"`
	Desc      string
	AttrId    int       `gorm:""` //属性id
	Name      string    `gorm:"not null"`
	Pid       int       `gorm:"not null"` //父节点id
	Ctime     time.Time `gorm:"default:current_timestamp"`
	Utime     time.Time
	Path      string //`gorm:"not null"` 需要先插入记录再更新path
	//Level     int    `gorm:"not null"`
}

func (Group) TableName() string {
	return "group"
}

type Attribute struct {
	Id    int    `gorm:"primary_key;AUTO_INCREMENT"`
	No    string `gorm:"unique;not null"`
	Desc  string
	Name  string    `gorm:"not null"`
	Ctime time.Time `gorm:"default:current_timestamp"`
	Utime time.Time
}

func (Attribute) TableName() string {
	return "attribute"
}

type Operation struct {
	//todo
}

func (Operation) TableName() string {
	return "operation"
}

type GroupOperation struct {
	Id        int    `gorm:"primary_key;AUTO_INCREMENT"`
	Groups    string `gorm:"type:jsonb;not null"`
	EndTime   time.Time
	BeginTime time.Time
	IsFuture  int `gorm:"default:0"` // 0:历史组织树，1:还未生效的组织树
}

func (GroupOperation) TableName() string {
	return "group_operation"
}

type HistoryGroup struct {
	Pk        int       `gorm:"primary_key;AUTO_INCREMENT"`
	EndTime   time.Time `gorm:"not null"`
	Id        int       `gorm:"not null"`
	No        string    `gorm:"unique;not null"`
	AdminId   int       `gorm:"not null"`
	CreatorId int       `gorm:"not null"`
	Desc      string
	AttrId    int       `gorm:"not null"` //属性id
	Name      string    `gorm:"not null"`
	Pid       int       `gorm:"not null"` //父节点id
	Ctime     time.Time `gorm:"default:current_timestamp"`
	Utime     time.Time
	Path      string //`gorm:"not null"`
}

func (HistoryGroup) TableName() string {
	return "history_group"
}
