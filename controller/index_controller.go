package controller

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/julienschmidt/httprouter"
	_ "github.com/lib/pq"
)

func Index(w http.ResponseWriter, r *http.Request, _ httprouter.Params) {
	fmt.Fprint(w, "Welcome to the shop around the corner! :) \n")
}

func checkError(err error) {
	if err != nil {
		log.Fatal(err)
	}
}

type ResponseStruct struct {
	Message string      `json:"message"`
	Code    int         `json:"status_code"`
	Results interface{} `json:"results"`
}

func FormResponse(message string, code int, results interface{}) []byte {
	response := ResponseStruct{Message: message, Code: code, Results: results}
	jData, err := json.Marshal(response)
	if err != nil {
		panic(err)
	}
	return jData
}
