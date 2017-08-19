package controller

import (
	"fmt"
	"net/http"
	"strings"
	"strconv"

	"database/sql"
	"github.com/girishkoundinya/SAC_Server/database"
	"github.com/julienschmidt/httprouter"
)

// http://localhost:3006/search?tagId=3&latitude=12.969&longitude=80.24865
func Search(w http.ResponseWriter, r *http.Request, ps httprouter.Params) {
	queryValues := r.URL.Query()
	tagId := queryValues.Get("tagId")
	latitude := queryValues.Get("latitude")
	longitude := queryValues.Get("longitude")

	w.Header().Set("Content-Type", "application/json")
	result := searchShops(tagId, latitude, longitude)
	if len(result) > 0 {
		w.Write(FormResponse("Success", 200, result))
	} else {
		w.Write(FormResponse("No shops found :(", 404, result))
	}
}

// http://localhost:3006/search_suggest?search_text=Tea
func SearchSuggestions(w http.ResponseWriter, r *http.Request, key httprouter.Params) {
	queryValues := r.URL.Query()
	searchText := queryValues.Get("search_text")

	w.Header().Set("Content-Type", "application/json")
	result := searchTag(searchText)
	if len(result) > 0 {
		w.Write(FormResponse("Success", 200, result))
	} else {
		w.Write(FormResponse("No tags found :(", 404, result))
	}
}

// http://localhost:3006/search_chrome_extension?search_text=Tea&latitude=12.969&longitude=80.24865
func SearchChromeExtension(w http.ResponseWriter, r *http.Request, key httprouter.Params) {
	queryValues := r.URL.Query()
	searchText := queryValues.Get("search_text")

	w.Header().Set("Content-Type", "application/json")
	result := searchTag(searchText)
	if len(result) > 0 {
		tags := ""
		for _, tag := range result {
			tags = tags + strconv.Itoa(tag.ID) + ","
		}
		tags = tags[:len(tags)-1]

		latitude := queryValues.Get("latitude")
		longitude := queryValues.Get("longitude")
		w.Header().Set("Content-Type", "application/json")
		result := searchShops(tags, latitude, longitude)
		if len(result) > 0 {
			w.Write(FormResponse("Success", 200, result))
		} else {
			w.Write(FormResponse("Unable to find '" + searchText + "' near you", 404, result))
		}
	} else {
		w.Write(FormResponse("Unable to find '" + searchText + "' near you", 404, result))
	}
}

type Tag struct {
	ID        int     `json:"id"`
	Name      string  `json:"name"`
}

func searchShops(tags string, latitude string, longitude string) []Shop {
	columns := "DISTINCT(shops.id), shops.name, shops.phone, shops.latitude, shops.longitude, shops.address"
	where_conditions := `where ST_Distance(shops.location_geom, ST_Transform(ST_SetSRID(ST_MakePoint(` + latitude + `, ` + longitude + `),4326),2163)) < 40000 and tag_id in (` + tags + `)`
	query := `SELECT ` + columns + ` FROM shops LEFT JOIN shop_tags on shops.id = shop_tags.shop_id ` + where_conditions + ` limit 25`

	rows, err := database.DB.Query(query);
	checkError(err);

	var shop Shop
	var shops []Shop

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

func searchTag(searchText string) []Tag {
	query := `SELECT id, name FROM tags WHERE lower(name) LIKE '%` + strings.ToLower(searchText) + `%'`

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
