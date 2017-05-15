package wxisv

import (
	"bytes"
	"common/lib/util"
	"crypto"
	"crypto/md5"
	"crypto/rand"
	"crypto/rsa"
	"crypto/sha1"
	"crypto/tls"
	"crypto/x509"
	"encoding/base64"
	"encoding/xml"
	"errors"
	"fmt"
	"io/ioutil"
	"net/http"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/clbanning/mxj"
)

const (
	GATEWAY      = "https://api.mch.weixin.qq.com/"
	DEBUGGATEWAY = "https://api.mch.weixin.qq.com/sandboxnew/"
	DEBUGKEY     = "ABCDEFGHIJKLMNOPQRSTUVWXYZ123476"
)

var DefaultClient *Client

type Client struct {
	//=======【基本信息设置】=====================================
	//
	/**
	 * TODO: 修改这里配置为您自己申请的商户信息
	 * 微信公众号信息配置
	 *
	 * APPID：绑定支付的APPID（必须配置，开户邮件中可查看）
	 *
	 * MCHID：商户号（必须配置，开户邮件中可查看）
	 *
	 * KEY：商户支付密钥，参考开户邮件设置（必须配置，登录商户平台自行设置）
	 * 设置地址：https://pay.weixin.qq.com/index.php/account/api_cert
	 *
	 * APPSECRET：公众帐号secert（仅JSAPI支付的时候需要配置， 登录公众平台，进入开发者中心可设置），
	 * 获取地址：https://mp.weixin.qq.com/advanced/advanced?action=dev&t=advanced/dev&token=2005451881&lang=zh_CN
	 * @var string
	 */
	GateWay    string
	AppId      string
	MchId      string
	Key        string
	AppSecret  string
	NotifyUrl  string
	TlsConfig  *tls.Config
	PublicKey  *rsa.PublicKey
	PrivateKey *rsa.PrivateKey
	LocalIP    string
}

func NewClient(appId, mchId, key, appSecret, notifyUrl, wechatCA, wechatCert, wechatKey string, isDebug bool) (cli *Client, err error) {
	// load cert
	cert, err := tls.X509KeyPair([]byte(wechatCert), []byte(wechatKey))
	if err != nil {
		return
	}
	pool := x509.NewCertPool()
	ok := pool.AppendCertsFromPEM([]byte(wechatCA))
	if !ok {
		err = fmt.Errorf("create certs error")
		return
	}
	tlsConfig := &tls.Config{
		Certificates: []tls.Certificate{cert},
		RootCAs:      pool,
	}

	ip, _ := util.GetLocalIP()

	gateway := GATEWAY
	if isDebug {
		gateway = DEBUGGATEWAY
		key = DEBUGKEY
	}

	cli = &Client{gateway, appId, mchId, key, appSecret, notifyUrl, tlsConfig, nil, nil, ip}
	if DefaultClient == nil {
		DefaultClient = cli
	}
	return
}

var signTypeMap = map[string]string{}

func setSignType(cmd, method string) {
	method = strings.ToUpper(method)
	if method != "RSA" {
		method = "MD5"
	}
	signTypeMap[cmd] = method
}

func getSignType(cmd string) string {
	if m, ok := signTypeMap[cmd]; ok {
		return m
	}
	return "md5"
}

func needSecure(cmd string) bool {
	return strings.Index(cmd, "secapi") == 0
}

type ICommonReply interface {
	IsError() bool
	GetError() error
}

type CommonReply struct {
	ReturnCode  string `xml:"return_code"`
	ReturnMsg   string `xml:"return_msg"`
	ResultCode  string `xml:"result_code"`
	ErrCode     string `xml:"err_code"`
	ErrCodeDesc string `xml:"err_code_des"`
}

func (c *CommonReply) IsError() bool {
	if c == nil {
		return true
	}
	return c.ResultCode != "SUCCESS"
}

func (c *CommonReply) GetError() error {
	if c == nil {
		return errors.New("empty reply")
	}
	if c.ReturnCode != "SUCCESS" {
		return fmt.Errorf("[%s-%s]", c.ReturnCode, c.ReturnMsg)

	}
	if c.ResultCode != "SUCCESS" {
		return fmt.Errorf("%s-[%s-%s]", c.ResultCode, c.ErrCode, c.ErrCodeDesc)
	}
	return nil
}

// Submit the order to weixin pay and return the prepay id if success,
// Prepay id is used for app to start a payment
// If fail, error is not nil, check error for more information
func (cli *Client) sendCommand(cmd string, arg map[string]string, reply ICommonReply) (err error) {
	signType := getSignType(cmd)
	odrInXml := cli.signedOrderRequestXmlString(signType, arg)
	resp, err := cli.doHttpPost(cli.GateWay+cmd, []byte(odrInXml), needSecure(cmd))
	if err != nil {
		return
	}

	err = xml.Unmarshal(resp, reply)
	if err != nil {
		return
	}

	if reply.IsError() {
		fmt.Println(string(resp))
		err = reply.GetError()
		return
	}

	if cmd == "payitil/report" {
		return
	}

	// 验签
	mv, err := mxj.NewMapXml(resp)
	if err != nil {
		return
	}
	mv, ok := mv["xml"].(map[string]interface{})
	if !ok {
		err = fmt.Errorf("parse xml error")
		return
	}
	mp := make(map[string]string)
	for k, v := range mv {
		mp[k] = fmt.Sprintf("%v", v)
	}
	if signType == "RSA" {
		err = checkRsaSign(mp, cli.PublicKey, mp["sign"])
	} else {
		if signMd5(mp, cli.Key) != mp["sign"] {
			err = fmt.Errorf("签名验证失败 %s, %+v", mp["sign"], mp)
		}
	}

	fmt.Println(mp)

	return
}

// doRequest post the order in xml format with a sign
func (cli *Client) doHttpPost(targetUrl string, body []byte, needSecure bool) ([]byte, error) {
	req, err := http.NewRequest("POST", targetUrl, bytes.NewBuffer([]byte(body)))
	if err != nil {
		return []byte(""), err
	}
	req.Header.Add("Content-type", "application/x-www-form-urlencoded;charset=UTF-8")

	tr := &http.Transport{
		TLSClientConfig: &tls.Config{InsecureSkipVerify: false},
	}
	if needSecure {
		tr.TLSClientConfig = cli.TlsConfig
	}

	client := &http.Client{Transport: tr, Timeout: time.Second * 10}

	resp, err := client.Do(req)
	if err != nil {
		return []byte(""), err
	}

	defer resp.Body.Close()
	respData, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return []byte(""), err
	}

	return respData, nil
}
func (cli *Client) signedOrderRequestXmlString(signType string, param map[string]string) string {
	param = cli.newOrderRequest(signType, param)
	var sign string
	if signType == "RSA" {
		sign = signRsa(param, cli.PrivateKey)
	} else {
		sign = signMd5(param, cli.Key)
	}

	param["sign"] = sign
	return toXmlString(param)
}

func (cli *Client) newOrderRequest(signType string, param map[string]string) map[string]string {
	param["appid"] = cli.AppId
	param["mch_id"] = cli.MchId
	param["sign_type"] = strings.ToUpper(signType)
	param["nonce_str"] = newNonceString()
	return param
}

// SortAndConcat sort the map by key in ASCII order,
// and concat it in form of "k1=v1&k2=2"
func sortAndConcat(param map[string]string) string {
	var keys []string
	for k := range param {
		keys = append(keys, k)
	}

	var sortedParam []string
	sort.Strings(keys)
	for _, k := range keys {
		// fmt.Println(k, "=", param[k])
		sortedParam = append(sortedParam, k+"="+param[k])
	}

	return strings.Join(sortedParam, "&")
}

// Sign the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func signMd5(param map[string]string, key string) string {
	newMap := make(map[string]string)
	for k, v := range param {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		newMap[k] = v
	}

	preSignStr := sortAndConcat(newMap)
	preSignWithKey := preSignStr + "&key=" + key

	return fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey)))
}

// SignRSA the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func signRsa(param map[string]string, key *rsa.PrivateKey) string {
	newMap := make(map[string]string)
	for k, v := range param {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		newMap[k] = v
	}

	preSignStr := sortAndConcat(newMap)
	hashed := sha1.Sum([]byte(preSignStr))
	signed, err := rsa.SignPKCS1v15(rand.Reader, key, crypto.SHA1, hashed[:])
	if err != nil {
		fmt.Println("rsa.SignPKCS1v15", err)
		return ""
	}
	sign := base64.StdEncoding.EncodeToString(signed)
	return sign
}

// Sign the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func checkRsaSign(param map[string]string, key *rsa.PublicKey, signed string) error {
	newMap := make(map[string]string)
	for k, v := range param {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		newMap[k] = v
	}

	preSignStr := sortAndConcat(newMap)

	sig := make([]byte, base64.StdEncoding.DecodedLen(len(signed)))
	n, err := base64.StdEncoding.Decode(sig, []byte(signed))
	if err != nil {
		return err
	}
	sig = sig[:n]

	hashed := sha1.Sum([]byte(preSignStr))

	err = rsa.VerifyPKCS1v15(key, crypto.SHA1, hashed[:], sig)
	return err
}

// ToXmlString convert the map[string]string to xml string
func toXmlString(param map[string]string) string {
	xml := "<xml>"
	for k, v := range param {
		xml = xml + fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xml = xml + "</xml>"

	return xml
}

// NewNonceString return random string in 32 characters
func newNonceString() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

const ChinaTimeZoneOffset = 8 * 60 * 60 //Beijing(UTC+8:00)

// NewTimestampString return
func newTimestampString() string {
	return fmt.Sprintf("%d", time.Now().Unix()+ChinaTimeZoneOffset)
}
