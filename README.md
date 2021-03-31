# <img alt="formdata" src="https://cdn.statically.io/gh/neox5/go-formdata/main/formdata_logo.svg" width="300" />

`formdata` is a simple and idomatic Go library for `multipart/form-data`.

The main focus for this library is parsing, validating and accessing
form-data from HTTP requests. The core element of this libary is `FormData`,
which wraps the multipart.Form object and adds additional validation
capabilities.

Validation is written to enable chaining and therefore improve code readability.

## Features

- **Documentation** - Great documentation and real world examples.

- **Independent** - No external dependencies besides the Go standard library,
  meaning it won't bloat your project.

- **Parsing** - Directly parsing http.Requests into `FormData`. A wrapper which
  extends `multipart.Form` with additional validation capabilities. 

- **Chainable** - Easy and intuative validation with chainable functions 
  (examples below).

## Usage

```go
// I will copy my real life example in the next version
```