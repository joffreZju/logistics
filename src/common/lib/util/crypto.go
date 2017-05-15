package util

import (
	"bytes"
	"crypto/md5"
	"crypto/rc4"
	"crypto/sha1"
	"encoding/hex"
)

func Md5Cal2String(data []byte) string {
	f := md5.New()
	f.Write(data)
	md5str := hex.EncodeToString(f.Sum(nil))
	return md5str
}

func Md5Cal2Byte(data []byte) []byte {
	f := md5.New()
	f.Write(data)
	return f.Sum(nil)
}

func ValidateMd5(data []byte, sum []byte) bool {
	mysum := Md5Cal2Byte(data)
	return bytes.Equal(mysum, sum)
}

func Sha1Cal(data []byte) string {
	fsha1 := sha1.New()
	fsha1.Write(data)
	fileh := hex.EncodeToString(fsha1.Sum(nil))
	return fileh
}

func Rc4Crypt(data []byte, key []byte) []byte {
	rc, _ := rc4.NewCipher(key)
	dst := make([]byte, len(data))
	rc.XORKeyStream(dst, data)
	return dst
}
