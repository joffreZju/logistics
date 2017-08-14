package email

import (
	"github.com/astaxie/beego"
	gomail "gopkg.in/gomail.v2"
)

func SendEmail(targets []string, subject, body string) {
	if len(targets) == 0 {
		return
	}
	smtpHost := beego.AppConfig.String("emailAccount::smtp")
	smtpPort := beego.AppConfig.DefaultInt("emailAccount::port", 25)
	from := beego.AppConfig.String("emailAccount::from")
	password := beego.AppConfig.String("emailAccount::password")

	m := gomail.NewMessage()
	m.SetAddressHeader("From", from, "")

	tos := []string{}
	for _, v := range targets {
		tos = append(tos, m.FormatAddress(v, ""))

	}
	m.SetHeader("To", tos...)
	m.SetHeader("Subject", subject)
	m.SetBody("text/html", body)

	dial := gomail.NewPlainDialer(smtpHost, smtpPort, from, password)
	if e := dial.DialAndSend(m); e != nil {
		beego.Error("发送邮件失败:", e)

	} else {
		beego.Info("发送邮件成功:", targets)

	}

}
