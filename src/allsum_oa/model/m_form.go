package model

import (
	"time"
)

const (
	TplInit = iota
	TplAbled
	TplDisabled
)

const (
	ApproveDraft = iota
	ApproveCommited
	Approving
	ApproveAccessed
	ApproveNotAccessed
	ApproveCanceled
	//草稿->提交->审批中->审批通过（或不通过）->取消审批(除了审批通过或不通过，其他都可以取消)
)

const (
	ApproveOpinionAgree = iota
	ApproveOpinionRefuse
)

const (
	TreeFlowUpNo = iota
	TreeFlowUpYes
)

type Formtpl struct {
	No         string `gorm:"primary_key"`
	Name       string `gorm:"not null"`
	Type       string
	Desc       string
	Ctime      time.Time
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
	Desc       string
	Ctime      time.Time
	Content    string   `gorm:"type:jsonb;not null"`
	Attachment StrSlice `gorm:"type:text[]"`
}

func (Form) TableName() string {
	return "form"
}

type Approvaltpl struct {
	No         string `gorm:"primary_key"`
	Name       string `gorm:"not null"`
	Desc       string
	Ctime      time.Time
	FormtplNo  string    `gorm:"not null"`
	TreeFlowUp int       //是否按组织树向上流动 0:否，1:是
	RoleFlow   IntSlice  `gorm:"type:int[]"` //role_id 的组成的数组
	AllowRows  IntSlice  `gorm:"type:int[]"`
	BeginTime  time.Time `gorm:"not null"`
	Status     int       `gorm:"not null"`

	FormtplContent Formtpl `gorm:"-"`
}

func (Approvaltpl) TableName() string {
	return "approvaltpl"
}

type Approval struct {
	No       string `gorm:"primary_key"`
	Name     string `gorm:"not null"`
	Desc     string
	Ctime    time.Time
	FormNo   string   `gorm:"not null"`
	UserFlow IntSlice `gorm:"type:int[]"` //创建审批单时确定具体审批人
	UserId   int      `gorm:"not null"`
	RoleId   int      `gorm:"not null"`
	GroupId  int      `gorm:"not null"`
	Status   int      `gorm:"not null"`

	FormContent Form `gorm:"-"`
}

func (Approval) TableName() string {
	return "approval"
}

type ApproveFlow struct {
	Id         int    `gorm:"AUTO_INCREMENT,primary_key"`
	ApprovalNo string `gorm:"not null"`
	UserId     int
	//RoleId     int `gorm:"not null"`
	Opinion int `gorm:"not null"`
	Comment string
}

func (ApproveFlow) TableName() string {
	return "approve_flow"
}
