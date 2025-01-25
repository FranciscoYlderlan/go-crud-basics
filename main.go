package main

import (
	"encoding/json"
	"fmt"
	"html"
	"log"
	"net/http"
)

type Contact struct {
	Id    int    `json:"id"`
	Name  string `json:"name"`
	Email string `json:"email"`
	Phone string `json:"phone"`
}

type ContactStore struct {
	Contacts map[int]Contact
}

func (c *ContactStore) Create(w http.ResponseWriter, r *http.Request) {
	var newContact Contact

	err := json.NewDecoder(r.Body).Decode(&newContact)

	if err != nil {
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := len(c.Contacts) + 1

	newContact.Id = id

	c.Contacts[id] = newContact

	w.Header().Set("Context-Type", "application/json")

	json.NewEncoder(w).Encode(newContact)

	w.WriteHeader(http.StatusCreated)

}

func main() {

	// service := &ContactStore{Contacts: make(map[int]Contact)}

	mux := http.NewServeMux()

	http.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
	})

	log.Fatal(http.ListenAndServe(":8080", mux))

}
