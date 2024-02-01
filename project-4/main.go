package main

import (
	"os"
)

func main() {

	app := App{}

	app.Initialize(
		"localhost",
		3452,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	app.run(":3001")
}
