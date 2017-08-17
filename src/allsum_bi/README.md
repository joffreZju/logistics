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


# 目标

- 一个报表快速开发工具

- 简易的etl抽取工具和kettle抽取工具的混合引入

- 链接OA实现权限控制


# 开发进度
- [x] 多数据库链接管理 
- [x] go etltransport的引入
- [x] 需求流程业务开发
- [x] 清洗任务
- [x] 数据录入
- [x] 报表设置
- [x] 各功能api开发 
- [x] kettle工具引入
- [x] 通用报表支持( 抽取操作涉及到超大表操作要谨慎)
- [x] 数据库访问权限赋予
- [x] demandlist查看权限
- [x] 任务进程池和进程队列
- [x] 邮件提醒 
- [ ] 用户报表管理（订阅 等）





修改记录
开发阶段：

测试阶段：
 8/12: mv util folder to  sevices & etl skip compare with string & add xmin in dest table for transporter
 8/14: 给etl抽取的地方加上了 最后一个修改的开发人员的字段 方便定位责任人。  抽取任务错误达到上限时 自动终止任务 增加了给开发人员的邮件提醒，测试部分 测试人员反馈和开发人员确认反馈的地方加上了邮件提醒
修复了并发读写map时未上锁所引起的奔溃bug
8/15: 压测etl时，我创建了很多名字相近的表，发现我用的transporter工具在从源表获取表名时用的模糊匹配的方法， 我将其修改了让其指定源表表名，然后目标表名可以任意定义，这样就更符合咱们的业务场景
8/16: 发现数据库更新struct的时候，如果int类型的字段为0时默认不更新， 这里个问题修复一下
8/17: Go transporter 的抽取在做小表同步的时候很方便，但是做大表同步性能不佳 大表建议用kettle
