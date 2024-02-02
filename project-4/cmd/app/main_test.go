package main_test

import (
	"log"
	"os"
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

func TestMain(m *testing.M) {

	app.Initialize(
		"localhost",
		5432,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"))

	ensureTableExists()

	code := m.Run()

	clearTable()

	// Go does not use an integer return value from main to indicate exit status.
	// If you’d like to exit with a non-zero status you should use os.Exit.
	os.Exit(code)
}

func ensureTableExists() {

	if _, err := app.DB.Exec(tableCreationQuery); err != nil {
		log.Fatal(err)
	}
}

func clearTable() {
	app.DB.Exec("DELETE FROM products")
	app.DB.Exec("ALTER SEQUENCE products_id_seq RESTART WITH 1")
}
