package dataload

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/aggregation"
	"allsum_bi/services/util"
	"encoding/json"
	"fmt"
	"strings"

	"github.com/astaxie/beego"
)

func AddDataLoad(dataload map[string]string) (err error) {
	dataload_db, err := models.GetDataLoadByUuid(dataload["uuid"])
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(dataload_db.Owner)
	table_name := dataload["table_name"]
	schema_table := schema + "." + table_name
	create_script := dataload["create_script"]
	isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema_table)
	beego.Debug("check table exist ", isexsit)
	if !isexsit {
		if create_script == "" {
			return fmt.Errorf("miss create script!")
		}
		create_script_real := strings.Replace(create_script, util.SCRIPT_TABLE, schema_table, util.SCRIPT_LIMIT)
		if !db.SchemaExist(util.BASEDB_CONNID, schema) {
			err = db.CreateSchema(schema)
			if err != nil {
				return err
			}
		}
		err = db.Exec(util.BASEDB_CONNID, create_script_real)
		if err != nil {
			return
		}
	}
	alter_script := dataload["alter_script"]
	if alter_script != "" && alter_script != "null" {
		alter_script_real := strings.Replace(alter_script, util.SCRIPT_TABLE, schema_table, util.SCRIPT_LIMIT)
		//TODO not sure is able exec multiple sql need check
		err = db.Exec(util.BASEDB_CONNID, alter_script_real)
		if err != nil {
			return
		}
		new_create_sql, err := db.GetTableDesc(util.BASEDB_CONNID, schema, table_name, schema, table_name)
		if err != nil {
			return err
		}
		create_script = strings.Replace(new_create_sql, schema_table, util.SCRIPT_TABLE, util.SCRIPT_LIMIT)
	}
	flush_script := dataload["flush_script"]
	aggregationId := 0
	cron := dataload["cron"]
	documents := dataload["documents"]
	if flush_script != "" {

		aggregate_ops, err := aggregation.AddAggregateByDataload(dataload_db.Uuid, dataload_db.Owner, schema_table, flush_script, cron, documents)
		if err != nil {
			return err
		}
		aggregationId = aggregate_ops.Id
	}
	dataload_db.CreateScript = create_script
	dataload_db.AlterScript = "null"
	//	dataload_db.FlushScript = flush_script
	//	dataload_db.Cron = cron
	dataload_db.Basetable = table_name
	dataload_db.Documents = documents
	dataload_db.Aggregateid = aggregationId
	dataload_db.Status = util.DATALOAD_STARTED
	dataload_db.WebPath = dataload["webpath"]
	columnmaps, err := db.GetTableColumes(util.BASEDB_CONNID, table_name, schema)
	if err != nil {
		beego.Error("get table columes err:", err)
		return
	}
	columns, err := json.Marshal(columnmaps)
	if err != nil {
		beego.Error("get columnmaps err", err)
		return
	}
	dataload_db.Columns = string(columns)
	dataload_db.Name = dataload["name"]

	err = models.UpdateDataLoad(dataload_db, "create_script", "alter_script", "documents", "aggregateid", "status", "columns", "insert_script", "name", "web_path")
	if err != nil {
		return
	}
	return
}

func TestCreateScript(dataload_uuid string, table_name string, create_script string) (err error) {
	dataload, err := models.GetDataLoadByUuid(dataload_uuid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(dataload.Owner)
	isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema+"."+table_name)
	if isexsit {
		return fmt.Errorf("table is exsit ", schema+"."+table_name)
	}
	table_name_test := schema + "." + table_name + "_test"
	create_script_real := strings.Replace(create_script, util.SCRIPT_TABLE, table_name_test, util.SCRIPT_LIMIT)
	if !db.SchemaExist(util.BASEDB_CONNID, schema) {
		err = db.CreateSchema(schema)
		if err != nil {
			return err
		}
	}

	err = db.Exec(util.BASEDB_CONNID, create_script_real)
	if err != nil {
		return
	}
	defer func() {
		db.DeleteSchemaTable(util.BASEDB_CONNID, table_name_test)
	}()
	dataload.CreateScript = create_script
	err = models.UpdateDataLoad(dataload, "create_script")
	return
}

func TestAlterScript(dataload_uuid string, table_name string, alter_script string) (err error) {
	dataload, err := models.GetDataLoadByUuid(dataload_uuid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(dataload.Owner)
	//	isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema+"."+table_name)
	//	if !isexsit {
	//		return fmt.Errorf("table is not exist ", schema+"."+table_name)
	//	}
	if !db.SchemaExist(util.BASEDB_CONNID, schema) {
		err = db.CreateSchema(schema)
		if err != nil {
			return err
		}
	}

	table_name_test := schema + "." + table_name + "_test"
	create_script_real := strings.Replace(dataload.CreateScript, util.SCRIPT_TABLE, table_name_test, util.SCRIPT_LIMIT)
	err = db.Exec(util.BASEDB_CONNID, create_script_real)
	if err != nil {
		return
	}
	defer func() {
		db.DeleteSchemaTable(util.BASEDB_CONNID, table_name_test)
	}()
	alter_script_real := strings.Replace(alter_script, util.SCRIPT_TABLE, table_name_test, util.SCRIPT_LIMIT)
	//TODO not sure is able exec multiple sql need check
	err = db.Exec(util.BASEDB_CONNID, alter_script_real)
	if err != nil {
		return
	}
	//	new_create_sql, err := db.GetTableDesc(util.BASEDB_CONNID, schema, table_name+"_test", schema, table_name+"_test")
	//	if err != nil {
	//		return
	//	}
	//	new_create_sql_format := strings.Replace(new_create_sql, table_name_test, util.SCRIPT_TABLE, 1)
	//	dataload.CreateScript = new_create_sql_format
	dataload.AlterScript = alter_script
	err = models.UpdateDataLoad(dataload, "create_script", "alter_script")
	return
}

func InsertNewData(uuid string, fields []string, data []map[string]interface{}) (err error) {
	dataload, err := models.GetDataLoadByUuid(uuid)
	if err != nil {
		return
	}
	tablename := dataload.Basetable
	insertSql, params, err := db.MakeInsetSql(tablename, fields, data)
	if err != nil {
		return
	}
	err = db.Exec(util.BASEDB_CONNID, insertSql, params...)
	return
}

func UpdateData(uuid string, fields []string, datas []map[string]interface{}) (err error) {
	dataload, err := models.GetDataLoadByUuid(uuid)
	if err != nil {
		return
	}
	tablename := dataload.Basetable
	schema_table := strings.Split(tablename, ".")
	pks, err := db.GetTablePk(util.BASEDB_CONNID, schema_table[0], schema_table[1])
	if err != nil {
		return
	}
	for _, data := range datas {
		condition := map[string]interface{}{}
		for key, _ := range pks {
			condition[key] = data[key]
		}
		sql, params, err := db.MakeUpdateSql(tablename, condition, fields, data)
		if err != nil {
			return err
		}
		beego.Debug("sql:", sql, params)
		err = db.Exec(util.BASEDB_CONNID, sql, params...)
		if err != nil {
			return err
		}
	}
	return
}

func GetData(uuid string, conditions []map[string]interface{}, limit int) (columns []map[string]string, datas []map[string]interface{}, err error) {
	dataload, err := models.GetDataLoadByUuid(uuid)
	if err != nil {
		return
	}
	table_name := dataload.Basetable
	//	schema_table := strings.Split(table_name, ".")
	err = json.Unmarshal([]byte(dataload.Columns), &columns)
	if err != nil {
		return
	}

	datas, err = db.ListTableData(util.BASEDB_CONNID, table_name, conditions, limit)
	if err != nil {
		return
	}
	return
}
