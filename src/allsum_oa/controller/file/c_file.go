package file

import (
	"allsum_oa/controller/base"
	"common/lib/errcode"
	"common/lib/ossfile"
	"github.com/astaxie/beego"
	"io/ioutil"
	"strings"
)

type Controller struct {
	base.Controller
}

func (c *Controller) UploadFile() {
	prefix := c.UserComp //todo
	f, h, e := c.GetFile("file")
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	defer f.Close()
	data, e := ioutil.ReadAll(f)
	if e != nil {
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	url, e := ossfile.PutFile(prefix, h.Filename, data)
	if e != nil {
		c.ReplyErr(errcode.ErrUploadFileFailed)
		return
	}
	c.ReplySucc(map[string]string{"url": url})
}

func (c *Controller) DownloadFile() {
	url := c.GetString("url")
	data, e := ossfile.GetFile(url)
	if e != nil {
		beego.Error(e)
		c.ReplyErr(errcode.ErrDownloadFileFailed)
		return
	}
	var filename string
	s := strings.Split(url, "/")
	if len(s) == 2 {
		filename = s[1]
	} else {
		filename = s[0]
	}
	c.ReplyFile("", filename, data)
}