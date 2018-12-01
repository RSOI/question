package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/valyala/fasthttp"
)

// IndexGET returns usage statistic
func IndexGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response
	r.Status = 200
	r.Data, err = model.GetUsageStatistic(string(ctx.Host()))
	if err != nil {
		r.Status, r.Error = errToResponse(err)
		r.Data = nil
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
