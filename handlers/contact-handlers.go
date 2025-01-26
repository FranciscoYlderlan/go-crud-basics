package handlers

import (
	"net/http"
	"strconv"

	"go-crud-basics/store"
)

func HandleGetRequest(w http.ResponseWriter, r *http.Request, store *store.ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.FindContactById(w, r, contactId)
	} else {
		store.ListContacts(w, r)
	}
}

func HandlePostRequest(w http.ResponseWriter, r *http.Request, store *store.ContactStore) {
	store.CreateContact(w, r)
}

func HandlePutRequest(w http.ResponseWriter, r *http.Request, store *store.ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.UpdateContact(w, r, contactId)
	} else {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
	}
}

func HandleDeleteRequest(w http.ResponseWriter, r *http.Request, store *store.ContactStore) {
	queryParams := r.URL.Query()

	hasContactID := queryParams.Get("id") != ""
	if hasContactID {
		contactId, _ := strconv.Atoi(queryParams.Get("id"))
		store.DeleteContact(w, r, contactId)
	} else {
		http.Error(w, "Contact ID is required", http.StatusBadRequest)
	}
}
