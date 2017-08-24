## 代码结构
- controller所有的控制器
    - admin:allsum管理员相关路由处理 
    - api:系统间调用路由
    - file:文件上传
    - firm:注册公司管理员相关操作
    - form:表单以及审批单相关操作
    - group:组织树相关操作
    - msg:用户消息相关操作
    - public:一些公共接口
    - role:角色管理相关接口
    - user:基本用户相关接口
    
- service复杂的数据库操作逻辑，以及一些定时任务等
    - s_initAllsum 部署应用时，自动初始化allsum所在的schema
    - s_ticker 一些定时任务，包括组织树修改生效，审批单设定生效等
    
- model数据层
    - m_init 初始化数据库连接并使用AutoMigrate初始化数据表。
    - m_common 封装了StrSlice,IntSlice,JsonMap对应postgres的数组和jsonb，使用更方便。

## 部署
- 部署方式同样可以采用nohup & 或者 systemd，应用启动之前，先要执行inniFunctions.sql来初始化功能配置表


## 一些系统逻辑
- 用户注册一个公司后，状态为待审批，管理员审核通过后，系统会自动创建该公司对应的schema，并建立管理员角色赋予其所有的功能权限。
- allsum管理员可以管理配置系统功能菜单。
- 用户登录时会将其公司信息，组织信息，角色信息，权限信息存放在redis中。
- 公司的管理员可以对角色进行管理。
- 公司的管理员可以对组织树进行管理，对组织树的修改可以设置生效时间，同时还要保存每一次修改的快照，快照以json的格式保存在group_operation表中，
然后由定时任务将其扫描到group表中,未生效的组织树修改操作可以取消。

- 公司管理员可以管理公司员工的信息。

## redis相关
- 每个用户登录之后会在redis中存在两条记录
```
1. key是token的加密字符串，value是token的值，其中SingleID是UserId

key:vEiFDZIkmq7Mr64c31fU8mQkuEn/l+N87wEbSGliQV8=
value:{
  "ClientID": "stowage_user",
  "SingleID": "145",
  "GroupID": "",
  "Value": "vEiFDZIkmq7Mr64c31fU8mQkuEn/l+N87wEbSGliQV8=",
  "BizInfo": "",
  "DeadLine": 1503759573
}


2. key是UserId拼接其token加密串，value是用户的公司，组织，角色，功能权限集，注意字符串的拼接方式（首尾都有 - ），可以快速查找

key:145-vEiFDZIkmq7Mr64c31fU8mQkuEn/l+N87wEbSGliQV8=
value:{
    company:C0726105137846,
    roles:-1-2-3-,
    groups:-1-2-3-,
    functions:-serviceUrl1-serviceUrl2-
}
```
- redis中的信息会在用户新的登录操作时被覆盖，用户注销操作时会清空。


## 定时任务
- 每小时会扫描一次，扫描每一个schema中的group_operation,formtpl,approvaltpl三张表，如果其设置生效时间距现在不足1.5小时，那么就改为将其生效。


## 审批流逻辑
- 审批流程设定增加组织选定和是否必审功能。
- 选择组织：选择角色后，选择该角色对应的组织。审批时寻找该组织下的角色进行审批。不选组织，默认在发起人的组织树寻找该角色。
- 是否必审：勾选为必审时，流程中的哪个角色发起都要走到这个角色进行审批。审批流以外的都要从头开始审批。
- 审批走向：从第一个角色开始，在发起人的同一个组织树里，进行一步一步审批。
- 在发起的角色之前的角色，没有必审时不用审批直接跳过，往发起角色后面的审批角色进行审批。
- 审批流里面有两个相同角色的情况下，如果允许，那么当一个此角色的用户发起时从前面一个开始审批；如果不允许，那么前端验证下就可以。
- 必审的一步，一般要指定组织。如果没有指定组织，先从发起人上级找，没找到的话再从发起人下级找。


## 系统常量定义

```
User.Status{1:正常,2:锁定}
User.Gender{1:男,2:女}
User.UserType{1:普通用户}

Company.Status{1:注册,2:管理员审核通过,3:审核不通过,4:删除}

表单状态
Formtpl.Status{1:初始化,2:启用,3:禁用}

审批单模板状态
Approvaltpl.Status{1:初始化,2:启用,3:禁用}
Approvaltpl.EmailMsg{1:No,2:Yes}

流程是否必审
ApprovaltplFlow.Necessary{1:不必须,2:必审}

审批单状态
Approval.Status{1：正在审批，2：审批通过，3：审批不通过，4：审批取消，5：审批停止，无法进行下去（没有审批人）}
Approval.ApproveFlow.Status{1：正在审批，2：审批通过，3：审批不通过}
Approval.EmailMsg{1:No,2:Yes}

组织树操作记录
GroupOperation.Status{1:历史记录,2:未生效记录}

审批单查询可选字段
beginTime:{"2017-07-01T14:47:00+08:00"}
conditon:{"approving","finished"}

AppVersion{
    Environment: 1:开发2:测试3:预发布4:生产
    DownloadUrl: []string, 多个下载地址
    UpgradeType: 1:透明2:友好提示3:强制升级
}
```
## 部分errcode定义
```
token相关错误码
ERR_InvalidateToken = ValidationError{Code: 40001, Msg: "Invalid token"}
ERR_TokenEmpty      = ValidationError{Code: 41001, Msg: "Token is empty"}
ERR_TokenExpired    = ValidationError{Code: 42001, Msg: "Token is expired"}

验证码相关错误码
ErrAuthCodeError            = &CodeError{20103, "验证码错误"}
ErrAuthCodeExpired          = &CodeError{20104, "验证码已经失效"}
ErrUserCodeHasAlreadyExited = &CodeError{20106, "验证码已经发送，请60秒后重试"}
```
## 接口修改记录

#### 2017-08-11
- 增加审批流模板，编辑审批流模板两个接口都增加了EmailMsg字段,该字段对应常量值见上。


#### 2017-07-17

- MatchUsers修改为1-2-3, "-"拼接

- allsum_oa表单审批中的所有接口地址都修改了，其他接口地址会逐步改。表单模板部分只改了地址，接口内容没有改。

- 新增审批流模板时，根据用户已选的roleId,调用接口获取可选的所有组织列表，

- 新增和编辑审批流模板，删除skipBlankRole，treeflowup，roleflow字段，增加FlowContent

- 发起审批单删掉approval.Status 字段。

- 删掉approval.{TreeFlowUp, SkipBlankRole, RoleFlow, CurrentRole},增加CurrentFlow

- 审批单详情，返回所有的流程id正序排列，id小于approval.CurrentFlow的是已经走过的流程。


#### long long ago

- 注册加了username字段，加了获取历史消息接口,获取（发起，收到）的审批单列表

- 所有上传的文件名，系统统一在前面拼接了长度为36的uuid字符串，需要展示文件名的地方，下载文件后直接截掉即可。

- 系统所有desc（Desc）字段都更改为descrp，涉及修改的接口：
    - 修改传入参数：增加组织属性，更新组织属性，组织树修改的五个接口，增加角色，更新角色
    - 修改返回参数：接口返回数据涉及到以下数据结构的，Desc更改为Descrp
    ```
        user.Descrp
        function.Descp 
        company.Descrp 
        form.descrp 
        formtpl.descrp 
        approval.descrp 
        approvaltpl.descrp 
        role.descrp 
        attibute.descrp
    ```

- 组织树修改的五个接口全部加上descrp和beginTime(可选)，删除组织树：post方式，同样加字段。

- 公司增加AdminId字段，是公司管理员（老板）的userId。

- 上传公司资质文件改为 修改公司信息。url改为http://allsum.com:8094/v2/firm/update_firm_info

- 增加用户修改个人信息接口

- 审批单不能保存草稿，直接提交


## System Change Log

### 2017-08-11
- 所有schema下oa_approval, oa_approvaltpl增加email_msg字段

- 前端发送password字段加上了md5加密，密文传输。