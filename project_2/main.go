package main

import (
	"encoding/json"
	"log"
	"net/http"
	"regexp"
	"strconv"
	"strings"
)

//############  DATA

type Person struct {
	ID      uint64            `json:"id"`
	Name    string            `json:"name"`
	Surname string            `json:"surname"`
	Phone   uint              `json:"phone"`
	Address map[string]string `json:"address"`
}

var contacts_data []Person

//############  CONTACTS RESOURCE "/api/v1/contacts"

// ###  URL Checkers
var (
	contacts_uri = regexp.MustCompile(`^\/api\/v1\/contacts[\/]*$`)
	contact_uri  = regexp.MustCompile(`^\/api\/v1\/contacts\/(\d+)$`)
)

// ###  "api/v1/contacts" Handler
type contacts struct{}

func (contacts *contacts) GetContacts(resp http.ResponseWriter, req *http.Request) {
	resp_dat, err := json.Marshal(contacts_data)

	if err != nil {
		http.Error(resp, "Error encoding JSON", http.StatusInternalServerError)
	} else {
		resp.Write(resp_dat)
	}
}

func (contacts *contacts) CreateContact(resp http.ResponseWriter, req *http.Request) {
	var contact Person

	err := json.NewDecoder(req.Body).Decode(&contact)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)

	} else {
		contacts_data = append(contacts_data, contact)
		resp.WriteHeader(http.StatusOK)
	}
}

func (contacts *contacts) DeleteContacts(resp http.ResponseWriter, req *http.Request) {
	contacts_data = make([]Person, 0)
}

func (contacts *contacts) DeleteContact(resp http.ResponseWriter, req *http.Request) {

	contact_id, err := strconv.ParseUint(strings.Split(req.URL.Path, "/")[4], 10, 64)

	if err != nil {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
	}

	var exist bool = false
	var i int
	contacts_data_size := len(contacts_data)

	for i = 0; i <= contacts_data_size; i++ {
		if contact_id == contacts_data[i].ID {
			exist = true
			break
		}
	}

	auxSlice := append(contacts_data[:i], contacts_data[i+1:contacts_data_size]...)
	contacts_data = auxSlice

	if exist != false {
		resp.WriteHeader(http.StatusOK)
	} else {
		http.Error(resp, "There is an internal error", http.StatusInternalServerError)
	}
}

func (contacts *contacts) UpdateContact(resp http.ResponseWriter, req *http.Request) {
	var p Person
	err := json.NewDecoder(req.Body).Decode(&p)

	if err != nil {
		http.Error(resp, "Error decoding JSON", http.StatusInternalServerError)
	}

	var exist bool = false
	var i int
	contacts_data_size := len(contacts_data)

	for i = 0; i <= contacts_data_size; i++ {
		if p.ID == contacts_data[i].ID {
			exist = true
			break
		}
	}

	if exist == false {
		contacts_data = append(contacts_data, p)
	} else {
		contacts_data[i] = Person{
			ID:      contacts_data[i].ID,
			Name:    p.Name,
			Surname: p.Surname,
			Phone:   p.Phone,
			Address: p.Address,
		}

		resp.WriteHeader(http.StatusOK)
	}
}

func (contacts *contacts) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

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

func main() {

	contacts_data = append(
		contacts_data,
		Person{
			ID:      0,
			Name:    "Mike",
			Surname: "Smith",
			Phone:   12345678,
			Address: map[string]string{
				"address":  "12 Main Street",
				"postcode": "DN22 4BN",
				"city":     "Nottingham",
			},
		},
	)

	contacts_data = append(
		contacts_data,
		Person{
			ID:      1,
			Name:    "Paul",
			Surname: "Lewis",
			Phone:   87654321,
			Address: map[string]string{
				"address":  "12 London Road",
				"postcode": "NM22 9JK",
				"city":     "York",
			},
		},
	)

	mux := http.NewServeMux()

	mux.Handle("/api/v1/contacts", &contacts{})
	mux.Handle("/api/v1/contacts/", &contacts{})

	log.Fatal(http.ListenAndServe(":3001", mux))
}
