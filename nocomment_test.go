// Copyright 2016 Joel Scoble
// Licensed under the Apache License, Version 2.0 (the "License"); you may not
// use this file except in compliance with the License. You may obtain a copy
// of the License at
//
// http://www.apache.org/licenses/LICENSE-2.0
//
// Unless required by applicable law or agreed to in writing, software
// distributed under the License is distributed on an "AS IS" BASIS, WITHOUT
// WARRANTIES OR CONDITIONS OF ANY KIND, either express or implied. See the
// License for the specific language governing permissions and limitations
// under the License.

package nocomment

import (
	"testing"
)

type stripperTest struct {
	name        string
	ignoreHash  bool
	ignoreSlash bool
	input       []byte
	output      string
}

var stripperTests = []stripperTest{
	{"ignoreBothEmpty", true, true, []byte(""), ""},
	{"ignoreBoth", true, true, []byte("//this is a comment\rHello World# another comment\r"),
		"//this is a comment\rHello World# another comment\r"},
	{"ignoreNeither", false, false, []byte("//this is a comment\rHello World# another comment\r"),
		"Hello World"},
	{"ignoreNeitherNoTrailingNL", false, false, []byte("//this is a comment\rHello World# another comment"),
		"Hello World"},
	{"ignoreSlash", false, true, []byte("//this is a comment\rHello World# another comment\r"),
		"//this is a comment\rHello World"},
	{"ignoreHash", true, false, []byte("//this is a comment\rHello World# another comment\r"),
		"Hello World# another comment\r"},
	{"blockComments", false, false, []byte("/* block comment\r\n*/\r\nHello World"),
		"\r\nHello World"},
}

type cleanTest struct {
	name   string
	input  []byte
	output string
}

var cleanTests = []cleanTest{
	{"Empty", []byte(""), ""},
	{"line comments", []byte("//this is a comment\rHello World# another comment\r"), "Hello World"},
	{"blockComments", []byte("/* block comment\r\n*/\r\nHello World"), "\r\nHello World"},
}

func TestStripper(t *testing.T) {
	for _, test := range stripperTests {
		s := NewStripper()
		s.SetIgnoreHash(test.ignoreHash)
		s.SetIgnoreSlash(test.ignoreSlash)
		result := s.Clean(test.input)
		if string(result) != test.output {
			t.Errorf("%s: expected %s got %s\n", test.name, test.output, string(result))
		}
	}
}

func TestClean(t *testing.T) {
	for _, test := range cleanTests {
		result := Clean(test.input)
		if string(result) != test.output {
			t.Errorf("%s: expected %s got %s\n", test.name, test.output, string(result))
		}
	}
}
