package main

type Person struct {
	ID      uint64            `json:"id"`
	Name    string            `json:"name"`
	Surname string            `json:"surname"`
	Phone   uint              `json:"phone"`
	Address map[string]string `json:"address"`
}

/*
WARNING!
The make()function is the right way to create an empty map.
If you make an empty map in a different way and write to it, it will causes a runtime panic.
So avoid this: var Contacts_data map[uint64]Person
*/
var Contacts_data map[uint64]Person = make(map[uint64]Person)
