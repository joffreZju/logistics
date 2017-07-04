package main

import (
	"allsum_bi/db"
	_ "allsum_bi/routers"
	"allsum_bi/services/etl"
	"common/filter"
	"common/lib/cache"
	"common/lib/push"
	"common/lib/redis"
	"math/rand"
	"net/http"
	"os"
	"os/signal"
	"time"

	_ "allsum_bi/models"

	"github.com/astaxie/beego"
	"github.com/ysqi/tokenauth2beego"
)

func main() {
	beego.LoadAppConfig("ini", "conf/allsum_bi.conf")

	//pprf 工具
	go pprof()

	//db 初始化
	db.InitDb()

	//启动etl
	etl.Start()
	//	etl.TestETL()
	go beego.Run()

	//INIT
	Init()

	signal_f()
}

func pprof() {
	beego.Debug(http.ListenAndServe("0.0.0.0:6060", nil))
}

func signal_f() {
	c := make(chan os.Signal, 1)
	signal.Notify(c, os.Interrupt, os.Kill)

	<-c
}

func Init() {
	rand.Seed(time.Now().UnixNano())

	key := beego.AppConfig.String("seed")
	// init log
	err := InitLog()
	if err != nil {
		beego.Error("init log failed : ", err)
		return
	}

	// init tokenauth
	err = tokenauth2beego.Init(key)
	if err != nil {
		beego.Error("init token auth failed : ", err)
		return
	}
	// init redis cache
	err = cache.Init(key)
	if err != nil {
		beego.Error("init cache failed : ", err)
		return
	}

	// init redis cache
	err = redis.Init(key)
	if err != nil {
		beego.Error("init redis client failed : ", err)
		return
	}

	// init push
	err = push.Init()

}

func InitLog() (err error) {
	filter.LoadLogFilter()
	typ := beego.AppConfig.String("log::type")
	cons := beego.AppConfig.String("log::params")
	err = beego.SetLogger(typ, cons)
	beego.SetLogFuncCall(true)
	//beego.BeeLogger.SetLogFuncCallDepth(4)
	return
}
