package service

import (
	"allsum_oa/model"
	"common/lib/keycrypt"
	"github.com/astaxie/beego"
	"time"
)

func InitAllsum() {
	allsumNo := beego.AppConfig.String("allsum::companyNo")
	allsumAdmin := beego.AppConfig.String("allsum::adminUser")
	db := model.NewOrm()
	count := 0
	e := db.Table(model.Company{}.TableName()).Where("no=?", "allsum").Count(&count).Error
	if e != nil {
		beego.Error("初始化allsum出错:", e)
		return
	} else if count != 0 {
		beego.Info("allsum已经初始化")
		return
	} else {
		beego.Info("开始初始化allsum")
	}

	pwdEncode := keycrypt.Sha256Cal("123456")
	adminUser := &model.User{
		Tel:      allsumAdmin,
		No:       model.UniqueNo("U"),
		UserName: "李龙峰",
		Gender:   1,
		Password: pwdEncode,
		UserType: model.UserTypeAdmin,
		Status:   model.UserStatusOk,
		Ctime:    time.Now(),
	}
	err := model.CreateUser("public", adminUser)
	if err != nil {
		beego.Error("初始化allsum admin user出错:", err)
		return
	} else {
		beego.Info("初始化allsum admin user成功")
	}

	comp := &model.Company{
		No:       allsumNo,
		FirmName: "杭州壹算科技有限公司",
		FirmType: "科技公司",
		Creator:  adminUser.Id,
		AdminId:  adminUser.Id,
		Status:   model.CompanyStatApproveWait,
	}
	err = model.CreateCompany(comp)
	if err != nil {
		beego.Error("初始化allsum company出错", err)
		return
	} else {
		beego.Info("初始化allsum company成功")
	}

	err = model.AddUserToCompany(comp.No, adminUser.Id)
	if err != nil {
		beego.Error("初始化allsum add user出错", err)
		return
	} else {
		beego.Info("初始化allsum add user成功")
	}

	e = AuditCompany(allsumNo, adminUser, model.CompanyStatApproveAccessed, "")
	if err != nil {
		beego.Error("初始化allsum，模拟审核出错", err)
		return
	} else {
		beego.Info("初始化allsum，模拟审核成功")
	}

	beego.Info("初始化allsum成功")
}
