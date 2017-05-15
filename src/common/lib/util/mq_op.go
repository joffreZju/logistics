package util

import (
	"crypto/hmac"
	"crypto/sha1"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"github.com/astaxie/beego"
	"io/ioutil"
	"net/http"
	"strings"
	"time"
)

var (
	//MQ_URL                       = beego.AppConfig.String("MQ_URL")
	//MQ_TOPIC_PRODUCER            = beego.AppConfig.String("MQ_TOPIC_PRODUCER")
	//MQ_TOPIC_CONSUMER            = beego.AppConfig.String("MQ_TOPIC_CONSUMER")
	//MQ_PRODUCER_ID               = beego.AppConfig.String("MQ_PRODUCER_ID")
	//MQ_CONSUMER_ID               = beego.AppConfig.String("MQ_CONSUMER_ID")
	ALI_ACCESS_KEY_ID     string = "LTAIysfw9MWnCZFk"
	ALI_ACCESS_KEY_SECRET string = "hAuLM27EkdVVxtfvbYHgq5XPDRvial"
)

//producer And consumer use MqSign
func MqSigh(signStr, accessKey string) string {
	mac := hmac.New(sha1.New, []byte(accessKey))
	mac.Write([]byte(signStr))
	s := base64.StdEncoding.EncodeToString([]byte(mac.Sum(nil)))
	strings.TrimRight(s, " ")
	return s
}

//define MqMsg for MQ consumer to unMarshal json from MQ
type MqMsg struct {
	Body      string
	MsgHandle string
	MsgId     string
}

func Producer(bodyStr string) error {
	Topic := beego.AppConfig.String("MQ_TOPIC_PRODUCER")
	URL := beego.AppConfig.String("MQ_URL")
	ProducerID := beego.AppConfig.String("MQ_PRODUCER_ID")
	newline := "\n"
	content := Md5Cal2String([]byte(bodyStr))
	date := fmt.Sprintf("%d", time.Now().UnixNano())[0:13]
	signStr := Topic + newline + ProducerID + newline + content + newline + date
	sign := MqSigh(signStr, ALI_ACCESS_KEY_SECRET)
	client := &http.Client{}
	req, err := http.NewRequest("POST", URL+"/message/?topic="+Topic+"&time="+date+"&tag=http"+"&key=http", strings.NewReader(bodyStr))
	if err != nil {
		beego.Error(err)
		return fmt.Errorf("MQ Producer error: %v", err)
	}

	req.Header.Set("Signature", sign)
	req.Header.Set("AccessKey", ALI_ACCESS_KEY_ID)
	req.Header.Set("ProducerID", ProducerID)
	req.Header.Set("Content-Type", "text/html;charset=UTF-8")

	resp, err := client.Do(req)
	if err != nil {
		beego.Error(err)
		return fmt.Errorf("MQ Producer error: %v", err)
	}
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("read MQ response body error: ", err)
		return err
	}
	var respMsg MqMsg
	json.Unmarshal(body, &respMsg)
	beego.Debug("MQ producer status", respMsg.MsgId, resp.Status)

	if resp.StatusCode == 201 {
		return nil
	} else {
		beego.Error(err)
		return fmt.Errorf("MQ Producer error: %v", resp.Status)
	}
}
