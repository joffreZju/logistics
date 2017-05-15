package push

import (
	"encoding/json"
	"fmt"
	"github.com/GiterLab/aliyun-sms-go-sdk/sms"
	"github.com/astaxie/beego"
)

func Init() (err error) {
	//alidayu.AppKey = "LTAIwxFn7egYfvra"
	//alidayu.AppKey = "23441572"
	//alidayu.AppSecret = "nBfpqo4StRZv9JreRsLQpFaZKKUT1h"
	//alidayu.AppSecret = "8d890383b59c43ccdeb50fe9d0074087"
	return
}

//
//func SendMsgWithDayuToMobile(mobile, code string, product string) bool {
//	sucess, response := alidayu.SendSMS(mobile, "登录验证", "SMS_13400735", fmt.Sprintf(`{"code":"%v","product":"%v"}`, code, product))
//	beego.Debug(response)
//	return sucess
//}
//
//func SendSMSWithDayu(mobile, name, tplId string, params map[string]string) bool {
//	args, _ := json.Marshal(params)
//	sucess, response := alidayu.SendSMS(mobile, name, tplId, string(args))
//	beego.Debug(response)
//	fmt.Println(response, sucess)
//	return sucess
//}

//define sms template code
const (
	ALI_ACCESS_KEY_ID       = "LTAIysfw9MWnCZFk"
	ALI_ACCESS_KEY_SECRET   = "hAuLM27EkdVVxtfvbYHgq5XPDRvial"
	SMS_SIGN_NAME           = "壹算科技"
	SMS_TEMPLATE_WEB        = "SMS_58265055"
	SMS_TEMPLATE_WHEN_ERROR = "SMS_63875806"
)

func SendSmsCodeToMobile(mobile, code string) error {
	param := make(map[string]string)
	param["smscode"] = code
	c := sms.New(ALI_ACCESS_KEY_ID, ALI_ACCESS_KEY_SECRET)
	str, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("send smscode failed,%v", err)
	}
	e, err := c.SendOne(mobile, SMS_SIGN_NAME, SMS_TEMPLATE_WEB, string(str))
	if err != nil {
		return fmt.Errorf("send sms failed,%v,%v", err, e.Error())
	}
	return nil
}

func SendErrorSms(mobile, content string) error {
	param := make(map[string]string)
	param["content"] = content
	c := sms.New(ALI_ACCESS_KEY_ID, ALI_ACCESS_KEY_SECRET)
	str, err := json.Marshal(param)
	if err != nil {
		return fmt.Errorf("send sms failed,%v", err)
	}
	e, err := c.SendOne(mobile, SMS_SIGN_NAME, SMS_TEMPLATE_WHEN_ERROR, string(str))
	if err != nil {
		return fmt.Errorf("send sms failed,%v,%v", err, e.Error())
	}
	beego.Info("calculate failed:", content)
	return nil
}
