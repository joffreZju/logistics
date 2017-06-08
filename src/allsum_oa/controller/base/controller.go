package base

import (
	mycache "common/lib/cache"
	"common/lib/redis"
	"fmt"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/cache"
	"strconv"
	"strings"
)

type Controller struct {
	beego.Controller
	ControllerName string
	ActionName     string
	IsFailed       bool

	Cache       cache.Cache         // 共享缓存
	LocalCache  cache.Cache         // 本机缓存
	RedisClient *redis.RedisManager // 本机缓存

	UserID     int64  // 用户ID
	UserComp   string // 用户公司
	UserGroups string // 用户组织
	UserRoles  string // 用户角色
	appName    string // app 名称
	appOS      string // app 系统
	appVer     string // app 版本号
}

//func (c *Controller) GetAppName() string {
//	if len(c.appName) > 0 {
//		return c.appName
//	}
//	return "app"
//}
//
//func (c *Controller) GetAppOS() string {
//	if strings.Index(c.appOS, "iphone") >= 0 {
//		return "ios"
//	}
//	if strings.Index(c.appOS, "android") >= 0 {
//		return "android"
//	}
//	return "unkown"
//}
//
//func (c *Controller) GetOSVersion() string {
//	if strings.Index(c.appOS, "iphone") >= 0 {
//		return strings.TrimSpace(c.appOS[6:])
//	} else if strings.Index(c.appOS, "android") >= 0 {
//		return strings.TrimSpace(c.appOS[7:])
//	}
//	return "unknown"
//}
//
//func (c *Controller) GetAppVersion() string {
//	return c.appVer
//}
//
//func (c *Controller) GetAppMainVersion() string {
//	fs := strings.Split(c.appVer, ".")
//	if len(fs) >= 2 {
//		return fmt.Sprintf("%s.%s", fs[0], fs[1])
//	}
//	return c.appVer
//}

func (c *Controller) Prepare() {
	strs := strings.Split(c.Ctx.Request.URL.Path, "/")
	if len(strs) > 2 {
		c.ControllerName = strs[len(strs)-2]
		c.ActionName = strs[len(strs)-1]
	}

	//perfcounter.Add(beego.BConfig.AppName+".request.total", 1)
	//perfcounter.Add(fmt.Sprintf("%s.%s.%s.request.total", beego.BConfig.AppName,
	//	c.ControllerName, c.ActionName), 1)
	c.LocalCache = mycache.LocalCache
	c.Cache = mycache.Cache
	c.RedisClient = redis.Client

	// 获取客户端版本号
	c.appName = strings.ToLower(strings.TrimSpace(c.GetString("source")))
	c.appOS = strings.ToLower(strings.TrimSpace(c.GetString("os")))
	c.appVer = strings.ToLower(strings.TrimSpace(c.GetString("ver")))

	// 从 access_token 中获取uid, 客户端可不传uid
	uid := c.Ctx.Request.Header.Get("uid")
	if len(uid) > 0 {
		c.Ctx.Input.SetParam("uid", uid)
		c.UserID, _ = strconv.ParseInt(uid, 10, 64)
	}
	c.UserComp = c.Ctx.Request.Header.Get("cno")
	m, e := c.RedisClient.Hmget(uid, []string{"roles", "groups"})
	if e != nil {
		beego.Error(e)
	}
	c.UserGroups = m["groups"]
	c.UserRoles = m["roles"]
}

func (c *Controller) Finish() {
	if c.IsFailed {
		fmt.Println("request failed")
		//perfcounter.Add(beego.BConfig.AppName+".request.failed", 1)
		//perfcounter.Add(fmt.Sprintf("%s.%s.%s.request.failed", beego.BConfig.AppName,
		//c.ControllerName, c.ActionName), 1)
	} else {
		fmt.Println("request success")
		//perfcounter.Add(beego.BConfig.AppName+".request.success", 1)
		//perfcounter.Add(fmt.Sprintf("%s.%s.%s.request.success", beego.BConfig.AppName,
		//	c.ControllerName, c.ActionName), 1)
	}
}
