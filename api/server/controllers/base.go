package controllers

import (
	"encoding/json"

	. "netstatd"

	"github.com/astaxie/beego/context"
)

var D *Netstatd

func Init(d *Netstatd) {
	D = d
}

func SetSuccessOutput(ctx *context.Context, statusCode int, data interface{}) {
	ctx.Output.SetStatus(statusCode)
	if data != nil {
		b, _ := json.Marshal(data)
		ctx.Output.Body(b)
	} else {
		ctx.Output.Body([]byte(""))
	}
}

func SetErrorOutput(ctx *context.Context, statusCode int, err error) {
	ctx.Output.SetStatus(statusCode)
	body := make(map[string]interface{})
	body["message"] = err.Error()
	b, _ := json.Marshal(body)
	ctx.Output.Body(b)
}
