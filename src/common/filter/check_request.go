package filter

import (
	"common/lib/errcode"
	"crypto/md5"
	"encoding/hex"
	"encoding/json"
	"fmt"
	"io/ioutil"
	"math"
	"sort"
	"strconv"
	"strings"
	"time"

	"github.com/astaxie/beego"
	"github.com/astaxie/beego/context"
)

var seedMap = map[string]string{}

func generateMD5Hash(text string) string {
	hasher := md5.New()
	hasher.Write([]byte(text))
	return hex.EncodeToString(hasher.Sum(nil))
}

func AddURLCheckSeed(key, seed string) {
	seedMap[key] = seed
}

func getURLCheckSeed(key string) string {
	seed, ok := seedMap[key]
	if ok {
		return seed
	}
	return seedMap["default"]
}

func CheckRequestFilter() beego.FilterFunc {
	isCheck := beego.BConfig.RunMode == "prod"
	AddURLCheckSeed("default", beego.AppConfig.String("seed"))
	return func(ctx *context.Context) {
		params := make(map[string]string)
		body, _ := ioutil.ReadAll(ctx.Request.Body)
		ctx.Request.Body.Close()
		json.Unmarshal(body, &params)
		// beego.Debug("request path : " + ctx.Request.URL.Path)
		// 上传图片不做md5校验
		if strings.Index(ctx.Request.URL.Path, "/v2/pic/") == 0 ||
			strings.Index(ctx.Request.URL.Path, "upload_pic") >= 0 ||
			strings.Index(ctx.Request.URL.Path, "upload_report_pic") >= 0 {
			return
		}

		for k := range ctx.Request.PostForm {
			params[k] = ctx.Request.PostForm[k][0]
		}
		for k := range ctx.Request.Form {
			params[k] = ctx.Request.Form[k][0]
		}

		// 查找传过来的md5
		md := ""
		for k, v := range params {
			if strings.ToLower(k) == "md5" {
				md = v
				delete(params, k)
			}
		}
		if md == "" {
			md = ctx.Request.Header.Get("md5")
		}

		str := SortAndConcat(params)
		str = strings.TrimSpace(str + getURLCheckSeed(params["source"]) + ctx.Request.URL.Path)
		sum := generateMD5Hash(str)
		if sum != strings.ToLower(md) {
			beego.Info("request url:", str)
			beego.Error("request md5 not match:", sum, md)
			if isCheck {
				ctx.Output.JSON(errcode.ErrCheckRequestFailed, false, false)
				return
			} else {
				ctx.Output.JSON(errcode.New(100001, fmt.Sprintf("URL不合法  %s  %s", str, sum)), false, false)
				return
			}
		}

		var tm int64
		for k, v := range params {
			ctx.Request.Form.Add(k, v)
			if k == "t" {
				tm, _ = strconv.ParseInt(v, 10, 64)
			}
		}

		now := time.Now().Unix()
		if math.Abs(float64(now-tm)) > 60*15 {
			beego.Info("request time:", now, tm)
			if isCheck {
				ctx.Output.JSON(errcode.ErrRequestExpired, false, false)
				return
			}
		}

		// Print cookies
		//for _, ck := range ctx.Request.Cookies() {
		//	beego.Debug(fmt.Sprintf("Cookie %s : %s", ck.Name, ck.Value))
		//}
	}
}

// SortAndConcat sort the map by key in ASCII order,
// and concat it in form of "k1=v1&k2=2"
func SortAndConcat(param map[string]string) string {
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

func RequestFilter() beego.FilterFunc {
	return func(ctx *context.Context) {
		// beego.Debug("request path : " + ctx.Request.URL.Path)
		// 上传文件不做md5校验
		if strings.Index(ctx.Request.URL.Path, "file_add") >= 0 {
			return
		}
		beego.Info("reqeust url:", ctx.Input.URI()) //for test
		params := make(map[string]interface{})
		body, _ := ioutil.ReadAll(ctx.Request.Body)
		beego.Info("request body string:", body) //for test
		ctx.Request.Body.Close()
		json.Unmarshal(body, &params)

		if ctx.Request.Method == "POST" {
			beego.Info("post request postform data:", ctx.Request.PostForm) //for test
			for k := range ctx.Request.PostForm {
				params[k] = ctx.Request.PostForm[k][0]
			}
		} else {
			beego.Info("get request form data:", ctx.Request.Form) //for test
			for k := range ctx.Request.Form {
				params[k] = ctx.Request.Form[k][0]
			}
		}

		for k, v := range params {
			if str_v, ok := v.(string); ok {
				ctx.Request.Form.Add(k, str_v)
			} else {
				b, _ := json.Marshal(v)
				ctx.Request.Form.Add(k, string(b))
			}
		}
		//ctx.Request.Form.Add("body",string(body))
		beego.Info("request params:", params)
		bbb, _ := json.Marshal(params)
		beego.Info("request all params json string:", string(bbb)) //for test
	}
}
