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

package formdata

import (
	"net/http"
	"strings"
)

const (
	// DefaultParseMaxMemory is the default maxMemory value for
	// ParseMultipartForm used by the Parse method.
	DefaultParseMaxMemory int64 = 2 << 19 // 1 MiB = 1,048,576 Bytes
)

// Parse envokes ParseMax(r, DefaultParseMaxMemory).
func Parse(r *http.Request) (*FormData, error) {
	return ParseMax(r, DefaultParseMaxMemory)
}

// ParseMax parses a request body as multipart/form-data. The whole
// request body is parsed and up to a total of maxMemory bytes of its file parts
// are stored in memory, with the remainder stored on disk in temporary files.
//
// To limit the size of the incoming request set http.MaxBytesReader before
// parsing.
func ParseMax(r *http.Request, maxMemory int64) (*FormData, error) {
	contentType := r.Header.Get("Content-Type")

	if !strings.HasPrefix(contentType, "multipart/form-data") {
		return nil, ErrNotMultipartFormData
	}

	if err := r.ParseMultipartForm(maxMemory); err != nil {
		return nil, err
	}

	return &FormData{
		r.MultipartForm,
		make([]*ValidationError, 0),
	}, nil
}
