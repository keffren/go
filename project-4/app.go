package main

import (
	"database/sql"

	"github.com/gorilla/mux"
)

// The APP is comprise by:
//	- Router (Gorilla/mux)
//	- DB 	 (PostgrSQL)

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

// To be useful and testable, App will need two methods:
//	- initialize
//	-  run

func (a *App) Initialize(h string, p int, u string, pw string, dbName string) {

}

func (a *App) run(p string) {

}
