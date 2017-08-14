package filter

import (
	"common/lib/errcode"
	"common/lib/redis"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"strings"
)

func CheckApiPermission(exemptPrefix []string) beego.FilterFunc {
	return func(ctx *context.Context) {
		path := ctx.Request.URL.Path
		for _, v := range exemptPrefix {
			if strings.Contains(path, v) {
				return
			}
		}
		uid := ctx.Request.Header.Get("uid")
		tokenStr := ctx.Request.Header.Get("access_token")
		userKey := fmt.Sprintf("%s-%s", uid, tokenStr)
		m, e := redis.Client.Hmget(userKey, []string{"functions"})
		if e != nil {
			beego.Error(e)
			ctx.Output.JSON(errcode.ErrGetLoginInfo, false, false)
		} else if !strings.Contains(m["functions"], fmt.Sprintf("-%s-", path)) {
			beego.Error("后端权限验证失败")
			ctx.Output.JSON(errcode.ErrUrlPermission, false, false)
		} else {
			beego.Info("userId:", uid, "path:", path, "后端权限验证通过")
		}
	}
}
