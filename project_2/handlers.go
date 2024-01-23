package main

import (
	"encoding/json"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

var (
	contacts_uri = regexp.MustCompile(`^\/api\/v1\/contacts[\/]*$`)
	contact_uri  = regexp.MustCompile(`^\/api\/v1\/contacts\/(\d+)$`)
)

type Contacts struct{}

func (contacts *Contacts) GetContacts(resp http.ResponseWriter, req *http.Request) {
	resp_dat, err := json.Marshal(Contacts_data)

	if err != nil {
		http.Error(resp, "Error encoding JSON", http.StatusInternalServerError)
	} else {
		resp.Write(resp_dat)
	}
}

func (contacts *Contacts) CreateContact(resp http.ResponseWriter, req *http.Request) {
	var contact Person

	err := json.NewDecoder(req.Body).Decode(&contact)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)

	} else {
		Contacts_data = append(Contacts_data, contact)
		resp.WriteHeader(http.StatusOK)
	}
}

func (contacts *Contacts) DeleteContacts(resp http.ResponseWriter, req *http.Request) {
	Contacts_data = make([]Person, 0)
}

func (contacts *Contacts) DeleteContact(resp http.ResponseWriter, req *http.Request) {

	contact_id, err := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if err != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
	}

	var exist bool = false
	var i int
	contacts_data_size := len(Contacts_data)

	for i = 0; i <= contacts_data_size; i++ {
		if contact_id == Contacts_data[i].ID {
			exist = true
			break
		}
	}

	auxSlice := append(Contacts_data[:i], Contacts_data[i+1:contacts_data_size]...)
	Contacts_data = auxSlice

	if exist != false {
		resp.WriteHeader(http.StatusOK)
	} else {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
	}
}

func (contacts *Contacts) UpdateContact(resp http.ResponseWriter, req *http.Request) {
	var p Person
	err := json.NewDecoder(req.Body).Decode(&p)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
	}

	var exist bool = false
	var i int
	contacts_data_size := len(Contacts_data)

	for i = 0; i <= contacts_data_size; i++ {
		if p.ID == Contacts_data[i].ID {
			exist = true
			break
		}
	}

	if exist == false {
		Contacts_data = append(Contacts_data, p)
	} else {
		Contacts_data[i] = Person{
			ID:      Contacts_data[i].ID,
			Name:    p.Name,
			Surname: p.Surname,
			Phone:   p.Phone,
			Address: p.Address,
		}

		resp.WriteHeader(http.StatusOK)
	}
}

func (contacts *Contacts) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("content-type", "application/json")

	// CRUD operations
	if req.Method == http.MethodGet && contacts_uri.MatchString(req.URL.Path) {
		contacts.GetContacts(resp, req)

	} else if req.Method == http.MethodPost && contacts_uri.MatchString(req.URL.Path) {
		contacts.CreateContact(resp, req)

	} else if req.Method == http.MethodDelete && contacts_uri.MatchString(req.URL.Path) {
		contacts.DeleteContacts(resp, req)

	} else if req.Method == http.MethodDelete && contact_uri.MatchString(req.URL.Path) {
		contacts.DeleteContact(resp, req)

	} else if req.Method == http.MethodPut && contact_uri.MatchString(req.URL.Path) {
		contacts.UpdateContact(resp, req)
	}
}
