# <img alt="formdata" src="https://cdn.statically.io/gh/neox5/go-formdata/main/formdata_logo.svg" width="300" />

[![GoDoc](https://img.shields.io/badge/go-documentation-blue.svg?style=flat-square)](https://pkg.go.dev/mod/github.com/neox5/go-formdata)

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
    r.Body = http.MaxBytesReader(w, r.Body, formdata.DefaultParseMaxMemory) // 1MiB

    // PARSE formdata from request
    fd, err := formdata.Parse(r)
    if err == formdata.ErrNotMultipartFormData {
      // ...handle unsupported media type
      return
    }
    if err != nil {
      // ...handle internal server error
      return
    }

    // VALIDATE formdata
    fd.Validate("from").Required().HasN(1)
    fd.Validate("subject").Required().HasN(1)
    fd.Validate("body").Required().HasN(1)
    fd.Validate("to").Required().HasNMin(1).MatchAllEmail()

    if fd.HasErrors() {
      message := fmt.Sprintf("validation errors: %s", strings.Join(fd.Errors(), "; "))
      // ...handle bad request
      return
    }

    // ACCESS formdata values
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
          // ...handle invalid attachment
          return
        }
        msg.AddReaderAttachment(file.Filename, reader)
      }
    }

    s.sendMessage(msg)
  }
}
```

## Parsing

### Methods

- **Parse** - envokes ParseMax(r, DefaultParseMaxMemory), DefaultParseMaxMemory = 1MiB
  
- **ParseMax** - parses a request body as multipart/form-data. The whole request
  body is parsed and up to a total of `maxMemory` bytes of its files parts are
  stored in memory, with the rmainder stored on disk in temporary files.


## FormData

`FormData` is the central element of formdata and it is returned from either 
[Parse](#parsing) or [ParseMax](#parsing).

```go
// FormData extends multipart.Form with additional validation capabilities.
type FormData struct {
	*multipart.Form
	errors []*ValidationError
}
```

### FormData Methods
- **Validate** - returns a Value Validation on the given key
- **ValidateFile** - returns a File Validation on the given key
- **HasErrors** - checks if FormData has validation errors
- **Errors** - returns validation errors as `[]string`
- **Exists** - checks if key exists in FormData.Value
- **FileExists** - checks if key exists in FormData.File
- **Get** - returns [FormDataValue](#formdatavalue) for given key
- **GetFile** - returns [FormDataFile](#formdatafile) for given key

## Validation

### Global Validation
- **Required** - add required validation, checks if key exists in FormData
- **HasN** - checks if Value/File has `N` elements
- **HasNMin** - cheks if Value/File has minimum `N` elements

### Value Validation
- **Match** - validates if the first element matches a regular expression
- **MatchAll** - validates if all elements match a given regular expression
- **MatchEmail** - validates if the first element matches an email
- **MatchAllEmail** - validates if all elements are matching an email 

### File Validation
...still loading

## FormDataValue

`FormDataValue` is the returned type of the [Get](#formdata-methods) method on the 
[FormData](#formdata) type.

### Methods
- **At** - gets element of FormDataValue at the given index
- **First** - gets the first element of FormDataValue

## FormDataFile

`FormDataFile` is the returned type of the [GetFile](#formdata-methods) method on the 
[FormData](#formdata) type.

### Methods
- **At** - gets element of FormDataFile at the given index
- **First** - gets the first element of FormDataFile

## Inspiration

This library is conceptually similar to [albrow/forms](https://github.com/albrow/forms), with the following major behavioral differences:

- Focusing only on multipart/form-data.
- Wrapping multipart.Form without redefining a new data struct.
- Support for chainable validation.
- Support for multiple files per key.

## License

`formdata` is licensed under the MIT License. See the LICENSE file for more information.