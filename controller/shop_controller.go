package controller

import (
	"database/sql"
	"fmt"
	"net/http"

	"github.com/girishkoundinya/SAC_Server/database"

	"github.com/julienschmidt/httprouter"
)

func ShopCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	fmt.Fprint(w, "ShopCreate!\n")

}

func ShopDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	shopID := ps.ByName("shopid")
	result := fetchShopDetail(shopID)
	w.Header().Set("Content-Type", "application/json")
	w.Write(FormResponse("Success", 200, result))

}

type Shop struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Address   string  `json:"address"`
}

var shops []Shop

func fetchShopDetail(id string) []Shop {
	DB, err := database.GetDatabase()
	checkError(err)
	rows, err := DB.Query("SELECT id,name,phone,latitude,longitude,address FROM shops WHERE id = $1", id)
	checkError(err)
	var shop Shop
	for rows.Next() {
		switch err := rows.Scan(&shop.ID, &shop.Name, &shop.Phone, &shop.Latitude, &shop.Longitude, &shop.Address); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			shops = append(shops, shop)
		default:
			checkError(err)
		}
	}
	return shops
}

func ProductCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Param) {

	fmt.Fprint(w, "Product create!\n")
}

func ProductDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	productID := ps.ByName("productid")
	fmt.Fprintf(w, "ProductDetail, %s!\n", productID)

}
