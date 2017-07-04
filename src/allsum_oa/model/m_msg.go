package model

import "time"

const (
	MsgAppTypeAll = iota + 1
	MsgAppTypeWeb
	MsgAppTypeAndroid
	MsgAppTypeIPhone
)

const (
	MsgTypeSystem = iota + 1
	MsgTypeNeedApprove
)

//for users
type Message struct {
	Id        int `gorm:"primary_key;auto_increment"`
	Title     string
	MsgType   int
	AppType   int
	Content   string
	UserId    int
	CompanyNo string
	Ctime     time.Time `gorm:"default:current_timestamp"`
}

func (Message) TableName() string {
	return "allsum_message"
}
