package service

import (
	"allsum_oa/model"
	"bytes"
	"common/lib/push"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/jinzhu/gorm"
	"net"
	"net/http"
	"time"
)

func CreateMsg(m *model.Message) (err error) {
	err = model.NewOrm().Create(m).Error
	return
}

func GetLatestMsg(company string, uid, maxId int) (msgs []*model.Message, err error) {
	msgs = []*model.Message{}
	err = model.NewOrm().Where("id > ? and user_id = ? and company_no = ?", maxId, uid, company).Find(&msgs).Error
	if err == gorm.ErrRecordNotFound {
		err = nil
	}
	return
}

func DelMsgById(msgId int) (err error) {
	m := &model.Message{Id: msgId}
	err = model.NewOrm().Delete(m).Error
	return nil
}

func DelMsgByType(company string, uid, tp int) (err error) {
	m := &model.Message{
		CompanyNo: company,
		UserId:    uid,
		MsgType:   tp,
	}
	err = model.NewOrm().Where(m).Delete(m).Error
	return
}

func SaveAndSendMsg(m *model.Message) (err error) {
	err = CreateMsg(m)
	if err != nil {
		beego.Error("save an message failed:", err)
		return
	}
	content, err := json.Marshal(m)
	if err != nil {
		return err
	}
	alias := fmt.Sprintf("%s_%d", m.CompanyNo, m.UserId)
	//go sendMsgToWeb(m)
	go push.JPushMsgByAlias([]string{alias}, string(content))

	return nil
}

func sendMsgToWeb(m *model.Message) {
	client := &http.Client{
		Transport: &http.Transport{
			Dial: func(netw, addr string) (net.Conn, error) {
				deadline := time.Now().Add(time.Second * 10)
				c, err := net.DialTimeout(netw, addr, 5*time.Second)
				if err != nil {
					return nil, err
				}
				c.SetDeadline(deadline)
				return c, nil
			},
		},
	}

	webpush := beego.AppConfig.String("webpush")
	body, err := json.Marshal(m)
	if err != nil {
		beego.Error("json marshal error:", err)
		return
	}
	request, err := http.NewRequest("POST", webpush, bytes.NewReader(body))
	if err != nil {
		beego.Error("http.NewRequest:", err)
		return
	}
	_, err = client.Do(request)
	if err != nil {
		beego.Error(err)
	}
}
