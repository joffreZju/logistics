package util

import (
	"errors"
	"math/rand"
	"net"
	"strings"
	"time"
)

const (
	KC_RAND_KIND_NUM   = 0 // 纯数字
	KC_RAND_KIND_LOWER = 1 // 小写字母
	KC_RAND_KIND_UPPER = 2 // 大写字母
	KC_RAND_KIND_ALL   = 3 // 数字、大小写字母
)

// 随机字符串
func Krand(size int, kind int) []byte {
	ikind, kinds, result := kind, [][]int{[]int{10, 48}, []int{26, 97}, []int{26, 65}}, make([]byte, size)
	is_all := kind > 2 || kind < 0
	rand.Seed(time.Now().UnixNano())
	for i := 0; i < size; i++ {
		if is_all { // random ikind
			ikind = rand.Intn(3)
		}
		scope, base := kinds[ikind][0], kinds[ikind][1]
		result[i] = uint8(base + rand.Intn(scope))
	}
	return result
}

//获取本地ip
func GetLocalIP() (localip string, err error) {
	localip = "-"
	addrs, e := net.InterfaceAddrs()
	if e != nil {
		err = e
		return
	}
	for _, addr := range addrs {
		switch addr.(type) {
		case *net.IPNet:
			ip := addr.(*net.IPNet).IP
			if ip.IsGlobalUnicast() && !ip.Equal(ip.Mask(addr.(*net.IPNet).Mask)) {
				if strings.Index(ip.String(), "10.") == 0 {
					localip = ip.String()
				} else if strings.Index(ip.String(), "192") == 0 {
					localip = ip.String()
				}
			}
		}
		// 优先返回10开头的内网IP
		if strings.Index(localip, "10.") == 0 {
			break
		}
	}
	if localip == "-" {
		err = errors.New("get local ip failed")
	}
	return
}
