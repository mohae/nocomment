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
	name              string
	keepCComments     bool
	keepCPPComments   bool
	keepShellComments bool
	input             string
	output            string
	err               string
}

var stripperTests = []stripperTest{
	{"empty", false, false, false, "", "", ""},
	{"keepAllEmpty", true, true, true, "", "", ""},
	{"basic line", false, false, false, "Hello World", "Hello World", ""},
	{
		"remove all", false, false, false, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"Hello World", "",
	},
	{
		"keepC", true, false, false, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"/* this is a c comment */Hello World", "",
	},

	{
		"keepCPP", false, true, false, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"// this is a C++ comment\nHello World", "",
	},
	{
		"keepShell", false, false, true, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"Hello World# this is a shell comment\n", "",
	},
	{
		"keepCCPP", true, true, false, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"/* this is a c comment */// this is a C++ comment\nHello World", "",
	},
	{
		"keepCShell", true, false, true, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"/* this is a c comment */Hello World# this is a shell comment\n", "",
	},
	{
		"keepCPPShell", false, true, true, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"// this is a C++ comment\nHello World# this is a shell comment\n", "",
	},

	{
		"keepAll", true, true, true, "/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n", "",
	},
	{
		"quotedC", false, false, false, "\"/* this is a c comment */\"// this is a C++ comment\nHello World# this is a shell comment\n",
		"\"/* this is a c comment */\"Hello World", "",
	},
	{
		"quotedCPP", false, false, false, "/* this is a c comment */\"// this is a C++ comment\n\"Hello World# this is a shell comment\n",
		"\"// this is a C++ comment\n\"Hello World", "",
	},
	{
		"quotedShell", false, false, false, "/* this is a c comment */// this is a C++ comment\nHello World\"# this is a shell comment\n\"",
		"Hello World\"# this is a shell comment\n\"", "",
	},
	{
		"quotedCCPP", false, false, false, "\"/* this is a c comment */\"\"// this is a C++ comment\n\"Hello World# this is a shell comment\n",
		"\"/* this is a c comment */\"\"// this is a C++ comment\n\"Hello World", "",
	},

	{
		"quotedCShell", false, false, false, "\"/* this is a c comment */\"// this is a C++ comment\nHello World\"# this is a shell comment\n\"",
		"\"/* this is a c comment */\"Hello World\"# this is a shell comment\n\"", "",
	},
	{
		"quotedCPPShell", false, false, false, "/* this is a c comment */\"// this is a C++ comment\n\"Hello World\"# this is a shell comment\n\"",
		"\"// this is a C++ comment\n\"Hello World\"# this is a shell comment\n\"", "",
	},
	{
		"quotedAll", false, false, false, "\"/* this is a c comment */\"\"// this is a C++ comment\n\"Hello World\"# this is a shell comment\n\"",
		"\"/* this is a c comment */\"\"// this is a C++ comment\n\"Hello World\"# this is a shell comment\n\"", "",
	},
	{
		"brokenBlockQuote", false, false, false, "/* this is a c comment // this is a C++ comment\n\"Hello World\"# this is a shell comment\n\"",
		"", "index 0: unclosed block comment",
	},
	{
		"unclosedQuote", false, false, false, "hello \"/* this is a c comment */// this is a C++ comment\nHello World# this is a shell comment\n",
		"", "index 6: unterminated quoted string",
	},
}

func TestStripper(t *testing.T) {
	var s Stripper
	for _, test := range stripperTests {
		s.KeepCComments = test.keepCComments
		s.KeepCPPComments = test.keepCPPComments
		s.KeepShellComments = test.keepShellComments
		result, err := s.Clean([]byte(test.input))
		if err != nil {
			if err.Error() != test.err {
				t.Errorf("%s: got %q want %q", test.name, err, test.err)
			}
			continue
		}
		if err == nil && test.err != "" {
			t.Errorf("%s: got no error; wanted %q", test.name, test.err)
			continue
		}
		if string(result) != test.output {
			t.Errorf("%s: got %q want %q\n", test.name, string(result), test.output)
		}
	}
}

// This tests that a bug that resulted in unterminated quotes string error
// when a windows path separator was right before a quote. Since \ can denote
// an escape sequence, we need to figure out if it is an escape or not. If it
// is an escape and it escapes a quote, then it's not the end of a quote.
func TestUnterminatedQuotedStringBug(t *testing.T) {
	tests := []struct {
		val      string
		expected string
	}{
		{`"inline=dir c:\\"`, `"inline=dir c:\\"`},
		// if this wasn't handled properly, it would be considers an unterminated quoted string.
		{`"inline=dir \"  "`, `"inline=dir \"  "`},
	}

	var s Stripper
	for i, test := range tests {
		result, err := s.Clean([]byte(test.val))
		if err != nil {
			t.Fatalf("%d: unexpected error: %s", i, err)
		}
		if string(result) != test.expected {
			t.Errorf("%d: got %q want %q\n", i, string(result), test.expected)
		}
	}
}
