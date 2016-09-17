package main

import (
	"database/sql"
	"log"
	"os"

	"github.com/joho/godotenv"
	_ "github.com/lib/pq"
)

var dbConn *(sql.DB)

func loadDotenv() error {
	err := godotenv.Load()
	if err != nil {
		log.Fatal(err)
		return err
	}
	return nil
}

func initDBConn(user, dbname, passwd string) error {
	db, err := sql.Open("postgres", os.Getenv("POSTGRES_INFO"))
	if err != nil {
		log.Fatal(err)
		return err
	}
	dbConn = db
	return nil
}
