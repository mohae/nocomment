# nocomment
removes comments

Nocomment removes comments. Comments can either be line comments or block comments.

## Comments
Comments will be elided from the text if the beginning comment delimiter is found and is not within quoted text: quoted text starts with a `"` and ends with a `"`.

_Support for single-quotes, `'`, and/or raw quotes ` may be added._

### Line comment
For line comments, by default, nocomment interprets `#` and `//` as the beginning of a line comment. Line comments are terminated when a new line is encountered: `\r`, `\n`, or `\r\n`.

There are toggles for whether or not an octothorpe (hash), `#`, or a double slash, `//` should be accepted as the beginning of a line comment.  By default, both types of line comments are enabled.  Disabling both types of line comments will result in no line comments, which is probably not desirable.

Configuration of line comments can only occur if `NewLex()` is called.  If `NewLex()` is called; e.g. `l := NewLex(name input)`, `l.Run()` must be called using a Go routine. 

Calling `Lex` results in the lexing of the received input.
### Block comment
Nocomment uses C style block comments, `/* */`. Block comments may span new lines.

## Usage
Input is expected to be `[]byte` and the cleaned input is returned as `[]byte`.

    import github.com/mohae/nocomment

To elide comments from text:

    cleaned := Clean(input)

If you want to configure what is accepted as line-comments, the `Stripper` struct must be used.  This struct has exported methods, that accept bools, that enable the toggling of line comment types: `Stripper.SetIgnoreHash()` and `Stripper.SetIgnoreSlash()`.

    s := NewStripper()
    s.SetIgnoreHash(true) // do not accept # as the start of a line comment.

Once set, the input is cleaned with the `Clean()` method, which also accepts and returns `[]byte`

    cleaned := s.Clean(input)

## Notes:
This is based on Rob Pike's lexer design, though greatly simplified. There is code in the package that has been copied from the source code within `http://golang.org/src/text/template/parse/`. The files containing any copied code also have the original copyright notices.

For license and copyright information on this package, please read the LICENSE file.
