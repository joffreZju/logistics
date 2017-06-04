package model

import (
	"time"
)

type User struct {
	Id  int    `gorm:"primary_key" ` // 用户id,继承自public
	No  string `gorm:"size:64"`
	Tel string `gorm:"unique_index;size:15;not null" json:",omitempty"`
	//Password   string    `json:"-"`                                           // 密码
	UserName   string    `gorm:"size:64" json:",omitempty"`
	Icon       string    `gorm:"size:64" json:",omitempty"`
	Descp      string    `gorm:"" json:",omitempty"`
	Gender     int8      `gorm:"default:1" json:",omitempty"`
	Address    string    `gorm:"size:64" json:",omitempty"`
	CreateTime time.Time `gorm:"default:current_timestamp" json:",omitempty"`
	LoginTime  time.Time `gorm:"timestamp" json:",omitempty"`
	Mail       string    `gorm:"size:64" json:",omitempty"`
	Status     int       `gorm:"not null" json:",omitempty"`
	UserType   int       `gorm:"default:1" json:",omitempty"` //1 普通用户
	// Companys []Company `orm:"-" json:",omitempty"`                          // 用户的所在组织
	Roles  []Role  `gorm:"-"`
	Groups []Group `gorm:"-"`
}

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
	AttrId    int       `gorm:"not null"` //属性id
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
