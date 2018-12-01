package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
	"github.com/valyala/fasthttp"
)

// RemoveDELETE remove question
func RemoveDELETE(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response

	var QuestionToRemove model.Question
	err = json.Unmarshal(ctx.PostBody(), &QuestionToRemove)
	if err != nil {
		r.Status, r.Error = errToResponse(err)
	} else {
		validate, f := view.ValidateDeleteQuestion(QuestionToRemove)
		if !validate {
			r.Status = 400
			r.Error = "one of next parameters are required: " + f
		} else {
			switch f {
			case "id":
				err = QuestionToRemove.DeleteByID()
			case "author_id":
				err = QuestionToRemove.DeleteByAuthorID()
			}
			r.Status, r.Error = errToResponse(err)
		}
	}

	model.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
