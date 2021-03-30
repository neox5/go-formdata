/*
 * Created on Tue Mar 30 2021
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
	"regexp"
	"strings"
	"testing"
)

var specialFilemap = map[string][]string{
	"empty": {},
}

var specialValuemap = map[string][]string{
	"empty":       {},
	"emptystring": {""},
}

func populatedFormData() *FormData {
	fd := emptyFormData()

	for key, files := range filemap {
		for _, f := range files {
			fd.File[key] = append(fd.File[key], &multipart.FileHeader{Filename: f})
		}
	}

	for key, values := range valuemap {
		fd.Value[key] = append(fd.Value[key], values...)
	}

	// special cases
	for key, files := range specialFilemap {
		for _, f := range files {
			fd.File[key] = append(fd.File[key], &multipart.FileHeader{Filename: f})
		}
	}

	for key, values := range specialValuemap {
		fd.Value[key] = append(fd.Value[key], values...)
	}

	return fd
}

func TestRequired(t *testing.T) {
	fd := populatedFormData()

	// No errors after initialization
	if fd.HasErrors() {
		t.Errorf("Errors before validation: expected: 0, got: %d", len(fd.Errors()))
	}

	// Positive required validation
	fd.Validate("emails").Required()
	if fd.HasErrors() {
		t.Errorf("Required validation error: got: %s", strings.Join(fd.Errors(), " "))
	}
	fd.ValidateFile("documents").Required()
	if fd.HasErrors() {
		t.Errorf("Required file validation error: got:%s", strings.Join(fd.Errors(), " "))
	}

	// Negative reqired validation
	fd.Validate("invalidkey").Required()
	if !fd.HasErrors() {
		t.Errorf("Required validation error missing for \"invalidkey\"")
	}
	fd.ValidateFile("invalidkey").Required()
	if !fd.HasErrors() || len(fd.Errors()) != 2 {
		t.Errorf("Required file validation error missing for \"invalidkey\"")
	}
}

func TestHasN(t *testing.T) {
	fd := populatedFormData()

	// No errors after initialization
	if fd.HasErrors() {
		t.Errorf("Errors before validation: expected: 0, got: %d", len(fd.Errors()))
	}

	// Positive hasN validation
	fd.Validate("emails").HasN(4)
	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Errorf("HasN error: got: %s", got)
	}
	fd.ValidateFile("photos").HasN(3)
	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Errorf("HasN error: got: %s", got)
	}

	// Negative hasN validation
	fd.Validate("emails").HasN(1)
	expected := "'emails': Invalid number of elements: expected: 1, got: 4"
	got := strings.Join(fd.Errors(), " ")
	if !fd.HasErrors() {
		t.Errorf("No HasN validation error!")
	}
	if got != expected {
		t.Errorf("Invalid HasN error message: expected: \"%s\", got: \"%s\"", expected, got)
	}
	// reset formdata
	fd = populatedFormData()
	fd.ValidateFile("documents").HasN(1)
	expected = "'documents': Invalid number of elements: expected: 1, got: 4"
	got = strings.Join(fd.Errors(), " ")
	if !fd.HasErrors() {
		t.Errorf("No HasN file validation error!")
	}
	if got != expected {
		t.Errorf("Invalid file HasN error message: expected: \"%s\", got: \"%s\"", expected, got)
	}
}

func TestHasNMin(t *testing.T) {
	fd := populatedFormData()

	// No errors after initialization
	if fd.HasErrors() {
		t.Errorf("Errors before validation: expected: 0, got: %d", len(fd.Errors()))
	}

	// Positive hasN validation
	fd.Validate("emails").HasNMin(2)
	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Errorf("HasNMin error: got: %s", got)
	}
	fd.ValidateFile("photos").HasNMin(1)
	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Errorf("HasNMin error: got: %s", got)
	}

	// Negative hasN validation
	fd.Validate("emails").HasNMin(5)
	expected := "'emails': Invalid number of elements: expected: >=5, got: 4"
	got := strings.Join(fd.Errors(), " ")
	if !fd.HasErrors() {
		t.Errorf("No HasNMin validation error!")
	}
	if got != expected {
		t.Errorf("Invalid HasNMin error message: expected: \"%s\", got: \"%s\"", expected, got)
	}
	// reset formdata
	fd = populatedFormData()
	fd.ValidateFile("documents").HasNMin(5)
	expected = "'documents': Invalid number of elements: expected: >=5, got: 4"
	got = strings.Join(fd.Errors(), " ")
	if !fd.HasErrors() {
		t.Errorf("No HasNMin file validation error!")
	}
	if got != expected {
		t.Errorf("Invalid file HasNMin error message: expected: \"%s\", got: \"%s\"", expected, got)
	}
}

func TestMatch(t *testing.T) {
	testcases := []struct {
		key       string
		regex     *regexp.Regexp
		hasErrors bool
	}{
		{"games", regexp.MustCompile(`.*\.(gameboy|sega|ps)$`), false},
		{"emails", regexp.MustCompile(`^[^@]*$`), true},
	}

	for _, testcase := range testcases {
		fd := populatedFormData()

		fd.Validate(testcase.key).Match(testcase.regex)
		got := fd.HasErrors()
		expected := testcase.hasErrors
		if got != expected {
			t.Logf("first elmement: %s", fd.Get(testcase.key).First())
			t.Errorf("Invalid testcase for key \"%s\" and regex \"%v\": expected: %v, got %v", testcase.key, testcase.regex, expected, got)
		}
	}
}

func TestMatchFilePanic(t *testing.T) {
	fd := populatedFormData()
	expected := "Match is not supported for file validation!"
	defer func() {
		r := recover()
		if r == nil {
			t.Errorf("Match did not panic on validating files!")

			t.Logf("recoverd from: %s", r.(string))
		}
		if r.(string) != expected {
			t.Errorf("Match paniced with invalid message: expected: \"%s\", got \"%s\"", expected, r.(string))
		}
	}()

	fd.ValidateFile("documents").Match(regexp.MustCompile(`.*\.(pdf|md|docx)$`))

}

func TestMatchAll(t *testing.T) {
	testcases := []struct {
		key       string
		regex     *regexp.Regexp
		hasErrors bool
	}{
		{"games", regexp.MustCompile(`.*\.(gameboy|sega|ps)$`), false},
		{"emails", regexp.MustCompile(`^[^@]*$`), true},
	}

	for _, testcase := range testcases {
		fd := populatedFormData()

		fd.Validate(testcase.key).MatchAll(testcase.regex)
		got := fd.HasErrors()
		expected := testcase.hasErrors
		if got != expected {
			t.Logf("first elmement: %s", fd.Get(testcase.key).First())
			t.Errorf("Invalid testcase for key \"%s\" and regex \"%v\": expected: %v, got %v", testcase.key, testcase.regex, expected, got)
		}
	}
}

func TestMatchEmail(t *testing.T) {
	fd := populatedFormData()

	fd.Validate("emails").MatchEmail()

	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Logf("first elmement: %s", fd.Get("emails").First())
		t.Errorf("Error match email: got: %s", got)
	}

	fd.Validate("games").MatchEmail()
	if !fd.HasErrors() {
		t.Errorf("No errors matching games as emails!")
	}
}

func TestMatchAllEmail(t *testing.T) {
	fd := populatedFormData()

	fd.Validate("emails").MatchAllEmail()

	if fd.HasErrors() {
		got := strings.Join(fd.Errors(), " ")
		t.Logf("first elmement: %s", strings.Join(fd.Get("emails"), " "))
		t.Errorf("Error matchall email: got: %s", got)
	}
}
