# REST API WITH POSTGRESQL **[WIP]**

TODO

## Project Layout

```
├── cmd
│   └── app
│       └── main.go
│       └── main_test.go
├── internal
│   └── rest
│       └── handlers.go
├── pkg
│   └── database
│       ├── model.go
│       └── schema.sql
├── test
│   └── main_test.go
├── go.mod
├── go.sum
```

### Go directories

- **`/cmd`**
Main applications for this project
- **`/internal`**
Private application and library code
- **`/pkg`**
Library code that's ok to use by external applications
- **`/test`**
Additional external test apps and test data

## Testing Go package

Package testing provides support for automated testing of Go packages. It is intended to be used in concert with the `go test` command, which automates execution of any function of the form. 

Testing types:

- **`*testing.T`**
    - T stands for "testing."
    - This type is the most commonly used in Go testing. It represents a testing context and provides methods to report test failures and log messages.
- **`*testing.B`**
    - B stands for "benchmark."
    - This type is used for benchmarking functions. It provides methods for measuring the time taken by a function.
- **`*testing.M`**
    - M stands for "main."
    - This type is used for writing custom test binaries with a main function. It allows you to customize the behavior of the go test command.
- **`*testing.F`**
    - F stands for "flag."
    - This type is not commonly used in regular tests; it is used for custom test flags.

Useful resources:

- [Testing package](https://pkg.go.dev/testing)
- [Quick tutorial](https://github.com/golang-standards/project-layout)

## Extra resources

- [Postgres in GO](https://hevodata.com/learn/golang-postgres/)
- [Testing package](https://www.youtube.com/watch?v=FjkSJ1iXKpg&ab_channel=GolangDojo)
