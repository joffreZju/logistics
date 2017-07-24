--
-- PostgreSQL database dump
--

-- Dumped from database version 9.6.2
-- Dumped by pg_dump version 9.6.2

SET statement_timeout = 0;
SET lock_timeout = 0;
SET idle_in_transaction_session_timeout = 0;
SET client_encoding = 'UTF8';
SET standard_conforming_strings = on;
SET check_function_bodies = false;
SET client_min_messages = warning;
SET row_security = off;

--
-- Name: bi_system; Type: SCHEMA; Schema: -; Owner: user_logistic
--

CREATE SCHEMA bi_system;


ALTER SCHEMA bi_system OWNER TO user_logistic;

--
-- Name: SCHEMA bi_system; Type: COMMENT; Schema: -; Owner: user_logistic
--

COMMENT ON SCHEMA bi_system IS 'standard public schema';


--
-- Name: manager; Type: SCHEMA; Schema: -; Owner: user_logistic
--

CREATE SCHEMA manager;


ALTER SCHEMA manager OWNER TO user_logistic;

--
-- Name: plpgsql; Type: EXTENSION; Schema: -; Owner: 
--

CREATE EXTENSION IF NOT EXISTS plpgsql WITH SCHEMA pg_catalog;


--
-- Name: EXTENSION plpgsql; Type: COMMENT; Schema: -; Owner: 
--

COMMENT ON EXTENSION plpgsql IS 'PL/pgSQL procedural language';


SET search_path = bi_system, pg_catalog;

SET default_tablespace = '';

SET default_with_oids = false;

--
-- Name: aggregate_log; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE aggregate_log (
    id bigint NOT NULL,
    aggregateid bigint NOT NULL,
    reportid bigint NOT NULL,
    error text,
    res character varying(256),
    "timestamp" timestamp without time zone NOT NULL
);


ALTER TABLE aggregate_log OWNER TO user_logistic;

--
-- Name: COLUMN aggregate_log.id; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log.id IS '序列id';


--
-- Name: COLUMN aggregate_log.aggregateid; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log.aggregateid IS '同步id';


--
-- Name: COLUMN aggregate_log.reportid; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log.reportid IS '报表id';


--
-- Name: COLUMN aggregate_log.error; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log.error IS '报错信息';


--
-- Name: COLUMN aggregate_log.res; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log.res IS '聚合结果信息';


--
-- Name: COLUMN aggregate_log."timestamp"; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_log."timestamp" IS '时间';






--
-- Name: aggregate_ops_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE aggregate_ops_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE aggregate_ops_id_seq OWNER TO user_logistic;

--
-- Name: aggregate_ops; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE aggregate_ops (
    id bigint DEFAULT nextval('aggregate_ops_id_seq'::regclass) NOT NULL,
    reportid bigint NOT NULL,
    uuid character varying(64) NOT NULL,
    name character varying(64) NOT NULL,
    create_script text,
    alter_script text,
    script text,
    script_type character varying(16),
    dest_table character varying(256),
    cron character varying(256),
    documents text,
    status character varying(16) NOT NULL,
    aggregate_use_tables character varying(1024)
);


ALTER TABLE aggregate_ops OWNER TO user_logistic;

--
-- Name: COLUMN aggregate_ops.id; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.id IS '序列id';


--
-- Name: COLUMN aggregate_ops.reportid; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.reportid IS '报表id';


--
-- Name: COLUMN aggregate_ops.create_script; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.create_script IS '初始化脚本用‘；’分割';


--
-- Name: COLUMN aggregate_ops.script; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.script IS '脚本定期执行';


--
-- Name: COLUMN aggregate_ops.script_type; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.script_type IS '为数据库内部自己做定期执行的脚本不定期触发';


--
-- Name: COLUMN aggregate_ops.dest_table; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.dest_table IS '目标表‘，’分割';


--
-- Name: COLUMN aggregate_ops.cron; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.cron IS '停服重启首次运行的钟点／0 为马上启动 0～24';


--
-- Name: COLUMN aggregate_ops.documents; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.documents IS '描述';


--
-- Name: COLUMN aggregate_ops.status; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.status IS '使用状态';


--
-- Name: COLUMN aggregate_ops.name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.name IS '名称';


--
-- Name: COLUMN aggregate_ops.aggregate_use_tables; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN aggregate_ops.aggregate_use_tables IS '清洗中用到的表';


--
-- Name: test_info_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE test_info_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE test_info_id_seq OWNER TO user_logistic;

--
-- Name: test_info; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE test_info (
    id bigint DEFAULT nextval('test_info_id_seq'::regclass) NOT NULL,
    uuid character varying(64) NOT NULL,
    reportid bigint,
    documents text,
    filepaths varchar(128)[],
    status int2 
);


ALTER TABLE test_info OWNER TO user_logistic;


COMMENT ON COLUMN test_info.filepaths IS '文件路径列表';




--
-- Name: data_load_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE data_load_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE data_load_id_seq OWNER TO user_logistic;

--
-- Name: data_load; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE data_load (
    id bigint DEFAULT nextval('data_load_id_seq'::regclass) NOT NULL,
    uuid character varying(64) NOT NULL,
    name character varying(125),
    owner character varying(128) NOT NULL,
    columns character varying,
    create_script text,
    alter_script text,
    basetable character varying(256),
    documents text,
    status character varying(16) NOT NULL,
    web_path character varying(356),
    aggregateid bigint,
    webfile_name character varying
);


ALTER TABLE data_load OWNER TO user_logistic;

--
-- Name: COLUMN data_load.owner; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN data_load.owner IS '所属组织';


--
-- Name: demand_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE demand_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE demand_id_seq OWNER TO user_logistic;

--
-- Name: demand; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE demand (
    id bigint DEFAULT nextval('demand_id_seq'::regclass) NOT NULL,
    uuid character varying(64) NOT NULL,
    owner character varying(128) NOT NULL,
    owner_name character varying(256) NOT NULL,
    exhibitor character varying(128) NOT NULL,
    reportid bigint NOT NULL,
    description text NOT NULL,
    price double precision,
    resultcode text,
    assignetime timestamp without time zone,
    complettime timestamp without time zone,
    assigner_authority text,
    assigner_name character varying(64),
    handler_name character varying(64),
    inittime timestamp without time zone,
    doc_url character varying(256),
    doc_name character varying(64),
    contactid bigint NOT NULL,
    assignerid bigint,
    handlerid bigint NOT NULL,
    deadline timestamp without time zone,
    status smallint DEFAULT 0 NOT NULL
);


ALTER TABLE demand OWNER TO user_logistic;

--
-- Name: COLUMN demand.reportid; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.reportid IS '对应的报表id';


--
-- Name: COLUMN demand.description; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.description IS '描述备注';


--
-- Name: COLUMN demand.price; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.price IS '价格';


--
-- Name: COLUMN demand.resultcode; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.resultcode IS '处理结果描述';


--
-- Name: COLUMN demand.assignetime; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.assignetime IS '指派时间';


--
-- Name: COLUMN demand.complettime; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.complettime IS '完成时间';


--
-- Name: COLUMN demand.assigner_authority; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.assigner_authority IS '指派的db和schema 信息';


--
-- Name: COLUMN demand.owner_name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.owner_name IS '需求发起方';


--
-- Name: COLUMN demand.assigner_name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.assigner_name IS '指派人名称';


--
-- Name: COLUMN demand.handler_name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.handler_name IS '处理人名称';


--
-- Name: COLUMN demand.doc_url; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.doc_url IS '需求文档路径';


--
-- Name: COLUMN demand.doc_name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN demand.doc_name IS '需求文档名称';


--
-- Name: report_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE report_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE report_id_seq OWNER TO user_logistic;

--
-- Name: report; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE report (
    id bigint DEFAULT nextval('report_id_seq'::regclass) NOT NULL,
    demandid bigint NOT NULL,
    uuid character varying(64) NOT NULL,
    name character varying(256),
    owner character varying(125),
    reporttype character varying(16),
    description text,
    level smallint,
    status smallint NOT NULL,
    grouppath character varying(256)
);


ALTER TABLE report OWNER TO user_logistic;

--
-- Name: COLUMN report.grouppath; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report.grouppath IS '报表的分组设置';


--
-- Name: COLUMN report.reporttype; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report.reporttype IS '报表类型';


--
-- Name: report_group; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE report_group (
    id bigint NOT NULL,
    name character varying(256) NOT NULL,
    owner character varying(16) NOT NULL,
    description text,
    superior bigint NOT NULL
);


ALTER TABLE report_group OWNER TO user_logistic;

--
-- Name: COLUMN report_group.name; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report_group.name IS '组名称';


--
-- Name: COLUMN report_group.owner; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report_group.owner IS '组所属（组织id）';


--
-- Name: COLUMN report_group.description; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report_group.description IS '描述';


--
-- Name: COLUMN report_group.superior; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN report_group.superior IS '上级关系 0为根级';


--
-- Name: report_set_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE report_set_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE report_set_id_seq OWNER TO user_logistic;

--
-- Name: report_set; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE report_set (
    id bigint DEFAULT nextval('report_set_id_seq'::regclass) NOT NULL,
    reportid bigint NOT NULL,
    uuid character varying(64) NOT NULL,
    dbid varchar(128) NOT NULL,
    script text,
    resttype integer,
    conditions character varying(2048),
    enable_event_types character varying(1024) DEFAULT NULL::character varying,
    status character varying(64) NOT NULL,
    web_path character varying(256),
    webfile_name character varying(256)
);


ALTER TABLE report_set OWNER TO user_logistic;

--
-- Name: synchronous_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE synchronous_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE synchronous_id_seq OWNER TO user_logistic;

--
-- Name: synchronous; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE synchronous (
    id bigint DEFAULT nextval('synchronous_id_seq'::regclass) NOT NULL,
    uuid character varying(256) NOT NULL,
    owner character varying(256) NOT NULL,
    create_script text,
    param_script text,
    script text,
    source_db_id character varying(256),
    source_table character varying(128),
    dest_db_id character varying(256),
    dest_table character varying(128),
    cron character varying(256),
    documents text,
    error_limit integer DEFAULT 0,
    lasttime timestamp without time zone,
    status character varying(16) NOT NULL
);


ALTER TABLE synchronous OWNER TO user_logistic;

--
-- Name: COLUMN synchronous.owner; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous.owner IS '在pgsql里面 默认就当作schema';


--
-- Name: COLUMN synchronous.error_limit; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous.error_limit IS '错误上限';


--
-- Name: COLUMN synchronous.lasttime; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous.lasttime IS '最后一次执行时间';


--
-- Name: synchronous_log_id_seq; Type: SEQUENCE; Schema: bi_system; Owner: user_logistic
--

CREATE SEQUENCE synchronous_log_id_seq
    START WITH 1
    INCREMENT BY 1
    NO MINVALUE
    NO MAXVALUE
    CACHE 1;


ALTER TABLE synchronous_log_id_seq OWNER TO user_logistic;

--
-- Name: synchronous_log; Type: TABLE; Schema: bi_system; Owner: user_logistic
--

CREATE TABLE synchronous_log (
    id bigint DEFAULT nextval('synchronous_log_id_seq'::regclass) NOT NULL,
    syncid bigint NOT NULL,
    errormsg text,
    res text,
    "timestamp" timestamp without time zone NOT NULL,
    status smallint NOT NULL
);


ALTER TABLE synchronous_log OWNER TO user_logistic;

--
-- Name: COLUMN synchronous_log.id; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous_log.id IS '序列id';


--
-- Name: COLUMN synchronous_log.syncid; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous_log.syncid IS '同步id';


--
-- Name: COLUMN synchronous_log.errormsg; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous_log.errormsg IS '失败信息';


--
-- Name: COLUMN synchronous_log.res; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous_log.res IS '同步结果';


--
-- Name: COLUMN synchronous_log."timestamp"; Type: COMMENT; Schema: bi_system; Owner: user_logistic
--

COMMENT ON COLUMN synchronous_log."timestamp" IS '发生时间';


SET search_path = manager, pg_catalog;

--
-- Name: authority_group; Type: TABLE; Schema: manager; Owner: user_logistic
--

CREATE TABLE authority_group (
    id bigint NOT NULL,
    projectergroupid bigint NOT NULL,
    dbid character varying(16) NOT NULL,
    schema character varying(16) NOT NULL,
    rules character varying(256) NOT NULL,
    dbpasswd character varying(256) NOT NULL,
    starttime timestamp without time zone NOT NULL,
    usetime bigint
);


ALTER TABLE authority_group OWNER TO user_logistic;

--
-- Name: authority_user; Type: TABLE; Schema: manager; Owner: user_logistic
--

CREATE TABLE authority_user (
    "projecterId" bigint NOT NULL,
    groupid bigint NOT NULL,
    starttime timestamp without time zone NOT NULL,
    usetime integer
);


ALTER TABLE authority_user OWNER TO user_logistic;

--
-- Name: database_manager; Type: TABLE; Schema: manager; Owner: user_logistic
--

CREATE TABLE database_manager (
    dbid character varying(128) NOT NULL,
    name character varying(128) NOT NULL,
    dbname character varying(128) NOT NULL,
    dbtype character varying(20) NOT NULL,
    host character varying(128) NOT NULL,
    port integer NOT NULL,
    dbuser character varying(125) NOT NULL,
    password character varying(128) NOT NULL,
    params character varying(1024)
);


ALTER TABLE database_manager OWNER TO user_logistic;

SET search_path = bi_system, pg_catalog;

--
-- Name: aggregate_log aggregate_log_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY aggregate_log
    ADD CONSTRAINT aggregate_log_pkey PRIMARY KEY (id);


--
-- Name: aggregate_ops aggregate_ops_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY aggregate_ops
    ADD CONSTRAINT aggregate_ops_pkey PRIMARY KEY (id);


--
-- Name: data_load data_load_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY data_load
    ADD CONSTRAINT data_load_pkey PRIMARY KEY (id);


--
-- Name: demand demand_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY demand
    ADD CONSTRAINT demand_pkey PRIMARY KEY (id);


--
-- Name: report_set get_report_ops_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY report_set
    ADD CONSTRAINT get_report_ops_pkey PRIMARY KEY (id);


--
-- Name: data_load name; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY data_load
    ADD CONSTRAINT name UNIQUE (name);


--
-- Name: report_group report_group_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY report_group
    ADD CONSTRAINT report_group_pkey PRIMARY KEY (id);


--
-- Name: report report_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY report
    ADD CONSTRAINT report_pkey PRIMARY KEY (id);


--
-- Name: synchronous_log synchronous_log_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY synchronous_log
    ADD CONSTRAINT synchronous_log_pkey PRIMARY KEY (id);


--
-- Name: synchronous synchronous_pkey; Type: CONSTRAINT; Schema: bi_system; Owner: user_logistic
--

ALTER TABLE ONLY synchronous
    ADD CONSTRAINT synchronous_pkey PRIMARY KEY (id);


SET search_path = manager, pg_catalog;

--
-- Name: authority_group authority_group_pkey; Type: CONSTRAINT; Schema: manager; Owner: user_logistic
--

ALTER TABLE ONLY authority_group
    ADD CONSTRAINT authority_group_pkey PRIMARY KEY (id);


--
-- Name: authority_user authority_user_pkey; Type: CONSTRAINT; Schema: manager; Owner: user_logistic
--

ALTER TABLE ONLY authority_user
    ADD CONSTRAINT authority_user_pkey PRIMARY KEY ("projecterId", groupid);


--
-- Name: database_manager database_manager_pkey; Type: CONSTRAINT; Schema: manager; Owner: user_logistic
--

ALTER TABLE ONLY database_manager
    ADD CONSTRAINT database_manager_pkey PRIMARY KEY (dbid);



--
--kettle_job
--
-- ---------------------
CREATE TABLE "bi_system"."kettle_job" (
	"id" int4 NOT NULL DEFAULT nextval('kettle_job_id_seq'::regclass),
	"name" varchar(256) NOT NULL COLLATE "default",
	"cron" varchar(64) COLLATE "default",
	"kjbpath" varchar(256) COLLATE "default",
	"ktrpaths" varchar(256) COLLATE "default",
	"status" int2 NOT NULL,
	"uuid" varchar(128) NOT NULL COLLATE "default",
	"lock" varchar(512) COLLATE "default"
)
WITH (OIDS=FALSE);
ALTER TABLE "bi_system"."kettle_job" OWNER TO "user_logistic";

-- ----------------------------
--  Primary key structure for table kettle_job
-- ----------------------------
ALTER TABLE "bi_system"."kettle_job" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;

CREATE SEQUENCE "bi_system"."kettle_job_id_seq"  
START WITH 1  
INCREMENT BY 1  
NO MINVALUE  
NO MAXVALUE  
CACHE 1;

alter table "bi_system"."kettle_job" alter column id set default nextval('kettle_job_id_seq' );


--
--
--kettle_job_log
--
-- --------------------------------
CREATE TABLE "bi_system"."kettle_job_log" (
	"id" int8 NOT NULL DEFAULT nextval('kettle_job_log_id_seq'::regclass),
	"kettle_job_id" int8 NOT NULL,
	"error_info" text,
	"status" int2 NOT NULL
)
WITH (OIDS=FALSE);
ALTER TABLE "bi_system"."kettle_job_log" OWNER TO "user_logistic";

-- ----------------------------
--  Primary key structure for table kettle_job_log
-- ----------------------------
ALTER TABLE "bi_system"."kettle_job_log" ADD PRIMARY KEY ("id") NOT DEFERRABLE INITIALLY IMMEDIATE;

CREATE SEQUENCE "bi_system"."kettle_job_log_id_seq"  
START WITH 1  
INCREMENT BY 1  
NO MINVALUE  
NO MAXVALUE  
CACHE 1;

alter table "bi_system"."kettle_job_log" alter column id set default nextval('kettle_job_log_id_seq' );



--
-- PostgreSQL database dump complete
--

