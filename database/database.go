package database

import (
	"database/sql"

	_ "github.com/lib/pq"
)

func GetDatabase() (*sql.DB, error) {
	return sql.Open("postgres", "user=postgres dbname=sac host=pq-server sslmode=disable")
}
