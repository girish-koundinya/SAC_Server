package router

import "github.com/julienschmidt/httprouter"
import "gitlab.com/shoparoundthecorner_backend/controller"

func InitRouter() *httprouter.Router {
	router := httprouter.New()
	mapRoutes(router)
	return router
}

func mapRoutes(router *httprouter.Router) {
	router.GET("/", controller.Index)
}
