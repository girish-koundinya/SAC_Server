package router

import (
	"github.com/girishkoundinya/SAC_Server/controller"
	"github.com/julienschmidt/httprouter"
)

func InitRouter() *httprouter.Router {
	router := httprouter.New()
	mapRoutes(router)
	return router
}

func mapRoutes(router *httprouter.Router) {
	router.GET("/", controller.Index)
}
