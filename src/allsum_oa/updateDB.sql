-- 2017-08-11 所有schema下面的oa_approval, oa_approvaltpl增加字段：email_msg

ALTER TABLE "every_schema".oa_approval ADD COLUMN email_msg INTEGER;
ALTER TABLE "every_schema".oa_approvaltpl ADD COLUMN email_msg INTEGER;