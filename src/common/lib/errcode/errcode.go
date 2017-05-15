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
	ErrAuthCreateFailed             = &CodeError{20101, "出问题了，稍后再试吧~"}
	ErrAuthCheckFailed              = &CodeError{20102, "出问题了，稍后再试吧~"}
	ErrAuthCodeError                = &CodeError{20103, "验证码错误"}
	ErrAuthCodeExpired              = &CodeError{20104, "验证码已经失效"}
	ErrUserCodeHasAlreadyExited     = &CodeError{20106, "验证码已经发送，请60秒后重试"}
	ErrUserPremissionError          = &CodeError{20107, "您没有足够的权限查看该数据！"}
	ErrAgentCreatFailed             = &CodeError{20120, "新建代理商失败"}
	ErrAgentNotExisted              = &CodeError{20121, "代理商不存在"}

	ErrCreateOrderFailed       = &CodeError{20131, "创建订单失败"}
	ErrCreateOrderStatusFailed = &CodeError{20132, "创建订单状态失败"}
	ErrNoOrder                 = &CodeError{20133, "没有此订单"}
	ErrWXPay                   = &CodeError{20133, "微信支付暂不可用"}
	ErrALIPay                  = &CodeError{20134, "支付宝支付暂不可用"}
	ErrGetBillFailed           = &CodeError{20140, "获取账单失败"}
	ErrCreateBillFailed        = &CodeError{20141, "创建账单失败"}
	ErrAccountFundLow          = &CodeError{20142, "账户余额不足"}

	ErrUploadFileFailed = &CodeError{20150, "文件上传失败"}
	ErrFileNotExist     = &CodeError{20151, "文件不存在"}
	ErrUploadDocFailed  = &CodeError{20152, "文档上传失败"}

	ErrCouponExist    = &CodeError{20155, "存在重复编号"}
	ErrCouponNo       = &CodeError{20156, "号段错误"}
	ErrCouponVerify   = &CodeError{20157, "核销码错误"}
	ErrCouponUsed     = &CodeError{20158, "代金券已使用"}
	ErrCouponIllegal  = &CodeError{20159, "非法券"}
	ErrCouponNotExist = &CodeError{20160, "券不存在"}

	ErrNoTpl            = &CodeError{20170, "暂无匹配模板"}
	ErrTplIsNull        = &CodeError{20171, "模板字段为空或有误"}
	ErrCalResultIsNull  = &CodeError{20172, "计算结果还未返回"}
	ErrWrongCalNo       = &CodeError{20173, "计算记录号(CalNo)有误"}
	ErrCalPayFailed     = &CodeError{20173, "账户余额不足,计算扣费失败"}
	ErrNoFrequentCars   = &CodeError{20174, "暂无最近使用车辆"}
	ErrWrongJson        = &CodeError{20175, "cars,goods json格式有误"}
	ErrWrongCarsGoods   = &CodeError{20176, "计算数据有误"}
	ErrCalNoUserNoMatch = &CodeError{20177, "CalNo和用户不匹配"}
)

func ParseError(err error) (code int, msg string) {
	if e, ok := err.(*CodeError); ok {
		return e.Code, e.Msg
	}
	return ErrServerError.Code, ErrServerError.Msg
}
