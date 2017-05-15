package wxisv

import (
	"fmt"
	"math/rand"
	"strconv"
	"testing"
	"time"
)

var pay *Pay

func init() {
	caString := `-----BEGIN CERTIFICATE-----
MIIDIDCCAomgAwIBAgIENd70zzANBgkqhkiG9w0BAQUFADBOMQswCQYDVQQGEwJV
UzEQMA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2Vy
dGlmaWNhdGUgQXV0aG9yaXR5MB4XDTk4MDgyMjE2NDE1MVoXDTE4MDgyMjE2NDE1
MVowTjELMAkGA1UEBhMCVVMxEDAOBgNVBAoTB0VxdWlmYXgxLTArBgNVBAsTJEVx
dWlmYXggU2VjdXJlIENlcnRpZmljYXRlIEF1dGhvcml0eTCBnzANBgkqhkiG9w0B
AQEFAAOBjQAwgYkCgYEAwV2xWGcIYu6gmi0fCG2RFGiYCh7+2gRvE4RiIcPRfM6f
BeC4AfBONOziipUEZKzxa1NfBbPLZ4C/QgKO/t0BCezhABRP/PvwDN1Dulsr4R+A
cJkVV5MW8Q+XarfCaCMczE1ZMKxRHjuvK9buY0V7xdlfUNLjUA86iOe/FP3gx7kC
AwEAAaOCAQkwggEFMHAGA1UdHwRpMGcwZaBjoGGkXzBdMQswCQYDVQQGEwJVUzEQ
MA4GA1UEChMHRXF1aWZheDEtMCsGA1UECxMkRXF1aWZheCBTZWN1cmUgQ2VydGlm
aWNhdGUgQXV0aG9yaXR5MQ0wCwYDVQQDEwRDUkwxMBoGA1UdEAQTMBGBDzIwMTgw
ODIyMTY0MTUxWjALBgNVHQ8EBAMCAQYwHwYDVR0jBBgwFoAUSOZo+SvSspXXR9gj
IBBPM5iQn9QwHQYDVR0OBBYEFEjmaPkr0rKV10fYIyAQTzOYkJ/UMAwGA1UdEwQF
MAMBAf8wGgYJKoZIhvZ9B0EABA0wCxsFVjMuMGMDAgbAMA0GCSqGSIb3DQEBBQUA
A4GBAFjOKer89961zgK5F7WF0bnj4JXMJTENAKaSbn+2kmOeUJXRmm/kEd5jhW6Y
7qj/WsjTVbJmcVfewCHrPSqnI0kBBIZCe/zuf6IWUrVnZ9NA2zsmWLIodz2uFHdh
1voqZiegDfqnc1zqcPGUIWVEX/r87yloqaKHee9570+sB3c4
-----END CERTIFICATE-----`
	certString := `-----BEGIN CERTIFICATE-----
MIIEYjCCA8ugAwIBAgIDcU1XMA0GCSqGSIb3DQEBBQUAMIGKMQswCQYDVQQGEwJD
TjESMBAGA1UECBMJR3Vhbmdkb25nMREwDwYDVQQHEwhTaGVuemhlbjEQMA4GA1UE
ChMHVGVuY2VudDEMMAoGA1UECxMDV1hHMRMwEQYDVQQDEwpNbXBheW1jaENBMR8w
HQYJKoZIhvcNAQkBFhBtbXBheW1jaEB0ZW5jZW50MB4XDTE2MTIyMzA4MjAxNVoX
DTI2MTIyMTA4MjAxNVowgZIxCzAJBgNVBAYTAkNOMRIwEAYDVQQIEwlHdWFuZ2Rv
bmcxETAPBgNVBAcTCFNoZW56aGVuMRAwDgYDVQQKEwdUZW5jZW50MQ4wDAYDVQQL
EwVNTVBheTEnMCUGA1UEAxQe5p2t5bee5rGC5Zyj56eR5oqA5pyJ6ZmQ5YWs5Y+4
MREwDwYDVQQEEwgxNzMzMjY1NTCCASIwDQYJKoZIhvcNAQEBBQADggEPADCCAQoC
ggEBANUYTLfx10GKHTwWpSklg++RMJdULC29RajpwXve8oYV4WzjNm0nnc6W6/rd
VOCK5cHyxpvXQCiS1fm/DxWqjoDjnjfb6xo56Fldm12JjZjTMinSsmUXsgwJpQG3
G2a/FT4CnZb3C7UFqS3Lnj3709WR1qAAh3suuj+P+iL6FKSTXT8Wn6I1hkb+eyJg
mANwEBQL/UywAMqsclmfHijq7YDGd3Wv8SZeBzeAB/jT6TknNzMm80h2SoWJWWb3
tJQjX7RzN4Erh+yr6GFWn15457aA2YetBa+Ybja1rGWRJKm/g+El9lka48uUDRGe
OXxP9PN/LOEqMlDn+AL5Qn6pRJkCAwEAAaOCAUYwggFCMAkGA1UdEwQCMAAwLAYJ
YIZIAYb4QgENBB8WHSJDRVMtQ0EgR2VuZXJhdGUgQ2VydGlmaWNhdGUiMB0GA1Ud
DgQWBBSl1wQsy+pN+qg2hlOLP3Apc2/w6jCBvwYDVR0jBIG3MIG0gBQ+BSb2ImK0
FVuIzWR+sNRip+WGdKGBkKSBjTCBijELMAkGA1UEBhMCQ04xEjAQBgNVBAgTCUd1
YW5nZG9uZzERMA8GA1UEBxMIU2hlbnpoZW4xEDAOBgNVBAoTB1RlbmNlbnQxDDAK
BgNVBAsTA1dYRzETMBEGA1UEAxMKTW1wYXltY2hDQTEfMB0GCSqGSIb3DQEJARYQ
bW1wYXltY2hAdGVuY2VudIIJALtUlyu8AOhXMA4GA1UdDwEB/wQEAwIGwDAWBgNV
HSUBAf8EDDAKBggrBgEFBQcDAjANBgkqhkiG9w0BAQUFAAOBgQAzdccNEznEYxf1
r3fxbgWvDLnMX4mH9FTFE98stdG9xrwdyQB8tqJJV9lyyscEgfIo4ie87byuhkNL
d0RN8WY6o2SMPgrj1/8H/oMWnP5J1qgY3vNaK5S+oiv1b6WUxgPeavtoiHvG4h20
kXtjIR4qj0+VImSrUfMmDVZVMiIvAg==
-----END CERTIFICATE-----`
	keyString := `-----BEGIN PRIVATE KEY-----
MIIEvwIBADANBgkqhkiG9w0BAQEFAASCBKkwggSlAgEAAoIBAQDVGEy38ddBih08
FqUpJYPvkTCXVCwtvUWo6cF73vKGFeFs4zZtJ53Oluv63VTgiuXB8sab10AoktX5
vw8Vqo6A45432+saOehZXZtdiY2Y0zIp0rJlF7IMCaUBtxtmvxU+Ap2W9wu1Bakt
y549+9PVkdagAId7Lro/j/oi+hSkk10/Fp+iNYZG/nsiYJgDcBAUC/1MsADKrHJZ
nx4o6u2Axnd1r/EmXgc3gAf40+k5JzczJvNIdkqFiVlm97SUI1+0czeBK4fsq+hh
Vp9eeOe2gNmHrQWvmG42taxlkSSpv4PhJfZZGuPLlA0Rnjl8T/TzfyzhKjJQ5/gC
+UJ+qUSZAgMBAAECggEBAMQgD0wlO7bIhUuuk+gg7SNq/8vn3pliYGCsdDWr5o7e
SJHNNWSVV7qyURKc7ueTLw+ogH8iR5yQOHwaCqooRev+krpaoDGNJnpJmxsl5LrJ
dpvjnelJO8e0gLfpbUDNkaF3Cs/NJGtBgInzo/rscfVYuq6cjhUj1qt1ugTDIois
oupSE2mNIs5TsI21h0L9ng+3LDtXENDDIXgCM9p9OvIqD0LVCcZ4NbhMjL6HrCjX
BARA7C8DufntDk+c+rYYMDi0+QYK4qFfuOmIimaLrkSA6OVVmyZfaf7T4bj+3pBN
oZ66CSOWaYQnRV5uQRBm1EcXgC73GLZqBG7rP3fsXgECgYEA7AIT4xmR3wi75L0B
v9EqHTvyHgU/hNGpOCXqQlBvk+EBomZ1jVXcHcvMp77uUzsbRtFCeBlDzdRtegim
SRydiUKvRpnXF8PqZoPdxRiBhzfXAGbKvhhacUuqP3fV6txfe94Wkkp2LRPI1MoA
XClInqxrmQagK06VPTLuZKiiiOECgYEA5yVXSj1hGDvdQiSp2hZKu0hCxO3W4xKn
99osPUwhhTSxT78JlE98FFpDtRAFpLnYphZGOypgT7sE1YtpQPG3mn7zgseviD1S
V7tteoDQz5kS/u87wwC/9+rgGUFbI7xuepGzIbA86d/y56Z1YsIiEF4zHWWW/Y4d
V8rKHaV9mrkCgYB8dnQKdjepia+dZ9f+Us6E8FI1Zssivnchd01df5H4SNdVz/b6
fGdDB3F8nYKOPkOaS01kjN5nNDov+1PGhuLFunc5InR+wgFh4vUXtl7I8rfeLFeL
fMhlq2OzaP1ViLaKWotIxyAfkal+HrGl6Ne1ZnSwFQBvFlg6GBwE1bIxwQKBgQCn
nk+Hma0gasEPpxC5AvNcjpFEx4jOEAhYVxE/vkaMl2KBhvKGZ4F+LNruoVjGVLMD
9iEl5JwFFYTy6m8AVokjcy5ZRz9GV9mvn05LyMAj20iIMKowxglv2hZ6mgdiidG/
9oplQq1ZmDpIvFBhtpAHOJhul+3/nyAuOvOIviqwmQKBgQCEY2UTvs0ED6yC28Fq
JOeFsq3l0ACY758LMHf/lDa7Hxwn/r0/gAJsQup2KMFWqQRZ/OSXfHkzFrzprgyk
mDM7n5chLAS6o78y3edGOenOrJteyzep4KbSARhT48nevrrZM8S+2AOwSYERKXxx
73tJGIYp5FaIqAmOgaOrV7quig==
-----END PRIVATE KEY-----`

	client, err := NewClient("wx64e8fcc77a6999af", "1453829202", "2720c6250e4c35247ec02bd1e67f1a2d", "",
		"http://106.14.38.116:8081/notify/wxpay", caString, certString, keyString, false)
	if err != nil {
		panic(err)
	}
	pay = NewPay(client)
	rand.Seed(time.Now().UnixNano())
}

func TestMicroPay(t *testing.T) {
	authCode := "288230866698343074"
	orderNo, desc, totalAmount := strconv.Itoa(100000+rand.Intn(900000)), "测试支付", int64(1)
	reply, err := pay.MicroPay("4", "1428339002", orderNo, desc, totalAmount, authCode)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

func TestQrPay(t *testing.T) {
	orderNo, desc, totalAmount := strconv.Itoa(10000+rand.Intn(90000)), "AllSum服务预付费", int64(42)
	reply, err := pay.QrPay(orderNo, desc, totalAmount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

func TestAppPay(t *testing.T) {
	orderNo, desc, totalAmount := strconv.Itoa(10000+rand.Intn(90000)), "AllSum服务预付费", int64(242)
	reply, err := pay.AppPay("7", "", orderNo, desc, totalAmount)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

func TestWeChatOfficialAccountsPay(t *testing.T) {
	orderNo, desc, totalAmount, openId := strconv.Itoa(10000+rand.Intn(90000)), "AllSum服务预付费", int64(242), "oKrTMwRp8WFRZ9Cc5odtR67B8QI8"
	reply, err := pay.WeChatOfficialAccountsPay(orderNo, desc, totalAmount, openId)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

func TestQueryOrder(t *testing.T) {
	orderNo := "451064"
	reply, err := pay.QueryOrder("1428339002", orderNo)
	if err != nil {
		t.Fatal(err)
	}
	fmt.Printf("%+v\n", reply)
}

func TestReportTil(t *testing.T) {
	err := pay.ReportTil("https://api.mch.weixin.qq.com/pay/batchreport/micropay/total", []interface{}{})
	if err != nil {
		t.Fatal(err)
	}
}
