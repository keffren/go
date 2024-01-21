package main

import (
	"encoding/json"
	"log"
	"net/http"
)

type Person struct {
	Name    string `json:"name"`
	Surname string `json:"surname"`
	Phone   uint   `json:"phone"`
}

// "/" handler
// Handler parameters:(     response  ,     request     )
func homeHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("This is the Home page."))
	//resp.Write([]byte{116, 104, 105, 115})

}

// "/api" handler
func apiHandler(resp http.ResponseWriter, req *http.Request) {

	html_index := `
		<!DOCTYPE html>
		<html>
		<head>
			<title>HTML Response</title>
		</head>
		<body>
			<h1>Hello, World!</h1>
		</body>
		</html>
	`
	// Set the Content-Type header to text/html
	resp.Header().Set("Content-Type", "text/html")

	// Write the HTML response
	resp.Write([]byte(html_index))
}

// "/api/people" handler
func peopleResourceHandler(resp http.ResponseWriter, req *http.Request) {

	people_data := []Person{
		{
			Name:    "Jhon",
			Surname: "Smith",
			Phone:   12345678,
		},
		{
			Name:    "Mike",
			Surname: "William",
			Phone:   87654321,
		},
	}

	// Encode people_data to JSON Object
	//json.Marshal() returns a JSON encode in a Bytes Slice ([]byte)
	res, err := json.Marshal(people_data)

	if err == nil {
		resp.Header().Set("Content-Type", "application/json")
		resp.Write(res)
	} else {
		http.Error(resp, "Error encoding JSON", http.StatusInternalServerError)
	}
}

type aboutHandler struct{}

func (ah aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to About Section!"))
}

func main() {

	// Create a new instance of ServeMux (mux)
	mux := http.NewServeMux()

	//Register the association of: endpoint-handler using a Handle Function
	mux.HandleFunc("/", homeHandler)
	mux.HandleFunc("/api", apiHandler)
	mux.HandleFunc("/api/people", peopleResourceHandler)

	//Register the association of: endpoint-handler using a Handle
	aboutHandler := aboutHandler{}
	mux.Handle("/about", aboutHandler)

	//Let's start the webServer
	log.Fatal(http.ListenAndServe(":3001", mux))
}
