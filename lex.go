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
	return fmt.Sprintf("%q", t.value)
}

const lineCommentSlash = "//"
const lineCommentHash = "#"
const blockCommentBegin = "/*"
const blockCommentEnd = "*/"
const cr = '\r'
const nl = '\n'
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
	tokenQuotedText       // "
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

var CommentStrings = [...]string{
	none:         "none",
	commentSlash: "line comment: //",
	commentHash:  "line comment: #",
	commentBlock: "block comment",
}

const eof = -1

type stateFn func(*lexer) stateFn

type lexer struct {
	name       string     // the name of the imput; used for errors
	input      string     // the string being scanned
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

func NewLexer(name, input string ) *lexer {
	return &lexer{
		name: name,
		input: input,
		state:  lexText,
		tokens: make(chan token, 2),
	}
}

// SetIgnoreHash set's whether or not # is the start of a line comment
func (l *lexer) SetIgnoreHash(b bool) {
		l.ignoreHash = b
}

// SetIgnoreSlash set's whether or not // is the start of a line comment
func (l *lexer) SetIgnoreSlash(b bool) {
		l.ignoreSlash = b
}

// accept consumes the next rune if it's from the valid set.
func (l *lexer) accept(valid string) bool {
	if strings.IndexRune(valid, l.next()) >= 0 {
		return true
	}
	l.backup()
	return false
}

// acceptRun consumes a run of runes from the valid set
func (l *lexer) acceptRun(valid string) {
	for strings.IndexRune(valid, l.next()) >= 0 {
	}
	l.backup()
}

// backup steps back one rune. Can be called only once per call of next.
func (l *lexer) backup() {
	l.pos -= l.width
}

// emit passes an item back to the client.
func (l *lexer) emit(t tokenType) {
	l.tokens <- token{t, l.start, l.input[l.start:l.pos]}
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
	r, w := utf8.DecodeRuneInString(l.input[l.pos:])
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
func (l *lexer) Run() {
	for state := lexText; state != nil; {
		state = state(l)
	}
	close(l.tokens) // No more tokens will be delivered
}

func lex(name, input string) *lexer {
	l := &lexer{
		name:   name,
		input:  input,
		state:  lexText,
		tokens: make(chan token, 2),
	}
	go l.Run() // concurrently run state machine
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
	i := strings.Index(l.input[l.pos:], blockCommentEnd)
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
			if strings.HasPrefix(l.input[l.pos:], lineCommentSlash) {
				l.commentTyp = commentSlash
				return lexLineComment // next state
			}
		}
		if !l.ignoreHash {
			if strings.HasPrefix(l.input[l.pos:], lineCommentHash) {
				l.commentTyp = commentHash
				return lexLineComment
			}
		}
		if strings.HasPrefix(l.input[l.pos:], blockCommentBegin) {
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
