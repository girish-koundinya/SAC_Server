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
	router.GET("/search", controller.Search)
	router.GET("/search_suggest", controller.SearchSuggestions)
	router.GET("/shop/:shopid", controller.ShopDetail)
	router.GET("/shop/:shopid/product/:productid", controller.ProductDetail)

	router.PUT("/shop/:shopid", controller.ShopUpdate)
	router.POST("/shop", controller.ShopCreate)
}
