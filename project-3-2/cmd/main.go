package main

func main() {

	// Load the AWS Configuration profile (~/.aws/config)
	/*cfg, err := config.LoadDefaultConfig(context.TODO())
	if err != nil {
		log.Fatal(err)
	}

	// Using the Config value, create the DynamoDB client
	svc := dynamodb.NewFromConfig(cfg)

	Table := MyDynamoDBTable{
		DynamoDbClient: svc,
		TableName:      "contacts",
	}

	//table.Init()
	log.Println(Table) */

	/* newContact := Contact{
		Id:      2,
		Name:    "Mateo",
		Surname: "Zuluaga",
		Phone:   123456788,
		City:    "NewYork",
	}

	table.AddItem(newContact)

	items_ls, err := table.GetItems()
	log.Println(err)
	log.Println(app.ParseListItemsToListContacts(items_ls))*/

}
