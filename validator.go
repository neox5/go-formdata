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
	"fmt"
	"regexp"
)

var emailRegex = regexp.MustCompile("^[a-zA-Z0-9.!#$%&'*+\\/=?^_`{|}~-]+@[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?(?:\\.[a-zA-Z0-9](?:[a-zA-Z0-9-]{0,61}[a-zA-Z0-9])?)*$")

type Validation struct {
	data   *FormData
	key    string
	isFile bool
}

// Required checks if a key exists in the form-data
func (v *Validation) Required() *Validation {
	if v.isFile {
		if _, exists := v.data.File[v.key]; !exists {
			v.addRequiredError(v.key)
		}
		return v
	}

	if _, exists := v.data.Value[v.key]; !exists {
		v.addRequiredError(v.key)
	}

	return v
}

// HasN validates if a value has the given number of elements.
func (v *Validation) HasN(count int) *Validation {
	if v.isFile {
		got := len(v.data.GetFile(v.key))
		if got != count {
			msg := fmt.Sprintf("Invalid number of elements: expected: %d, got: %d", count, got)
			v.addError(v.key, msg)
		}
		return v
	}

	got := len(v.data.Get(v.key))
	if got != count {
		msg := fmt.Sprintf("Invalid number of elements: expected: %d, got: %d", count, got)
		v.addError(v.key, msg)
	}
	return v
}

// HasNMin validates if a value has minimal N number of elements.
func (v *Validation) HasNMin(count int) *Validation {
	if v.isFile {
		got := len(v.data.GetFile(v.key))
		if got < count {
			msg := fmt.Sprintf("Invalid number of elements: expected: >=%d, got: %d", count, got)
			v.addError(v.key, msg)
		}
		return v
	}

	got := len(v.data.Get(v.key))
	if got < count {
		msg := fmt.Sprintf("Invalid number of elements: expected: >=%d, got: %d", count, got)
		v.addError(v.key, msg)
	}
	return v
}

// Match validates if the first element of the value matches the given regular
// expression.
func (v *Validation) Match(regex *regexp.Regexp) *Validation {
	if v.isFile {
		panic("Match is not supported for file validation!")
	}

	if !regex.MatchString(v.data.Get(v.key).First()) {
		v.addMatchError(regex)
	}
	return v
}

// MatchAll validates if all elements of the value matching the given regular
// expresssion.
func (v *Validation) MatchAll(regex *regexp.Regexp) *Validation {
	for i, el := range v.data.Get(v.key) {
		if !regex.MatchString(el) {
			v.addMatchAtIndexError(i, regex)
		}
	}
	return v
}

// MatchEmail validates if the frist element of the value matches an email.
func (v *Validation) MatchEmail() *Validation {
	v.Match(emailRegex)
	return v
}

// MatchAllEmail validates if all elements of the value matching an email.
func (v *Validation) MatchAllEmail() *Validation {
	v.MatchAll(emailRegex)
	return v
}
