package main

import (
	"allsum_oa/controller/file"
	"allsum_oa/controller/form"
	"allsum_oa/controller/group"
	"allsum_oa/controller/msg"
	"allsum_oa/controller/role"
	"allsum_oa/controller/user"
	"common/filter"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/plugins/cors"
)

const (
	ExemptPrefix      string = "/exempt"
	FilePrefix        string = "/v2/file"
	PublicPrefix      string = "/v2/public"
	UserPrefix        string = "/v2/user"
	FirmInfoPrefix    string = "/v2/firm/info"
	FirmUserPrefix    string = "/v2/firm/user"
	FirmGroupPrefix   string = "/v2/firm/group"
	FirmRolePrefix    string = "/v2/firm/role"
	FormtplPrefix     string = "/v2/firm/formtpl"
	ApprovaltplPrefix string = "/v2/firm/approvaltpl"
	ApprovalPrefix    string = "/v2/firm/approval"

	AdminPrefix string = "/v2/admin"
)

func LoadRouter() {
	//文件上传下载
	beego.Router(FilePrefix+"/upload", &file.Controller{}, "Post:UploadFile")

	// user相关
	beego.Router("/test", &user.Controller{}, "*:Test")
	beego.Router(ExemptPrefix+PublicPrefix+"/appversion/latest/get", &user.Controller{}, "*:GetLatestAppVersion")
	beego.Router(ExemptPrefix+PublicPrefix+"/functions/get", &user.Controller{}, "*:GetFunctionsTree")
	beego.Router(ExemptPrefix+PublicPrefix+"/smscode", &user.Controller{}, "*:GetCode")
	beego.Router(ExemptPrefix+UserPrefix+"/register", &user.Controller{}, "*:UserRegister")
	beego.Router(ExemptPrefix+UserPrefix+"/login", &user.Controller{}, "*:UserLogin")
	beego.Router(ExemptPrefix+UserPrefix+"/login_phone", &user.Controller{}, "*:UserLoginPhone")
	beego.Router(ExemptPrefix+UserPrefix+"/login_out", &user.Controller{}, "*:LoginOut")
	beego.Router(ExemptPrefix+UserPrefix+"/forgetpwd", &user.Controller{}, "*:Forgetpwd")
	beego.Router(UserPrefix+"/info/get", &user.Controller{}, "*:GetUserInfo")
	beego.Router(UserPrefix+"/info/update", &user.Controller{}, "*:UpdateUserInfo")
	beego.Router(UserPrefix+"/pwd/reset", &user.Controller{}, "*:Resetpwd")
	beego.Router(UserPrefix+"/company/switch", &user.Controller{}, "*:SwitchCurrentFirm")
	beego.Router(UserPrefix+"/msg/history/get", &msg.Controller{}, "*:GetHistoryMsg")
	beego.Router(UserPrefix+"/msg/latest/get", &msg.Controller{}, "*:GetLatestMsg")
	beego.Router(UserPrefix+"/msg/page/get", &msg.Controller{}, "*:GetMsgsByPage")
	beego.Router(UserPrefix+"/msg/del", &msg.Controller{}, "*:DelMsgById")

	//allsum管理员审核公司
	beego.Router(AdminPrefix+"/firms/get", &user.Controller{}, "*:AdminGetFirmList")
	beego.Router(AdminPrefix+"/firm/audit", &user.Controller{}, "*:AdminFirmAudit")
	beego.Router(AdminPrefix+"/function/add", &user.Controller{}, "Post:AdminAddFunction")
	beego.Router(AdminPrefix+"/function/update", &user.Controller{}, "Post:AdminUpdateFunction")
	beego.Router(AdminPrefix+"/function/del", &user.Controller{}, "*:AdminDelFunction")
	beego.Router(AdminPrefix+"/appversion/list/get", &user.Controller{}, "*:GetAppVersionList")
	beego.Router(AdminPrefix+"/appversion/add", &user.Controller{}, "*:AddAppVersion")

	//公司管理员相关
	beego.Router(FirmInfoPrefix+"/update", &user.Controller{}, "*:UpdateFirmInfo")
	beego.Router(FirmUserPrefix+"/list/get", &user.Controller{}, "*:FirmGetUserList")
	beego.Router(FirmUserPrefix+"/search", &user.Controller{}, "*:FirmSearchUsersByName")
	beego.Router(FirmUserPrefix+"/add", &user.Controller{}, "Post:FirmAddUser")
	beego.Router(FirmUserPrefix+"/profile/update", &user.Controller{}, "Post:FirmUpdateUserProfile")
	beego.Router(FirmUserPrefix+"/rolegroup/update", &user.Controller{}, "Post:FirmUpdateUserRoleAndGroup")
	beego.Router(FirmUserPrefix+"/control", &user.Controller{}, "Post:FirmControlUserStatus")
	//管理组织树
	beego.Router(FirmGroupPrefix+"/attr/add", &group.Controller{}, "*:AddAttr")
	beego.Router(FirmGroupPrefix+"/attr/list/get", &group.Controller{}, "*:GetAttrList")
	beego.Router(FirmGroupPrefix+"/attr/update", &group.Controller{}, "*:UpdateAttr")
	beego.Router(FirmGroupPrefix+"/attr/del", &group.Controller{}, "*:DelAttr")
	beego.Router(FirmGroupPrefix+"/add", &group.Controller{}, "Post:AddGroup")
	beego.Router(FirmGroupPrefix+"/merge", &group.Controller{}, "Post:MergeGroups")
	beego.Router(FirmGroupPrefix+"/move", &group.Controller{}, "*:MoveGroup")
	beego.Router(FirmGroupPrefix+"/del", &group.Controller{}, "*:DelGroup")
	beego.Router(FirmGroupPrefix+"/update", &group.Controller{}, "Post:UpdateGroup")
	beego.Router(FirmGroupPrefix+"/list/get", &group.Controller{}, "*:GetGroupList")
	beego.Router(FirmGroupPrefix+"/operation/list/get", &group.Controller{}, "*:GetGroupOpList")
	beego.Router(FirmGroupPrefix+"/operation/search", &group.Controller{}, "*:SearchGroupOpsByTime")
	beego.Router(FirmGroupPrefix+"/operation/detail/get", &group.Controller{}, "*:GetGroupOpDetail")
	beego.Router(FirmGroupPrefix+"/operation/cancel", &group.Controller{}, "*:CancelGroupOp")
	beego.Router(FirmGroupPrefix+"/users/get", &group.Controller{}, "*:GetUsersOfGroup")
	beego.Router(FirmGroupPrefix+"/users/add", &group.Controller{}, "Post:AddUsersToGroup")
	beego.Router(FirmGroupPrefix+"/users/del", &group.Controller{}, "Post:DelUsersFromGroup")
	//管理角色
	beego.Router(FirmRolePrefix+"/list/get", &role.Controller{}, "*:GetRoleList")
	beego.Router(FirmRolePrefix+"/add", &role.Controller{}, "Post:AddRole")
	beego.Router(FirmRolePrefix+"/update", &role.Controller{}, "Post:UpdateRole")
	beego.Router(FirmRolePrefix+"/del", &role.Controller{}, "*:DelRole")
	beego.Router(FirmRolePrefix+"/users/get", &role.Controller{}, "*:GetUsersOfRole")
	beego.Router(FirmRolePrefix+"/users/add", &role.Controller{}, "Post:AddUsersToRole")
	beego.Router(FirmRolePrefix+"/users/del", &role.Controller{}, "Post:DelUsersFromRole")

	//审批相关
	//管理表单模板
	beego.Router(FormtplPrefix+"/add", &form.Controller{}, "Post:AddFormtpl")
	beego.Router(FormtplPrefix+"/get", &form.Controller{}, "*:GetFormtplList")
	beego.Router(FormtplPrefix+"/update", &form.Controller{}, "Post:UpdateFormtpl")
	beego.Router(FormtplPrefix+"/control", &form.Controller{}, "*:ControlFormtpl")
	beego.Router(FormtplPrefix+"/del", &form.Controller{}, "*:DelFormtpl")
	//审批单模板操作
	beego.Router(ApprovaltplPrefix+"/role/groups/match", &form.Controller{}, "*:GetMatchGroupsOfRole")
	beego.Router(ApprovaltplPrefix+"/add", &form.Controller{}, "Post:AddApprovaltpl")
	beego.Router(ApprovaltplPrefix+"/get", &form.Controller{}, "*:GetApprovaltplList")
	beego.Router(ApprovaltplPrefix+"/detail", &form.Controller{}, "*:GetApprovaltplDetail")
	beego.Router(ApprovaltplPrefix+"/update", &form.Controller{}, "Post:UpdateApprovaltpl")
	beego.Router(ApprovaltplPrefix+"/control", &form.Controller{}, "*:ControlApprovaltpl")
	beego.Router(ApprovaltplPrefix+"/del", &form.Controller{}, "*:DelApprovaltpl")
	//审批流相关
	beego.Router(ApprovalPrefix+"/add", &form.Controller{}, "*:AddApproval")
	beego.Router(ApprovalPrefix+"/cancel", &form.Controller{}, "*:CancelApproval")
	beego.Router(ApprovalPrefix+"/approve", &form.Controller{}, "*:Approve")
	beego.Router(ApprovalPrefix+"/from_me/get", &form.Controller{}, "*:GetApprovalsFromMe")
	beego.Router(ApprovalPrefix+"/to_me/get", &form.Controller{}, "*:GetApprovalsToMe")
	beego.Router(ApprovalPrefix+"/detail", &form.Controller{}, "*:GetApprovalDetail")

	beego.InsertFilter("*", beego.BeforeRouter, cors.Allow(&cors.Options{
		//AllowOrigins:     []string{"http://localhost:8090", "http://www.suanpeizaix.comw", "http://www.suanpeizaix.com:8090"},
		AllowAllOrigins:  true,
		AllowHeaders:     []string{"access_token", "Authorization", "X-Requested-With", "Content-Type", "Origin", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowMethods:     []string{"GET", "DELETE", "PUT", "PATCH", "POST", "OPTIONS"},
		ExposeHeaders:    []string{"Authorization", "Content-Type", "Content-Length", "Access-Control-Allow-Origin", "Access-Control-Allow-Headers"},
		AllowCredentials: true,
	}))

	// 非登录态列表
	notNeedAuthList := []string{
		ExemptPrefix + PublicPrefix + "/appversion/latest/get",
		ExemptPrefix + PublicPrefix + "/functions/get",
		ExemptPrefix + PublicPrefix + "/smscode",
		ExemptPrefix + UserPrefix + "/register",
		ExemptPrefix + UserPrefix + "/login",
		ExemptPrefix + UserPrefix + "/login_phone",
		ExemptPrefix + UserPrefix + "/login_out",
		ExemptPrefix + UserPrefix + "/forgetpwd",
	}

	// 请求合法性验证 这个要放在第一个
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
