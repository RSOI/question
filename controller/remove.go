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
			r.Error = "required fields are empty: " + f
		} else {
			switch f {
			case "id":
				err = QuestionModel.DeleteQuestionByID(QuestionToRemove)
			case "author_id":
				err = QuestionModel.DeleteQuestionByAuthorID(QuestionToRemove)
			}
			r.Status, r.Error = errToResponse(err)
		}
	}

	QuestionModel.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
