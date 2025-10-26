package postgres

import (
	"database/sql"
	"log"

	_ "github.com/lib/pq"
)

func DatabaseConnect(databaseUrl string) (*sql.DB, error) {
	db, err := sql.Open("postgres", databaseUrl)
	if err != nil {
		log.Fatal("error while connecting to database" + err.Error())
		return nil, err
	}
	return db, nil
}
