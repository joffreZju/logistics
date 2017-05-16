package image

import (
	"bytes"
	"fmt"

	"github.com/aliyun/aliyun-oss-go-sdk/oss"
	"github.com/astaxie/beego"
)

var client *oss.Client

var (
	endpoint_prod = "oss-cn-hangzhou-internal.aliyuncs.com"
	endpoint_dev  = "oss-cn-hangzhou.aliyuncs.com"
	accessID      = "LTAIysfw9MWnCZFk"
	accessKey     = "hAuLM27EkdVVxtfvbYHgq5XPDRvial"
	bucketName    = "allsum-images"
	//imageHost     = "image.allsum.com"
	imageHost = "oss-cn-hangzhou.aliyuncs.com"
)

func Init() (err error) {
	// New Client
	endpoint := endpoint_dev
	if beego.BConfig.RunMode == "prod" {
		endpoint = endpoint_prod
	}
	client, err = oss.New(endpoint, accessID, accessKey)
	if err != nil {
		return
	}
	return
}

func CreateImage(group string, name string, data []byte) (url string, err error) {
	bucket, err := client.Bucket(bucketName)
	if err != nil {
		return
	}
	err = bucket.PutObject(group+"/"+name, bytes.NewReader(data))
	if err != nil {
		return
	}
	url = generateImageURL(group, name)
	return
}

func generateImageURL(group, name string) (url string) {
	if len(group) > 0 {
		return fmt.Sprintf("http://%s/%s/%s", imageHost, group, name)
	}
	return fmt.Sprintf("http://%s/%s", imageHost, group)
}
