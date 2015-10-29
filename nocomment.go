// Package nocomment removes line and block comments from the provided bytes.
//
// Line comments start with either // or # and end when an EOL is encountered.
// What to accept as line comments can be set.
//
// Block comments start with /* and end with */ and can span lines.
//
// Anything in quotes, "", is ignored.
package nocomment

// Stripper handles the elision of comments from text.
type Stripper struct {
	*lexer
}

// NewStripper returns a Stripper.
func NewStripper() *Stripper {
	return &Stripper{lexer: newLexer([]byte(""))}
}

// Clean removes comments from the input.
func (s *Stripper) Clean(input []byte) []byte {
	// make output the same cap as input
	output := make([]byte, 0, len(input))
	s.lexer.input = input
	go s.lexer.run()
	for {
		token := s.lexer.nextToken()
		if token.typ == tokenEOF || token.typ == tokenError {
			break
		}
		output = append(output, token.String()...)
	}
	return output
}

// SetIgnoreHash sets whether or not the hash (octothorpe), `#`, should be
// ignored as comments.
func (s *Stripper) SetIgnoreHash(b bool) {
	s.lexer.ignoreHash = b
}

// SetIgnoreSlash sets whether or not double slashes, `//`, should be
// ignored as comments.
func (s *Stripper) SetIgnoreSlash(b bool) {
	s.lexer.ignoreSlash = b
}

// Clean cleans the input of comments using nocomment's defaults.
func Clean(input []byte) []byte {
	s := NewStripper()
	return s.Clean(input)
}
