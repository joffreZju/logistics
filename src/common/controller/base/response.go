package base

import (
	"common/lib/errcode"
	"encoding/json"
	"fmt"
	"net/http"
)

type Response struct {
	Code int         `json:"code"`
	Msg  string      `json:"msg"`
	Data interface{} `json:"data"`
}

func (c *Controller) ReplySucc(data interface{}) {
	resp := Response{0, "", data}
	callback := c.Ctx.Input.Query("callback")
	if callback == "" {
		c.Ctx.Output.JSON(resp, false, false)
	} else {
		c.Ctx.Output.JSONP(resp, false)
	}
}

func (c *Controller) ReplyErr(err error) {
	c.IsFailed = true
	code, msg := errcode.ParseError(err)
	resp := &Response{
		Code: code,
		Msg:  msg,
	}
	//beego.BeeLogger.SetLogFuncCallDepth(4)
	//beego.Error(resp)
	//beego.BeeLogger.SetLogFuncCallDepth(3)
	callback := c.Ctx.Input.Query("callback")
	if callback == "" {
		c.Ctx.Output.JSON(resp, false, false)
	} else {
		c.Ctx.Output.JSONP(resp, false)
	}
}

func (c *Controller) ReplyCache(data []byte) {
	c.Ctx.Output.Body(data)
}

func (c *Controller) ReplyText(data interface{}) {
	c.Ctx.Output.Header("Content-Type", "text/plain; charset=utf-8")
	content, err := json.Marshal(Response{0, "", data})
	if err != nil {
		http.Error(c.Ctx.ResponseWriter, err.Error(), http.StatusInternalServerError)
		return
	}
	c.Ctx.Output.Body(content)
}

func (c *Controller) ReplyFile(mime, name string, data []byte) {
	if mime == "" {
		mime = "application/octet-stream"
	}
	c.Ctx.Output.Header("Content-Type", fmt.Sprintf("%s; charset=utf-8", mime))
	c.Ctx.Output.Header("Content-Disposition", fmt.Sprintf("attachment;filename=%s", name))
	c.Ctx.Output.Body(data)
}
