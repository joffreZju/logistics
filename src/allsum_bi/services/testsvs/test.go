package testsvs

import (
	"allsum_bi/models"
	"allsum_bi/services/util"
	"common/lib/email"
	"common/lib/service_client/oaclient"
)

func SendMailTestInfo(testinfo models.TestInfo) (err error) {
	var toid int
	var status_str string
	if testinfo.Status == util.IS_OPEN {
		toid = testinfo.Testerid
		status_str = "打开"
	} else {
		toid = testinfo.Handlerid
		status_str = "关闭"
	}
	userinfo, err := oaclient.GetUserInfo(toid)
	if err != nil {
		return
	}
	subject := "#BI测试#" + testinfo.Title
	body := testinfo.Documents + "\n状态: " + status_str
	email.SendEmail([]string{userinfo["Mail"].(string)}, subject, body)
	return
}
