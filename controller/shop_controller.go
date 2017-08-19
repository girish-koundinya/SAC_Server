package controller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func ShopCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ShopCreate!\n")

}

func ShopDetail(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "ShopDetail!\n")

}

func ProductCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Param) {
	fmt.Fprint(w, "Product create!\n")

}

func ProductDetail(w http.ResponseWriter, r *http.Request, _ httprouter.Param) {
	fmt.Fprint(w, "ProductDetail\n")
}
