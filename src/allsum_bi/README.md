BI SYSTEM
程序部分：
controller:    API
db：           数据库链接和一些数据库基本操作
models：       数据库业务层操作
routers：      路由层
services:      数据服务层

数据服务层4大模块儿:
 需求业务：
       demand相关字眼的代码
 数据录入：
       dataload相关字眼的代码
 抽取业务：
       ETL. sync, kettle 相关代码
清洗业务：
       aggragate相关代码

报表数据获取：
       reportset相关代码 











修改记录
开发阶段：

测试阶段：
 8/12: mv util folder to  sevices & etl skip compare with string & add xmin in dest table for transporter
 8/14: 给etl抽取的地方加上了 最后一个修改的开发人员的字段 方便定位责任人。  抽取任务错误达到上限时 自动终止任务 增加了给开发人员的邮件提醒，测试部分 测试人员反馈和开发人员确认反馈的地方加上了邮件提醒
修复了并发读写map时未上锁所引起的奔溃bug
