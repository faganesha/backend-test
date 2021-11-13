package connection

import (
	"database/sql"
	"log"
)

func Connection() *sql.DB {
	db, errDb := sql.Open("mysql", "root:@tcp(127.0.0.1:3306)/backend_test")
	if errDb != nil {
		log.Fatal(errDb)
	}
	return db
}
