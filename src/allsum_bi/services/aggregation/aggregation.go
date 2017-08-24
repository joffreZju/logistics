package aggregation

import (
	"allsum_bi/db"
	"allsum_bi/models"
	"allsum_bi/services/util"
	"common/lib/service_client/oaclient"
	"fmt"
	"strings"
)

func AddAggregateByDataload(name string, reportid int, aggregateid int, owner string, tablename string, flush_script string, cron string, documents string) (aggregate models.AggregateOps, err error) {
	aggregate = models.AggregateOps{
		Name:      "data_load_" + name,
		DestTable: tablename,
		Reportid:  reportid,
		Script:    flush_script,
		Cron:      cron,
		Documents: documents,
		Status:    util.AGGREGATE_STARTED,
	}
	if aggregateid == 0 {
		err = models.InsertAggregateReturnAggregate(&aggregate)
		if err != nil {
			return
		}
	} else {
		aggregate.Id = aggregateid
		err = models.UpdateAggregate(aggregate)
		if err != nil {
			return
		}
	}

	schema := db.GetCompanySchema(owner)
	schema_table, _ := db.EncodeTableSchema(util.BASEDB_CONNID, schema, tablename)
	flush_script_real := strings.Replace(flush_script, util.SCRIPT_TABLE, schema_table, -1)
	flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, -1)

	err = AddCronWithFlushScript(aggregate.Id, cron, flush_script_real)

	return
}

func AddAggregate(uuid string, table_name string, create_script string, alter_script string, flush_script string, cron string, documents string) (err error) {
	aggregate, err := models.GetAggregateOpsByUuid(uuid)
	if err != nil {
		return
	}
	demand, err := models.GetReportDemand(aggregate.Reportid)
	if err != nil {
		return
	}
	//TODO add common report check
	report, err := models.GetReport(aggregate.Reportid)
	if err != nil {
		return
	}
	var schemas []string
	if report.Reporttype == util.REPORT_TYPE_COMMON {
		schemas, err = oaclient.GetAllCompanySchema()
		if err != nil {
			return
		}
	} else {
		schemas = []string{db.GetCompanySchema(demand.Owner)}

	}

	for _, schema := range schemas {
		err = db.CreateSchema(schema)
		if err != nil {
			return
		}
		schema_table := schema + "." + table_name
		isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema_table)
		if !isexsit {
			create_script_real := strings.Replace(create_script, util.SCRIPT_TABLE, schema_table, -1)
			create_script_real = strings.Replace(create_script_real, util.SCRIPT_SCHEMA, schema, -1)
			err = db.Exec(util.BASEDB_CONNID, create_script_real)
			if err != nil {
				return
			}
		}
		if aggregate.AlterScript != "" {
			alter_script_real := strings.Replace(aggregate.AlterScript, util.SCRIPT_TABLE, schema_table, -1)
			alter_script_real = strings.Replace(alter_script_real, util.SCRIPT_SCHEMA, schema, -1)
			//TODO not sure is able exec multiple sql need check
			err = db.Exec(util.BASEDB_CONNID, alter_script_real)
			if err != nil {
				return
			}
		}
		new_create_sql, err := db.GetTableDesc(util.BASEDB_CONNID, schema, table_name, schema, table_name)
		if err != nil {
			return err
		}
		new_create_sql_format := strings.Replace(new_create_sql, schema_table, util.SCRIPT_TABLE, -1)

		flush_script_real := strings.Replace(flush_script, util.SCRIPT_TABLE, schema_table, -1)
		flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, -1)

		err = AddCronWithFlushScript(aggregate.Id, cron, flush_script_real)

		aggregate.CreateScript = new_create_sql_format
		aggregate.AlterScript = ""
		aggregate.Cron = cron
		aggregate.Script = flush_script
		aggregate.Documents = documents
		aggregate.DestTable = table_name
		aggregate.Status = util.AGGREGATE_STARTED
		err = models.UpdateAggregate(aggregate, "dest_table", "create_script", "alter_script", "flush_script", "cron", "documents", "status")
	}
	return
}

func TestCreateScript(uuid string, table_name string, create_script string) (err error) {
	aggregate, err := models.GetAggregateOpsByUuid(uuid)
	if err != nil {
		return
	}
	demand, err := models.GetReportDemand(aggregate.Reportid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(demand.Owner)
	if !db.SchemaExist(util.BASEDB_CONNID, schema) {
		err = db.CreateSchema(schema)
		if err != nil {
			return err
		}
	}

	schema_table := schema + "." + table_name

	isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema_table)
	if isexsit {
		return fmt.Errorf("table is exsit ", schema_table)
	}
	table_name_test := schema_table + "_test"
	create_script_real := strings.Replace(create_script, util.SCRIPT_TABLE, table_name_test, -1)
	err = db.Exec(util.BASEDB_CONNID, create_script_real)
	if err != nil {
		return
	}
	defer func() {
		db.DeleteSchemaTable(util.BASEDB_CONNID, table_name_test)
	}()
	aggregate.CreateScript = create_script
	err = models.UpdateAggregate(aggregate, "create_script")
	return
}

func TestAlterScript(uuid string, table_name string, alter_script string) (err error) {
	aggregate, err := models.GetAggregateOpsByUuid(uuid)
	if err != nil {
		return
	}
	demand, err := models.GetReportDemand(aggregate.Reportid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(demand.Owner)
	schema_table := schema + "." + table_name
	isexsit := db.CheckTableExist(util.BASEDB_CONNID, schema_table)
	if !isexsit {
		return fmt.Errorf("table is not exist ", schema+"."+table_name)
	}
	table_name_test := schema_table + "_test"
	if !db.SchemaExist(util.BASEDB_CONNID, schema) {
		err = db.CreateSchema(schema)
		if err != nil {
			return err
		}
	}

	create_script_real := strings.Replace(aggregate.CreateScript, util.SCRIPT_TABLE, table_name_test, -1)
	err = db.Exec(util.BASEDB_CONNID, create_script_real)
	if err != nil {
		return
	}
	defer func() {
		db.DeleteSchemaTable(util.BASEDB_CONNID, table_name_test)
	}()
	alter_script_real := strings.Replace(alter_script, util.SCRIPT_TABLE, table_name_test, -1)
	//TODO not sure is able exec multiple sql need check
	err = db.Exec(util.BASEDB_CONNID, alter_script_real)
	if err != nil {
		return
	}
	//new_create_sql, err := db.GetTableDesc(util.BASEDB_CONNID, schema, table_name+"_test", schema, table_name+"_test")
	//if err != nil {
	//	return
	//}
	//new_create_sql_format := strings.Replace(new_create_sql, table_name_test, util.SCRIPT_TABLE, 1)
	//aggregate.CreateScript = new_create_sql_format
	aggregate.AlterScript = alter_script
	err = models.UpdateAggregate(aggregate, "create_script", "alter_script")
	return
}

func TestFlushScript(uuid string, table_name string, flush_script string, cron string) (err error) {
	aggregate, err := models.GetAggregateOpsByUuid(uuid)
	if err != nil {
		return
	}
	demand, err := models.GetReportDemand(aggregate.Reportid)
	if err != nil {
		return
	}
	schema := db.GetCompanySchema(demand.Owner)
	if !db.SchemaExist(util.BASEDB_CONNID, schema) {
		err = db.CreateSchema(schema)
		if err != nil {
			return err
		}
	}

	schema_table := schema + "." + table_name
	table_name_test := schema_table + "_test"
	create_script_real := strings.Replace(aggregate.CreateScript, util.SCRIPT_TABLE, table_name_test, -1)
	err = db.Exec(util.BASEDB_CONNID, create_script_real)
	if err != nil {
		return
	}
	defer func() {
		db.DeleteSchemaTable(util.BASEDB_CONNID, table_name_test)
	}()
	flush_script_real := strings.Replace(flush_script, util.SCRIPT_TABLE, table_name_test, -1)
	flush_script_real = strings.Replace(flush_script_real, util.SCRIPT_SCHEMA, schema, -1)
	err = db.Exec(util.BASEDB_CONNID, flush_script_real)
	if err != nil {
		return
	}
	//TODO cron test
	err = TestAddCronWithFlushScript(cron, flush_script_real)
	if err != nil {
		return
	}
	aggregate.Script = flush_script
	aggregate.Cron = cron
	err = models.UpdateAggregate(aggregate, "flush_script", "cron")
	return
}
