package filter

import (
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func AddAllowedCrossDomain(domain string) beego.FilterFunc {
	return func(ctx *context.Context) {
		ctx.Output.Header("Access-Control-Allow-Origin", domain)
		ctx.Output.Header("Access-Control-Allow-Methods", "POST,GET,OPTIONS,DELETE")
		ctx.Output.Header("Access-Control-Allow-Headers", "x-requested-with,content-type")
		ctx.Output.Header("Access-Control-Allow-Credentials", "true")
	}
}
