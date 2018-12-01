package controller

import (
	"encoding/json"

	"github.com/RSOI/question/model"
	"github.com/RSOI/question/view"
	"github.com/valyala/fasthttp"
)

// AskPOST new question
func AskPUT(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response

	var NewQuestion model.Question
	err = json.Unmarshal(ctx.PostBody(), &NewQuestion)
	if err != nil {
		r.Status, r.Error = errToResponse(err)
	} else {
		validate, missingFields := view.ValidateNewQuestion(NewQuestion)
		if !validate {
			r.Status = 400
			r.Error = "required fields are empty: " + missingFields
		} else {
			r.Data, err = model.AddQuestion(NewQuestion)
			if err == nil {
				r.Status = 201
			} else {
				r.Data = nil
				r.Status, r.Error = errToResponse(err)
			}
		}
	}

	model.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
