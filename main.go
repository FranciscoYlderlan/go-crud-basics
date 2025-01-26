package main

import (
	"log"
	"net/http"

	"go-crud-basics/store"
)

func main() {
	contactStore := store.NewContactStore()
	mux := http.NewServeMux()

	RegisterRoutes(mux, contactStore)

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
