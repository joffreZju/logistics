package model

import "time"

type Role struct {
	Id    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"not null"`
	Desc  string
	Ctime time.Time `gorm:"not null"`
}

func (Role) TableName() string {
	return "role"
}

type RoleFunc struct {
	Id     int       `gorm:"primary_key;AUTO_INCREMENT"`
	RoleId int       `gorm:"not null"`
	FuncId int       `gorm:"not null"`
	Ctime  time.Time `gorm:"not null"`
}

func (RoleFunc) TableName() string {
	return "role_func"
}

type Function struct {
	Id    int    `gorm:"primary_key;AUTO_INCREMENT"`
	Name  string `gorm:"not null"`
	Desc  string
	Pid   int       `gorm:"not null"`
	Ctime time.Time `gorm:"not null"`
	Path  string    `gorm:"not null"`
	//Level int       `gorm:"not null"`
}

func (Function) TableName() string {
	return "function"
}

type UserRole struct {
	Id     int       `gorm:"primary_key;AUTO_INCREMENT"`
	RoleId int       `gorm:"not null"`
	UserId int       `gorm:"not null"`
	Ctime  time.Time `gorm:"not null"`
}

func (UserRole) TableName() string {
	return "user_role"
}
