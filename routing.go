package main

import (
	"encoding/json"

	"github.com/RSOI/question/controller"
	"github.com/RSOI/question/ui"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func sendResponse(ctx *fasthttp.RequestCtx, r ui.Response) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)
	controller.LogStat(ctx.Path(), r.Status, r.Error)

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

func askPUT(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response

	r.Data, err = controller.AskPUT(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	if r.Status == 200 {
		r.Status = 201 // REST :)
	}
	sendResponse(ctx, r)
}

func questionGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response

	id := ctx.UserValue("id").(string)
	r.Data, err = controller.QuestionGET(id)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func questionsGET(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response

	aid := ctx.UserValue("authorid").(string)
	r.Data, err = controller.QuestionsGET(aid)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func updatePATCH(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response

	r.Data, err = controller.UpdatePATCH(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func removeDELETE(ctx *fasthttp.RequestCtx) {
	var err error
	var r ui.Response

	err = controller.RemoveDELETE(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func initRoutes() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", controller.IndexGET)
	router.PUT("/ask", askPUT)
	router.GET("/question/id:id", questionGET)
	router.GET("/questions/author:authorid", questionsGET)
	router.PATCH("/update", updatePATCH)
	router.DELETE("/delete", removeDELETE)

	return router
}
