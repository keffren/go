package internal

import (
	"log"
	"strconv"

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

func (c Contact) GetKey() map[string]types.AttributeValue {

	contactKey := map[string]types.AttributeValue{
		"Id": &types.AttributeValueMemberN{
			Value: strconv.FormatUint(uint64(c.Id), 10),
		},
	}

	return contactKey
}
