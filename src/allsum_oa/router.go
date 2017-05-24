package main

import (
	"common/filter"
	"github.com/astaxie/beego"
)

const (
	ExemptPrefix string = "/exempt"
	UserPrefix   string = "/v2/user"
	ManagePrefix string = "/v2/admin"
)

func LoadRouter() {

	// user 相关
	//beego.Router(ExemptPrefix+"/user/getcode", &user.Controller{}, "*:GetCode")
	//beego.Router(ExemptPrefix+"/user/register", &user.Controller{}, "*:UserRegister")
	//beego.Router(ExemptPrefix+"/user/login", &user.Controller{}, "*:UserLogin")
	//beego.Router(ExemptPrefix+"/user/login_phone", &user.Controller{}, "*:UserLoginPhoneCode")
	//beego.Router(UserPrefix+"/login_out", &user.Controller{}, "*:LoginOut")
	//beego.Router(UserPrefix+"/info", &user.Controller{}, "*:GetUserInfo")
	//beego.Router(UserPrefix+"/passwd/modify", &user.Controller{}, "*:Resetpwd")
	//beego.Router(UserPrefix+"/edit_profile", &user.Controller{}, "*:EditProfile")
	//
	////文档
	//beego.Router(ManagePrefix+"/doc/add", &doc.Controller{}, "POST:AddDocument")         //文档上传
	//beego.Router(UserPrefix+"/doc/view", &doc.Controller{}, "GET:GetDocUsing")           //文档查看
	//beego.Router(ManagePrefix+"/doc/list", &doc.Controller{}, "GET:GetDocList")          //文档列表
	//beego.Router(ManagePrefix+"/doc/set_status", &doc.Controller{}, "Post:SetDocStatus") //文档列表
	//
	//beego.Router(UserPrefix+"/doc/file_add", &doc.Controller{}, "POST:AddFile")      //文件上传
	//beego.Router(UserPrefix+"/doc/file_down", &doc.Controller{}, "GET:FileDownload") //文件下载

	// 非登录态列表
	notNeedAuthList := []string{
		// aliyun check
		"/",
		// user
		ExemptPrefix + "/user/getcode", ExemptPrefix + "/user/register", ExemptPrefix + "/user/login",
	}

	// add filter
	// 请求合法性验证 这个要放在第一个
	//beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckRequestFilter())
	//filter.AddURLCheckSeed("wxapp", "bFvKYrlnHdtSaaGk7B1t") // 添加URLCheckSeed
	beego.InsertFilter("/v2/*", beego.BeforeRouter, filter.CheckAuthFilter("stowage_user", notNeedAuthList))
	beego.InsertFilter("/*", beego.BeforeRouter, filter.RequestFilter())
}
