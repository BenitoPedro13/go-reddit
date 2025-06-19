package main

import (
	"log"
	"net/http"

	"github.com/benitopedro13/go-reddit/postgres"
	"github.com/benitopedro13/go-reddit/web"
)

func main() {
	store, err := postgres.NewStore("postgres://postgres:password@localhost:5432/postgres?sslmode=disable")
	if err != nil {
		log.Fatal(err)
	}

	h := web.NewHandler(store)
	log.Println("Starting server on port http://localhost:3000")
	log.Fatal(http.ListenAndServe(":3000", h))
}
