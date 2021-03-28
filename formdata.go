/*
 * Created on Sun Mar 28 2021
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

// package formdata provides helpers for multipart/form-data requests with no
// external dependencies.
//
// The main focus for this package is parsing, validating and accessing
// form-data from HTTP requests. The core element of this libary is FormData,
// which wraps the multipart.Form object and adds additional validation
// capabilities.
//
// Validation is written to enable chaining and therefore improve code
// readability.
package formdata

import "mime/multipart"

// FormData extends multipart.Form with additional validation capabilities.
type FormData struct {
	*multipart.Form
	errors []*ValidationError
}

// HasErrors checks if form-data has any validation errors.
func (fd *FormData) HasErrors() bool {
	return len(fd.errors) > 0
}

// Errors returns validation errors as []string. If there are no validation
// errors an empty []string is returned.
func (fd *FormData) Errors() []string {
	errors := []string{}
	if len(fd.errors) == 0 {
		return errors
	}
	for _, err := range fd.errors {
		errors = append(errors, err.String())
	}
	return errors
}

// Exists checks if FormData.Value has given key.
func (fd *FormData) Exists(key string) bool {
	_, exists := fd.Value[key]
	return exists
}

// Get gets the value array associated with the given key. If there are no
// values associated with the key, Get returns an empty []string.
func (fd *FormData) Get(key string) FormDataValue {
	if fd.Value == nil {
		return []string{}
	}
	v := fd.Value[key]
	if len(v) == 0 {
		return []string{}
	}
	return v
}

// FileExists checks if FormData.File has given key.
func (fd *FormData) FileExists(key string) bool {
	_, exists := fd.File[key]
	return exists
}

// GetFile returns the *multipart.FileHeader array associated with the given key.
// If there are no value associated with the key, GetFile returns an empty
// []*multipart.FileHeader.
func (fd *FormData) GetFile(key string) FormDataFile {
	if fd.Value == nil {
		return []*multipart.FileHeader{}
	}
	f := fd.File[key]
	if len(f) == 0 {
		return []*multipart.FileHeader{}
	}
	return f
}
