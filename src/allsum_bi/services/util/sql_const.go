package util

//SQL

const CREATE_USER_AUTHOIRTY = "CREATE SEQUENCE {SCHEMA_NAME}.user_authority_id_seq START WITH 1 INCREMENT BY 1 NO MINVALUE NO MAXVALUE CACHE 1; " +
	"ALTER TABLE {SCHEMA_NAME}.user_authority_id_seq OWNER TO user_logistic;" +
	"CREATE TABLE {SCHEMA_NAME}.user_authority " +
	"(id bigint DEFAULT nextval('{SCHEMA_NAME}.user_authority_id_seq'::regclass) NOT NULL, " +
	"roleid bigint NOT NULL, " +
	"rolename varchar(128) NOT NULL, " +
	"reportid bigint, " +
	"reportsetids bigint[], " +
	"createtime timestamp, " +
	"limittime int); " +
	"ALTER TABLE {SCHEMA_NAME}.user_authority OWNER TO user_logistic; " +
	"COMMENT ON COLUMN {SCHEMA_NAME}.user_authority.roleid IS '角色id'; " +
	"COMMENT ON COLUMN {SCHEMA_NAME}.user_authority.limittime IS '角色的权限的有效时间 单位为天 为0 时表示无限期'; "
