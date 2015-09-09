# nocomment
removes comments

Nocomment removes comments. Comments can either be line comments or block comments.

## Comments
If a comment delimiter exists outside of quoted text. Quoted text starts with a `"` and ends with a `"`. Comments cannot exist within quoted text.

Support for single-quotes, `'`, and/or raw quotes `\`` may be added.

### Line comment
For line comments, by default, nocomment interprets `#` and `//` as the beginning of a line comment. Line comments are terminated when a new line is encountered.

There are toggles for whether or not an octothorpe (hash), `#`, or a double slash, `//` should be accepted as the beginning of a line comment.  By default, both types of line comments are enabled.  Disabling both types of line comments will result in no line comments, which is probably not desirable.

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
