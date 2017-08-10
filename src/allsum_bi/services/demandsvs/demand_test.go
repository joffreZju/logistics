package demandsvs

import (
	"fmt"
	"testing"

	"github.com/astaxie/beego"
)

func Test_getoauser(b *testing.T) {
	beego.AppConfig.Set("service_client::oa_host", "allsum.com:8094")
	users, err := GetHandlerUserFromOA()
	fmt.Println("user:", users, err)
}
