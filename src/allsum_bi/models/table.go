package models

import (
	"allsum_bi/util"
)

func GetDatabaseManagerTableName() string {
	return util.BI_MANAGER + ".database_manager"
}

func GetSynchronousTableName() string {
	return util.BI_SCHEMA + ".synchronous"
}

func GetSynchronousLogTableName() string {
	return util.BI_SCHEMA + ".synchronous_log"
}

func GetDemandTableName() string {
	return util.BI_SCHEMA + ".demand"
}

func GetReportTableName() string {
	return util.BI_SCHEMA + ".report"
}

func GetDataLoadTableName() string {
	return util.BI_SCHEMA + ".data_load"
}

func GetAggregateOpsTableName() string {
	return util.BI_SCHEMA + ".aggregate_ops"
}

func GetReportSetTableName() string {
	return util.BI_SCHEMA + ".report_set"
}

func GetTestInfoTableName() string {
	return util.BI_SCHEMA + ".test_info"
}

func GetUserAuthorityTableName(companyid string) string {
	return util.BI_COMMENT_PREFIX + companyid + ".user_authority"
}
