package model

import "time"

type Role struct {
	Id     int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name   string `gorm:"not null"`
	Descrp string
	Ctime  time.Time  `gorm:"default:current_timestamp"`
	Funcs  []Function `gorm:"-" json:",omitempty"`
}

func (Role) TableName() string {
	return "role"
}

type RoleFunc struct {
	Id     int       `gorm:"primary_key;AUTO_INCREMENT"`
	RoleId int       `gorm:"not null"`
	FuncId int       `gorm:"not null"`
	Ctime  time.Time `gorm:"default:current_timestamp"`
}

func (RoleFunc) TableName() string {
	return "role_func"
}

type UserRole struct {
	Id     int       `gorm:"primary_key;AUTO_INCREMENT"`
	RoleId int       `gorm:"not null"`
	UserId int       `gorm:"not null"`
	Ctime  time.Time `gorm:"default:current_timestamp"`
}

func (UserRole) TableName() string {
	return "user_role"
}
