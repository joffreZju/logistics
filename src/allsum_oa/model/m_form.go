package model

import (
	"time"
)

//审批/表单模板的状态
const (
	TplInit = iota + 1
	TplAbled
	TplDisabled
)

//审批单生命周期状态
//审批中，审批通过，不通过，取消审批，停止(无法进行下午)(审批未完成之前都可以取消)
const (
	ApprovalStatWaiting = iota + 1
	ApprovalStatAccessed
	ApprovalStatNotAccessed
	ApprovalStatCanceled
	ApprovalStatStop
)

//审批单是否沿组织树流动
const (
	TreeFlowUpNo = iota + 1
	TreeFlowUpYes
)

//审批单是否跳过没有用户的角色
const (
	SkipBlankRoleNo = iota + 1
	SkipBlankRoleYes
)

//获取 正在审批的 或 已完成的 审批单
const (
	GetApprovalApproving = "approving"
	GetApprovalFinished  = "finished"
)

type Formtpl struct {
	No         string `gorm:"primary_key"`
	Name       string `gorm:"not null"`
	Type       string
	Descrp     string
	Ctime      time.Time `gorm:"default:current_timestamp"`
	Content    string    `gorm:"type:jsonb;not null"`
	Attachment StrSlice  `gorm:"type:text[]"`
	BeginTime  time.Time `gorm:"not null"`
	Status     int       `gorm:"not null"`
}

func (Formtpl) TableName() string {
	return "formtpl"
}

type Form struct {
	No         string `gorm:"primary_key"`
	Name       string `gorm:""`
	Type       string
	Descrp     string
	Ctime      time.Time `gorm:"default:current_timestamp"`
	Content    string    `gorm:"type:jsonb;not null"`
	Attachment StrSlice  `gorm:"type:text[]"`
}

func (Form) TableName() string {
	return "form"
}

type Approvaltpl struct {
	No            string `gorm:"primary_key"`
	Name          string `gorm:"not null"`
	Descrp        string
	Ctime         time.Time `gorm:"default:current_timestamp"`
	FormtplNo     string    `gorm:"not null"`
	TreeFlowUp    int       //是否按组织树向上流动 1:否，2:是
	SkipBlankRole int       //是否跳过空角色 1:否，2:是
	RoleFlow      IntSlice  `gorm:"type:int[]"` //role_id 的组成的数组
	AllowRoles    IntSlice  `gorm:"type:int[]"`
	BeginTime     time.Time `gorm:"not null"`
	Status        int       `gorm:"not null"`

	FormtplContent *Formtpl `gorm:"-"`
}

func (Approvaltpl) TableName() string {
	return "approvaltpl"
}

type Approval struct {
	No            string `gorm:"primary_key"`
	Name          string `gorm:"not null"`
	Descrp        string
	Ctime         time.Time `gorm:"default:current_timestamp"`
	FormNo        string    `gorm:"not null"`
	TreeFlowUp    int       //1:否，2:是
	SkipBlankRole int       //1:否，2:是
	RoleFlow      IntSlice  `gorm:"type:int[]"`
	CurrentRole   int       //当前正在审批的角色id(current_user是pg的关键字)
	UserId        int       `gorm:"not null"`
	RoleId        int       `gorm:""`
	GroupId       int       `gorm:""`
	UserName      string    `gorm:""`
	RoleName      string    `gorm:""`
	GroupName     string    `gorm:""`
	Status        int       `gorm:"not null"`

	FormContent  *Form          `gorm:"-"`
	ApproveFLows []*ApproveFlow `gorm:"-"`
}

func (Approval) TableName() string {
	return "approval"
}

//流转时创建flow（roleId，RoleName赋值），用户审批时更新userId，userName，status，comment字段
type ApproveFlow struct {
	Id         int    `gorm:"AUTO_INCREMENT,primary_key"`
	ApprovalNo string `gorm:"not null"`
	MatchUsers string // 满足条件的用户id拼接1_2_3_
	UserId     int
	UserName   string
	RoleId     int
	RoleName   string
	Status     int //只有三种状态1：正在审批，2：审批通过，3：审批不通过
	Comment    string
	Ctime      time.Time `gorm:"default:current_timestamp"`
}

func (ApproveFlow) TableName() string {
	return "approve_flow"
}
