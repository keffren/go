[![Link-Checker](https://github.com/keffren/go/actions/workflows/link_checker.yml/badge.svg)](https://github.com/keffren/go/actions/workflows/link_checker.yml)

# REST API  WITH STANDARD GO LIBRARY

The previous project focused on creating a web server using common Go packages. Now, let's apply that knowledge to build a REST API.

## CRUD Operations

CRUD stands for **C**reate, **R**ead, **U**pdate, and **D**elete, which are the four basic operations that can be performed on an API resource. These operations are fundamental in the context of data management and are often associated with database systems and web applications. 

| **CRUD** | **HTTP** | **REST** | **Behavior**
| --- | --- | --- | --- |
| Create | POST | /api/movies | Create a movie within movies collection resource |
| Read | GET | /api/movies | Read the list of movies |
| Update | PUT | /api/movies/{id} | Update a movie |
| Delete | DELETE | /api/movies/{id} | Delete a movie from its ID |

### How implement CRUD Operations

As it was commented in the last project: A mux, created using ServeMux, has a method called `Handle()`. Which has the goal of associate/register a handler function with a specific route pattern (endpoint).

Hence, within the `mux.handle()` function, the CRUD operations will be implemented. The next code represents a simple way to do it:

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
```

However, the above sample is not the recommended way about code readability principles. A better approach could be:

```
type contacts struct{}

func (c *contacts) Create(resp http.ResponseWriter, req *http.Request){
	//TODO
}
func (c *contacts) Get(resp http.ResponseWriter, req *http.Request){
	//TODO
}
func (c *contacts) Update(resp http.ResponseWriter, req *http.Request){
	//TODO
}
func (c *contacts) Delete(resp http.ResponseWriter, req *http.Request){
	//TODO
}

func (c *contacts) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("content-type", "application/json")

	// CRUD operations
	if req.Method == http.MethodGet {
		c.Get(resp,req)
	} else if req.Method == http.MethodPost {
		c.Create(resp,req)
	} else if req.Method == http.MethodDelete {
		c.Delete(resp,req)
	} else if req.Method == http.MethodPut {
		c.Update(resp,req)
	}
}
```

### Avoid handling errors in endpoints.

The API needs to route the request to the proper handler.

```
var (
	contacts_uri = regexp.MustCompile(`^\/api\/v1\/contacts[\/]*$`)
	contact_uri  = regexp.MustCompile(`^\/api\/v1\/contacts\/(\d+)$`)
)

type contacts struct{} 

func (c *contacts) ServeHTTP(resp http.ResponseWriter, req *http.Request) {

	resp.Header().Set("content-type", "application/json")

	// CRUD operations
	if req.Method == http.MethodGet && contacts_uri.MatchString(req.URL.Path){
		c.Get(resp,req)
	}
}
```

## Type Request

The `Request struct` represents an HTTP request in Go's standard library: It encapsulates various components of an HTTP request, including the request method, URL, protocol version, headers, body, and other relevant details.

The most relevant `Request fields` are:

- **Method string** - HTTP method
- **URL \*url.URL** -  URL for the request URL
- **Header Header** - Header for request headers
- **Body io.ReadCloser** - Body for the request body
- **Response \*Response** -  Information about the redirect response in the case of client redirects. 

### How print the body of req *http.Request

```
body, err := io.ReadAll(req.Body)
if err != nil {
    fmt.Println("Error reading request body:", err)
    return
}
fmt.Println("Request Body:", string(body))
```

## Test the REST API

I have used the [POSTMAN](https://www.postman.com/) desktop app to perform HTTP requests to the API.