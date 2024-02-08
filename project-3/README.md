[![Link-Checker](https://github.com/keffren/go/actions/workflows/link_checker.yml/badge.svg)](https://github.com/keffren/go/actions/workflows/link_checker.yml)

# REST API WITH GORILLA/MUX

## Gorilla/Mux

**Gorilla/mux implements a request router** and dispatcher for matching incoming requests to their respective handler.

Like the standard `http.ServeMux`, `mux.Router` matches incoming requests against a list of registered routes and calls a handler for the route that matches the URL or other conditions. 

The main Gorilla/mux features are:

- It implements the `http.Handler` interface so it is compatible with the standard `http.ServeMux`.
- Requests can be matched based on URL host, path, path prefix, schemes, header and query values, HTTP methods or using custom matchers.
- Routes can be used as **subrouters**: nested routes are only tested if the parent route matches.

Gorilla/mux package official documentation: [Gorilla/mux repository](https://github.com/gorilla/mux?tab=readme-ov-file)

### What Does Router means?

A router is a component responsible for directing incoming HTTP requests to the appropriate handler functions based on the request's route or URL. The router examines the URL of an incoming HTTP request and dispatches the request to the corresponding code that should handle it.

The difference between using a router package and a custom multiplexer might not be noticeable if you're implementing a small-scale application; however, the repetitive development can be frustrating when the application grows bigger.

## Difference between Custom Mux and Gorilla/Mux

### Custom Mux

```
type contacts struct{}

func (contacts *contacts) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("content-type", "application/json")

	// CRUD operations
	if req.Method == http.MethodGet {
		//TODO
	} else if req.Method == http.MethodPost {
		//TODO
	} else if req.Method == http.MethodDelete {
		//TODO
	} else if req.Method == http.MethodPut {
		//TODO
	}
}

func main(){

    mux := http.NewServeMux()
	mux.Handle("/api/v1/contacts", &Contacts{})
    http.ListenAndServe(":3001", mux)
}
```

### Gorilla/Mux

```
func GetContacts(resp http.ResponseWriter, req *http.Request) {}

func GetContact(resp http.ResponseWriter, req *http.Request) {}

func CreateContact(resp http.ResponseWriter, req *http.Request) {}

func UpdateContact(resp http.ResponseWriter, req *http.Request) {}

func DeleteContacts(resp http.ResponseWriter, req *http.Request) {}

func DeleteContact(resp http.ResponseWriter, req *http.Request) {}

func main() {
	router := mux.NewRouter()

	router.HandleFunc("/api/v1/contacts", GetContacts).Methods("GET")
	router.HandleFunc("/api/v1/contacts/{id}", GetContact).Methods("GET")
	router.HandleFunc("/api/v1/contacts", CreateContact).Methods("POST")
	router.HandleFunc("/api/v1/contacts/{id}", UpdateContact).Methods("PUT")
	router.HandleFunc("/api/v1/contacts", DeleteContact).Methods("DELETE")
	router.HandleFunc("/api/v1/contacts", DeleteContacts).Methods("DELETE")
}
```