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
// under the License..
//
// the design of these tests is based on the tests in
// http://golang.org/src/text/template/parse/lex_test.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO.LICENSE file.

package nocomment

import (
	"testing"
)

type lexTest struct {
	name   string
	input  []byte
	tokens []token
}

var tEOF = token{tokenEOF, 0, ""}

var lexTests = []lexTest{
	{"empty", []byte(""), []token{tEOF}},
	{"justText", []byte("hello world"), []token{{tokenText, 0, "hello world"}, tEOF}},
	{"simpleLineCommentCPPNL", []byte("//this is a comment\nHello World\n"),
		[]token{{tokenText, 0, "Hello World\n"}, tEOF}},
	{"simpleLineCommentCPPCRNL", []byte("//this is a comment\r\nHello World\r\n"),
		[]token{{tokenText, 0, "Hello World\r\n"}, tEOF}},
	{"prePostLineCommentCPPNL", []byte("//this is a comment\nHello World// another comment\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},

	{"prePostLineCommentCPPCRNL", []byte("//this is a comment\r\nHello World// another comment\r\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"simpleLineCommentShellNL", []byte("#this is a comment\nHello World\n"),
		[]token{{tokenText, 0, "Hello World\n"}, tEOF}},
	{"simpleLineCommentShellCRNL", []byte("#this is a comment\r\nHello World\r\n"),
		[]token{{tokenText, 0, "Hello World\r\n"}, tEOF}},
	{"prePostLineCommentShellL", []byte("#this is a comment\nHello World# another comment\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentShellCRNL", []byte("#this is a comment\r\nHello World# another comment\r\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},

	{"prePostLineCommentShellHashNL", []byte("//this is a comment\nHello World# another comment\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentShellHashCRNL", []byte("//this is a comment\r\nHello World# another comment\r\n"),
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"simpleCCommentNL", []byte("/*this is a comment*/\nHello World\n"),
		[]token{{tokenText, 0, "\nHello World\n"}, tEOF}},
	{"simpleCCommentCRNL", []byte("/*this is a comment*/\r\nHello World\r\n"),
		[]token{{tokenText, 0, "\r\nHello World\r\n"}, tEOF}},
	{"prePostCCommentNL", []byte("/*this is a comment\n*/Hello World/* another comment*/\n"),
		[]token{{tokenText, 0, "Hello World"}, {tokenText, 0, "\n"}, tEOF}},

	{"prePostCCommentCRNL", []byte("/*this is a comment\r\n*/Hello World/* another comment*/\r\n"),
		[]token{{tokenText, 0, "Hello World"}, {tokenText, 0, "\r\n"}, tEOF}},
	{"simpleCCommentMultiLineNL", []byte("/*this\n is a\n comment\n*/Hello World\n"),
		[]token{{tokenText, 0, "Hello World\n"}, tEOF}},
	{"simpleCCommentMultiLineCRNL", []byte("/*this\r\n is a\r\n comment\r\n*/Hello World\r\n"),
		[]token{{tokenText, 0, "Hello World\r\n"}, tEOF}},
	{"noCommentQuotedText", []byte(`This is some text. "#This is not a comment // neither is this /* or this */" sooo, no comments!`),
		[]token{{tokenText, 0, "This is some text. "}, {tokenQuotedText, 0, `"#This is not a comment // neither is this /* or this */"`}, {tokenText, 0, " sooo, no comments!"}, tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest, left, right string) (tokens []token) {
	l := lex(t.input)
	for {
		tkn := l.nextToken()
		tokens = append(tokens, tkn)
		if tkn.typ == tokenEOF || tkn.typ == tokenError {
			break
		}
	}
	return tokens
}

func equal(t *testing.T, i int, i1, i2 []token) {
	if len(i1) != len(i2) {
		t.Errorf("%d: got %d tokens want %d", i, len(i1), len(i2))
		t.Errorf("%d: got\t%#v\nwant:\t%#v\n", i, i1, i2)
		return
	}
	// pos isn't checked
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			t.Errorf("%d:%d:typ: got %v want %v\ttoken: %#v", i, k, i1[k].typ, i2[k].typ, i1[k])
			continue
		}
		if i1[k].value != i2[k].value {
			t.Errorf("%d:%d:value: got %q want %q\ttoken: %#v", i, k, i1[k].value, i2[k].value, i1[k])
			continue
		}
	}
}

// test comment lexing
func TestLex(t *testing.T) {
	for i, test := range lexTests {
		tokens := collect(&test, "", "")
		equal(t, i, tokens, test.tokens)
	}
}

/*
// test enabling/disabling different line comment types
func TestLineLex(t *testing.T) {
	for _, test := range lineLexTests {
		tokens := collectLineTest(&test)
		if !equal(tokens, test.tokens, false) {
			t.Errorf("%s: got \n\t%+v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}
*/
