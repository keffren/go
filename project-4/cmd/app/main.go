package main

import (
	"os"

	"github.com/keffren/go/project-4/internal/rest"
)

func main() {

	app := rest.App{}

	app.Initialize(
		"localhost",
		3452,
		os.Getenv("DB_USERNAME"),
		os.Getenv("DB_PASSWORD"),
		os.Getenv("DB_NAME"),
	)

	app.Run(":3001")
}
