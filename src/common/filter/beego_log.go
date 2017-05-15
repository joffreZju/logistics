package filter

import (
	"path"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

func LoadLogFilter() {
	beego.DefaultAccessLogFilter = &myLogFilter{beego.DefaultAccessLogFilter}
}

type myLogFilter struct {
	parent beego.FilterHandler
}

func (l *myLogFilter) Filter(ctx *context.Context) bool {
	if l.parent != nil && l.parent.Filter(ctx) {
		return true
	}
	// 过滤阿负载均衡请求
	requestPath := path.Clean(ctx.Request.URL.Path)
	if ctx.Request.Method == "HEAD" && (requestPath == "/" || requestPath == "/health") {
		return true
	}
	return false

}
