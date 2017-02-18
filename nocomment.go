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

// Package nocomment removes line and block comments from the provided bytes.
//
// Line comments start with either // or # and end when an EOL is encountered.
// What to accept as line comments is configurable.
//
// Block comments start with /* and end with */ and can span lines.
//
// Anything within quotes, "", is ignored.
package nocomment

// Stripper handles the elision of comments from text. The style of comments to
// elide is configurable: all supported styles are elided by default.
type Stripper struct {
	// KeepCComments: do not elide C style comments (/* */).
	KeepCComments bool
	// KeepCPPComments: do not elide C++ style comments (//).
	KeepCPPComments bool
	// KeepShellComments: do not elide C style comments (#).
	KeepShellComments bool
}

// Clean removes comments from the input.
func (s *Stripper) Clean(input []byte) (b []byte, err error) {
	// make output the same cap as input
	b = make([]byte, 0, len(input))
	l := lex(input)
	for {
		t := l.nextToken()
		switch t.typ {
		case tokenCComment:
			if !s.KeepCComments { // if C comments are to be elided, don't append this token
				continue
			}
		case tokenCPPComment:
			if !s.KeepCPPComments { // if C++ comments are to be elided, don't append this token
				continue
			}
		case tokenShellComment:
			if !s.KeepShellComments { // if shell comments are to be elided, don't append this token
				continue
			}
		case tokenEOF:
			goto done
		case tokenError:
			return b, t
		}
		b = append(b, t.String()...)
	}

done:
	return b, nil
}
