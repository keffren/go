package main

type Person struct {
	ID      uint64            `json:"id"`
	Name    string            `json:"name"`
	Surname string            `json:"surname"`
	Phone   uint              `json:"phone"`
	Address map[string]string `json:"address"`
}

var Contacts_data []Person
