package push

import (
	"encoding/json"
	"fmt"
	"github.com/GiterLab/aliyun-sms-go-sdk/sms"
	"github.com/astaxie/beego"
	jpushclient "github.com/ylywyn/jpush-api-go-client"
	"time"
)

//define sms template code
const (
	ALI_ACCESS_KEY_ID       = "LTAIysfw9MWnCZFk"
	ALI_ACCESS_KEY_SECRET   = "hAuLM27EkdVVxtfvbYHgq5XPDRvial"
	SMS_SIGN_NAME           = "壹算科技"
	SMS_TEMPLATE_WEB        = "SMS_58265055"
	SMS_TEMPLATE_WHEN_ERROR = "SMS_63875806"
	JPUSH_DEVKEY            = "eb3143316015c58253ea734a"
	JPUSH_DEVSECRET         = "9289bfeb4c50cc0b424e77a1"
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

func JPushMsgByAlias(usersAlias []string, msgContent string, noticeInfo ...map[string]interface{}) {
	payload := jpushclient.NewPushPayLoad()
	//Platform
	var pf jpushclient.Platform
	pf.Add(jpushclient.ANDROID)

	//Audience
	var ad jpushclient.Audience
	if len(usersAlias) == 0 || usersAlias == nil {
		ad.All()
	} else {
		ad.SetAlias(usersAlias)
	}

	//Message
	var msg jpushclient.Message
	msg.Content = msgContent

	//Notice
	var notice jpushclient.Notice
	if len(noticeInfo) != 0 {
		notice.SetAlert("您有一条新消息")
		payload.SetNotice(&notice)
	}

	payload.SetPlatform(&pf)
	payload.SetAudience(&ad)
	payload.SetMessage(&msg)
	bytes, _ := payload.ToBytes()

	t := time.Now().Unix()
	fmt.Printf("极光推送,%d,%s\n", t, string(bytes))
	//push
	c := jpushclient.NewPushClient(JPUSH_DEVSECRET, JPUSH_DEVKEY)
	ret, err := c.Send(bytes)
	if err != nil {
		fmt.Printf("推送失败,%d,%s\n", t, err.Error())
	} else {
		fmt.Printf("推送成功,%d,%s\n", t, ret)
	}
}
