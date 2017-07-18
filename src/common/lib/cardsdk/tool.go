package cardsdk

import (
	"crypto/hmac"
	"crypto/md5"
	"crypto/sha256"
	"encoding/base64"
	"fmt"
	"io"
	"net"
	"net/http"
	"strconv"
	"strings"
	"time"

	uuid "github.com/satori/go.uuid"
)

func Sha256(data []byte, secret []byte) (res []byte, err error) {
	h := hmac.New(sha256.New, secret)
	h.Write(data)
	res = h.Sum(nil)
	return
}
func GetRemoteIp() (ip string, err error) {
	conn, err := net.Dial("udp", "suanpeizai.com:80")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	defer conn.Close()
	ip = strings.Split(conn.LocalAddr().String(), ":")[0]
	return
}

func AddApiAuth(Req *http.Request, AppKey string, AppSecret string, reqbody []byte, querystr string) (reqres *http.Request, err error) {
	HTTPMethod := Req.Method
	Accept := "application/json"
	ContentType := "application/octet-stream; charset=utf-8"
	X_Ca_Request_Mode := "debug"
	X_Ca_Version := "1"
	X_Ca_Stage := "RELEASE" // TEST、PRE、RELEASE
	X_Ca_Key := AppKey
	X_Ca_Signature_Headers := strings.Join([]string{"X-Ca-Key", "X-Ca-Request-Mode", "X-Ca-Stage", "X-Ca-Timestamp", "X-Ca-Version"}, ",")
	X_Ca_Timestamp := strconv.Itoa(int(time.Now().UnixNano()))
	X_Ca_Timestamp = X_Ca_Timestamp[0:13]
	X_Ca_Nonce := uuid.NewV4().String()

	hash := md5.New()
	io.WriteString(hash, string(reqbody))
	md5body := hash.Sum(nil)
	Content_MD5 := base64.StdEncoding.EncodeToString(md5body[:])
	Headers := "X-Ca-Key:" + X_Ca_Key + "\n" +
		"X-Ca-Request-Mode:" + X_Ca_Request_Mode + "\n" +
		"X-Ca-Stage:" + X_Ca_Stage + "\n" +
		"X-Ca-Timestamp:" + X_Ca_Timestamp + "\n" +
		"X-Ca-Version:" + X_Ca_Version + "\n"

	url := Req.URL.Path + "?" + querystr
	stringToSign := HTTPMethod + "\n" + Accept + "\n" + Content_MD5 + "\n" + ContentType + "\n" + "\n" + Headers + url
	fmt.Println("stringToSign:", stringToSign)
	sign, err := Sha256([]byte(stringToSign), []byte(AppSecret))
	if err != nil {
		return
	}
	X_Ca_Signature := base64.StdEncoding.EncodeToString(sign)

	Req.Header.Set("X-Ca-Key", X_Ca_Key)
	Req.Header.Set("X-Ca-Version", X_Ca_Version)
	Req.Header.Set("X-Ca-Stage", X_Ca_Stage)
	Req.Header.Set("X-Ca-Timestamp", X_Ca_Timestamp)
	Req.Header.Set("X-Ca-Request-Mode", X_Ca_Request_Mode)
	Req.Header.Set("Accept", Accept)
	Req.Header.Set("Content-Type", ContentType)
	Req.Header.Set("X-Ca-Signature-Headers", X_Ca_Signature_Headers)
	Req.Header.Set("X-Ca-Timestamp", X_Ca_Timestamp)
	Req.Header.Set("X-Ca-Nonce", X_Ca_Nonce)
	Req.Header.Set("Content-MD5", Content_MD5)
	Req.Header.Set("X-Ca-Signature", X_Ca_Signature)

	reqres = Req
	return
}
