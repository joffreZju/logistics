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