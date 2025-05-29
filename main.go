package main

import (
	"database/sql"
	"insider-case/router" // router klasörünü import et
	"log"
	"net/http"

	_ "github.com/mattn/go-sqlite3"
)

func main() {
	db, err := sql.Open("sqlite3", "./league.db")
	if err != nil {
		log.Fatal(err)
	}
	defer db.Close()

	router := router.NewRouter(db)
	http.ListenAndServe(":8080", router.SetupRoutes())
}
