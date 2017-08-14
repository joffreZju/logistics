package email

import (
	"testing"

	"github.com/astaxie/beego"
)

func Test_email(b *testing.T) {
	beego.AppConfig.Set("emailAccount::smtp", "smtp.mxhichina.com")
	beego.AppConfig.Set("emailAccount::prot", "25") //should be int
	beego.AppConfig.Set("emailAccount::from", "joffre@suanpeizai.com")
	beego.AppConfig.Set("emailAccount::password", "Wang1234")

	SendEmail([]string{"bl.he@suanpeizai.com"}, "邮件测试", "邮件测试")
}
