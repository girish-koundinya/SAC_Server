package controller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ShopCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {

	fmt.Fprint(w, "ShopCreate!\n")

}

func ShopDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	shopID := ps.ByName("shopid")
	fmt.Fprintf(w, "ShopCreate, %s!\n", shopID)
}

func ProductCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Param) {

	fmt.Fprint(w, "Product create!\n")

}

func ProductDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	productID := ps.ByName("productid")
	fmt.Fprintf(w, "ProductDetail, %s!\n", productID)

}
