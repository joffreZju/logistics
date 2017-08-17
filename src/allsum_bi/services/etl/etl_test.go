package etl

import (
	"allsum_bi/db"
	"allsum_bi/services/util"
	"fmt"
	"sync"
	"testing"

	"github.com/astaxie/beego"
)

//向测试库注入 200 个schema
//func Test_mkschematable(b *testing.T) {
//
//	beego.LoadAppConfig("ini", "../../../../conf/allsum_bi.conf")
//	db.InitDb()
//
//	//这是一个 测试数据库 bi拥有其创建schema的权限
//	dbid := "5e972f07-10c2-4f62-8778-a8ec1d8281ff"
//	for i := 0; i < 200; i++ {
//		schemaname := fmt.Sprintf("test__%v", i)
//		db.CreateSchemaWithDbid(dbid, schemaname)
//		createsql, err := db.GetTableDescFromSource(util.BASEDB_CONNID, "public", "route_base", schemaname, "route_base")
//		if err != nil {
//			fmt.Println("createsql :", err, createsql)
//			return
//		}
//		err = db.Exec(dbid, createsql)
//		fmt.Println(err)
//	}
//}

//向目标schmema table 注入测试数据
func Test_cleandata(b *testing.T) {
	beego.LoadAppConfig("ini", "../../../../conf/allsum_bi.conf")
	db.InitDb()
	dbid := "5e972f07-10c2-4f62-8778-a8ec1d8281ff"
	var wg sync.WaitGroup
	for i := 0; i < 200; i++ {
		schemaname := fmt.Sprintf("test__%v", i)
		wg.Add(1)
		sql := "truncate " + schemaname + ".route_base"
		go func() {
			db.Exec(dbid, sql)
			wg.Done()
		}()
	}
	wg.Wait()
}

func Test_mktabledata(b *testing.T) {
	beego.LoadAppConfig("ini", "../../../../conf/allsum_bi.conf")
	db.InitDb()
	dbid := "5e972f07-10c2-4f62-8778-a8ec1d8281ff"
	var wg sync.WaitGroup
	for i := 0; i < 20; i++ {
		schemaname := fmt.Sprintf("test__%v", i)
		wg.Add(1)
		sourcejs, _ := MakeSourceJs(util.BASEDB_CONNID)
		sinkjs, _ := MakeSinkJsWithID(dbid)
		transportjs, _ := MakeTransportJs(dbid, "public", "route_base", schemaname, "route_base", "", "")
		runjs := MakeRunJs(sourcejs, sinkjs, transportjs)
		beego.Debug("runjs:", runjs)
		go func() {
			DoETL(i, []byte(runjs))
			wg.Done()
		}()
	}
	wg.Wait()
}
