package store

import (
	"encoding/json"
	"log"
	"net/http"

	"go-crud-basics/models"
)

type ContactStore struct {
	Contacts map[int]models.Contact
}

func NewContactStore() *ContactStore {
	return &ContactStore{Contacts: make(map[int]models.Contact)}
}

func (store *ContactStore) ListContacts(w http.ResponseWriter, r *http.Request) {
	w.Header().Set("Content-Type", "application/json")
	var contactList []models.Contact

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
	var newContact models.Contact

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
	var updatedContact models.Contact

	isInvalidJSON := json.NewDecoder(r.Body).Decode(&updatedContact) != nil
	if isInvalidJSON {
		log.Println("Error decoding JSON")
		http.Error(w, "Invalid JSON format", http.StatusBadRequest)
		return
	}

	contactExists := store.Contacts[contactId] != (models.Contact{})
	if contactExists {
		updatedContact.Id = contactId
		store.Contacts[contactId] = updatedContact
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}

func (store *ContactStore) DeleteContact(w http.ResponseWriter, r *http.Request, contactId int) {
	w.Header().Set("Content-Type", "application/json")

	contactExists := store.Contacts[contactId] != (models.Contact{})
	if contactExists {
		delete(store.Contacts, contactId)
		w.WriteHeader(http.StatusOK)
	} else {
		http.Error(w, "Contact not found", http.StatusNotFound)
	}
}
