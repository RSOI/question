package controller

import (
	"encoding/json"

	"github.com/RSOI/question/ui"
	"github.com/valyala/fasthttp"
)

// IndexGET returns usage statistic
func IndexGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response
	r.Status = 200
	r.Data, err = QuestionModel.GetUsageStatistic(string(ctx.Host()))
	if err != nil {
		r.Status, r.Error = ui.ErrToResponse(err)
		r.Data = nil
	}

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

// LogStat stores service usage
func LogStat(path []byte, status int, err string) {
	QuestionModel.LogStat(path, status, err)
}
