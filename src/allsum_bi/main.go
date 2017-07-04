package main

import (
	"allsum_bi/db"
	_ "allsum_bi/routers"
	"allsum_bi/services/etl"
	"net/http"
	"os"
	"os/signal"

	_ "allsum_bi/models"

	"github.com/astaxie/beego"
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
