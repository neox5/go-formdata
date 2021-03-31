# <img alt="formdata" src="https://cdn.statically.io/gh/neox5/go-formdata/main/formdata_logo.svg" width="300" />

`formdata` is a simple and idomatic Go library for `multipart/form-data`.

The main focus for this library is parsing, validating and accessing
form-data from HTTP requests. The core element of this libary is `FormData`,
which wraps the multipart.Form object and adds additional validation
capabilities.

Validation is written to enable chaining and therefore improve code readability.

## Features

- **Parsing** - Directly parsing http.Requests into `FormData`. A wrapper which
  extends `multipart.Form` with additional validation capabilities. 

- **Chainability** - Easy and intuative validation with chainable functions 
  (examples below).

- **Independent** - No external dependencies besides the Go standard library,
  meaning it won't bloat your project.

- **Documentation** - With real world examples.

## Install

`go get -u github.com/neox5/go-formdata`

## Usage
Example shows how `formdata` helps handling a request for an email endpoint:

```go
func (s *Server) handleMailRequestV1() http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
    // set http.MaxBytesReader before parsing to limit the size of the incoming
    // request.
    r.Body = http.MaxBytesReader(w, r.Body, formdata.DefaultParseMaxMemory) // 1MB

    fd, err := formdata.Parse(r)
    if err == formdata.ErrNotMultipartFormData {
      // handle unsupported media type
      return
    }
    if err != nil {
      // handle internal server error
      return
    }

    fd.Validate("from").Required().HasN(1)
    fd.Validate("subject").Required().HasN(1)
    fd.Validate("body").Required().HasN(1)
    fd.Validate("to").Required().HasNMin(1).MatchAllEmail()

    if fd.HasErrors() {
      message := fmt.Sprintf("validation errors: %s", strings.Join(fd.Errors(), "; "))
      // handle bad request
      return
    }

    from := fd.Get("from").First()
    subject := fd.Get("subject").First()
    body := fd.Get("body").First()

    msg := NewMail(from, subject, body)

    to := fd.Get("to")

    for _, recipient := range to {
      msg.AddRecipient(recipient)
    }

    if fd.FileExists("attachment") {
      for _, file := range fd.GetFile("attachment") {
        reader, err := file.Open()
        if err != nil {
          // handle invalid attachment
          return
        }
        msg.AddReaderAttachment(file.Filename, reader)
      }
    }

    s.sendMessage(msg)
  }
}
```