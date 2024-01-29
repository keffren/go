package internal

import (
	"log"

	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

type Contact struct {
	Id      uint64 `json:"id"`
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   uint64 `json:"phone"`
	City    string `json:"city"`
}

func (c Contact) ParseToItem() (map[string]types.AttributeValue, error) {
	res := make(map[string]types.AttributeValue)

	res, err := attributevalue.MarshalMap(c)

	if err != nil {
		log.Printf("Error converting contact type to attributevalue type:\n%v", err)
	}

	return res, err
}

func (c Contact) GetKey() (map[string]types.AttributeValue, error) {
	key := map[string]uint64{
		"Id":    c.Id,
		"Phone": c.Phone,
	}

	keyConverted, err := attributevalue.MarshalMap(key)

	if err != nil {
		log.Printf("Error converting contact type to attributevalue type:\n%v", err)
	}

	return keyConverted, err
}

func ParseItemToContact(item map[string]types.AttributeValue) (Contact, error) {

	contact := Contact{}
	err := attributevalue.UnmarshalMap(item, &contact)

	return contact, err
}

func ParseListItemsToListContacts(item_ls []map[string]types.AttributeValue) ([]Contact, error) {

	contacts_ls := make([]Contact, 0)

	for _, item := range item_ls {

		contact, err := ParseItemToContact(item)

		if err != nil {
			return nil, err
		} else {
			contacts_ls = append(contacts_ls, contact)
		}
	}

	return contacts_ls, nil
}
