package main

import (
	"allsum_oa/controller/file"
	"allsum_oa/controller/form"
	"allsum_oa/controller/group"
	"allsum_oa/controller/role"
	"allsum_oa/controller/user"
	"common/filter"
	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
	FirmPrefix   string = "/v2/firm"
	FilePrefix   string = "/v2/file"
	AdminPrefix  string = "/v2/admin"
)

func LoadRouter() {

	// user相关
	beego.Router(ExemptPrefix+"/test", &user.Controller{}, "*:Test")
	beego.Router(ExemptPrefix+"/get_functions", &user.Controller{}, "*:GetFunctionsTree")
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "*:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "*:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "*:UserLogin")
	beego.Router(ExemptPrefix+"/user/login_phone", &user.Controller{}, "*:UserLoginPhone")
	beego.Router(ExemptPrefix+"/user/login_out", &user.Controller{}, "*:LoginOut")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "*:GetUserInfo")
	beego.Router(UserPrefix+"/resetpwd", &user.Controller{}, "*:Resetpwd")
	beego.Router(UserPrefix+"/switch_company", &user.Controller{}, "*:SwitchCurrentFirm")

	//allsum管理员审核公司
	beego.Router(AdminPrefix+"/firm_info", &user.Controller{}, "*:AdminGetFirmInfo")
	beego.Router(AdminPrefix+"/firm_list", &user.Controller{}, "*:AdminGetFirmList")
	beego.Router(AdminPrefix+"/firm_audit", &user.Controller{}, "*:AdminFirmAudit")

	//文件上传下载
	beego.Router(FilePrefix+"/upload", &file.Controller{}, "Post:UploadFile")
	//beego.Router(FilePrefix+"/download", &file.Controller{}, "*:DownloadFile")

	//comapny相关
	//beego.Router(FirmPrefix+"/register", &user.Controller{}, "Post:FirmRegister")
	//beego.Router(FirmPrefix+"/modify", &user.Controller{}, "Post:FirmModify")
	beego.Router(FirmPrefix+"/add_license", &user.Controller{}, "*:AddLicenseFile")
	beego.Router(FirmPrefix+"/add_user", &user.Controller{}, "Post:FirmAddUser")
	beego.Router(FirmPrefix+"/update_user", &user.Controller{}, "Post:UpdateUserProfile")
	//组织树
	beego.Router(FirmPrefix+"/add_attr", &group.Controller{}, "*:AddAttr")
	beego.Router(FirmPrefix+"/get_attrs", &group.Controller{}, "*:GetAttrList")
	beego.Router(FirmPrefix+"/update_attr", &group.Controller{}, "*:UpdateAttr")
	beego.Router(FirmPrefix+"/add_group", &group.Controller{}, "Post:AddGroup")
	beego.Router(FirmPrefix+"/merge_groups", &group.Controller{}, "Post:MergeGroups")
	beego.Router(FirmPrefix+"/move_group", &group.Controller{}, "*:MoveGroup")
	beego.Router(FirmPrefix+"/del_group", &group.Controller{}, "*:DelGroup")
	beego.Router(FirmPrefix+"/update_group", &group.Controller{}, "Post:UpdateGroup")
	beego.Router(FirmPrefix+"/get_groups", &group.Controller{}, "*:GetGroupList")
	beego.Router(FirmPrefix+"/getusers_ofgroup", &group.Controller{}, "*:GetUsersOfGroup")
	beego.Router(FirmPrefix+"/addusers_togroup", &group.Controller{}, "Post:AddUsersToGroup")
	beego.Router(FirmPrefix+"/delusers_fromgroup", &group.Controller{}, "Post:DelUsersFromGroup")
	//角色
	beego.Router(FirmPrefix+"/get_roles", &role.Controller{}, "*:GetRoleList")
	beego.Router(FirmPrefix+"/add_role", &role.Controller{}, "Post:AddRole")
	beego.Router(FirmPrefix+"/update_role", &role.Controller{}, "Post:UpdateRole")
	beego.Router(FirmPrefix+"/del_role", &role.Controller{}, "*:DelRole")
	beego.Router(FirmPrefix+"/getusers_ofrole", &role.Controller{}, "*:GetUsersOfRole")
	beego.Router(FirmPrefix+"/addusers_torole", &role.Controller{}, "Post:AddUsersToRole")
	beego.Router(FirmPrefix+"/delusers_fromrole", &role.Controller{}, "Post:DelUsersFromRole")

	//审批相关
	//表单模板操作
	beego.Router(FirmPrefix+"/add_formtpl", &form.Controller{}, "Post:AddFormtpl")
	beego.Router(FirmPrefix+"/get_formtpls", &form.Controller{}, "*:GetFormtplList")
	beego.Router(FirmPrefix+"/update_formtpl", &form.Controller{}, "Post:UpdateFormtpl")
	beego.Router(FirmPrefix+"/control_formtpl", &form.Controller{}, "*:ControlFormtpl")
	beego.Router(FirmPrefix+"/del_formtpl", &form.Controller{}, "*:DelFormtpl")
	//审批单模板操作
	beego.Router(FirmPrefix+"/add_atpl", &form.Controller{}, "Post:AddApprovaltpl")
	beego.Router(FirmPrefix+"/get_atpls", &form.Controller{}, "*:GetApprovaltplList")
	beego.Router(FirmPrefix+"/update_atpl", &form.Controller{}, "Post:UpdateApprovaltpl")
	beego.Router(FirmPrefix+"/control_atpl", &form.Controller{}, "*:ControlApprovaltpl")
	beego.Router(FirmPrefix+"/del_atpl", &form.Controller{}, "*:DelApprovaltpl")
	//审批流相关
	beego.Router(FirmPrefix+"/add_approval", &form.Controller{}, "*:AddApproval")
	beego.Router(FirmPrefix+"/update_approval", &form.Controller{}, "*:UpdateApproval")

	// 非登录态列表
	notNeedAuthList := []string{
		// aliyun check
		//"/",
		// user
		ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register",
		ExemptPrefix + "/user/login", ExemptPrefix + "/user/login_phone",
		ExemptPrefix + "/user/login_out", ExemptPrefix + "/test",
		ExemptPrefix + "/get_functions",
	}

	// add filter
	// 请求合法性验证 这个要放在第一个
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
