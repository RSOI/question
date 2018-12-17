package main

import (
	"encoding/json"
	"fmt"

	"github.com/RSOI/question/controller"
	"github.com/RSOI/question/ui"
	"github.com/RSOI/question/utils"
	"github.com/buaazp/fasthttprouter"
	"github.com/valyala/fasthttp"
)

func sendResponse(ctx *fasthttp.RequestCtx, r ui.Response, nolog ...bool) {
	ctx.Response.Header.Set("Content-Type", "application/json")
	ctx.Response.SetStatusCode(r.Status)
	utils.LOG(fmt.Sprintf("Sending response. Status: %d", r.Status))

	doLog := true
	if len(nolog) > 0 {
		doLog = !nolog[0]
	}

	if doLog {
		controller.LogStat(ctx.Path(), r.Status, r.Error)
	}

	content, _ := json.Marshal(r)
	ctx.Write(content)
}

func indexGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get service stats (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.IndexGET(ctx.Host())
	r.Status, r.Error = ui.ErrToResponse(err)

	nolog := true
	sendResponse(ctx, r, nolog)
}

func askPUT(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Ask new question (%s)", ctx.Path()))
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
	utils.LOG(fmt.Sprintf("Request: Get one question (%s)", ctx.Path()))
	var err error
	var r ui.Response

	id := ctx.UserValue("id").(string)
	r.Data, err = controller.QuestionGET(id)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func questionsGET(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Get questions by author (%s)", ctx.Path()))
	var err error
	var r ui.Response

	aid := ctx.UserValue("authorid").(string)
	r.Data, err = controller.QuestionsGET(aid)
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func updatePATCH(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Update question (%s)", ctx.Path()))
	var err error
	var r ui.Response

	r.Data, err = controller.UpdatePATCH(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func removeDELETE(ctx *fasthttp.RequestCtx) {
	utils.LOG(fmt.Sprintf("Request: Delete question (%s)", ctx.Path()))
	var err error
	var r ui.Response

	err = controller.RemoveDELETE(ctx.PostBody())
	r.Status, r.Error = ui.ErrToResponse(err)
	sendResponse(ctx, r)
}

func initRoutes() *fasthttprouter.Router {
	utils.LOG("Setup router...")
	router := fasthttprouter.New()
	router.GET("/", indexGET)
	router.PUT("/ask", askPUT)
	router.GET("/question/id:id", questionGET)
	router.GET("/questions/author:authorid", questionsGET)
	router.PATCH("/update", updatePATCH)
	router.DELETE("/delete", removeDELETE)

	return router
}
