package main

import (
	"context"
	"log"
	"reflect"
	"time"

	"github.com/aws/aws-sdk-go-v2/aws"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
	"github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb"
	"github.com/aws/aws-sdk-go-v2/service/dynamodb/types"
)

/* #########################################  DYNAMODB CONTROLLER     ########################################### */

type MyDynamoDBTable struct {
	DynamoDbClient *dynamodb.Client
	TableName      string
}

// Init function deploys the DB
func (table MyDynamoDBTable) Init() {

	tableExistsERR := table.TableExists()
	if tableExistsERR == nil {
		log.Printf("DynamoDB - Couldn't create the %v table because it already exists\n", table.TableName)
		return
	}

	creationERR := table.createContactsTable()

	if creationERR != nil {
		log.Fatal("DynamoDB - Error performing Init func:", creationERR)
	}

	//Wait table creation
	waiter := dynamodb.NewTableExistsWaiter(table.DynamoDbClient)
	waiterERR := waiter.Wait(
		context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: &table.TableName,
		},
		5*time.Minute,
	)

	if waiterERR != nil {
		log.Printf("DynamoDB - There was a error during the waiting:\n%v\n", waiterERR)
	}

	log.Printf("DynamoDB - The table %v has been created!", table.TableName)

	//Populate contacts table
	contact_one := Contact{
		Id:      0,
		Name:    "Mateo",
		Surname: "Serna",
		Phone:   123456789,
		City:    "NewYork",
	}

	putItemOneERR := table.AddItem(contact_one)
	if putItemOneERR != nil {
		log.Printf("DynamoDB - Error adding a contact:\n%v\n", putItemOneERR)
		return
	}

	contact_two := Contact{
		Id:      1,
		Name:    "Joseph",
		Surname: "Smith",
		Phone:   123456788,
		City:    "Madrid",
	}

	putItemTwoERR := table.AddItem(contact_two)
	if putItemTwoERR != nil {
		log.Printf("DynamoDB - Error adding a contact:\n%v\n", putItemTwoERR)
		return
	}

	if putItemOneERR != nil && putItemTwoERR != nil {
		log.Print("DynamoDB -  The init has been performed successfully!")
	}
}

// TableExists determines whether a DynamoDB table exists.
func (table MyDynamoDBTable) TableExists() error {

	_, err := table.DynamoDbClient.DescribeTable(
		context.TODO(),
		&dynamodb.DescribeTableInput{
			TableName: &table.TableName,
		},
	)

	return err
}

// CreateContactsTable creates a DynamoDB table with a composite primary key defined as
// a numeric sort key named `id`, and a numeric partition key named `phone`.
func (table MyDynamoDBTable) createContactsTable() error {

	tableInputs := dynamodb.CreateTableInput{
		AttributeDefinitions: []types.AttributeDefinition{
			{
				AttributeName: aws.String("Id"),
				AttributeType: types.ScalarAttributeTypeN,
			},
			{
				AttributeName: aws.String("Phone"),
				AttributeType: types.ScalarAttributeTypeN,
			},
		},
		KeySchema: []types.KeySchemaElement{
			{
				AttributeName: aws.String("Id"),
				KeyType:       types.KeyTypeHash,
			},
			{
				AttributeName: aws.String("Phone"),
				KeyType:       types.KeyTypeRange,
			},
		},
		TableName: &table.TableName,
		ProvisionedThroughput: &types.ProvisionedThroughput{
			ReadCapacityUnits:  aws.Int64(2),
			WriteCapacityUnits: aws.Int64(2),
		},
	}

	tableOutput, err := table.DynamoDbClient.CreateTable(
		context.TODO(),
		&tableInputs,
	)

	if err != nil {
		log.Printf("DynamoDB - Couldn't create table %v. Here's why: %v\n", table.TableName, err)
	} else {
		tableDesc := tableOutput.TableDescription
		log.Printf("DynamoDB - The %v table has been created successfully. Here's its details: %v\n", table.TableName, tableDesc)
	}

	return err
}

// ################## CRUD OPERATIONS

// Put contact Item into Contacts DynamoDB table
func (table MyDynamoDBTable) AddItem(c Contact) error {
	newItem, marshalERR := attributevalue.MarshalMap(c)

	if marshalERR != nil {
		log.Println("DynamoDB - Error converting contact type to attributevalue type:")
		log.Fatal(marshalERR)
	}

	newItemInput := dynamodb.PutItemInput{
		Item:      newItem,
		TableName: &table.TableName,
	}

	itemOutput, putItemERR := table.DynamoDbClient.PutItem(context.TODO(), &newItemInput)

	if putItemERR == nil {
		log.Printf("DynamoDB - The item has been added successfully! Details: %v\n", itemOutput)
	} else {
		log.Printf("DynamoDB - Couldn't add the item with ID:%v. Details: %v\n", c.Id, putItemERR)
	}

	return putItemERR
}

// Retrieve a contact Item from Contacts DynamoDB table based on its KEY
func (table MyDynamoDBTable) GetItem(c Contact) (map[string]types.AttributeValue, error) {

	contactKey, getKeyERR := c.GetKey()

	if getKeyERR != nil {
		log.Fatal(getKeyERR)
	}

	item := dynamodb.GetItemInput{
		Key:       contactKey,
		TableName: &table.TableName,
	}

	itemOutput, err := table.DynamoDbClient.GetItem(context.TODO(), &item)

	if err == nil {
		log.Printf("DynamoDB - Couldn't retrieve the item with ID:%v\n%v\n", c.Id, err)
		return nil, err
	} else {
		log.Printf("DynamoDB - Item with ID:%v has been gotten successfully!\n%v\n", c.Id, itemOutput.Item)
		return itemOutput.Item, err
	}
}

// Retrieve all items from Contacts table
func (table MyDynamoDBTable) GetItems() ([]map[string]types.AttributeValue, error) {

	// Create filter for the Expression to fill the input struct with
	filter := expression.Name("Id").GreaterThanEqual(expression.Value(0))

	// Get back the item's attributes
	proj := expression.NamesList(
		expression.Name("Id"),
		expression.Name("Name"),
		expression.Name("Surname"),
		expression.Name("Phone"),
		expression.Name("City"),
	)

	//Build the expression
	expr, exprERR := expression.NewBuilder().WithFilter(filter).WithProjection(proj).Build()

	if exprERR != nil {
		log.Printf("DynamoDB - Got error building expression for scan the table. Here's the details\n%v\n", exprERR)
		return nil, exprERR
	}

	// Build the query input parameters
	scanInput := dynamodb.ScanInput{
		TableName:                 &table.TableName,
		ExpressionAttributeNames:  expr.Names(),
		ExpressionAttributeValues: expr.Values(),
		FilterExpression:          expr.Filter(),
		ProjectionExpression:      expr.Projection(),
	}

	scanOutput, err := table.DynamoDbClient.Scan(context.TODO(), &scanInput)

	if err != nil {
		log.Printf("DynamoDB - Got error scanning the table:\n%s\n", err)
	} else {
		log.Print("DynamoDB - The Contacts table scan was successfull!\n")
	}

	return scanOutput.Items, err
}

// Delete a contact Item from Contacts DynamoDB table based on its ID
func (table MyDynamoDBTable) DeleteItem(c Contact) error {

	contactKey, getKeyERR := c.GetKey()

	if getKeyERR != nil {
		log.Fatal("DynamoDB - Couldn't retrieve the Item key to delete")
	}

	deleteItemInput := dynamodb.DeleteItemInput{
		Key:       contactKey,
		TableName: &table.TableName,
	}

	_, err := table.DynamoDbClient.DeleteItem(context.TODO(), &deleteItemInput)

	if err == nil {
		log.Printf("DynamoDB - Couldn't delete item with ID:%v\n%v\n", c.Id, err)
	} else {
		log.Printf("DynamoDB - Item with ID:%v deleted successfully!\n", c.Id)
	}

	return err
}

// Update one o more item attributes
func (table MyDynamoDBTable) UpdateIem(c Contact) error {

	contactKey, getKeyERR := c.GetKey()

	if getKeyERR != nil {
		log.Fatal("DynamoDB - Couldn't retrieve the Item key to delete")
	}

	/* To iterate over a custom struct (Contacts), it needs to use 'reflection package'
	   because Go can't do it itself.
	   	- reflect.Type represents the type of a Go expression.
	   	- reflect.Value represents the value of a Go expression.
	*/
	val := reflect.ValueOf(c)
	typ := reflect.TypeOf(c)

	var updateExpression expression.UpdateBuilder

	for i := 0; i < val.NumField(); i++ {
		if typ.Field(i).Name != "Id" && typ.Field(i).Name != "Phone" {
			if !val.Field(i).IsZero() {
				updateExpression = updateExpression.Set(
					expression.Name(typ.Field(i).Name),
					expression.Value(val.Field(i).Interface()),
				)
			}
		}
	}

	update, expBuilderERR := expression.NewBuilder().WithUpdate(updateExpression).Build()

	if expBuilderERR != nil {
		log.Printf("DynamoDB - Error building update expression:\n %v\n", expBuilderERR)
		return expBuilderERR
	}

	updateItemInput := dynamodb.UpdateItemInput{
		Key:                       contactKey,
		TableName:                 &table.TableName,
		UpdateExpression:          update.Update(),
		ExpressionAttributeNames:  update.Names(),
		ExpressionAttributeValues: update.Values(),
		ReturnValues:              types.ReturnValueAllNew,
	}

	updateItemOutput, err := table.DynamoDbClient.UpdateItem(context.TODO(), &updateItemInput)

	if err != nil {
		log.Printf("DynamoDB - Couldn't update the item with ID:%v\n%v\n", c.Id, err)
	} else {
		log.Printf("DynamoDB - Item with ID:%v updated successfully!\n Here's the details:\n%v\n", c.Id, updateItemOutput.Attributes)
	}

	return err
}
