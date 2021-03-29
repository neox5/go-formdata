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
	"bytes"
	"io"
	"math/rand"
	"mime/multipart"
	"net/http"
	"testing"
)

func TestParse(t *testing.T) {
	_, err := Parse(testRequestValidContentType(t))
	if err != nil {
		t.Error(err)
	}

	_, err = Parse(testRequestInvalidContentType(t))
	if _, ok := err.(*FormDataError); err == nil || !ok {
		t.Error(err)
	}
}

func testRequestValidContentType(t *testing.T) *http.Request {
	t.Helper()
	r, boundary := testRequestWithMultipartForm(t)
	r.Header.Add("Content-Type", "multipart/form-data; boundary="+boundary)
	return r
}

func testRequestInvalidContentType(t *testing.T) *http.Request {
	t.Helper()
	r, _ := testRequestWithMultipartForm(t)
	r.Header.Add("Content-Type", "application/json")
	return r
}

func testRequestWithMultipartForm(t *testing.T) (*http.Request, string) {
	t.Helper()

	body := bytes.NewBuffer([]byte{})
	multipartWriter := multipart.NewWriter(body)

	fieldData := map[string][]string{
		"from":    {"ZionTec <noreply@example.com>"},
		"subject": {"Updates on Example project"},
		"body":    {"this is an awesome mail body"},
		"to":      {"awesomedude@gmail.com", "chairman@example.com"},
	}

	// populate form
	for key, values := range fieldData {
		for _, value := range values {
			if err := multipartWriter.WriteField(key, value); err != nil {
				t.Fatalf("WriteField %s: %v", key, err)
			}
		}
	}

	// add files
	fileWriter, err := multipartWriter.CreateFormFile("attachment", "test_file.txt")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	if _, err := io.Copy(fileWriter, bytes.NewReader([]byte("This is my second test file"))); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	fileWriter, err = multipartWriter.CreateFormFile("attachment", "test_binary.bin")
	if err != nil {
		t.Fatalf("CreateFormFile: %v", err)
	}
	randomBinary := make([]byte, 50*1024)
	rand.Read(randomBinary)
	if _, err := io.Copy(fileWriter, bytes.NewReader(randomBinary)); err != nil {
		t.Fatalf("io.Copy: %v", err)
	}

	// close multipart writer
	if err := multipartWriter.Close(); err != nil {
		t.Fatalf("fileWriter.Close: %v", err)
	}

	r, err := http.NewRequest("POST", "/", body)
	if err != nil {
		t.Errorf("http.NewRequest: %v", err)
	}

	return r, multipartWriter.Boundary()
}

// func tmpFileWithContent(t *testing.T, content string) *os.File {
// 	t.Helper()

// 	tmpFile, err := ioutil.TempFile(os.TempDir(), "tmpFileWithContent-")
// 	if err != nil {
// 		t.Fatalf("Can't create tempfile: %v", err)
// 	}
// 	t.Cleanup(func() { os.Remove(tmpFile.Name()) })

// 	text := []byte(content)
// 	if _, err := tmpFile.Write(text); err != nil {
// 		t.Fatalf("Can't write to tmpFile: %v", err)
// 	}
// 	tmpFile.Close()

// 	return tmpFile
// }
