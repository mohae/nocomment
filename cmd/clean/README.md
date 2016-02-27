clean
=====

Clean removes comments from files: takes an input file, strips all comments (#, //, /*...*/), and writes the result to an output file.

## Usage

	  go get github.com/mohae/nocomment/cmd/clean

	  clean -i input.file -o output.file

## Help output

	  Usage of .lean:
	    -i string
		      input file: required (short)
	    -input string
	          input file: required
		-o string
		      output file: required (short)
	    -output string
	          output file: required


