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
// This code is based on
// http://golang.org/src/text/template/parse/lex.go
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the GO.LICENSE file.

package nocomment

import (
	"bytes"
	"fmt"
	"strings"
	"unicode/utf8"
)

// Pos is a byte position in the original input text.
type Pos int

type token struct {
	typ   tokenType
	pos   Pos
	value string
}

func (t token) String() string {
	switch {
	case t.typ == tokenEOF:
		return "EOF"
	case t.typ == tokenError:
		return t.value
	}
	return fmt.Sprintf("%s", t.value)
}

func (t token) Error() string {
	return fmt.Sprintf("index %d: %s", int(t.pos), t.value)
}

const (
	cppComment    = "//"
	shellComment  = "#"
	cCommentBegin = "/*"
	cCommentEnd   = "*/"
	cr            = '\r'
	nl            = '\n'
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenText          // anything that isn't one of the following
	tokenCPPComment    // //
	tokenShellComment  // #
	tokenCCommentStart // /*
	tokenCCommentEnd   // */
	tokenCComment      // /* */
	tokenQuotedText    // text that is quoted
	tokenDoubleQuote   // "
)

var key = map[string]tokenType{
	"//": tokenCPPComment,
	"#":  tokenShellComment,
	"/*": tokenCCommentStart,
	"*/": tokenCCommentEnd,
	"\"": tokenDoubleQuote,
}

type commentType int

const (
	none commentType = iota
	// C++ style comments
	CPPComment
	// shell style comments
	ShellComment
	// C style comments
	CComment
	// quotes aren't comment types but they are here as an unexported value
	// because we need to handle quoted text (which may have comment delimeters
	// in them which should not be processed as comments)
	// feels weird: -mohae
	doubleQuote // ""
)

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	input      []byte     // the string being scanned
	state      stateFn    // the next lexing function to enter
	pos        Pos        // current position of this item
	start      Pos        // start position of this item
	width      Pos        // width of last rune read from input
	lastPos    Pos        // position of most recent item returned by nextItem
	tokens     chan token // channel of scanned tokens
	parenDepth int        // nesting depth of () exprs <- probably not needed
}

func lex(input []byte) *lexer {
	l := lexer{
		input:  input,
		state:  lexText,
		tokens: make(chan token, 2),
	}
	go l.run()
	return &l
}

// run lexes the input by executing state functions until the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered
}

// next returns the next rune in the input.
func (l *lexer) next() rune {
	if int(l.pos) >= len(l.input) {
		l.width = 0
		return eof
	}
	r, w := utf8.DecodeRune(l.input[l.pos:])
	l.width = Pos(w)
	l.pos += l.width
	return r
}

// peek returns but does not consume the next rune in the input
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// backup steps back one rune. Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, string(l.input[l.start:l.pos])}
	l.start = l.pos
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// error returns an error token and terminates the scan by passing back a nil
// pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// nextToken returns the next token from the input.
func (l *lexer) nextToken() token {
	tkn := <-l.tokens
	l.lastPos = tkn.pos
	return tkn
}

// drain the channel so the lex go routine will exit: called by caller.
func (l *lexer) drain() {
	for range l.tokens {
	}
}

// stateFn to process input and tokenize things
func lexText(l *lexer) stateFn {
	for {
		is, typ := l.atComment()
		if is {
			if l.pos > l.start {
				l.emit(tokenText)
			}
			switch typ {
			case CPPComment:
				return lexCPPComment
			case ShellComment:
				return lexShellComment
			case CComment:
				return lexCComment
			//case quoteSingle:
			//return lexSingleQuote
			case doubleQuote:
				return lexDoubleQuote
				//case quoteTick:
				//return lexQuote
			}
		}
		if l.next() == eof {
			break
		}
	}
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(tokenText)
	}
	l.emit(tokenEOF) // Useful to make EOF a token
	return nil       // Stop the run loop.
}

// atComment returns if the next rune(s) are either a comment or a quote and if
// so, its type.
func (l *lexer) atComment() (is bool, typ commentType) {
	r, s := utf8.DecodeRune(l.input[l.pos:])
	// with one character, only ShellComment or quote will match, We check to see
	// which it was and if it was neither process an additional rune, which is
	// probably unnecessary.
	t, ok := key[string(r)]
	if ok {
		switch t {
		case tokenShellComment:
			return true, ShellComment
		case tokenDoubleQuote:
			return true, doubleQuote
		}
	}
	// otherwise get a second rune
	rr, _ := utf8.DecodeRune(l.input[int(l.pos)+s:])
	t, ok = key[string(r)+string(rr)]
	if !ok {
		return false, none
	}
	switch t {
	case tokenCPPComment:
		return true, CPPComment
	case tokenCCommentStart:
		return true, CComment
	}
	return false, none
}

// lexCPPComment handles lexing of C++ style comments: // to \n or eof
func lexCPPComment(l *lexer) stateFn {
	// scan until the comment is consumed: EOL is encountered
	for {
		r := l.next()
		if r == nl || r == eof {
			break
		}
	}
	// comment is done, ignore processed runes and continue lexing
	l.emit(tokenCPPComment)
	return lexText
}

// lexShellComment handles lexing of shell style comments: # to \n or eof
func lexShellComment(l *lexer) stateFn {
	// scan until the comment is consumed: EOL is encountered
	for {
		r := l.next()
		if r == nl || r == eof {
			break
		}
	}
	// comment is done, ignore processed runes and continue lexing
	l.emit(tokenShellComment)
	return lexText
}

// lexCComment handles the lexing of C style comments: they start with /* and
// end with */; they may span new lines
func lexCComment(l *lexer) stateFn {
	l.pos += Pos(len(cCommentBegin))
	// find end of comment or error if none
	i := bytes.Index(l.input[l.pos:], []byte(cCommentEnd))
	if i < 0 {
		return l.errorf("unclosed block comment")
	}
	l.pos += Pos(i + len(cCommentEnd))
	l.emit(tokenCComment)
	return lexText
	// comment is done, ignore processed runes and continue lexing
}

// lexQuote processes everything within ""
func lexDoubleQuote(l *lexer) stateFn {
	// consume the start quote
	l.next()
Loop:
	for {
		switch l.next() {
		case eof:
			return l.errorf("unterminated quoted string")
		case '\\':
			r := l.peek()
			// There are two things to look for: is the next char another \ or is it a quote?
			switch r {
			case '\\':
				l.next() // it doesn't start an escape so it is consumed
			case '"':
				l.next() // it's an escaped quote so it should be consumed as it's not the end of the quoted text.
			}
		case '"':
			break Loop
		}
	}
	l.emit(tokenQuotedText)
	return lexText
}
