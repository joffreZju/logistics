package ossfile

import (
	"bytes"
	"fmt"
	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
)

var client *oss.Client

var (
	endpoint_prod = "oss-cn-shanghai-internal.aliyuncs.com"
	endpoint_dev  = "oss-cn-shanghai.aliyuncs.com"
	accessID      = "LTAIysfw9MWnCZFk"
	accessKey     = "hAuLM27EkdVVxtfvbYHgq5XPDRvial"
	bucketName    = "allsum-images"
	//imageHost   = "image.allsum.com"
	//imageHost   = "oss-cn-hangzhou.aliyuncs.com"
)

func init() {
	// New Client
	endpoint := endpoint_dev
	if beego.BConfig.RunMode == "prod" {
		endpoint = endpoint_prod
	}
	var err error
	client, err = oss.New(endpoint, accessID, accessKey)
	if err != nil {
		panic("aliyun oss init failed")
	}
	return
}

func PutFile(prefix string, filename string, data []byte) (url string, err error) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return "", err
	}
	filepath := prefix + "/" + filename
	err = bucket.PutObject(filepath, bytes.NewReader(data))
	if err != nil {
		return "", err
	}
	url = fmt.Sprintf("http://%s.%s/%s", bucket.BucketName, client.Config.Endpoint, filepath)
	return url, nil
}

//func GetFile(url string) ([]byte, error) {
//	bucket, e := client.Bucket(bucketName)
//	if e != nil {
//		return nil, e
//	}
//	body, e := bucket.GetObject(url)
//	if e != nil {
//		return nil, e
//	}
//	data, e := ioutil.ReadAll(body)
//	body.Close()
//	if e != nil {
//		return nil, e
//	}
//	return data, nil
//}

func DelFiles(urls []string) error {
	bucket, e := client.Bucket(bucketName)
	if e != nil {
		return e
	}
	if len(urls) == 0 {
		return nil
	} else if len(urls) == 1 {
		e = bucket.DeleteObject(urls[0])
		return e
	} else {
		_, e = bucket.DeleteObjects(urls)
		return e
	}
}
