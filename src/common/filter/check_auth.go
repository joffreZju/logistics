package filter

import (
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
	"github.com/ysqi/tokenauth"
	"github.com/ysqi/tokenauth2beego/o2o"
)

var (
	NoAuthMap = make(map[string]bool)
)

func CheckAuthWebFilter(api bool, group string, notNeedAuthList []string) beego.FilterFunc {
	noAuthMap := map[string]bool{}
	for _, path := range notNeedAuthList {
		noAuthMap[path] = true
	}

	d := &tokenauth.DefaultProvider{}
	audience := &tokenauth.Audience{
		Name:        "CusSingleTokenCheck",
		ID:          group,
		TokenPeriod: tokenauth.TokenPeriod,
	}
	audience.Secret = d.GenerateSecretString(audience.ID)

	if o2o.Auth == nil {
		o2o.Auth = &o2o.O2OAutomatic{}
	}
	o2o.Auth.TokenFunc = d.GenerateTokenString
	o2o.Auth.Audience = audience

	return func(ctx *context.Context) {
		path := ctx.Request.URL.Path
		if token, err := o2o.Auth.CheckToken(ctx.Request); err != nil {
			if !noAuthMap[path] && !(strings.Index(path, ".") > 0 && noAuthMap[path[strings.Index(path, "."):]]) {
				if api || strings.Contains(path, "/v2/") {
					o2o.Auth.ReturnFailueInfo(err, ctx)
				} else {
					ctx.Redirect(302, "/login") // TODO: redirectPath
				}
			}
		} else {
			// beego.Debug("CheckAuthFilter token", token)
			ctx.Request.Header.Add("uid", token.SingleID)
			// 刷新过期时间
			token.DeadLine = time.Now().Unix() + int64(tokenauth.TokenPeriod)
			tokenauth.Store.FlushToken(token)
		}
	}
}

func CheckAuthFilter(group string, notNeedAuthList []string) beego.FilterFunc {
	noAuthMap := map[string]bool{}
	for _, path := range notNeedAuthList {
		noAuthMap[path] = true
	}

	d := &tokenauth.DefaultProvider{}
	audience := &tokenauth.Audience{
		Name:        "CusSingleTokenCheck",
		ID:          group,
		TokenPeriod: tokenauth.TokenPeriod,
	}
	audience.Secret = d.GenerateSecretString(audience.ID)

	if o2o.Auth == nil {
		o2o.Auth = &o2o.O2OAutomatic{}
	}
	o2o.Auth.TokenFunc = d.GenerateTokenString
	o2o.Auth.Audience = audience

	return func(ctx *context.Context) {
		path := ctx.Request.URL.Path
		if ctx.Request.Method == "OPTIONS" {
			return
		}
		if token, err := o2o.Auth.CheckToken(ctx.Request); err != nil {
			if !noAuthMap[ctx.Request.URL.Path] {
				beego.Debug("request to: ", path, ctx.Request.Method, err.Error(), "token:", token)
				o2o.Auth.ReturnFailueInfo(err, ctx)
			}
		} else {
			beego.Debug("CheckAuthFilter token", token)
			//for visit admin page
			//if strings.Contains(path, "/admin") {
			//	if token.GroupID != "1" {
			//		o2o.Auth.ReturnFailueInfo(err, ctx)
			//	}
			//} else {
			ctx.Request.Header.Add("uid", token.SingleID)
			ctx.Request.Header.Add("cno", token.GroupID)
			token.DeadLine = time.Now().Unix() + int64(tokenauth.TokenPeriod)
			tokenauth.Store.FlushToken(token)
			//}

			//strs := strings.Split(token.SingleID, "_")
			//if len(strs) == 1 {
			//	ctx.Request.Header.Add("uid", token.SingleID)
			//} else if len(strs) == 2 {
			//	ctx.Request.Header.Add("uid", strs[1])
			//	return
			//} else {
			//	beego.Debug(err.Error(), "token:", token)
			//	o2o.Auth.ReturnFailueInfo(err, ctx)
			//}
		}
	}
}

func GetAuthFilter() beego.FilterFunc {

	d := &tokenauth.DefaultProvider{}
	audience := &tokenauth.Audience{
		Name:        "CusSingleTokenCheck",
		ID:          "allsum_mobile",
		TokenPeriod: tokenauth.TokenPeriod,
	}
	audience.Secret = d.GenerateSecretString(audience.ID)

	if o2o.Auth == nil {
		o2o.Auth = &o2o.O2OAutomatic{}
	}
	o2o.Auth.TokenFunc = d.GenerateTokenString
	o2o.Auth.Audience = audience

	return func(ctx *context.Context) {
		if token, err := o2o.Auth.CheckToken(ctx.Request); err != nil {
			//beego.Debug("CheckAuthFilter notoken", noAuthMap[ctx.Request.URL.Path], ctx.Request.URL.Path)
			beego.Debug(err.Error(), "token:", token)
		} else {
			beego.Debug("CheckMobileWebFilter token", token.SingleID)
			strs := strings.Split(token.SingleID, "_")
			if len(strs) != 2 || strs[0] != "allsum" {
				return
			} else {
				ctx.Request.Header.Add("uid", strs[1])
			}
		}
	}
}
