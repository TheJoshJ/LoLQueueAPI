package initializers

import (
	"database/sql"
	"log"
)

func CreatePostgresConnect() {
	connect := "user=pqgotest dbname=pqgotest sslmode=verify-full"
	db, err := sql.Open("postgres", connect)
	if err != nil {
		log.Fatal(err)
	} else {
		log.Println("Connected to Postgres DB")
	}
	var Session *sql.DB = db
}
