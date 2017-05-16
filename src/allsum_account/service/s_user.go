package service

import (
	"allsum_account/model"
	"common/lib/errcode"
	"common/lib/util"
	"encoding/binary"
	"encoding/hex"

	"github.com/astaxie/beego"
)

var seedu = "asdf1234&*();.,"

func UserCreate(u *model.User) (err error) {
	//TODO add uid
	err = model.CreateUser(u)
	if err != nil {
		beego.Error("UserCreate error: ", err)
		err = errcode.ErrUserCreateFailed
		return
	}
	return
}
func UserUpdate(u *model.User, fileds ...string) (err error) {
	err = model.UpdateUser(u, fileds...)
	if err != nil {
		beego.Error("UserUpdate error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}
func GetUserByTel(tel string) (u *model.User, err error) {
	u, err = model.GetUserByTel(tel)
	if err != nil {
		beego.Error("GetUserByTel error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}

func GetUserInfo(id int) (u *model.User, err error) {
	u, err = model.GetUser(id)
	if err != nil {
		beego.Error("GetUserByTel error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}

func GenUserNo(u *model.User) (no string, err error) {
	bb := util.Md5Cal2Byte([]byte(u.Tel + seedu))
	binary.LittleEndian.PutUint32(bb[12:], uint32(u.Id))
	no = hex.EncodeToString(bb)
	return
}
