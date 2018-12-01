package main

import (
	"github.com/RSOI/question/controller"
	"github.com/buaazp/fasthttprouter"
)

func initRoutes() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", controller.IndexGET)
	router.PUT("/ask", controller.AskPUT)
	router.GET("/question/id:id", controller.QuestionGET)
	router.GET("/questions/author:authorid", controller.QuestionsGET)
	router.PATCH("/update", controller.UpdatePATCH)
	router.DELETE("/delete", controller.RemoveDELETE)

	return router
}
