package util

import (
	"encoding/hex"
	"fmt"
	"math/rand"
	"strconv"
	"strings"
	"time"
)

var (
	codeRange = []byte{'A', 'B', 'C', 'D', 'E', 'F', 'G', 'H', 'I', 'J', 'K', 'L', 'M', 'N', 'O', 'P', 'Q', 'I',
		'S', 'T', 'U', 'V', 'W', 'X', 'Y', 'Z', '1', '2', '3', '4', '5', '6', '7', '8', '9', '0'}
)

func GetTradeNo(tp int, id int) string {
	str := strings.Replace(time.Now().Format("0102150405.000"), ".", "", 1)
	str += strconv.Itoa(tp)
	str += fmt.Sprintf("%04d", id)
	return str
}

func RandomByte6(s int) string {
	r := rand.New(rand.NewSource(int64(s + 199)))
	var code [6]byte
	for i := 0; i < 6; i++ {
		j := r.Intn(36)
		code[i] = codeRange[j]
	}
	return string(code[:])
}

func RandomByte16() string {
	var code = make([]byte, 16)
	var r = rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < 16; i++ {
		code[i] = byte(r.Intn(255))
	}
	return hex.EncodeToString(code)
}

func UniqueRandom() string {
	tm := time.Now().UnixNano()
	tms := strconv.FormatInt(tm, 10)
	return tms
}
