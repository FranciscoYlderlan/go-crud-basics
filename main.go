package main

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
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

func (c *ContactStore) List(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contacts []Contact

	for _, ct := range c.Contacts {
		contacts = append(contacts, ct)
	}

	json.NewEncoder(w).Encode(contacts)

}

func (c *ContactStore) Find(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	if contact, ok := c.Contacts[id]; ok {
		json.NewEncoder(w).Encode(contact)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}

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

func (c *ContactStore) Update(w http.ResponseWriter, r *http.Request, id int) {

	w.Header().Set("Content-Type", "application/json")

	var contactUpdated Contact

	err := json.NewDecoder(r.Body).Decode(&contactUpdated)

	if err != nil {
		log.Println("Error decoding JSON:", err)
		http.Error(w, err.Error(), http.StatusBadRequest)
		return
	}

	if _, ok := c.Contacts[id]; ok {
		contactUpdated.Id = id
		c.Contacts[id] = contactUpdated

	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}

}

func (c *ContactStore) Delete(w http.ResponseWriter, r *http.Request, id int) {
	w.Header().Set("Content-Type", "application/json")

	if _, ok := c.Contacts[id]; ok {
		delete(c.Contacts, id)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}

}

func handleGetContacts(w http.ResponseWriter, r *http.Request, service *ContactStore) {
	q := r.URL.Query()

	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Find(w, r, id)
	} else {
		service.List(w, r)
	}

}

func handlePostContact(w http.ResponseWriter, r *http.Request, service *ContactStore) {
	service.Create(w, r)
}

func handlePutContact(w http.ResponseWriter, r *http.Request, service *ContactStore) {
	q := r.URL.Query()

	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Update(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func handleDeleteContact(w http.ResponseWriter, r *http.Request, service *ContactStore) {
	q := r.URL.Query()

	if q.Get("id") != "" {
		id, _ := strconv.Atoi(q.Get("id"))
		service.Delete(w, r, id)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}

}

func main() {

	service := &ContactStore{Contacts: make(map[int]Contact)}

	mux := http.NewServeMux()

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {

		switch r.Method {
		case http.MethodGet:
			handleGetContacts(w, r, service)
		case http.MethodPost:
			handlePostContact(w, r, service)
		case http.MethodPut:
			handlePutContact(w, r, service)
		case http.MethodDelete:
			handleDeleteContact(w, r, service)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}

	})
	log.Println("Server started on http://localhost:8080")

	log.Fatal(http.ListenAndServe(":8080", mux))

}
