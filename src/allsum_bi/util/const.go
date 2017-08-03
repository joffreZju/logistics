package util

//目前只支持pgsql
const MYSQL_DB_TYPE = "mysql"
const PG_DB_TYPE = "pgsql"
const MONGO_DB_TYPE = "mongo"

//basedb conn id

const BASEDB_CONNID = "bi_base_db"

//schema
const BI_SCHEMA = "bi_system"
const BI_MANAGER = "manager"
const BI_COMMENT_PREFIX = ""

//ETL
const TRANSFORM_PATH = "./transform_js/"
const JS_TEMPLATE = "var %s = postgres({\"uri\": \"postgres://%s:%d/%s?sslmode=disable&user=%s&password=%s\"})"
const JS_TRANSPORT = "t.Source( source, \"/^%s\\/\" ).%sSave( sink, \"/%s/\" )"
const TRANSPORTFORM_GOJA = "Transform(%s)."

const DEFAULT_TRANSPORT = "skip({\"field\": \"%s\", \"operator\": \"<=\", \"match\": %v})"
const DEFAULT_PARAMS_SQL = "select max(%v) as %v from %v"

const SYNC_ENABLE = 1
const SYNC_DISABLE = 0

const SYNC_NONE = "none"
const SYNC_BUILDING = "building"
const SYNC_STARTED = "started"
const SYNC_ERROR = "error"
const SYNC_STOP = "stop"

//role
const ROLETYPE_ASSIGNER = 1
const ROLETYPE_PROJECTOR = 2
const ROLETYPE_TESTER = 3

//action
const ACTION_LISTDEMAND_ASSIGNER = "list_demand_assigner"
const ACTION_LISTDEMAND_PROJECTOR = "list_demand_projector"
const ACTION_LISTDEMAND_TESTER = "list_demand_tester"

//demand
const DEMAND_STATUS_NO_ASSIGN = 0 //未指派
const DEMAND_STATUS_BUILDING = 1  //开发中
const DEMAND_STATUS_TESTING = 2   //测试中
const DEMAND_STATUS_REVIEW = 3    //审核中
const DEMAND_STATUS_RELEASE = 4   //已发布
const DEMAND_STATUS_REBACK = 5    //打回

//report
const REPORT_STATUS_ANALYS = 0  //分析中
const REPORT_STATUS_DEVELOP = 1 //开发中
const REPORT_STATUS_TEST = 2    //测试中
const REPORT_STATUS_REVIEW = 3  //审核中
const REPORT_STATUS_RELEASE = 4 //已发布
const REPORT_STATUS_DISABLE = 5 //失效

const REPORT_TYPE_COMMON = 0  //通用报表类型
const REPORT_TYPE_PRIVATE = 1 //个性化报表类型

//script
const SCRIPT_TABLE = "{TABLE_NAME}"
const SCRIPT_SCHEMA = "{SCHEMA_NAME}"
const SCRIPT_OWNER = "{OWNER}"
const SCRIPT_LIMIT = 50

//dataload
const DATALOAD_BUILDING = "building"
const DATALOAD_STARTED = "started"
const DATALOAD_STOP = "stop"
const IS_INSERT = "insert"
const IS_UPDATE = "update"

//aggregate
const AGGREGATE_NONE = "none"
const AGGREGATE_BUILDING = "building"
const AGGREGATE_STARTED = "started"
const AGGREGATE_ERROR = "error"
const AGGREGATE_STOP = "stop"

//reportset
//aggregate
const REPORTSET_NONE = "none"
const REPORTSET_BUILDING = "building"
const REPORTSET_STARTED = "started"
const REPORTSET_ERROR = "error"
const REPORTSET_STOP = "stop"

//test
const TEST_MAX_UPLOAD_IMAGE = 9

const IS_OPEN = 1
const IS_CLOSE = 0

//KETTLEJOB
const KETTLEJOB_RIGHT = 1
const KETTLEJOB_FAIL = 0
