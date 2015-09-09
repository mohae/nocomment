// Copyright 2015 Joel Scoble.
// This code is governed by the MIT license, please
// refer to the LICENSE file.
//
// This code is based on
// http://golang.org/src/text/template/parse/lex.go
//
// Copyright 2011 The Go Authors. All rights reserved.
// Use of this source code is governed by a BSD-style
// license that can be found in the LICENSE file.
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

const (
	lineCommentSlash = "//"
	lineCommentHash = "#"
	blockCommentBegin = "/*"
 	blockCommentEnd = "*/"
	cr = '\r'
	nl = '\n'
)

type tokenType int

const (
	tokenError tokenType = iota
	tokenEOF
	tokenText              // anything that isn't one of the following
	tokenCommentSlash      // //
	tokenCommentHash       // #
	tokenBlockCommentStart // /*
	tokenBlockCommentEnd   // */
	tokenNL                // \n
	tokenCR                // \r
	tokenQuotedText        // "
)

var key = map[string]tokenType{
	"//": tokenCommentSlash,
	"#":  tokenCommentHash,
	"/*": tokenBlockCommentStart,
	"*/": tokenBlockCommentEnd,
	"\n": tokenNL,
	"\r": tokenCR,

}

type commentType int

const (
	none commentType = iota
	commentSlash
	commentHash
	commentBlock
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
	commentTyp commentType
	ignoreHash bool
	ignoreSlash bool
	allowSingleQuote bool // whether or not `'` is supported as a quote char.
}

func newLexer(input []byte ) *lexer {
	return &lexer{
		input: input,
		state:  lexText,
		tokens: make(chan token, 2),
	}
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
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

// error returns an error token and terminates the scan by passing back a nil
// pointer that will be the next state, terminating l.run.
func (l *lexer) errorf(format string, args ...interface{}) stateFn {
	l.tokens <- token{tokenError, l.start, fmt.Sprintf(format, args...)}
	return nil
}

// ignore skips over the pending input before this point.
func (l *lexer) ignore() {
	l.start = l.pos
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

// nextToken returns the next token from the input.
func (l *lexer) nextToken() token {
	for {
		select {
		case token := <-l.tokens:
			return token
		default:
			l.state = l.state(l)
		}
	}
	panic("not reached")
}

// peek returns but does not consume the next rune in the input
func (l *lexer) peek() rune {
	r := l.next()
	l.backup()
	return r
}

// run lexes the input by executing state functions until the state is nil.
func (l *lexer) run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered
}

func lex(input []byte) *lexer {
	l := &lexer{
		input:  input,
		state:  lexText,
		tokens: make(chan token, 2),
	}
	go l.run() // concurrently run state machine
	return l
}

// lexLineComment handles scanning of line comments.
// Line comments start with either # or // and end with a new line.
func lexLineComment(l *lexer) stateFn {
	// based on type consume the start
	if l.pos > l.start {
		l.emit(tokenText)
	}
	switch l.commentTyp {
	case commentSlash:
		l.pos += Pos(len(lineCommentSlash))
	case commentHash:
		l.pos += Pos(len(lineCommentHash))
	}
	// scan until the comment is consumed: EOL is encountered
	for {
		t := l.next()
		if t == cr {
			// if this is \r\n consume the \n
			if l.peek() == nl {
				t = l.next()
			}
			break
		}
		if t == nl {
			break
		}
		if t == eof {
			break
		}
	}
	// comment is done, ignore processed runes and continue lexing
	l.ignore()
	return lexText
}

// lexBlockComment handles the scanning of block comments.
// Block comments start with a /* and end with */. They may
// span new lines
func lexBlockComment(l *lexer) stateFn {
	if l.pos > l.start {
		l.emit(tokenText)
	}
	l.pos += Pos(len(blockCommentBegin))
	// find end of comment or error if none
	i := bytes.Index(l.input[l.pos:], []byte(blockCommentEnd))
	if i < 0 {
		return l.errorf("unclosed block comment")
	}
	l.pos += Pos(i + len(blockCommentEnd))
	l.ignore()
	return lexText
	// comment is done, ignore processed runes and continue lexing
}

// lexReturn handles a carriage return, `\r`.
func lexReturn(l *lexer) stateFn {
	l.pos += Pos(len(string(cr)))
	l.emit(tokenCR)
	return lexText
}

// lexNewLine handles a new line, `\n`
func lexNewLine(l *lexer) stateFn {
	l.pos += Pos(len(string(nl)))
	l.emit(tokenNL)
	return lexText
}

// lexQuote processes everything within ""
func lexQuote(l *lexer ) stateFn {
	// consume the start quote
	l.next()
Loop:
	for {
		switch l.next() {
		case eof:
			return l.errorf("unterminated quoted string")
		case '\\':
			if l.peek() == '"' {
				l.next() // skip it
			}
		case '"':
			break Loop
		}
	}
	l.emit(tokenQuotedText)
	return lexText
}
// stateFn to process input and tokenize things
func lexText(l *lexer) stateFn {
	for {
		if !l.ignoreSlash {
			if bytes.HasPrefix(l.input[l.pos:], []byte(lineCommentSlash)) {
				l.commentTyp = commentSlash
				return lexLineComment // next state
			}
		}
		if !l.ignoreHash {
			if bytes.HasPrefix(l.input[l.pos:], []byte(lineCommentHash)) {
				l.commentTyp = commentHash
				return lexLineComment
			}
		}
		if bytes.HasPrefix(l.input[l.pos:], []byte(blockCommentBegin)) {
			l.commentTyp = commentBlock
			return lexBlockComment
		}
		switch l.peek() {
		case cr:
			if l.pos > l.start {
				l.emit(tokenText)
			}
			return lexReturn
		case nl:
			if l.pos > l.start {
				l.emit(tokenText)
			}
			return lexNewLine
		case '"':
			if l.pos > l.start {
				l.emit(tokenText)
			}
			return lexQuote
		case eof:
			goto Done
		}
		l.next()
	}
Done:
	// Correctly reached EOF.
	if l.pos > l.start {
		l.emit(tokenText)
	}
	l.emit(tokenEOF) // Useful to make EOF a token
	return nil       // Stop the run loop.
}
