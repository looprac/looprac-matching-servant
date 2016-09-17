package main

import (
	"log"
	"net/http"
	"os"
)

func main() {
	err := loadDotenv()
	if err != nil {
		log.Fatal("Error loading .env file")
		return
	}

	err = initDBConn(os.Getenv("POSTGRES_USER"), os.Getenv("POSTGRES_DBNAME"), os.Getenv("POSTGRES_PASSWD"))
	if err != nil {
		log.Fatal("Error creating connection to database")
		return
	}

	router := NewRouter()

	initClient()

	port := ":" + os.Getenv("HTTP_PORT")
	log.Fatal(http.ListenAndServe(port, router))
}
