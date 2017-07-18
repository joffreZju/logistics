package cardsdk

import (
	"bytes"
	"encoding/base64"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"net/url"
	"strings"
)

const (
	APPKEY    = "23574208"
	APPSECRET = "afadf1a2693983fb6b82d0802852c98b"
)

//身份证识别
func IdCard(databytes []byte, side string) (res map[string]interface{}, err error) {
	bt := base64.StdEncoding.EncodeToString(databytes)
	imagevalue := map[string]interface{}{
		"dataType":  50,
		"dataValue": bt,
	}
	configurevalue := map[string]interface{}{
		"dataType":  50,
		"dataValue": "{\"side\":\"" + side + "\"}",
	}
	inputvalue := []interface{}{
		map[string]interface{}{
			"image":     imagevalue,
			"configure": configurevalue,
		},
	}
	reqbody := map[string]interface{}{
		"inputs": inputvalue,
	}
	jsonbody, err := json.Marshal(reqbody)
	if err != nil {
		return
	}
	url := "http://dm-51.data.aliyun.com/rest/160601/ocr/ocr_idcard.json"
	client := &http.Client{}
	req, err := http.NewRequest("POST", url, bytes.NewReader(jsonbody))
	if err != nil {
		return
	}
	//TODO key
	req, err = AddApiAuth(req, APPKEY, APPSECRET, jsonbody, "")
	if err != nil {
		return
	}
	fmt.Println(req.Header)
	resp, err := client.Do(req)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	fmt.Println("body:", string(body), resp.Header)
	err = json.Unmarshal(body, &res)
	return
}

//驾驶证识别
func DriverCard(databytes []byte) (res map[string]interface{}, err error) {
	bt := base64.StdEncoding.EncodeToString(databytes)
	imagevalue := map[string]interface{}{
		"dataType":  50,
		"dataValue": bt,
	}
	inputvalue := []interface{}{
		map[string]interface{}{
			"image": imagevalue,
		},
	}
	reqbody := map[string]interface{}{
		"inputs": inputvalue,
	}
	jsonbody, err := json.Marshal(reqbody)
	if err != nil {
		return
	}
	url := "http://dm-52.data.aliyun.com/rest/160601/ocr/ocr_driver_license.json"
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonbody))
	if err != nil {
		return
	}
	req, err = AddApiAuth(req, APPKEY, APPSECRET, jsonbody, "")
	if err != nil {
		return
	}

	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return
}

//行驶证识别
func VehicleCard(databytes []byte) (res map[string]interface{}, err error) {
	bt := base64.StdEncoding.EncodeToString(databytes)
	imagevalue := map[string]interface{}{
		"dataType":  50,
		"dataValue": bt,
	}
	inputvalue := []interface{}{
		map[string]interface{}{
			"image": imagevalue,
		},
	}
	reqbody := map[string]interface{}{
		"inputs": inputvalue,
	}
	jsonbody, err := json.Marshal(reqbody)
	if err != nil {
		return
	}
	url := "http://dm-53.data.aliyun.com/rest/160601/ocr/ocr_vehicle.json"
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonbody))
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, jsonbody, "")
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return
}

//营业证识别
func BusinessCard(databytes []byte) (res map[string]interface{}, err error) {
	bt := base64.StdEncoding.EncodeToString(databytes)
	imagevalue := map[string]interface{}{
		"dataType":  50,
		"dataValue": bt,
	}
	inputvalue := []interface{}{
		map[string]interface{}{
			"image": imagevalue,
		},
	}
	reqbody := map[string]interface{}{
		"inputs": inputvalue,
	}
	jsonbody, err := json.Marshal(reqbody)
	if err != nil {
		return
	}
	url := "http://dm-58.data.aliyun.com/rest/160601/ocr/ocr_business_license.json"
	client := &http.Client{}

	req, err := http.NewRequest("POST", url, bytes.NewBuffer(jsonbody))
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, jsonbody, "")
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return
}

//工商信息查询。areaCode 地区代码 comName 公司名关键字 industryCode 国家标准行业代码 page 页码（必选）
func BusinessInfo(areaCode string, comName string, industryCode string, page string) (res map[string]interface{}, err error) {
	querystrs := []string{}
	if areaCode != "" {
		querystrs = append(querystrs, areaCode)
	}
	if comName != "" {
		querystrs = append(querystrs, comName)
	}
	if industryCode != "" {
		querystrs = append(querystrs, industryCode)
	}
	querystrs = append(querystrs, page)
	querystr := strings.Join(querystrs, "&")
	QueryValue, err := url.ParseQuery(querystr)
	if err != nil {
		return
	}
	url := "http://qianzhan1.market.alicloudapi.com/OperVague?" + QueryValue.Encode()
	client := &http.Client{}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, []byte{}, querystr)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return

}

//工商信息模糊查询
func VagueBusinessInfo(comName string, page string) (res string, err error) {
	querystr := "comName=" + comName + "&page=" + page
	QueryValue, err := url.ParseQuery(querystr)
	if err != nil {
		return
	}
	uri := "http://qianzhan1.market.alicloudapi.com/OpenAli/CommerceVague?" + QueryValue.Encode()
	fmt.Println("uri:", uri)
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, []byte{}, querystr)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode != 200 {
		err = fmt.Errorf("errorcode: %v  errinfo: %v\n, header :%v", resp.StatusCode, string(body), resp.Header)
		return
	}
	err = json.Unmarshal(body, &res)
	return
}

//工商信息精准查询
func AccurateBussinessInfo(comName string) (res map[string]interface{}, err error) {
	querystr := "comName=" + comName
	QueryValue, err := url.ParseQuery(querystr)
	if err != nil {
		return
	}
	uri := "http://qianzhan1.market.alicloudapi.com/OpenAli/CommerceAccurate?" + QueryValue.Encode()
	fmt.Println("uri:", uri)
	client := &http.Client{}
	req, err := http.NewRequest("GET", uri, nil)
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, []byte{}, querystr)
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	if err != nil {
		return
	}

	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	if resp.StatusCode/200 != 1 {
		err = fmt.Errorf("errorcode: %v  errinfo: %v\n, header :%v", resp.StatusCode, string(body), resp.Header)
		return
	}
	err = json.Unmarshal(body, &res)
	return

}

//名片识别
func BCard(databytes []byte) (res map[string]interface{}, err error) {
	bt := base64.StdEncoding.EncodeToString(databytes)
	localip, err := GetRemoteIp()
	if err != nil {
		return
	}
	reqbody := map[string]interface{}{
		"uid":   localip,
		"lang":  "auto",
		"color": "gray",
		"image": bt,
	}
	jsonbody, err := json.Marshal(reqbody)
	uri := "http://businesscard.aliapi.hanvon.com/rt/ws/v1/ocr/bcard/recg?code=cf22e3bb-d41c-47e0-aa44-a92984f5829d"

	client := &http.Client{}

	req, err := http.NewRequest("POST", uri, bytes.NewBuffer(jsonbody))
	if err != nil {
		return
	}

	req, err = AddApiAuth(req, APPKEY, APPSECRET, jsonbody, "")
	if err != nil {
		return
	}
	resp, err := client.Do(req)
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	err = json.Unmarshal(body, &res)
	return
}
