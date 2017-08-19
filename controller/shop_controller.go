package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/girishkoundinya/SAC_Server/database"
	"github.com/julienschmidt/httprouter"
)

func ShopCreate(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	//CREATE SHOP
	//CHECK IF TAG EXISTS
	//IF DOESN'T EXIST ADD TO TAG AND AND TO SHOP TAG LINK
	r.ParseForm()
	name := r.Form.Get("name")
	description := r.Form.Get("description")
	phone := r.Form.Get("phone")
	owner := r.Form.Get("owner")
	address := r.Form.Get("address")
	categoryID := r.Form.Get("category_id")
	lat := r.Form.Get("latitude")
	long := r.Form.Get("longitude")
	tags := r.Form.Get("tags")
	shopID := createShop(name, description, phone, owner, address, categoryID, lat, long)
	insertTags(tags, shopID, categoryID)
	w.Header().Set("Content-Type", "application/json")
	w.Write(FormResponse("Success", 200, shopID))
}

func insertTags(tags string, shopid int, categoryID string) {
	tagArr := strings.Split(tags, ",")
	for _, element := range tagArr {
		log.Println(element)
		query := fmt.Sprintf("SELECT id from tags where lower(name) like '%s';", element)
		log.Println(query)
		rows, err := database.DB.Query(query)
		checkError(err)
		log.Println("*")
		log.Println(rows)
		if !rows.Next() {
			insertTag(element, shopid, categoryID)
		}
	}
}

func insertTag(element string, shopid int, categoryID string) {
	queryString := fmt.Sprintf("INSERT INTO tags(name, category_id) VALUES('%s', %s)", element, categoryID)
	var id int
	database.DB.QueryRow(queryString).Scan(&id)
	queryString = fmt.Sprintf("INSERT INTO shop_tags VALUES(%d, %d);", id, shopid)
	var relationShipID int
	database.DB.QueryRow(queryString).Scan(&relationShipID)
}

func createShop(name, description, phone, owner, address, categoryID, lat, long string) int {
	queryString := fmt.Sprintf(`INSERT INTO shops(name, description, phone, owner, address, category_id, latitude, longitude, location_geom) VALUES('%s','%s','%s','%s','%s','%s',%s,%s,ST_TRANSFORM(ST_SetSRID(ST_MakePoint(%s,%s),4326),2163)) RETURNING id;`, name, description, phone, owner, address, categoryID, lat, long, lat, long)
	var id int
	database.DB.QueryRow(queryString).Scan(&id)
	return id
}

func AddTag(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	r.ParseForm()
	tags := r.Form.Get("tags")
	categoryID := r.Form.Get("category_id")
	shopID := ps.ByName("shopid")
	shopInt, _ := strconv.Atoi(shopID)
	insertTags(tags, shopInt, categoryID)
	w.Header().Set("Content-Type", "application/json")
	w.Write(FormResponse("Success", 200, nil))
}

func ShopDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shopID := ps.ByName("shopid")
	result := fetchShopDetail(shopID)
	if len(result) > 0 {
		log.Println("*****")
		log.Println(fetchTrends(result))
		copy(result[0].Trends, fetchTrends(result))
	}
	log.Println(result)
	w.Header().Set("Content-Type", "application/json")
	w.Write(FormResponse("Success", 200, result))
}

func fetchTrends(result []Shop) []Trend {
	shop := result[0]
	query := fmt.Sprintf("select tag_id, count(*) from search_requests where request_time > '2017-08-19 00:00:00' and request_time < '2017-08-19 23:59:59' and tag_id in (select tag_id from shop_tags where shop_id = %d) group by tag_id order by count(*) desc limit 10;", shop.ID)
	var trends []Trend
	rows, err := database.DB.Query(query)
	checkError(err)
	for rows.Next() {
		var tagID, count int
		switch err := rows.Scan(&tagID, &count); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			var trend Trend
			tag := getTag(tagID)
			trend.tag = tag
			trend.count = count
			trends = append(trends, trend)
		default:
			checkError(err)
		}

	}
	return trends
}

func getTag(tagid int) *Tag {
	query := fmt.Sprintf("select id, name from tags where id = %d", tagid)
	rows, err := database.DB.Query(query)
	checkError(err)
	var tag Tag
	for rows.Next() {
		err := rows.Scan(&tag.ID, &tag.Name)
		checkError(err)
		return &tag
	}
	return nil
}

type Trend struct {
	tag   *Tag `json:"tag"`
	count int  `json:"count"`
}

type Shop struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
	Phone     string  `json:"phone"`
	Latitude  float32 `json:"latitude"`
	Longitude float32 `json:"longitude"`
	Address   string  `json:"address"`
	Trends    []Trend `json:"trend"`
}

var shops []Shop

func fetchShopDetail(id string) []Shop {
	rows, err := database.DB.Query("SELECT id,name,phone,latitude,longitude,address FROM shops WHERE id = $1", id)
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
