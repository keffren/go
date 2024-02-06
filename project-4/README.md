# REST API WITH POSTGRESQL

In this project, I will perform the next tutorial, which I highly recommend to Gophers beginners: [Building and Testing a REST API in Go with Gorilla Mux and PostgreSQL using TDD](https://semaphoreci.com/community/tutorials/building-and-testing-a-rest-api-in-go-with-gorilla-mux-and-postgresql)

## Project Layout

```
├── cmd
│   └── app
│       └── main.go
│       └── main_test.go
├── internal
│   └── rest
│       └── handlers.go
├── pkg
│   └── database
│       ├── model.go
│       └── schema.sql
├── test
│   └── main_test.go
├── go.mod
├── go.sum
```

### Go directories

- **`/cmd`**
Main applications for this project
- **`/internal`**
Private application and library code
- **`/pkg`**
Library code that's ok to use by external applications
- **`/test`**
Additional external test apps and test data

## Run locally

- Run postgres as docker container
    ```
    docker run --name <container_name> -e POSTGRES_PASSWORD=<password> -p 5432:5432 -d postgres:16-alpine3.19
    ```
- Load the environment variables
    ```
    source .env
    ```
- Execute the App testing using the IDE
    ```
    go test -v ./cmd/app 
    ```
- Execute the App  using the IDE
    ```
    go run cmd/app/main.go  
    ```

## Test

### Testing GO package

```
Run source project-4/.env
2024/02/06 16:38:27 postgres DB Connection stablished with postgres user
=== RUN   TestEmptyTable
--- PASS: TestEmptyTable (0.00s)
=== RUN   TestGetNonExistentProduct
--- PASS: TestGetNonExistentProduct (0.00s)
=== RUN   TestCreateProduct
--- PASS: TestCreateProduct (0.00s)
=== RUN   TestGetProduct
--- PASS: TestGetProduct (0.00s)
=== RUN   TestUpdateProduct
--- PASS: TestUpdateProduct (0.00s)
=== RUN   TestDeleteProduct
--- PASS: TestDeleteProduct (0.00s)
PASS
ok  	github.com/keffren/go/project-4/cmd/app	0.041s
```

### POSTMAN

<details>
Here's a simple POSTMAN collection to validate the REST API
</details>
<summary>

```
{
	"info": {
		"_postman_id": "0d48581d-6386-498d-a0d9-8d29988ecac0",
		"name": "Go project n4",
		"schema": "https://schema.getpostman.com/json/collection/v2.1.0/collection.json",
		"_exporter_id": "31285487"
	},
	"item": [
		{
			"name": "127.0.0.1:3001/products",
			"request": {
				"method": "POST",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\n    \"name\": \"cup\",\n    \"price\":  34.12\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "127.0.0.1:3001/products",
					"host": [
						"127",
						"0",
						"0",
						"1"
					],
					"port": "3001",
					"path": [
						"products"
					]
				}
			},
			"response": []
		},
		{
			"name": "127.0.0.1:3001/products/1",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "127.0.0.1:3001/products/1",
					"host": [
						"127",
						"0",
						"0",
						"1"
					],
					"port": "3001",
					"path": [
						"products",
						"1"
					]
				}
			},
			"response": []
		},
		{
			"name": "127.0.0.1:3001/products",
			"request": {
				"method": "GET",
				"header": [],
				"url": {
					"raw": "127.0.0.1:3001/products",
					"host": [
						"127",
						"0",
						"0",
						"1"
					],
					"port": "3001",
					"path": [
						"products"
					]
				}
			},
			"response": []
		},
		{
			"name": "127.0.0.1:3001/products/1",
			"request": {
				"method": "PUT",
				"header": [],
				"body": {
					"mode": "raw",
					"raw": "{\n    \"name\": \"bottle\",\n    \"price\": 1.51\n}",
					"options": {
						"raw": {
							"language": "json"
						}
					}
				},
				"url": {
					"raw": "127.0.0.1:3001/products/1",
					"host": [
						"127",
						"0",
						"0",
						"1"
					],
					"port": "3001",
					"path": [
						"products",
						"1"
					]
				}
			},
			"response": []
		},
		{
			"name": "127.0.0.1:3001/products/1",
			"request": {
				"method": "DELETE",
				"header": [],
				"url": {
					"raw": "127.0.0.1:3001/products/4",
					"host": [
						"127",
						"0",
						"0",
						"1"
					],
					"port": "3001",
					"path": [
						"products",
						"4"
					]
				}
			},
			"response": []
		}
	]
}
```
</summary>

## Extra resources

- GO TESTING
    - [Testing package](https://pkg.go.dev/testing)
    - [Quick tutorial](https://github.com/golang-standards/project-layout)
- POSTGRES with Go
    - [PostgreSQL Basics](https://www.w3schools.com/postgresql/index.php)
    - [PostgreSQL in GO](https://hevodata.com/learn/golang-postgres/)
