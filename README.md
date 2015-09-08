# nocomment
removes comments

Nocomment removes comments. Comments can either be line comments or block comments. 

## Comments
If a comment delimiter exists outside of quoted text. Quoted text starts with a `"` and ends with a `"`. Comments cannot exist within quoted text.

Support for single-quotes, `'`, and/or raw quotes `\`` may be added.

### Line comment
For line comments, by default, nocomment interprets `#` and `//` as the beginning of a line comment. Line comments are terminated when a new line is encountered.

There are toggles for whether or not an octothorpe. `#`, or a double slash, `//` should be accepted as the beginning of a line comment. While it is possible to have both enabled, which is nocomment's default, it is not possible to turn both off. Toggling one off will automatically enable the other option.

### Block comment
Nocomment uses C style block comments, `/* */`. Block comments may span new lines.

## Notes:
This is based on Rob Pike's lexer design, though greatly simplified. There is code in the package that has been copied from the source code within `http://golang.org/src/text/template/parse/`. The files containing any copied code also have the original copyright notices.

For license and copyright information on this package, please read the LICENSE file.
