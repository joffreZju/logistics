package user

import (
	"allsum_oa/controller/base"
	"encoding/json"
	"errors"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
)

const host = "http://allsum.com:8093"

func (c *Controller) get_account(url string, params interface{}) (data *base.Response, e error) {
	url = host + url
	req := httplib.Get(url)
	ps, ok := params.(map[string]string)
	if !ok {
		return nil, errors.New("assert params to map[string]string failed")
	} else {
		for k, v := range ps {
			req.Param(k, v)
		}
	}
	resp, e := req.Response()
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	bodystr, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	data = &base.Response{}
	e = json.Unmarshal(bodystr, data)
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	return data, nil
}

func (c *Controller) post_account(url string, params interface{}) (data *base.Response, e error) {
	url = host + url
	req := httplib.Post(url)
	req.JSONBody(params)
	resp, e := req.Response()
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	bodystr, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	data = &base.Response{}
	e = json.Unmarshal(bodystr, data)
	if e != nil {
		beego.Error(e)
		return nil, e
	}
	return data, nil
}
