package main

import (
	"context"
	"log"
	"net/http"

	"github.com/aws/aws-sdk-go-v2/config"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/gorilla/mux"
	"github.com/keffren/go/project-3-2/internal"
	"github.com/keffren/go/project-3-2/pkg/rest"
)

func main() {

	// ######################################  DB Services

	// Load the AWS Configuration profile (~/.aws/config)
	cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	var contactsTable internal.MyDynamoDBTable = internal.MyDynamoDBTable{
		DynamoDbClient: svc,
		TableName:      "contacts",
	}

	contactsHandler := rest.ContactsHandler{
		Database: &contactsTable,
	}

	// ######################################  Web Server
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/contacts", contactsHandler.GetContacts).Methods("GET")
	/* router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", contactsHandler.GetContact).Methods("GET")
	router.HandleFunc("/api/v1/contacts", contactsHandler.CreateContact).Methods("POST")
	router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", contactsHandler.UpdateContact).Methods("PUT")
	router.HandleFunc("/api/v1/contacts", contactsHandler.DeleteContacts).Methods("DELETE")
	router.HandleFunc("/api/v1/contacts/{id:[0-9]+}", contactsHandler.DeleteContact).Methods("DELETE") */

	log.Print("Listening on port 3001")
	http.ListenAndServe(":3001", router)
}
