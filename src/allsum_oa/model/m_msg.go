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
	MsgTypeApprove
)

//for users
type Message struct {
	Id        int `gorm:"primary_key;auto_increment"`
	Title     string
	MsgType   int
	AppType   int
	Content   JsonMap `gorm:"default:null;type:jsonb"`
	UserId    int
	CompanyNo string
	Ctime     time.Time `gorm:"default:current_timestamp"`
}

func (Message) TableName() string {
	return "oa_message"
}
