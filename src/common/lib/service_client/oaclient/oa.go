package oaclient

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
)

func GetAllCompanySchema() (schemas []string, err error) {
	oahost := beego.AppConfig.String("service_client::oa_host")
	url := "http://" + oahost + "/api/schema/list/get"
	beego.Debug("url:", url)
	resp, err := http.Get(url)
	if err != nil {
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		return
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	code, ok := res["code"]
	if !ok || code.(float64) != 0 {
		err = fmt.Errorf("err code ", res)
		return
	}
	schema_interfaces, ok := res["data"]
	if !ok {
		err = fmt.Errorf("miss data ")
		return
	}
	schemastrs := []string{}
	for _, schema := range schema_interfaces.([]interface{}) {
		schemastrs = append(schemastrs, schema.(string))
	}
	return schemastrs, err
}

func GetAllRoleByCompany(companyNo string) (roleinfos []map[string]interface{}, err error) {
	oahost := beego.AppConfig.String("service_client::oa_host")

	url := "http://" + oahost + "/api/role/list/get?companyNo=" + companyNo
	beego.Debug("url:", url)
	resp, err := http.Get(url)
	if err != nil {
		beego.Error("http client err", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("read body err:", err)
		return
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		beego.Error("json unmarshal err:", err, string(body))
		return
	}
	code, ok := res["code"]
	if !ok || code.(float64) != 0 {
		err = fmt.Errorf("err code ", res)
		return
	}
	roles, ok := res["data"]
	if !ok {
		err = fmt.Errorf("miss data in response")
		return
	}
	for _, role := range roles.([]interface{}) {
		roleinfos = append(roleinfos, role.(map[string]interface{}))
	}
	return
}

func GetAllUserByRole(companyNo string, roleid int) (userinfos []map[string]interface{}, err error) {
	oahost := beego.AppConfig.String("service_client::oa_host")

	url := fmt.Sprintf("http://%s/api/role/list/get?companyNo=%s&roleId=%v", oahost, companyNo, roleid)
	beego.Debug("url:", url)
	resp, err := http.Get(url)
	if err != nil {
		beego.Error("http client err", err)
		return
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)
	if err != nil {
		beego.Error("read body err:", err)
		return
	}
	var res map[string]interface{}
	err = json.Unmarshal(body, &res)
	if err != nil {
		return
	}
	code, ok := res["code"]
	if !ok || code.(float64) != 0 {
		err = fmt.Errorf("err code ", res)
		return
	}
	users, ok := res["data"]
	if !ok {
		err = fmt.Errorf("miss data in response")
		return
	}
	for _, user := range users.([]interface{}) {
		userinfos = append(userinfos, user.(map[string]interface{}))
	}
	return
}
