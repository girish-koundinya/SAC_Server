package database

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

var DB *sql.DB

func init() {
	DB, _ = getDatabase()
	log.Println(DB)
}

func getDatabase() (*sql.DB, error) {
	return sql.Open("postgres", "user=postgres dbname=sac host=pq-server sslmode=disable")
}
func checkError(err error) {
	if err != nil {
		log.Println(err)
	}
}
