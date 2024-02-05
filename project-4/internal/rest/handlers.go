package rest

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"strconv"

	"github.com/gorilla/mux"
	_ "github.com/lib/pq"

	"github.com/keffren/go/project-4/pkg/database"
	_ "github.com/keffren/go/project-4/pkg/database"
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

	//Check params not empty
	if h == "" {
		log.Fatal("The database host is not specified")
	}
	if u == "" || pw == "" || dbName == "" {
		log.Fatal("The database env vars are not exported")
	}

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

	a.Router.HandleFunc("/products", a.getProducts).Methods("GET")
	a.Router.HandleFunc("/products/{id:[0-9]+}", a.getProduct).Methods("GET")
	a.Router.HandleFunc("/products", a.createProduct).Methods("POST")
	a.Router.HandleFunc("/products/{id:[0-9]+}", a.updateProduct).Methods("PUT")
	a.Router.HandleFunc("/products/{id:[0-9]+}", a.deleteProduct).Methods("DELETE")

	// Execute the web server
	log.Printf("Listening at port %s\n", addr)
	log.Fatal(http.ListenAndServe(addr, a.Router))

	//Close DB
}

func (a *App) getProducts(w http.ResponseWriter, r *http.Request) {

	var p []database.Product

	p, err := database.GetProducts(a.DB)

	if err != nil {
		respondWithError(w, http.StatusBadRequest, "Internal Error")
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) getProduct(w http.ResponseWriter, r *http.Request) {

	// Get product ID from request
	vars := mux.Vars(r)

	pID, parseERR := strconv.Atoi(vars["id"])

	if parseERR != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
	}

	// Retrieve product from DB
	p := database.Product{
		ID: pID,
	}

	// Send JSON response
	if dbERR := p.GetProduct(a.DB); dbERR != nil {
		switch dbERR {
		case sql.ErrNoRows:
			respondWithError(w, http.StatusNotFound, "Product not found")
		default:
			respondWithError(w, http.StatusInternalServerError, dbERR.Error())
		}
		return
	}

	respondWithJSON(w, http.StatusOK, p)
}

func (a *App) createProduct(w http.ResponseWriter, r *http.Request) {

	// Get product from request
	reqBody, _ := io.ReadAll(r.Body)

	p := database.Product{}
	json.Unmarshal(reqBody, &p)

	//Check product name is not empty
	if p.Name == "" {
		respondWithError(w, http.StatusBadRequest, "Product not added because the name was not provided")
		return
	}

	// Add product into db
	if dbERR := p.CreateProduct(a.DB); dbERR != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal server error")
	}

	respondWithJSON(w, http.StatusCreated, p)
}

func (a *App) updateProduct(w http.ResponseWriter, r *http.Request) {

	// Get product arguments, from request, to update
	reqBody, _ := io.ReadAll(r.Body)
	updateP := database.Product{}
	json.Unmarshal(reqBody, &updateP)

	if updateP.Name == "" && updateP.Price == 0 {
		respondWithError(w, http.StatusBadRequest, "Invalid product values")
		return
	}

	// Update the product
	if err := updateP.UpdateProduct(a.DB); err != nil {
		respondWithError(w, http.StatusInternalServerError, "Internal Error")
	}

	respondWithJSON(w, http.StatusOK, "Product updated")

}

func (a *App) deleteProduct(w http.ResponseWriter, r *http.Request) {
	// Get product ID
	vars := mux.Vars(r)
	pID, parseERR := strconv.Atoi(vars["id"])

	p := database.Product{
		ID: pID,
	}

	if parseERR != nil {
		respondWithError(w, http.StatusBadRequest, "Invalid product ID")
		return
	}

	if err := p.DeleteProduct(a.DB); err != nil {
		respondWithError(w, http.StatusBadRequest, "Error deleting product")
		return
	}

	respondWithJSON(w, http.StatusOK, "Product deleted")

}

func respondWithError(w http.ResponseWriter, sc int, m string) {
	respondWithJSON(w, sc, map[string]string{"error": m})
}

func respondWithJSON(w http.ResponseWriter, sc int, payload interface{}) {
	response, _ := json.Marshal(payload)

	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(sc)
	w.Write(response)
}
