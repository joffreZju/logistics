package oaclient

import (
	"encoding/json"
	"io/ioutil"
	"net/http"

	"github.com/astaxie/beego"
)

func GetAllCompanySchema() (schemas []string, err error) {
	oahost := beego.AppConfig.String("service_client::oa_host")
	url := "http://" + oahost + "/api/schema/list/get"
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
	if schemas, ok := res["data"]; ok {
		return schemas.([]string), err
	}
	schemas = []string{}
	return
}
