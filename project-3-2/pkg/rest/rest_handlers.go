package rest

import (
	"encoding/json"
	"log"
	"net/http"
	"strconv"
	"strings"

	"github.com/keffren/go/project-3-2/internal"
)

type ContactsHandler struct {
	Database *internal.MyDynamoDBTable
}

func sendJSON(data any, resp http.ResponseWriter) {

	jsonToSend, err := json.Marshal(data)

	if err != nil {
		http.Error(resp, "Error encoding JSON", http.StatusInternalServerError)
	} else {
		resp.Header().Set("content-type", "application/json")
		resp.Write(jsonToSend)
	}
}

func (contactsHandler ContactsHandler) GetContacts(resp http.ResponseWriter, req *http.Request) {

	//Get table Contacts items, all of them
	items, getItemsERR := contactsHandler.Database.GetItems()

	if getItemsERR != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	//Parse the list of items to a list of Contacts
	contacts, err := internal.ParseListItemsToListContacts(items)

	if err != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
	} else {
		sendJSON(contacts, resp)
	}
}

func (contactsHandler ContactsHandler) GetContact(resp http.ResponseWriter, req *http.Request) {
	//Obtain contact id from request path
	contactID, parseUintERR := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if parseUintERR != nil {
		log.Printf("Rest handler - Error parsing string to uint64:\n%v\n", parseUintERR)
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	// Get item from Database
	item, getItemERR := contactsHandler.Database.GetItem(contactID)
	if getItemERR != nil {
		log.Printf("Rest handler - Error getting a item:\n%v\n", getItemERR)
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	} else if len(item) == 0 {
		http.Error(resp, "Contact doesn't exist", http.StatusNotFound)
		return
	}

	// Parse the Item to Contact type
	contact, parseItemERR := internal.ParseItemToContact(item)
	if parseItemERR != nil {
		log.Printf("Rest handler - Error parsing item to contact:\n%v\n", parseItemERR)
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	//Send contact as JSON
	sendJSON(contact, resp)
}

func (contactsHandler ContactsHandler) CreateContact(resp http.ResponseWriter, req *http.Request) {

	// Decode Person obtained from request
	var contactToAdd internal.Contact
	decodeContactErr := json.NewDecoder(req.Body).Decode(&contactToAdd)

	if decodeContactErr != nil {
		log.Printf("Rest handler - Error parsing json to contact:\n%v\n", decodeContactErr)
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	//Get contact id
	contactID := contactToAdd.Id

	//Check whether it already exists
	item, exists := contactsHandler.Database.GetItem(contactID)

	if exists != nil {
		http.Error(resp, "Ups... Internal Error", http.StatusInternalServerError)
		return
	} else if len(item) > 0 {
		http.Error(resp, "The contact already exists", http.StatusConflict)
		return
	}

	//ADD Person to contacts
	addContactERR := contactsHandler.Database.AddItem(contactToAdd)

	if addContactERR != nil {
		log.Printf("Rest handler - Error adding a new contact:\n%v\n", addContactERR)
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
	} else {
		resp.WriteHeader(http.StatusOK)
	}

}

func (contactsHandler ContactsHandler) UpdateContact(resp http.ResponseWriter, req *http.Request) {
	// Decode Contact obtained from request
	var contactToUpdate internal.Contact
	decodeContactErr := json.NewDecoder(req.Body).Decode(&contactToUpdate)

	if decodeContactErr != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	//Check it already exists
	//Get contact id
	contactID := contactToUpdate.Id

	//Check whether it already exists
	item, exists := contactsHandler.Database.GetItem(contactID)

	if exists != nil {
		http.Error(resp, "Ups... Internal Error", http.StatusInternalServerError)
		return
	} else if len(item) == 0 {
		http.Error(resp, "Contact doesn't exist", http.StatusNotFound)
		return
	}

	//Add contact
	if updateContactERR := contactsHandler.Database.UpdateIem(contactToUpdate); updateContactERR == nil {
		sendJSON("Resource updated successfully!", resp)
	} else {
		log.Printf("Rest handler - Error updating contact:\n%v\n", updateContactERR)
		http.Error(resp, "Sorry ... Couldn't update the contact", http.StatusInternalServerError)
	}
}

func (contactsHandler ContactsHandler) DeleteContacts(resp http.ResponseWriter, req *http.Request) {
	sendJSON("This operation is disabled", resp)
}

func (contactsHandler ContactsHandler) DeleteContact(resp http.ResponseWriter, req *http.Request) {
	//Obtain contact id from request path
	contactID, parseUintERR := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if parseUintERR != nil {
		log.Printf("Rest handler - Error parsing contact id:\n%v\n", parseUintERR)
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	//Delete action
	deleteContactERR := contactsHandler.Database.DeleteItem(contactID)

	if deleteContactERR == nil {
		sendJSON("Resource deleted successfully!", resp)
	} else {
		http.Error(resp, "Ups... Internal Error", http.StatusInternalServerError)
	}
}
