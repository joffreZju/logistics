package push

import (
	"encoding/json"
	"testing"
)

//func init() {
//	Init()
//}
//

//func TestPushDayu(t *testing.T) {
//	result := SendSMSWithDayu("15158134537", "壹算科技", "SMS_37830073", map[string]string{
//		"cp_id":    "2222222",
//		"cp_code":  "12234555",
//		"cp_title": "测试优惠券"})
//	fmt.Println(result)
//}

func TestJPush(t *testing.T) {
	alias := []string{"C0607145711618_21"}
	msg := &struct {
		CompanyNo string
		UserId    int
		MsgType   int
		Title     string
		Content   map[string]string
	}{
		"C0607145711618",
		21,
		2,
		`来自$` + "王俊" + `$的审批消息`,
		map[string]string{
			"ApprovalNo": "12345",
		},
	}
	b, e := json.Marshal(msg)
	if e != nil {
		t.Log("failed")
	}
	JPushMsgByAlias(alias, string(b), map[string]interface{}{})
	t.Log("hello world\n")
}
