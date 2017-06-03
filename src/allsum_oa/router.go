package main

import (
	"allsum_oa/controller/file"
	"allsum_oa/controller/user"
	"common/filter"
	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
	FirmPrefix   string = "/v2/admin"
	FilePrefix   string = "/v2/file"
)

func LoadRouter() {

	// user相关
	beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "*:GetCode")
	beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "*:UserRegister")
	beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "*:UserLogin")
	//beego.Router(ExemptPrefix+"/user/login_phone", &user.Controller{}, "*:UserLoginPhoneCode")
	//beego.Router(UserPrefix+"/login_out", &user.Controller{}, "*:LoginOut")
	beego.Router(UserPrefix+"/info", &user.Controller{}, "*:GetUserInfo")
	beego.Router(UserPrefix+"/passwd/modify", &user.Controller{}, "*:Resetpwd")
	beego.Router(UserPrefix+"/edit_profile", &user.Controller{}, "*:EditProfile")

	//comapny相关
	beego.Router(FirmPrefix+"/info", &user.Controller{}, "Get:GetFirmInfo")
	beego.Router(FirmPrefix+"/register", &user.Controller{}, "Post:FirmRegister")
	beego.Router(FirmPrefix+"/list", &user.Controller{}, "Get:GetFirmList")
	beego.Router(FirmPrefix+"/modify", &user.Controller{}, "Post:FirmModify")
	beego.Router(FirmPrefix+"/audit", &user.Controller{}, "Post:FirmAudit")
	beego.Router(FirmPrefix+"/add_user", &user.Controller{}, "Post:FirmAddUser")
	beego.Router(FirmPrefix+"/del_user", &user.Controller{}, "Post:FirmDelUser")

	//文件上传下载
	beego.Router(FilePrefix+"/upload", &file.Controller{}, "Post:UploadFile")
	beego.Router(FilePrefix+"/download", &file.Controller{}, "*:DownloadFile")

	// 非登录态列表
	//notNeedAuthList := []string{
	//	// aliyun check
	//	"/",
	//	// user
	//	ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register", ExemptPrefix + "/user/login",
	//}

	// add filter
	// 请求合法性验证 这个要放在第一个
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
