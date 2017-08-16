package etl

import (
	"allsum_bi/db"
	"allsum_bi/services/util"
	"fmt"
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
func Test_mktabledata(b *testing.T) {
	beego.LoadAppConfig("ini", "../../../../conf/allsum_bi.conf")
	db.InitDb()
	dbid := "5e972f07-10c2-4f62-8778-a8ec1d8281ff"
	for i := 0; i < 1; i++ {
		schemaname := fmt.Sprintf("test__%v", i)

		sourcejs, _ := MakeSourceJs(util.BASEDB_CONNID)
		sinkjs, _ := MakeSinkJsWithID(dbid)
		transportjs, _ := MakeTransportJs(dbid, "public", "route_base", schemaname, "route_base", "", "")
		runjs := MakeRunJs(sourcejs, sinkjs, transportjs)
		beego.Debug("runjs:", runjs)
		err := DoETL(i, []byte(runjs))
		if err != nil {
			break
		}
	}
}
