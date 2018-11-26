package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/valyala/fasthttp"
)

// IndexGET returns usage statistic
func IndexGET(ctx *fasthttp.RequestCtx) {
	ServiceResponse := model.GetUsageStatistic()

	ServiceResponse.Address = string(ctx.Host())
	content, _ := json.Marshal(ServiceResponse)
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(200)
	ctx.Write(content)
}
