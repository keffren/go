# Create a Web Server with Standard Library

## What is the standard library?

The standard library is a compilation of Go packages that offer common functionalities such as:

- String manipulation
- I/O, math
- HTTP client/server. 

It provides two benefits:

1. It furnishes a consistent API usable across Go programs.
1. It has been battle-tested by the Go community.

Consequently, ***Go programs employing the standard library should be more reliable and easier to maintain***, as they don't depend on external libraries.

## What is a MUX?

In Go, a "mux" typically refers to a router or multiplexer used in web development to handle HTTP requests. The term "mux" is short for "multiplexer."

> Multiplexer is a device that enables the simultaneous transmission of several messages or signals over one communications channel.

![](https://www.electronics-tutorials.ws/wp-content/uploads/2018/05/combination-multiplexer1.gif)

In the context of Go's standard library, the "net/http" package includes a basic ServeMux type, which is a simple HTTP request multiplexer. It allows you to associate different handler functions with different URL patterns.

### How does a Mux works?

Let's deep a bit more. A Mux is defined in Go as a `ServeMux` type, which is a `struct` that holds a couple of different data.

```
type ServeMux struct {
    mu    sync.RWMutex
    m     map[string]muxEntry
    es    []muxEntry
    hosts bool
}

type muxEntry struct {
    h       Handler
    pattern string
}
```

*What Do ServeMux fields represent?*

- `ServeMux.mu` A synchronization mutex to safely handle concurrent access to the ServeMux instance.
- `ServeMux.m` holds a map of URL pattern that is paired with muxEntry. This map is used to efficiently look up the handler for a given URL pattern.
- `ServeMux.es` is an ordered list of muxEntry objects used for pattern matching.
- `ServeMux.hosts` A boolean indicating whether the ServeMux is capable of handling requests based on the host header. If hosts is true, it means that the ServeMux considers the host information when routing requests.

*What Does muxEntry struct means?*

`muxEntry` is a struct that **represents an entry** in the routing table of the ServeMux. Each entry associates a URL pattern with a corresponding handler (h).

### ServeMux Handlers Function

A mux, created using ServeMux, has a method called `HandleFunc()`. Which is used to associate/register a handler function with a specific route pattern (endpoint). 
This takes in two parameters: the target endpoint, and a handler.

```
mux := http.NewServeMux()

mux.HandleFunc("/", homeHandler)

func homeHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("This is the Home page."))
}
```

The `resp http.ResponseWriter` parameter can also write any data to the response:

- A text (slice of type Byte)
    ```
    func homeHandler(resp http.ResponseWriter, req *http.Request) {
        resp.Write([]byte("This is the Home page."))
    }
    ```
- simple HTML
    ```
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
    ```
- JSON data.
    ```
    type Person struct {
        Name    string `json:"name"`
        Surname string `json:"surname"`
        Phone   uint   `json:"phone"`
    }

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
    ```

### ServeMux Handlers

A mux, created using ServeMux, has a method called `Handle()`. Which has the same goal as `HandleFunc()`: associate/register a handler function with a specific route pattern (endpoint). 

```
type aboutHandler struct{}

func (ah aboutHandler) ServeHTTP(w http.ResponseWriter, r *http.Request) {
	w.Write([]byte("Welcome to About Section!"))
}

func main() {
	mux := http.NewServeMux()

	aboutHandler := aboutHandler{}
	mux.Handle("/about", aboutHandler)
}
```
The difference is `mux.Handle` has a different parameter: a Handler parameter which uses `ServeHTTP` method. Which is an associated method.

> What Does associated method is? It refers to a method that is associated with a specific type. Unlike some other object-oriented programming languages, Go doesn't have traditional classes. Instead, Go uses a struct-based approach to define types and allows methods to be associated with these types.

### Creating a custom mux

Here's a brief example of using a mux in Go:

```
package main

import (
	"net/http"
)

func main() {
    // Create a new ServeMux instance (mux)
    mux := http.NewServeMux()

    //Register the associations of: endpoint-handler using a Handle Function
    mux.HandleFunc("/", homeHandler)

    //Let's start the webServer
    http.ListenAndServe(":3001", mux)
}

// Define a handler for the root ("/") path
func homeHandler(resp http.ResponseWriter, req *http.Request) {
	resp.Write([]byte("This is the Home page."))
}
```

## Last

Lastly, I highly recommend reading this [post](https://dev.to/jpoly1219/what-even-is-a-mux-4fng) from *Dev.to* Which explains the above concepts very well.