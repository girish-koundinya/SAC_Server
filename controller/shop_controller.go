package controller

import (
	"database/sql"
	"fmt"
	"log"
	"net/http"
	"strconv"
	"strings"
	"time"

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

	w.Header().Set("Content-Type", "application/json")

	jsonResponse := []struct{
		ShopDetails 		Shop `json:"shop"`
		InventoryTrends []Trend `json:"inventory_trends"`
		SuggestedTrends []Trend `json:"suggested_trends"`
	}{
		{result[0], fetchTrends(result[0], "inventory"), fetchTrends(result[0], "suggestion")} }
	w.Write(FormResponse("Success", 200, jsonResponse))
}

func fetchTrends(shop Shop, trend_type string) []Trend {
	yday_date := time.Now().Local().Add(-24*time.Hour).Format("2006-01-02")
	time_condition := `request_time > '` + yday_date + ` 00:00:00' and request_time < '` + yday_date + ` 23:59:59'`

	var tagid_condition string
	if trend_type == "inventory" {
		tagid_condition = `tag_id in (select tag_id from shop_tags where shop_id = ` + strconv.Itoa(shop.ID) + `)`
	} else {
		tagid_condition = `tag_id not in (select tag_id from shop_tags where shop_id = ` + strconv.Itoa(shop.ID) + `) and category_id = ` + strconv.Itoa(shop.CategoryID)
	}

	latitude := strconv.FormatFloat(shop.Latitude, 'f', -1, 32)
	longitude := strconv.FormatFloat(shop.Longitude, 'f', -1, 32)
	location_condition := `ST_Distance(search_requests.request_geom, ST_Transform(ST_SetSRID(ST_MakePoint(` + latitude + `, ` + longitude + `),4326),2163)) < 40000`

	where_conditions := time_condition + ` AND ` + tagid_condition + ` AND ` + location_condition
	aggr_query := `SELECT tag_id, count(*) AS "tag_count" FROM search_requests WHERE ` + where_conditions + ` group by tag_id order by count(*) desc limit 10`

	query := `select tags.name, aggr.tag_count from (` + aggr_query + `) AS aggr join tags on aggr.tag_id = tags.id order by aggr.tag_count desc`
	fmt.Println(query);

	rows, err := database.DB.Query(query)
	checkError(err)

	var trends []Trend
	var trend Trend

	for rows.Next() {
		switch err := rows.Scan(&trend.TagName, &trend.Count); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			trends = append(trends, trend)
		default:
			checkError(err)
		}

	}

	return trends
}

type Trend struct {
	TagName   string `json:"tag_name"`
	Count 		int  `json:"count"`
}

type Shop struct {
	ID        	int     `json:"id"`
	Name      	string  `json:"name"`
	Phone     	string  `json:"phone"`
	Latitude  	float64 `json:"latitude"`
	Longitude 	float64 `json:"longitude"`
	Address   	string  `json:"address"`
	CategoryID  int			`json:"category_id"`
}

var shops []Shop

func fetchShopDetail(id string) []Shop {
	rows, err := database.DB.Query("SELECT id,name,phone,latitude,longitude,address,category_id FROM shops WHERE id = $1", id)
	checkError(err)
	var shop Shop
	for rows.Next() {
		switch err := rows.Scan(&shop.ID, &shop.Name, &shop.Phone, &shop.Latitude, &shop.Longitude, &shop.Address, &shop.CategoryID); err {
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

func ShopTags(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	shopID := ps.ByName("shopid")

	w.Header().Set("Content-Type", "application/json")
	result := fetchShopTags(shopID)
	if len(result) > 0 {
		w.Write(FormResponse("Success", 200, result))
	} else {
		w.Write(FormResponse("No tags found :(", 404, result))
	}
}

func fetchShopTags(shopId string) []Tag{
	query := `SELECT tags.id, tags.name FROM tags LEFT JOIN shop_tags ON shop_tags.tag_id = tags.id WHERE shop_tags.shop_id = ` + shopId

	rows, err := database.DB.Query(query);
	checkError(err);

	var tag Tag
	var tags []Tag

	for rows.Next() {
		switch err := rows.Scan(&tag.ID, &tag.Name); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			tags = append(tags, tag)
		default:
			checkError(err)
		}
	}

	return tags
}

func ProductCreate(w http.ResponseWriter, r *http.Request, _ httprouter.Param) {
	fmt.Fprint(w, "Product create!\n")

}

func ProductDetail(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {

	productID := ps.ByName("productid")

	result := fetchTagTrend(productID);

	w.Header().Set("Content-Type", "application/json")
	if len(result) > 0 {
		w.Write(FormResponse("Success", 200, result))
	} else {
		w.Write(FormResponse("No trend found :(", 404, result))
	}
}

func fetchTagTrend(tagId string) []Trend {
	yday_date := time.Now().Local().Add(-24*time.Hour).Format("2006-01-02")
	where_conditions := `request_time > '` + yday_date + ` 00:00:00' and request_time < '` + yday_date + ` 23:59:59' and tag_id = ` + tagId
	query := `SELECT time_bucket('15 minute', request_time) AS time_group, count(*) FROM search_requests WHERE ` + where_conditions + ` GROUP BY time_group ORDER BY time_group ASC`
	rows, err := database.DB.Query(query)
	checkError(err)

	var trends []Trend
	var trend Trend

	for rows.Next() {
		switch err := rows.Scan(&trend.TagName, &trend.Count); err {
		case sql.ErrNoRows:
			fmt.Println("No rows were returned!")
		case nil:
			trends = append(trends, trend)
		default:
			checkError(err)
		}
	}

	return trends
}
