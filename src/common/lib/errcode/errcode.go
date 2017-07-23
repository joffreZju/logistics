package errcode

import "fmt"

type CodeError struct {
	Code int    `json:"code"`
	Msg  string `json:"msg"`
}

func (ce *CodeError) Error() string {
	return fmt.Sprintf("%d: %s", ce.Code, ce.Msg)
}

func New(code int, msg string) *CodeError {
	return &CodeError{code, msg}
}

var (
	ErrParams             = &CodeError{10000, "参数错误"}
	ErrCheckRequestFailed = &CodeError{10001, "URL请求不合法"}
	ErrRequestExpired     = &CodeError{10002, "URL请求过期"}
	ErrServerError        = &CodeError{10003, "服务器繁忙，请稍后重试"}
	ErrParamTime          = &CodeError{10004, "时间格式错误"}
	// user
	ErrGetUserInfoFailed            = &CodeError{20000, "获取用户信息失败"}
	ErrUserNotExisted               = &CodeError{20001, "用户不存在"}
	ErrUserAlreadyExisted           = &CodeError{20002, "用户已经存在"}
	ErrUserPasswordError            = &CodeError{20003, "用户密码错误~"}
	ErrBindTelFailed                = &CodeError{20004, "绑定手机号失败"}
	ErrUserUpdateFailed             = &CodeError{20005, "更新用户信息失败"}
	ErrUserCreateFailed             = &CodeError{20006, "新建用户失败"}
	ErrUserUploadPicFailed          = &CodeError{20007, "上传图片失败"}
	ErrOurUserTelHasAlreadyRegisted = &CodeError{20008, "手机号已经注册,请登录"}
	ErrGetOurUserByTelFailed        = &CodeError{20009, "未查询到用户请先联系管理员注册"}
	ErrOurUserGetAuthFailed         = &CodeError{20010, "获取签名信息失败，请重试"}
	ErrSendSMSMsgError              = &CodeError{20011, "发送短消息失败，请稍后重试"}
	ErrUserNeedInit                 = &CodeError{20012, "用户xu yao"}
	ErrUserLocked                   = &CodeError{20013, "用户被锁定了"}
	ErrGroupOfUser                  = &CodeError{20014, "用户组织错误"}
	ErrGetLoginInfo                 = &CodeError{20015, "获取登录信息失败,请重新登录"}
	ErrUpdateGroupTree              = &CodeError{20016, "当前组织树有未生效的修改"}
	ErrRoleOfUser                   = &CodeError{20017, "用户角色错误"}
	ErrStatOfApproval               = &CodeError{20018, "审批单状态错误"}
	ErrInfoOfUser                   = &CodeError{20019, "用户信息错误"}

	ErrAuthCreateFailed         = &CodeError{20101, "出问题了，稍后再试吧~"}
	ErrAuthCheckFailed          = &CodeError{20102, "出问题了，稍后再试吧~"}
	ErrAuthCodeError            = &CodeError{20103, "验证码错误"}
	ErrAuthCodeExpired          = &CodeError{20104, "验证码已经失效"}
	ErrUserCodeHasAlreadyExited = &CodeError{20106, "验证码已经发送，请60秒后重试"}
	ErrUserPremissionError      = &CodeError{20107, "您没有足够的权限查看该数据！"}
	ErrFirmCreateFailed         = &CodeError{20120, "新建企业失败"}
	ErrFirmNotExisted           = &CodeError{20121, "企业不存在"}
	ErrFirmUpdateFailed         = &CodeError{20122, "更新企业信息失败"}

	ErrCreateOrderFailed       = &CodeError{20131, "创建订单失败"}
	ErrCreateOrderStatusFailed = &CodeError{20132, "创建订单状态失败"}
	ErrGetBillFailed           = &CodeError{20140, "获取账单失败"}
	ErrCreateBillFailed        = &CodeError{20141, "创建账单失败"}

	ErrUploadFileFailed   = &CodeError{20150, "文件上传失败"}
	ErrDownloadFileFailed = &CodeError{20153, "文件下载失败"}
	ErrFileNotExist       = &CodeError{20151, "文件不存在"}
	ErrUploadDocFailed    = &CodeError{20152, "文档上传失败"}

	ErrCouponExist    = &CodeError{20155, "存在重复编号"}
	ErrCouponNo       = &CodeError{20156, "号段错误"}
	ErrCouponVerify   = &CodeError{20157, "核销码错误"}
	ErrCouponUsed     = &CodeError{20158, "代金券已使用"}
	ErrCouponIllegal  = &CodeError{20159, "非法券"}
	ErrCouponNotExist = &CodeError{20159, "券不存在"}

	//bi
	ErrActionGetReport           = &CodeError{20500, "获取报表信息出错"}
	ErrActionGetAggregate        = &CodeError{20501, "获取清洗信息出错"}
	ErrActionPutAggregate        = &CodeError{20502, "插入清洗信息出错"}
	ErrActionGetDataload         = &CodeError{20503, "获取数据录入出错"}
	ErrActionPutDataload         = &CodeError{20504, "添加数据录入项出错"}
	ErrActionInputData           = &CodeError{20505, "录入数据出错"}
	ErrActionGetDbMgr            = &CodeError{20506, "获取数据库信息失败"}
	ErrActionPutDbMgr            = &CodeError{20507, "插入据库信息失败"}
	ErrActionCreateConn          = &CodeError{20508, "创建数据库链接失败"}
	ErrActionDeleteDbMgr         = &CodeError{20509, "删除数据库链接失败"}
	ErrActionNoAuthority         = &CodeError{20510, "无权执行此操作"}
	ErrActionGetDemand           = &CodeError{20511, "获取需求信息失败"}
	ErrActionPutDemand           = &CodeError{20512, "添加需求信息失败"}
	ErrActionPutReport           = &CodeError{20513, "添加报表信息失败"}
	ErrActionPutSycn             = &CodeError{20514, "添加同步任务失败"}
	ErrActionGetSycn             = &CodeError{20515, "获取同步任务信息失败"}
	ErrActionGetSchemaTable      = &CodeError{20516, "获取表信息失败"}
	ErrActionGetReportSet        = &CodeError{20517, "获取报表设置信息失败"}
	ErrActionPutReportSet        = &CodeError{20518, "添加报表设置信息失败"}
	ErrActionGetReportData       = &CodeError{20519, "获取报表数据失败"}
	ErrActionPutTestData         = &CodeError{20520, "添加测试数据失败"}
	ErrActionGetTestInfo         = &CodeError{20521, "获取测试数据失败"}
	ErrActionGetJobInfo          = &CodeError{20522, "获取KETTLE任务数据失败"}
	ErrActionAddJobInfo          = &CodeError{20523, "添加KETTLE任务数据失败"}
	ErrActionSetJobNum           = &CodeError{20524, "设置KETTLEJOB数据出错"}
	ErrActionAddUserAuthority    = &CodeError{20525, "设置用户报表权限错误"}
	ErrActionGetUserAuthority    = &CodeError{20526, "获取用户报表权限"}
	ErrActionDeleteUserAuthirity = &CodeError{20527, "删除用户报表权限"}
)

func ParseError(err error) (code int, msg string) {
	if e, ok := err.(*CodeError); ok {
		return e.Code, e.Msg
	}
	return ErrServerError.Code, ErrServerError.Msg
}
