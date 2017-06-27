package main

import (
	"allsum_bi/db"
	_ "allsum_bi/routers"
	"allsum_bi/services/etl"
	"net/http"

	"github.com/astaxie/beego"
)

func main() {
	beego.LoadAppConfig("ini", "conf/bi_web.conf")

	//pprf 工具
	go pprof()

	//db 初始化
	db.InitDb()

	//启动etl
	etl.Start()
	etl.TestETL()

	beego.Run()
}

func pprof() {
	beego.Debug(http.ListenAndServe("0.0.0.0:6060", nil))
}
