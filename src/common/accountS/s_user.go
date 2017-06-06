package accountS

import (
	"common/accountM"
	"common/lib/errcode"
	"common/lib/util"
	"encoding/binary"
	"encoding/hex"
	"github.com/astaxie/beego"
	"strings"
)

var seedu = "asdf1234&*();.,"

func CreateUser(u *accountM.User) (err error) {
	//TODO add uid
	err = accountM.CreateUser(u)
	if err != nil {
		beego.Error("UserCreate error: ", err)
		if strings.Contains(err.Error(), "duplicate key") {
			err = errcode.ErrUserAlreadyExisted
		} else {
			err = errcode.ErrUserCreateFailed
		}
		return
	}
	return
}
func UpdateUser(u *accountM.User, fileds ...string) (err error) {
	err = accountM.UpdateUser(u, fileds...)
	if err != nil {
		beego.Error("UserUpdate error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}
func GetUserByTel(tel string) (u *accountM.User, err error) {
	u, err = accountM.GetUserByTel(tel)
	if err != nil {
		beego.Error("GetUserByTel error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}

func GetUserById(id int) (u *accountM.User, err error) {
	u, err = accountM.GetUserById(id)
	if err != nil {
		beego.Error("GetUserByTel error: ", err)
		err = errcode.ErrGetUserInfoFailed
		return
	}
	return
}

func GenUserNo(u *accountM.User) (no string, err error) {
	bb := util.Md5Cal2Byte([]byte(u.Tel + seedu))
	binary.LittleEndian.PutUint32(bb[12:], uint32(u.Id))
	no = hex.EncodeToString(bb)
	return
}
