package model

import "time"

//for users
type Message struct {
	Id      int `gorm:primary_key;auto_increment`
	Title   string
	Mtype   int
	Content string
	Ctime   time.Time `gorm:"default:current_timestamp"`
	UserId  int
}

func (Message) TableName() string {
	return "allsum_message"
}

func InsertMessage(m *Message) (err error) {
	err = NewOrm().Create(m).Error
	return
}

func DeleteMessage(id int) (err error) {
	m := &Message{Id: id}
	err = NewOrm().Delete(m).Error
	return
}

func GetLatestMessage(id int, uid int) (msgs []Message, err error) {
	err = NewOrm().Where("id > ? and user_id = ?", id, uid).Find(&msgs).Error
	return
}

func DeleteMessageByType(tp, uid int) (err error) {
	err = NewOrm().Where("mtype = ? and user_id = ?", tp, uid).Delete(Message{}).Error
	return
}
