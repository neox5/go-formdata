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
	"testing"
)

var filemap = map[string][]string{
	"documents": {"invoice.pdf", "passport.pdf", "readme.md", "recipe.docx"},
	"photos":    {"CFA00833.ARW", "CFA00833.JPG", "CFA00833.jpg"},
}

func sampleDocuments() FormDataFile {
	formDataFile := FormDataFile{}
	for _, f := range filemap["documents"] {
		formDataFile = append(formDataFile, &multipart.FileHeader{Filename: f})
	}

	return formDataFile
}
func samplePhotos() FormDataFile {
	formDataFile := FormDataFile{}
	for _, f := range filemap["photos"] {
		formDataFile = append(formDataFile, &multipart.FileHeader{Filename: f})
	}

	return formDataFile
}

func TestFileAt(t *testing.T) {
	f := sampleDocuments()

	validTestcases := []struct {
		index    int
		expected string
	}{
		{0, "invoice.pdf"},
		{1, "passport.pdf"},
		{2, "readme.md"},
		{3, "recipe.docx"},
	}

	invalidTestcases := []struct {
		index    int
		expected interface{}
	}{
		{-1, nil},
		{4, nil},
		{65535, nil},
	}

	for _, testcase := range validTestcases {
		got := f.FileAt(testcase.index).Filename
		if got != testcase.expected {
			t.Errorf("Did not get correct file at %d: expected: %s, got: %s", testcase.index, testcase.expected, got)
		}
	}

	for _, testcase := range invalidTestcases {
		got := f.FileAt(testcase.index)
		if got != nil {
			t.Errorf("Should return <nil> at index %d: got: %v", testcase.index, got)
		}
	}
}

func TestFirstFile(t *testing.T) {
	d := sampleDocuments()
	p := samplePhotos()

	got := d.FirstFile()
	if got != d[0] {
		t.Errorf("First file did not match: expected: %s, got: %s", d[0].Filename, got.Filename)
	}

	got = p.FirstFile()
	if got != p[0] {
		t.Errorf("First file did not match: expected: %s, got: %s", p[0].Filename, got.Filename)
	}
}
