// Copyright 2015 Joel Scoble.
// This code is governed by the MIT license, please
// refer to the LICENSE file.
//
// the design of these tests is based on the tests in
// http://golang.org/src/text/template/parse/lex_test.go
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
package nocomment

import (
	"testing"
)

type lexTest struct {
	name   string
	input  string
	tokens []token
}

var tEOF = token{tokenEOF, 0, ""}
var tNL = token{tokenNL, 0, "\n"}
var tCR = token{tokenCR, 0, "\r"}

var lexTests = []lexTest{
	{"empty", "", []token{tEOF}},
	{"justText", "hello world", []token{{tokenText, 0, "hello world"}, tEOF}},
	{"simpleLineCommentSlashNL", "//this is a comment\nHello World\n",
		[]token{{tokenText, 0, "Hello World"}, tNL, tEOF}},
	{"simpleLineCommentSlashCR", "//this is a comment\rHello World\r",
		[]token{{tokenText, 0, "Hello World"}, tCR, tEOF}},
	{"simpleLineCommentSlashCRNL", "//this is a comment\r\nHello World\r\n",
		[]token{{tokenText, 0, "Hello World"}, tCR, tNL, tEOF}},
	{"prePostLineCommentSlashNL", "//this is a comment\nHello World// another comment\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentSlashCR", "//this is a comment\rHello World// another comment\r",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentSlashCRNL", "//this is a comment\r\nHello World// another comment\r\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"simpleLineCommentHashNL", "#this is a comment\nHello World\n",
		[]token{{tokenText, 0, "Hello World"}, tNL, tEOF}},
	{"simpleLineCommentHashCR", "#this is a comment\rHello World\r",
		[]token{{tokenText, 0, "Hello World"}, tCR, tEOF}},
	{"simpleLineCommentHashCRNL", "#this is a comment\r\nHello World\r\n",
		[]token{{tokenText, 0, "Hello World"}, tCR, tNL, tEOF}},
	{"prePostLineCommentHashNL", "#this is a comment\nHello World# another comment\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentHashCR", "#this is a comment\rHello World# another comment\r",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentHashCRNL", "#this is a comment\r\nHello World# another comment\r\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentSlashHashNL", "//this is a comment\nHello World# another comment\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentSlashHashCR", "//this is a comment\rHello World# another comment\r",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"prePostLineCommentSlashHashCRNL", "//this is a comment\r\nHello World# another comment\r\n",
		[]token{{tokenText, 0, "Hello World"}, tEOF}},
	{"simpleBlockCommentNL", "/*this is a comment*/\nHello World\n",
		[]token{tNL, {tokenText, 0, "Hello World"}, tNL, tEOF}},
	{"simpleBlockCommentCR", "/*this is a comment*/\rHello World\r",
		[]token{tCR, {tokenText, 0, "Hello World"}, tCR, tEOF}},
	{"simpleBlockCommentCRNL", "/*this is a comment*/\r\nHello World\r\n",
		[]token{tCR, tNL, {tokenText, 0, "Hello World"}, tCR, tNL, tEOF}},
	{"prePostBlockCommentNL", "/*this is a comment\n*/Hello World/* another comment*/\n",
		[]token{{tokenText, 0, "Hello World"}, tNL, tEOF}},
	{"prePostBlockCommentCR", "/*this is a comment\r*/Hello World/* another comment*/\r",
		[]token{{tokenText, 0, "Hello World"}, tCR, tEOF}},
	{"prePostBlockCommentCRNL", "/*this is a comment\r\n*/Hello World/* another comment*/\r\n",
		[]token{{tokenText, 0, "Hello World"}, tCR, tNL, tEOF}},
	{"simpleBlockCommentMultiLineNL", "/*this\n is a\n comment\n*/Hello World\n",
		[]token{{tokenText, 0, "Hello World"}, tNL, tEOF}},
	{"simpleBlockCommentMultiLineCR", "/*this\r is a\r comment\r*/Hello World\r",
		[]token{{tokenText, 0, "Hello World"}, tCR, tEOF}},
	{"simpleBlockCommentMultiLineCRNL", "/*this\r\n is a\r\n comment\r\n*/Hello World\r\n",
		[]token{{tokenText, 0, "Hello World"}, tCR, tNL, tEOF}},
	{"noCommentQuotedText", `This is some text. "#This is not a comment // neither is this /* or this */" sooo, no comments!`,
		[]token{{tokenText, 0, "This is some text. "}, {tokenQuotedText, 0, `"#This is not a comment // neither is this /* or this */"`}, {tokenText, 0, " sooo, no comments!"}, tEOF}},
}

type lineLexTest struct {
	name string
	ignoreHash bool
	ignoreSlash bool
	input string
	tokens []token
}

var lineLexTests = []lineLexTest{
		{"ignoreBothEmpty", true, true, "", []token{tEOF}},
		{"ignoreBoth", true, true, "//this is a comment\rHello World# another comment\r",
			[]token{{tokenText, 0, "//this is a comment"}, tCR, {tokenText, 0, "Hello World# another comment"}, tCR, tEOF}},
		{"ignoreNeither", false, false, "//this is a comment\rHello World# another comment\r",
			[]token{{tokenText, 0, "Hello World"}, tEOF}},
		{"ignoreSlash", false, true, "//this is a comment\rHello World# another comment\r",
			[]token{{tokenText, 0, "//this is a comment"}, tCR, {tokenText, 0, "Hello World"}, tEOF}},
		{"ignoreHash", true, false, "//this is a comment\rHello World# another comment\r",
			[]token{{tokenText, 0, "Hello World# another comment"}, tCR, tEOF}},
}

// collect gathers the emitted items into a slice.
func collect(t *lexTest, left, right string) (tokens []token) {
	l := lex(t.name, t.input)
	for {
		token := l.nextToken()
		tokens = append(tokens, token)
		if token.typ == tokenEOF || token.typ == tokenError {
			break
		}
	}
	return tokens
}

// collectLineTests handles testing of enabling/disabling of line comment styles
func collectLineTest(t *lineLexTest) (tokens []token) {
	l := NewLexer(t.name, t.input)
	l.SetIgnoreHash(t.ignoreHash)
	l.SetIgnoreSlash(t.ignoreSlash)
	go l.Run()
	for {
		token := l.nextToken()
		tokens = append(tokens, token)
		if token.typ == tokenEOF || token.typ == tokenError {
			break
		}
	}
	return tokens
}

func equal(i1, i2 []token, checkPos bool) bool {
	if len(i1) != len(i2) {
		return false
	}
	for k := range i1 {
		if i1[k].typ != i2[k].typ {
			return false
		}
		if i1[k].value != i2[k].value {
			return false
		}
		if checkPos && i1[k].pos != i2[k].pos {
			return false
		}
	}
	return true
}

// test comment lexing
func TestLex(t *testing.T) {
	for _, test := range lexTests {
		tokens := collect(&test, "", "")
		if !equal(tokens, test.tokens, false) {
			t.Errorf("%s: got \n\t%+v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}

// test enabling/disabling different line comment types
func TestLineLex(t *testing.T) {
	for _, test := range lineLexTests {
		tokens := collectLineTest(&test)
		if !equal(tokens, test.tokens, false) {
			t.Errorf("%s: got \n\t%+v\nexpected\n\t%v", test.name, tokens, test.tokens)
		}
	}
}
