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

type ValidationError struct {
	key     string
	message string
}

func (ve ValidationError) String() string {
	return fmt.Sprintf("'%s': %s", ve.key, ve.message)
}

func (v *Validation) addError(key, msg string) {
	err := &ValidationError{
		key:     key,
		message: msg,
	}
	v.data.errors = append(v.data.errors, err)
}

func (v *Validation) addRequiredError(key string) {
	v.addError(v.key, "is required")
}

func (v *Validation) addMatchError(rx *regexp.Regexp) {
	msg := fmt.Sprintf("does not match: %s", rx.String())
	v.addError(v.key, msg)
}

func (v *Validation) addMatchAtIndexError(index int, rx *regexp.Regexp) {
	msg := fmt.Sprintf("Element %d does not match: %s", index, rx.String())
	v.addError(v.key, msg)
}
