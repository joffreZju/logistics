package oaclient

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
)

func Test_oa(b *testing.T) {
	beego.AppConfig.Set("service_client::oa_host", "allsum.com:8094")
	company := "allsum"
	roleid := 1
	schemas, err := GetAllCompanySchema()
	fmt.Println("schemas:", schemas, err)
	roleinfos, err := GetAllRoleByCompany(company)
	fmt.Println("roles:", roleinfos, err)
	userinfos, err := GetAllUserByRole(company, roleid)
	fmt.Println("userinfo:", userinfos, err)
	userinfo, err := GetUserInfo(7)
	fmt.Println("userinfo:", userinfo, err)

}
