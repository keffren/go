# REST API WITH GORILLA/MUX AND DYNAMODB ***[WIP]***

The last [project](/project-3/) I have build a REST API using the Gorilla/Mux Go package. 

In this project, I will leverage the API and modify the data type to DynamoDB (NoSQL DB) instead of using data as a Go variable.

## AWS DYNAMODB

DynamoDB is a fully managed, highly available database with replication across three AZ.
It is NoSQL database (not relational database) and is one of the flagship product (star/referent product) of AWS.
DynamoDB scales to massive workload, distributed “server-less” database. That means that the user doesn’t provision any server, infrastructure.
It is fast and consistent in performance, and low latency retrieval(Single-digit millisecond).

### Data type

DynamoDB is a **key-value par** and document database.

![](/project-3-2/docs/dyanamodb_table.png)

It is similar than relational database table, but the difference are: It is not relational database (NoSQL), instead of a table is an item and it is schemaless.

## Core components of DynamoDB

- **Tables**: A table is a collection of items.
- **Items** (Rows): An item is a collection of attributes.
- **Attributes** (Columns): An attribute contains data related to each item.
- **Indexes**: An index is a hashed data structure that uniquely identifies and groups attributes.
    - There are two main types of indexes Primary and Secondary.
    - The primary function of indexes is to speed up access to our stored attributes and prevent full scans (the distinctive difference between a database and a file store).

## DATA Types

- User -> API
  - JSON -> Go Type -> DynamoDB type
- API -> User
  - DynamoDB type -> Go Type -> JSON

## AWS SDK for Go

SDK stands for *S*oftware *D*evelopment *K*it. Hence, the AWS SDK for Go simpliﬁes use of AWS services by providing a set of libraries that are consistent and familiar for Go developers. It supports higher level abstractions for simplified development.

Go AWS SDK documentation:
  - [AWS Developer Guide](https://aws.github.io/aws-sdk-go-v2/docs/getting-started/)
  - [Go API reference](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2)
  - [DynamoDB code examples for the SDK for Go V2](https://github.com/awsdocs/aws-doc-sdk-examples/tree/main/gov2/dynamodb)

The following commands show how to retrieve the standard set of SDK modules to use in the application.

```
go get github.com/aws/aws-sdk-go-v2
go get github.com/aws/aws-sdk-go-v2/config
go get github.com/aws/aws-sdk-go-v2/service/dynamodb
```

This will retrieve the core SDK module, the config module which is used for loading the AWS profile and the dynamoDB service.

### Configuration and credential file settings

I recommend reading [AWS Documentation](https://docs.aws.amazon.com/cli/latest/userguide/cli-configure-files.html), in case it's the first time setting up the AWS for the Go environment

### Loading AWS configuration profile

There are a number of ways to initialize a service API client, but the following is the most common pattern recommended to users.

To configure the SDK to use the AWS shared configuration use the following code:

```
import (
  "context"
  "log"
  "github.com/aws/aws-sdk-go-v2/config"
)

// ...

cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
  log.Fatalf("failed to load configuration, %v", err)
}
```

`config.LoadDefaultConfig(context.TODO())` will construct an [aws.Config](https://pkg.go.dev/github.com/aws/aws-sdk-go-v2/aws#Config) using the AWS configuration profile. 

Service clients can be constructed using the loaded` aws.Config`, providing a consistent pattern for constructing clients.

### Constructing a Service Client

To make calls to an AWS service, you must first construct a service client instance. A **service client** provides low-level access to every API action for that service. For example, I create an DynamoDB service client to make calls to Amazon DynamoDB APIs.

Service clients can be constructed using either the `New` or `NewFromConfig` functions available in service client’s Go package. Each function will return a `Client` struct type containing the methods for invoking the service APIs.

```
import (
	"context"
	"log"
	"github.com/aws/aws-sdk-go-v2/config"
    "github.com/aws/aws-sdk-go-v2/service/dynamodb"
)

// ...

// Load the AWS Configuration profile (~/.aws/config)
cfg, err := config.LoadDefaultConfig(context.TODO())
if err != nil {
    log.Fatal(err)
}

// Using the Config value, create the DynamoDB client
svc := dynamodb.NewFromConfig(cfg)

```

### type.AttributeValue

Represents the data for an DynamoDB Table Attribute.** Each attribute value is described as a name-value pair**. The name is the data type, and the value is the data itself.

AttributeValue types:
- AttributeValueMember**B**
  - An attribute of type Binary. For example: "B": "dGhpcyB0ZXh0IGlzIGJhc2U2NC1lbmNvZGVk"
- AttributeValueMember**BOOL**
  - An attribute of type Boolean. For example: "BOOL": true
- AttributeValueMember**BS**
  - An attribute of type Binary Set. For example: "BS": ["U3Vubnk=", "UmFpbnk=", "U25vd3k="]
- AttributeValueMember**L**
  - An attribute of type List. For example: "L": [ {"S": "Cookies"} , {"S": "Coffee"}, {"N": "3.14159"}]
- AttributeValueMember**M**
  - An attribute of type Map. For example: "M": {"Name": {"S": "Joe"}, "Age": {"N": "35"}}
- AttributeValueMember**N**
  - An attribute of type Number. For example: "N": "123.45" Numbers are sent across the network to DynamoDB as strings
- AttributeValueMember**NS**
  - An attribute of type Number Set. For example: "NS": ["42.2", "-19", "7.5", "3.14"] 
- AttributeValueMember**NULL**
  - An attribute of type Null. For example: "NULL": true
- AttributeValueMember**S**
  - An attribute of type String. For example: "S": "Hello"
- AttributeValueMember**SS**
  - An attribute of type String Set. For example: "SS": ["Giraffe", "Hippo" ,"Zebra"]

Here's a example of a item:
```
newItem := dynamodb.PutItemInput{
  Item: map[string]types.AttributeValue{
    "id": &types.AttributeValueMemberN{
      Value: "0",
    },
    "phone": &types.AttributeValueMemberN{
      Value: "123456789",
    },
    "name": &types.AttributeValueMemberS{
      Value: "Mike",
    },
    "surname": &types.AttributeValueMemberS{
      Value: "Smith",
    },
    "city": &types.AttributeValueMemberS{
      Value: "lincoln",
    },
  },
  TableName: &table.TableName,
}
```

*How Can we convert the data of `AttributeValue` in standards Go types?*

Within AWS SDK for GO, there is a [feature](https://github.com/aws/aws-sdk-go-v2/tree/main/feature/dynamodb/attributevalue) that can be used to convert a DynamoDB attribute value back into a Go data type.

```
//Must be imported through terminal: go get go get "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue" 

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue"
```

### Expression

I have used the `expression` feature within update DynamoDB controller.

The **expression package**  provides functionality for building expressions that are used in various DynamoDB operations, such as querying and updating items. This package allows you to construct complex conditions and projections without manually crafting the entire DynamoDB expression.

```
//Must be imported through terminal: go get go get "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression" 

import "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"


func (table MyDynamoDBTable) UpdateIem(c Contact) error {
  // ...
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

  //...
}
```

## Go dependencies

```
go get "github.com/aws/aws-sdk-go-v2/config"
go get "github.com/aws/aws-sdk-go-v2/service/dynamodb"
go get "github.com/aws/aws-sdk-go-v2/feature/dynamodb/attributevalue" 
go get "github.com/aws/aws-sdk-go-v2/feature/dynamodb/expression"
```