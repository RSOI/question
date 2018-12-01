package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
	"github.com/valyala/fasthttp"
)

// UpdatePATCH remove question
func UpdatePATCH(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response

	var QuestionToUpdate model.Question
	err = json.Unmarshal(ctx.PostBody(), &QuestionToUpdate)
	if err != nil {
		r.Status, r.Error = errToResponse(err)
	} else {
		validate, f := view.ValidateUpdateQuestion(QuestionToUpdate)
		if !validate {
			r.Status = 400
			r.Error = "some of next parameters are required: " + f
		} else {
			r.Data, err = model.UpdateQuestion(QuestionToUpdate)
			r.Status, r.Error = errToResponse(err)
		}
	}

	model.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
