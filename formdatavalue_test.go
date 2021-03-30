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
	"testing"
)

var valuemap = map[string][]string{
	"emails": {"example@gmail.com", "elon@tesla.com", "richard@feynman.com", "charlie@berkshirehathaway.com"},
	"games":  {"wario_land_super_mario_land_3.gameboy", "sonic_2.sega", "tekken_3.ps"},
}

func sampleEmails() FormDataValue {
	formDataValue := FormDataValue{}
	for _, v := range valuemap["emails"] {
		formDataValue = append(formDataValue, v)
	}

	return formDataValue
}

func TestAt(t *testing.T) {
	emails := sampleEmails()

	validTestcases := []struct {
		index    int
		expected string
	}{
		{0, "example@gmail.com"},
		{1, "elon@tesla.com"},
		{2, "richard@feynman.com"},
		{3, "charlie@berkshirehathaway.com"},
	}

	invalidTestcases := []struct {
		index    int
		expected string
	}{
		{-1, ""},
		{4, ""},
		{65535, ""},
	}

	for _, testcase := range validTestcases {
		got := emails.At(testcase.index)
		if got != testcase.expected {
			t.Errorf("Did not get correct file at %d: expected: %s, got: %s", testcase.index, testcase.expected, got)
		}
	}

	for _, testcase := range invalidTestcases {
		got := emails.At(testcase.index)
		if got != "" {
			t.Errorf("Should return \"\" at index %d: got: %v", testcase.index, got)
		}
	}
}

func TestFirst(t *testing.T) {
	emails := sampleEmails()

	got := emails.First()
	if got != emails[0] {
		t.Errorf("First value did not match: expected: %s, got: %s", emails[0], got)
	}
}
