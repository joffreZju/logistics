
CREATE TABLE "public".function
(
  id       SERIAL  NOT NULL
    CONSTRAINT function_pkey
    PRIMARY KEY,
  name     TEXT    NOT NULL,
  descrp   TEXT,
  pid      INTEGER NOT NULL,
  ctime    TIMESTAMP WITH TIME ZONE DEFAULT now(),
  path     TEXT,
  icon     TEXT,
  sys_id   TEXT,
  services TEXT []
);
-- ----------------------------
-- 重新部署OA应用时，要先初始化public.functions表，使用如下的insert语句。
-------------------------------
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('根节点', null, 0, '2017-07-10 13:06:36.641041', '1', null, 'root', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('组织管理', null, 1, '2017-07-10 13:06:36.697595', '1-2', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('我的工作', null, 1, '2017-07-10 13:06:36.738168', '1-3', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('系统设置', null, 1, '2017-07-10 13:06:36.777182', '1-4', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('组织属性维护', null, 2, '2017-07-10 13:06:36.819808', '1-2-5', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('组织树管理', null, 2, '2017-07-10 13:06:36.866494', '1-2-6', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('组织树查询', null, 2, '2017-07-10 13:06:36.905274', '1-2-7', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('我的审批', null, 3, '2017-07-10 13:06:36.944919', '1-3-8', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('我的申请', null, 3, '2017-07-10 13:06:36.985430', '1-3-9', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('审批流设定', null, 4, '2017-07-10 13:06:37.027943', '1-4-10', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('用户管理', null, 4, '2017-07-10 13:06:37.070101', '1-4-11', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('角色管理', null, 4, '2017-07-10 13:06:37.111124', '1-4-12', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('审批表单设定', null, 4, '2017-07-10 13:06:37.149684', '1-4-13', null, 'oa', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('蜂群 BI', '', 1, '2017-07-11 13:01:05.138858', '1-14', '', 'bi_admin', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('报表管理', '', 14, '2017-07-11 13:02:55.398579', '1-14-15', '', 'bi_admin', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('报表开发', '', 14, '2017-07-11 13:03:01.767794', '1-14-16', '', 'bi_admin', null);
INSERT INTO public.function (name, descrp, pid, ctime, path, icon, sys_id, services) VALUES ('报表测试', '', 14, '2017-07-11 13:03:11.112201', '1-14-17', '', 'bi_admin', null);