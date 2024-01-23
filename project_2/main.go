package main

import (
	"log"
	"net/http"
)

func main() {

	Contacts_data = append(
		Contacts_data,
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

	Contacts_data = append(
		Contacts_data,
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

	mux.Handle("/api/v1/contacts", &Contacts{})
	mux.Handle("/api/v1/contacts/", &Contacts{})

	log.Fatal(http.ListenAndServe(":3001", mux))
}
