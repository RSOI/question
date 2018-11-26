package main

import (
	"github.com/RSOI/question/controller"
	"github.com/buaazp/fasthttprouter"
)

func initRoutes() *fasthttprouter.Router {
	router := fasthttprouter.New()
	router.GET("/", controller.IndexGET)

	return router
}
