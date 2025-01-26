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

func (store *ContactStore) ListContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contactList []Contact

	for _, contact := range store.Contacts {
		contactList = append(contactList, contact)
	}

	json.NewEncoder(w).Encode(contactList)
}

func (store *ContactStore) FindContactById(w http.ResponseWriter, r *http.Request, contactId int) {
	w.Header().Set("Content-Type", "application/json")

	contact, exists := store.Contacts[contactId]
	if exists {
		json.NewEncoder(w).Encode(contact)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func (store *ContactStore) CreateContact(w http.ResponseWriter, r *http.Request) {
	var newContact Contact

	isInvalidJSON := json.NewDecoder(r.Body).Decode(&newContact) != nil
	if isInvalidJSON {
		log.Println("Error decoding JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	newContact.Id = len(store.Contacts) + 1
	store.Contacts[newContact.Id] = newContact

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(http.StatusCreated)

	isEncodingError := json.NewEncoder(w).Encode(newContact) != nil
	if isEncodingError {
		log.Println("Error encoding JSON")
		http.Error(w, "Error processing the request", http.StatusInternalServerError)
	}
}

func (store *ContactStore) UpdateContact(w http.ResponseWriter, r *http.Request, contactId int) {
	w.Header().Set("Content-Type", "application/json")
	var updatedContact Contact

	isInvalidJSON := json.NewDecoder(r.Body).Decode(&updatedContact) != nil
	if isInvalidJSON {
		log.Println("Error decoding JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	contactExists := store.Contacts[contactId] != (Contact{})
	if contactExists {
		updatedContact.Id = contactId
		store.Contacts[contactId] = updatedContact
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func (store *ContactStore) DeleteContact(w http.ResponseWriter, r *http.Request, contactId int) {
	w.Header().Set("Content-Type", "application/json")

	contactExists := store.Contacts[contactId] != (Contact{})
	if contactExists {
		delete(store.Contacts, contactId)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func handleGetRequest(w http.ResponseWriter, r *http.Request, store *ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.FindContactById(w, r, contactId)
	} else {
		store.ListContacts(w, r)
	}
}

func handlePostRequest(w http.ResponseWriter, r *http.Request, store *ContactStore) {
	store.CreateContact(w, r)
}

func handlePutRequest(w http.ResponseWriter, r *http.Request, store *ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.UpdateContact(w, r, contactId)
	} else {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
	}
}

func handleDeleteRequest(w http.ResponseWriter, r *http.Request, store *ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.DeleteContact(w, r, contactId)
	} else {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
	}
}

func main() {
	contactStore := &ContactStore{Contacts: make(map[int]Contact)}
	mux := http.NewServeMux()

	mux.HandleFunc("/contacts", func(w http.ResponseWriter, r *http.Request) {
		switch r.Method {
		case http.MethodGet:
			handleGetRequest(w, r, contactStore)
		case http.MethodPost:
			handlePostRequest(w, r, contactStore)
		case http.MethodPut:
			handlePutRequest(w, r, contactStore)
		case http.MethodDelete:
			handleDeleteRequest(w, r, contactStore)
		default:
			http.Error(w, "Method not allowed", http.StatusMethodNotAllowed)
		}
	})

	log.Println("Server running on http://localhost:8080")
	log.Fatal(http.ListenAndServe(":8080", mux))
}
