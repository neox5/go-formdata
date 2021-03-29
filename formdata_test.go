/*
 * Created on Mon Mar 29 2021
 *
 * MIT License
 *
 * Copyright (c) 2021, Christian Faustmann / neox5, <faustmannchr@gmail.com>
 *
 * Permission is hereby granted, free of charge, to any person obtaining a copy
 * of this software and associated documentation files (the "Software"), to deal
 * in the Software without restriction, including without limitation the rights
 * to use, copy, modify, merge, publish, distribute, sublicense, and/or sell
 * copies of the Software, and to permit persons to whom the Software is
 * furnished to do so, subject to the following conditions:
 *
 * The above copyright notice and this permission notice shall be included in
 * all copies or substantial portions of the Software.
 *
 * THE SOFTWARE IS PROVIDED "AS IS", WITHOUT WARRANTY OF ANY KIND, EXPRESS OR
 * IMPLIED, INCLUDING BUT NOT LIMITED TO THE WARRANTIES OF MERCHANTABILITY,
 * FITNESS FOR A PARTICULAR PURPOSE AND NONINFRINGEMENT. IN NO EVENT SHALL THE
 * AUTHORS OR COPYRIGHT HOLDERS BE LIABLE FOR ANY CLAIM, DAMAGES OR OTHER
 * LIABILITY, WHETHER IN AN ACTION OF CONTRACT, TORT OR OTHERWISE, ARISING FROM,
 * OUT OF OR IN CONNECTION WITH THE SOFTWARE OR THE USE OR OTHER DEALINGS IN THE
 * SOFTWARE.
 */

package formdata

import (
	"mime/multipart"
	"testing"
)

func emptyFormData() *FormData {
	return &FormData{
		&multipart.Form{
			Value: make(map[string][]string),
			File:  make(map[string][]*multipart.FileHeader),
		},
		make([]*ValidationError, 0),
	}
}

func TestValidate(t *testing.T) {
	fd := emptyFormData()
	fd.Value["email"] = append(fd.Value["email"], "bigboss@example.com")

	if fd.HasErrors() {
		t.Errorf("Errors before validation: expected: 0, got: %d", len(fd.Errors()))
	}

	fd.Validate("email")
	if fd.HasErrors() {
		t.Errorf("Errors after initialize validation: expected: 0, got: %d", len(fd.Errors()))
	}

	fd.Validate("email").Required()
	if fd.HasErrors() {
		t.Errorf("Required validation error: expected: 0, got: %d", len(fd.Errors()))
	}
}

func TestExists(t *testing.T) {
	fd := emptyFormData()
	fd.Value["email"] = append(fd.Value["email"], "bigboss@example.com")

	if !fd.Exists("email") {
		t.Error("Email key does not exist in form-data")
	}

	if fd.Exists("address") {
		t.Error("Address should not exist in form-data")
	}
}

func TestErrors(t *testing.T) {
	fd := emptyFormData()
	fd.Value["email"] = append(fd.Value["email"], "example.com")

	errs := fd.Errors()
	got := len(errs)

	if got != 0 {
		t.Errorf("There should be no errors in empty form-data: got %d", got)
	}

	fd.Validate("email").MatchEmail()

	errs = fd.Errors()
	got = len(errs)

	if got != 1 {
		t.Errorf("Error count mismatch: expected: 1, got %d", got)
	}
}

func TestGet(t *testing.T) {
	fdWithNilValue := &FormData{
		&multipart.Form{
			Value: nil,
			File:  make(map[string][]*multipart.FileHeader),
		},
		make([]*ValidationError, 0),
	}

	// if Value is nil any check should return an empty string array
	if len(fdWithNilValue.Get("username")) != 0 || len(fdWithNilValue.Get("email")) != 0 {
		t.Errorf("Form-data with empty value field should always return empty string: got: %d, got: %d", len(fdWithNilValue.Get("username")), len(fdWithNilValue.Get("email")))
	}

	fd := emptyFormData()
	fd.Value["email"] = append(fd.Value["email"], "ceo@example.com")

	doesNotExist := fd.Get("username")
	if len(doesNotExist) != 0 {
		t.Errorf("Should not find any value: got: %s", doesNotExist)
	}

	doesExist := fd.Get("email")
	if len(doesExist) != 1 {
		t.Errorf("Array length mismatch: expected: 1, got: %d", len(doesExist))
	}
}
