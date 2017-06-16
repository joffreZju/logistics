### 用户相关接口

- 获取验证码（get，post）
    - request
        - tel:string（电话号码）
    - response
        - code:int
        - msg:string
        - data:json object 或 null
- 注册（post）
    - request
        - smscode:string(验证码，必填)
        - tel:string（电话，必填）
        - password:string（密码，必填）
        - addr:string
        - desc:string
        - gender:int(0/1)

- 获取用户公司列表（get，post）
    - 登录之前先请求获取用户的注册公司列表，然后选择其中一个来登录。
    - request
        - tel:string

- 密码登录（post）
    - 注意：需要登录之后的其他请求清携带此header中的token
    - request
        - companyNo:string(用户没有注册公司的话，选填)
        - tel:string
        - password:string
    - response
        - http header(Authorization:access_token boIcCYrQ9XQHS69fj4uQBG+5yGEyZ6PNdrFGTYDfDXw=)
        - code
        - msg
        - data:object

- 验证码登录（post）
    - request
        - companyNo:string(用户没有注册公司的话，选填)
        - tel:string
        - password:string
    - response
        - http header(Authorization:access_token boIcCYrQ9XQHS69fj4uQBG+5yGEyZ6PNdrFGTYDfDXw=)
        - code
        - msg
        - data:object（同上）
    
- 登出（直接携带token发get或post请求即可）

### 文件上传下载

- 上传（post，暂仅支持单个文件）
    - request
        - form-data(key:file,value:文件)
        - header(Content-Type:multipart/form-data; boundary=----WebKitFormBoundary7MA4YWxkTrZu0gW)
    
- 下载（get,post单个文件）
    - request
        - url:string(从后台获取到的url)
    - response
        - 开始下载或返回错误码


### 组织树相关接口

- 更新或增加组织属性
    - Id:int(更新)
    - Name:string
    - Desc:string

- 新增上下级
    - Group:{Name:"abc",Pid:int,AttrId:int}
    - Sons:[]int
    - 新增上级：sons为选中所有下一级子节点的id，同时这些节点的parent_id都一样且等于new_node.parent_id
    - 要修改所有sons的子孙节点的path
    - 新增下级：sons字段为空

- 合并
    - OldGroups:[]int
    - NewGroup:{Name:"abc",ParentId:int}
    - 同级，old_nodes必须有共同的parent_id,且等于new_node.parent_id
    
- 转让、升级
    - GroupId:int
    - NewPid:int

- 删除(软/硬)
    - GroupId:int
    - 有用户在该group的话不能删除
    - 有子节点：所有子节点自动升级
    - 没有子节点：直接删除

- 编辑
    - GroupId:int
    - NewName:string
    
- 向组织添加用户
    - GroupId:int
    - Users:[]int
    - 已经在该组织的用户会忽略掉

- 从组织批量删除用户
    - GroupId:int
    - Users:[]int


### 角色相关接口

- 添加角色
    - Name:string
    - Desc:string
    - FuncIds:[]int
    - 所选的功能树下所有的功能

- 修改角色
    - Id:int
    - Name:string
    - Desc:string
    - FuncIds:[]int
    - 所选的功能树下所有的功能

- 删除角色
    - Id:int
    - 如果该角色下还有用户，那么删除失败

- 向角色添加用户
    - RoleId:int
    - Users:[]int
    - 已经在该角色的用户会忽略掉

- 从角色中批量删除用户
    - RoleId:int
    - Users:[]int


### 管理员管理审批相关接口

##### 表单模板接口
- 添加表单模板
    - formtpl:{
        - Name:string
        - Type:string
        - Desc:string
        - Content:{json string}
        - Attachment:[]file (单独传？！)
        - BeginTime:time.Time string 2017-01-01 00:00:00
    }

- 更新表单模板
    - formtpl:{
        - No:string
        - Name:string
        - Type:string
        - Desc:string
        - Content:{json string}
        - Attachment:[]file ???
        - BeginTime:time.Time string 2017-01-01 00:00:00
    }

- 启用/禁用(1/2)表单模板
    - No:string
    - Status:1/2
    - 禁用表单模板，那么相对应的审批流模板approvaltpl也会被禁用
    
- 删除表单模板
    - No:string
    - 删除前需要确认是否有对应的审批流模板，如果删除formtpl,那么对用的approvalTpl会被禁用


##### 审批单模板接口
- 添加审批单模板
    - approvaltpl:{
        - Name:string
        - Type:string
        - Desc:string
        - FormtplNo:string
        - TreeFlowTag:int(是否按组织树向上流动0:no,1:yes)
        - RoleFlow:[]int
        - AllowRows:[]int
        - BeginTime:time.Time string 2017-01-01 00:00:00    
    }

- 更新审批单模板
    - approvaltpl:{
        - No:string
        - Name:string
        - Type:string
        - Desc:string
        - FormtplNo:string
        - TreeFlowTag:int(是否按组织树向上流动0:no,1:yes)
        - RoleFlow:[]int
        - AllowRows:[]int
        - BeginTime:time.Time string 2017-01-01 00:00:00    
    }

- 启用/禁用(1/2)审批单模板
    - No:string
    - Status:1/2
    
- 删除审批单模板
    - No:string

### 员工增删改审批单相关接口

- 增加审批单
    - approval:{
        - Name:string
        - Desc:string
        - UserFlow:[]int
        - UserId:int
        - RoleId:int
        - GroupId:int
        - Status:int(0草稿,1提交)
    }
    - form:{
        - Name:string
        - Type:string
        - Desc:string
        - Content:string
        - Attachment:[]file
    }
    
- 修改审批单（有限定）
    - approval:{
        - No:string
        - Name:string
        - Desc:string
        - UserFlow:[]int
        - UserId:int
        - RoleId:int
        - GroupId:int
        - Status:int(0草稿,1提交)
    }
    - form:{
        - No:string
        - Name:string
        - Type:string
        - Desc:string
        - Content:string
        - Attachment:[]file
    }
    
- 撤销审批单（有限定）
    - No:string
    
- 审批操作（同意或拒绝）
    - approve:{
        - ApprovalNo:string
        - UserId:int(审批人id)
        - Opinion:int(0:同意，1:拒绝,必填)
        - Comment:string(选填)
    }