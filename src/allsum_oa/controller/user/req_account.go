package user

import (
	"common/lib/errcode"
	"encoding/json"
	"github.com/astaxie/beego"
	"github.com/astaxie/beego/httplib"
	"io/ioutil"
)

const host = "http://allsum.com:8093"

func (c *Controller) get(url string, params ...string) (m map[string]interface{}, e error) {
	return nil, nil
}

func (c *Controller) post_account(url string, params interface{}) (m map[string]interface{}, ecode *errcode.CodeError) {
	url = host + url
	req := httplib.Post(url)
	req.JSONBody(params)
	resp, e := req.Response()
	ecode = &errcode.CodeError{Code: 99999}
	if e != nil {
		beego.Error(e)
		ecode.Msg = e.Error()
		return nil, ecode
	}
	bodystr, e := ioutil.ReadAll(resp.Body)
	if e != nil {
		beego.Error(e)
		ecode.Msg = e.Error()
		return nil, ecode
	}
	m = make(map[string]interface{})
	e = json.Unmarshal(bodystr, &m)
	if e != nil {
		beego.Error(e)
		ecode.Msg = e.Error()
		return nil, ecode
	}

	if m["code"].(float64) != 0 {
		e = json.Unmarshal(bodystr, ecode)
		if e != nil {
			beego.Error(e)
			ecode.Msg = e.Error()
			return nil, ecode
		}
		return nil, ecode
	}
	return m, nil
}
