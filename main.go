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
		log.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	id := len(c.Contacts) + 1

	newContact.Id = id

	c.Contacts[id] = newContact

	w.Header().Set("Content-Type", "application/json")

	w.WriteHeader(http.StatusCreated)

	err = json.NewEncoder(w).Encode(newContact)
	if err != nil {
		log.Println("Error encoding JSON:", err)
		http.Error(w, err.Error(), http.StatusInternalServerError)
		return
	}

	log.Println("Contact created:", newContact)

}

func main() {

	service := &ContactStore{Contacts: make(map[int]Contact)}

	mux := http.NewServeMux()

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		fmt.Println(r.URL)
		switch r.Method {
		case http.MethodGet:
			fmt.Fprintf(w, "Hello, %q", html.EscapeString(r.URL.Path))
		case http.MethodPost:
			service.Create(w, r)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
	log.Println("Server started on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
