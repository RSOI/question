package controller

import (
	"encoding/json"
	"strconv"

	"github.com/valyala/fasthttp"
)

// QuestionGET get question by id
func QuestionGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response

	qID, _ := strconv.Atoi(ctx.UserValue("id").(string))

	r.Data, err = QuestionModel.GetQuestionByID(qID)
	if err != nil {
		r.Data = nil
	}

	r.Status, r.Error = errToResponse(err)

	QuestionModel.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

// QuestionsGET get questions by author
func QuestionsGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r Response

	qAuthorID, _ := strconv.Atoi(ctx.UserValue("authorid").(string))

	r.Data, err = QuestionModel.GetQuestionsByAuthorID(qAuthorID)
	if err != nil {
		r.Data = nil
	}

	r.Status, r.Error = errToResponse(err)

	QuestionModel.LogStat(ctx.Path(), r.Status, r.Error)

	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}
