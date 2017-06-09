package service

import (
	"allsum_oa/model"

	"github.com/astaxie/beego"
)

func SaveOneMsg(title, content string, tp, uid int) error {
	m := new(model.Message)
	m.Title = title
	m.Mtype = tp
	m.UserId = uid
	m.Content = content
	err := model.InsertMessage(m)
	if err != nil {
		beego.Error("save an message failed:%v", err)
		return err
	}
	return nil
}

func DelMsg(msgId int) (err error) {
	err = model.DeleteMessage(msgId)
	if err != nil {
		beego.Error("delete message:%d failed:%v", msgId, err)
		return err
	}
	return nil
}

func DelMsgBatch(tp, uid int) (err error) {
	err = model.DeleteMessageByType(tp, uid)
	if err != nil {
		beego.Error("delete message by type:%d,%d failed:%v", tp, uid, err)
		return err
	}
	return nil
}

func GetMsgBatch(maxId int, uid int) (msgs []model.Message, err error) {
	msgs, err = model.GetLatestMessage(maxId, uid)
	if err != nil {
		beego.Error("user:%d get latest messages failed:%v", uid, err)
		return nil,err
	}
	return
}