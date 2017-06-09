package model

import (
	"time"
)

type UserGroup struct {
	Id      int       `gorm:"primary_key;AUTO_INCREMENT"`
	UserId  int       `gorm:"not null"`
	GroupId int       `gorm:"not null"`
	Ctime   time.Time `gorm:"not null"`
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
	Ctime     time.Time `gorm:"not null"`
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
	Ctime time.Time `gorm:"not null"`
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
	Ctime     time.Time `gorm:"not null"`
	Utime     time.Time
	Path      string //`gorm:"not null"`
	//Level     int    `gorm:"not null"`
}

func (HistoryGroup) TableName() string {
	return "history_group"
}
