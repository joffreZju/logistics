package main

import (
	"flag"
	"fmt"
	"os"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

var (
	Version    = "unknow"
	configPath = flag.String("c", "conf/allsum_account.conf", "config file path")
)

func main() {
	// 输出版本号
	if len(os.Args) > 1 && os.Args[1] == "version" {
		fmt.Println(Version)
		os.Exit(0)
	}

	// 加载配置文件
	flag.Parse()
	if len(*configPath) > 0 {
		err := beego.LoadAppConfig("ini", fmt.Sprintf("%s", *configPath))
		if err != nil {
			panic(err)
		}
	}

	// 初始化配置
	err := Init()
	if err != nil {
		panic(err)
	}

	// load router
	LoadRouter()

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowOrigins:     []string{"http://localhost:8090", "http://www.suanpeizaix.comw", "http://www.suanpeizaix.com:8090"},
		AllowOrigins:     []string{"*"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST"},
		AllowHeaders:     []string{"Origin"},
		ExposeHeaders:    []string{"Content-Length"},
		AllowCredentials: true,
	}))

	beego.Info("Init Server Begin..")
	beego.Run()
	beego.Info("Init Server End..")
}
