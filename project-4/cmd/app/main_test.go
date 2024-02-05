package main_test

import (
	"bytes"
	"encoding/json"
	"io"
	"log"
	"net/http"
	"net/http/httptest"
	"os"
	"strconv"
	"testing"

	"github.com/keffren/go/project-4/internal/rest"
)

var app rest.App

const tableCreationQuery = `CREATE TABLE IF NOT EXISTS products
(
    id SERIAL PRIMARY KEY,
    name TEXT NOT NULL,
    price NUMERIC(10,2) NOT NULL DEFAULT 0.00
)`

// ###############################  Integration TEST: main
func TestMain(m *testing.M) {

	app.Initialize(
		"localhost",
		5432,
		os.Getenv("TEST_DB_USERNAME"),
		os.Getenv("TEST_DB_PASSWORD"),
		os.Getenv("TEST_DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	// Go does not use an integer return value from main to indicate exit status.
	// If you’d like to exit with a non-zero status you should use os.Exit.
	os.Exit(code)
}

// ###############################   UNIT TESTS
func TestEmptyTable(t *testing.T) {
	clearTable()

	// Fetch data from products data
	req, _ := http.NewRequest("GET", "/products", nil)

	// Check if the table is empty frpm the http response
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusOK, resp.Code)
	if body := resp.Body.String(); body != "[]" {
		t.Errorf("Expected an empty array. Got %s", body)
	}
}

func TestGetNonExistentProduct(t *testing.T) {
	clearTable()

	//Fetch product with ID = 11
	req, _ := http.NewRequest("GET", "/products/11", nil)

	// Check product doesn't exist from http response
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, resp.Code)
	m := make(map[string]string)
	json.Unmarshal(resp.Body.Bytes(), &m)
	if m["error"] != "Product not found" {
		t.Errorf("Expected the 'error' key of the response to be set to 'Product not found'. Got '%s'", m["error"])
	}
}

func TestCreateProduct(t *testing.T) {

	clearTable()

	// Send POST request (add product)
	var jsonStr = []byte(`{"name":"test product", "price": 11.22}`)
	req, _ := http.NewRequest("POST", "/products", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	// Check the product has been added from the http response
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusCreated, resp.Code)

	m := make(map[string]interface{}) // can store values of different types and where the keys are strings.
	json.Unmarshal(resp.Body.Bytes(), &m)

	if m["name"] != "test product" {
		t.Errorf("Expected product name to be 'test product'. Got '%v'", m["name"])
	}

	if m["price"] != 11.22 {
		t.Errorf("Expected product price to be '11.22'. Got '%v'", m["price"])
	}

	// the id is compared to 1.0 because JSON unmarshaling converts numbers to
	// floats, when the target is a map[string]interface{}
	if m["id"] != 1.0 {
		t.Errorf("Expected product ID to be '1'. Got '%v'", m["id"])
	}
}

func TestGetProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	resp := executeRequest(req)

	checkResponseCode(t, http.StatusOK, resp.Code)
}

func TestUpdateProduct(t *testing.T) {

	clearTable()
	addProducts(1)

	// Get product with Id: 1
	req, _ := http.NewRequest("GET", "/products/1", nil)
	resp := executeRequest(req)

	originalProduct := make(map[string]interface{})
	json.Unmarshal(resp.Body.Bytes(), &originalProduct)

	// Update product with Id: 1
	var jsonStr = []byte(`{"name":"product-updated", "price": 11.22}`)
	req, _ = http.NewRequest("PUT", "/products/1", bytes.NewBuffer(jsonStr))
	req.Header.Set("Content-Type", "application/json")

	resp = executeRequest(req)

	// Check if the product with id:1 was updated
	checkResponseCode(t, http.StatusOK, resp.Code)

	m := make(map[string]interface{})
	aux, _ := io.ReadAll(resp.Body)
	json.Unmarshal(aux, &m)

	if m["id"] != originalProduct["id"] {
		t.Errorf("Expected the id to remain the same (%v). Got %v", originalProduct["id"], m["id"])
	}

	if m["name"] == originalProduct["name"] {
		t.Errorf("Expected the name to change from '%v' to '%v'. Got '%v'", originalProduct["name"], m["name"], m["name"])
	}

	if m["price"] == originalProduct["price"] {
		t.Errorf("Expected the price to change from '%v' to '%v'. Got '%v'", originalProduct["price"], m["price"], m["price"])
	}
}

func TestDeleteProduct(t *testing.T) {
	clearTable()
	addProducts(1)

	req, _ := http.NewRequest("GET", "/products/1", nil)
	resp := executeRequest(req)
	checkResponseCode(t, http.StatusOK, resp.Code)

	req, _ = http.NewRequest("DELETE", "/products/1", nil)
	resp = executeRequest(req)

	checkResponseCode(t, http.StatusOK, resp.Code)

	req, _ = http.NewRequest("GET", "/products/1", nil)
	resp = executeRequest(req)
	checkResponseCode(t, http.StatusNotFound, resp.Code)
}

// ###############################   AUX FUNCTIONS

func ensureTableExists() {

	if _, err := app.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM products")
	app.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}

func executeRequest(req *http.Request) *httptest.ResponseRecorder {
	resp := httptest.NewRecorder()
	app.Router.ServeHTTP(resp, req)

	return resp
}

func checkResponseCode(t *testing.T, expected, actual int) {
	if expected != actual {
		t.Errorf("Expected HTTP Status from response: %d. Got %d\n", expected, actual)
	}
}

func addProducts(count int) {
	if count < 1 {
		count = 1
	}

	for i := 0; i < count; i++ {
		app.DB.Exec("INSERT INTO products(name, price) VALUES($1, $2)", "Product "+strconv.Itoa(i), (i+1.0)*10)
	}
}
