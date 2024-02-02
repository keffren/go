package rest

import (
	"database/sql"
	"fmt"
	"log"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"
)

/* The APP is comprise by:
- Router (Gorilla/mux)
- DB 	 (PostgrSQL)
*/

type App struct {
	Router *mux.Router
	DB     *sql.DB
}

/* To be useful and testable, App will need two methods:
- initialize()
	- DB init
	- Router init
- run()
	- run the web server
*/

func (a *App) Initialize(h string, p int, u string, pw string, dbName string) {

	//Connection string
	psqlconn := fmt.Sprintf("host=%s port=%d user=%s password=%s dbname=%s sslmode=disable",
		h, p, u, pw, dbName)

	//Stablish DB connection: sql.OpenDB()
	db, err := sql.Open("postgres", psqlconn)

	if err != nil {
		log.Println("Error stablishing DB connection")
		log.Fatal(err)
	} else {
		log.Printf("%v DB Connection stablished with %v user\n", dbName, u)
	}
	a.DB = db
	a.Router = mux.NewRouter()
}

func (a *App) Run(addr string) {

}
