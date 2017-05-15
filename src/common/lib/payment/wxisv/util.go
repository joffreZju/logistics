package wxisv

import (
	"crypto/md5"
	"fmt"
	"reflect"
	"sort"
	"strconv"
	"strings"
	"time"
)

// map转xml
func Map2Xml(params map[string]string) string {
	xmlString := "<xml>"

	for k, v := range params {
		xmlString += fmt.Sprintf("<%s>%s</%s>", k, v, k)
	}
	xmlString += "</xml>"
	return xmlString
}

// xml转map
func Xml2Map(in interface{}) (map[string]string, error) {
	xmlMap := make(map[string]string)

	v := reflect.ValueOf(in)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	if v.Kind() != reflect.Struct {
		return nil, fmt.Errorf("xml2Map only accepts structs; got %T", v)
	}

	typ := v.Type()
	for i := 0; i < v.NumField(); i++ {
		fi := typ.Field(i)
		tagv := fi.Tag.Get("xml")

		if strings.Contains(tagv, ",") {
			tagvs := strings.Split(tagv, ",")

			switch tagvs[1] {
			case "innerXml":
				innerXmlMap, err := Xml2Map(v.Field(i).Interface())
				if err != nil {
					return nil, err
				}
				for k, v := range innerXmlMap {
					if _, ok := xmlMap[k]; !ok {
						xmlMap[k] = v
					}
				}
			}
		} else if tagv != "" && tagv != "xml" {
			xmlMap[tagv] = v.Field(i).String()
		}
	}
	return xmlMap, nil
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

// Sign the parameter in form of map[string]string with app key.
// Empty string and "sign" key is excluded before sign.
// Please refer to http://pay.weixin.qq.com/wiki/doc/api/app.php?chapter=4_3
func Sign(param map[string]string, key string) string {
	newMap := make(map[string]string)
	// fmt.Printf("%#v\n", param)
	for k, v := range param {
		if k == "sign" {
			continue
		}
		if v == "" {
			continue
		}
		newMap[k] = v
	}
	// fmt.Printf("%#v\n\n", newMap)

	preSignStr := SortAndConcat(newMap)
	preSignWithKey := preSignStr + "&key=" + key

	fmt.Println(preSignWithKey)

	return fmt.Sprintf("%X", md5.Sum([]byte(preSignWithKey)))
}

// NewNonceString return random string in 32 characters
func NewNonceString() string {
	nonce := strconv.FormatInt(time.Now().UnixNano(), 36)
	return fmt.Sprintf("%x", md5.Sum([]byte(nonce)))
}

// NewTimestampString return
func NewTimestampString() string {
	return fmt.Sprintf("%d", time.Now().Unix()+ChinaTimeZoneOffset)
}
