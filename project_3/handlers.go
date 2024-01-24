package main

import (
	"encoding/json"
	"fmt"
	"net/http"
	"strconv"
	"strings"
)

func sendJSON(data any, resp http.ResponseWriter) {

	jsonToSend, err := json.Marshal(data)

	if err != nil {
		http.Error(resp, "Error encoding JSON", http.StatusInternalServerError)
	} else {
		resp.Header().Set("content-type", "application/json")
		resp.Write(jsonToSend)
	}
}

func GetContacts(resp http.ResponseWriter, req *http.Request) {
	sendJSON(Contacts_data, resp)
}

func GetContact(resp http.ResponseWriter, req *http.Request) {

	//Obtain contact id from request path
	contactID, err := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if err != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	//Check if the contact resource exists
	//	If exists -> it'll be sent
	if _, ok := Contacts_data[contactID]; ok == true {
		sendJSON(Contacts_data[contactID], resp)
	} else {
		http.Error(resp, "This contact doesn't exist", http.StatusNotFound)
	}
}

func CreateContact(resp http.ResponseWriter, req *http.Request) {

	// Decode Person obtained from request
	var personToAdd Person
	err := json.NewDecoder(req.Body).Decode(&personToAdd)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	//Get contact id
	contactID := personToAdd.ID

	//Check whether it already exists
	if _, ok := Contacts_data[contactID]; ok == true {
		http.Error(resp, "This contact already exists", http.StatusConflict)
		return
	}

	//ADD Person to contacts
	Contacts_data[contactID] = personToAdd
	resp.WriteHeader(http.StatusOK)
}

func UpdateContact(resp http.ResponseWriter, req *http.Request) {
	// Decode Person obtained from request
	var personToUpdate Person
	err := json.NewDecoder(req.Body).Decode(&personToUpdate)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
		return
	}

	//Check it already exists
	contactID := personToUpdate.ID

	if _, ok := Contacts_data[contactID]; ok == true {

		Contacts_data[contactID] = Person{
			ID:      contactID,
			Name:    personToUpdate.Name,
			Surname: personToUpdate.Surname,
			Phone:   personToUpdate.Phone,
			Address: personToUpdate.Address,
		}
		sendJSON("Resource deleted successfully!", resp)
	} else {
		http.Error(resp, "The contact doesn't exist, please add into contacts", http.StatusNotFound)
	}
}

func DeleteContacts(resp http.ResponseWriter, req *http.Request) {
	Contacts_data = make(map[uint64]Person)
	sendJSON("Resource deleted successfully!", resp)
}

func DeleteContact(resp http.ResponseWriter, req *http.Request) {
	//Obtain contact id from request path
	contactID, err := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if err != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
		return
	}

	//Check it already exists
	if _, ok := Contacts_data[contactID]; ok == true {
		fmt.Print(contactID)
		delete(Contacts_data, contactID)
		sendJSON("Resource deleted successfully!", resp)
	} else {
		http.Error(resp, "The contact doesn't exist", http.StatusNotFound)
	}
}
