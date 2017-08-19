package controller

import (
	"fmt"
	"net/http"

	"github.com/julienschmidt/httprouter"
)

func Search(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Search!\n")
}

func SearchSuggestions(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "SearchSuggestions!\n")
}
