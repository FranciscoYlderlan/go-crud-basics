package main

import (
	"net/http"

	"go-crud-basics/handlers"
	"go-crud-basics/store"
)

func RegisterRoutes(mux *http.ServeMux, contactStore *store.ContactStore) {
	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handlers.HandleGetRequest(w, r, contactStore)
		case http.MethodPost:
			handlers.HandlePostRequest(w, r, contactStore)
		case http.MethodPut:
			handlers.HandlePutRequest(w, r, contactStore)
		case http.MethodDelete:
			handlers.HandleDeleteRequest(w, r, contactStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})
}
