package service

import (
	"allsum_oa/model"
	"bytes"
	"common/lib/push"
	"encoding/json"
	"net"
	"net/http"
	"strconv"
	"time"

	"github.com/astaxie/beego"
)

func SendMsg(m *model.Message) {
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
		beego.Error("json marshal error:%v", err)
		return
	}
	request, err := http.NewRequest("POST", webpush, bytes.NewReader(body))
	if err != nil {
		beego.Error("http.NewRequest:", err)
		return
	}
	_, err = client.Do(request)
	if err != nil {
		beego.Warn(err)
	}
}

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
	go SendMsg(m)
	go push.JPushCommonMsg([]string{strconv.Itoa(uid)}, content, map[string]interface{}{})

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
		return nil, err
	}
	return
}
