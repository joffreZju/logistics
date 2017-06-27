package push

import (
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
	tels := []string{"wy"}
	JPushCommonMsg(tels, "nihao", map[string]interface{}{})
	t.Log("hello world\n")
}
