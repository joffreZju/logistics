### 组织树相关接口

- 更新或增加组织属性
    - Update:"true" or "false"
    - No:string
    - Name:string
    - Desc:string

- 新增上下级
    - NewGroup:{Name:"abc",Pid:int,AttrId:int}
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


### 审批相关接口

##### 表单模板接口
- 添加表单模板
    - Name:string
    - Type:string
    - Desc:string
    - Content:{json string}
    - Attachment:[]file ???
    - BeginTime:time.Time string 2017-01-01 00:00:00

- 更新表单模板
    - No:string
    - Name:string
    - Type:string
    - Desc:string
    - Content:{json string}
    - Attachment:[]file ???
    - BeginTime:time.Time string 2017-01-01 00:00:00

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


- 启用/禁用(1/2)表单模板
    - No:string
    - Status:1/2
    - 禁用表单模板，那么相对应的审批流模板approvaltpl也会被禁用
    
- 删除表单模板
    - No:string
    - 删除前需要确认是否有对应的审批流模板，如果删除formtpl,那么对用的approvalTpl会被禁用

- 添加审批流模板

- 禁用/启用(1,2)审批流模板

- 删除审批流模板

- 更新审批流模板
