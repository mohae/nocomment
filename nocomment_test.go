// Copyright 2015 Joel Scoble.
// This code is governed by the MIT license, please
// refer to the LICENSE file.
package nocomment

import (
  "testing"
)

type stripperTest struct {
  name string
  ignoreHash bool
  ignoreSlash bool
  input []byte
  output string
}
var stripperTests = []stripperTest{
		{"ignoreBothEmpty", true, true, []byte(""), ""},
		{"ignoreBoth", true, true, []byte("//this is a comment\rHello World# another comment\r"),
			"//this is a comment\rHello World# another comment\r"},
		{"ignoreNeither", false, false, []byte("//this is a comment\rHello World# another comment\r"),
			"Hello World"},
		{"ignoreSlash", false, true, []byte("//this is a comment\rHello World# another comment\r"),
			"//this is a comment\rHello World"},
		{"ignoreHash", true, false, []byte("//this is a comment\rHello World# another comment\r"),
			"Hello World# another comment\r"},
    {"blockComments", false, false, []byte("/* block comment\r\n*/\r\nHello World"),
      "\r\nHello World"},
}

type stripperTest struct {
  name string
  input []byte
  output string
}
var stripperTests = []stripperTest{
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
