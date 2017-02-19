package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"
	"path/filepath"

	"github.com/mohae/nocomment"
)

var (
	app     = filepath.Base(os.Args[0])
	in, out string
)

func init() {
	flag.StringVar(&in, "input", "", "input file: required")
	flag.StringVar(&in, "i", "", "input file: required (short)")
	flag.StringVar(&out, "output", "", "output file: required")
	flag.StringVar(&out, "o", "", "output file: required (short)")
}

func main() {
	flag.Parse()

	if in == "" {
		fmt.Fprintf(os.Stderr, "%s: \"input\" is a required flag\n", app)
		flag.Usage()
		os.Exit(1)
	}

	// If the output wasn't specified, save
	if out == "" {
		fmt.Fprintf(os.Stderr, "\"output\" is a required flag\n")
		flag.Usage()
		os.Exit(1)
	}
	// read the input
	b, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	var s nocomment.Stripper
	b, err = s.Clean(b)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error removing comments: %s\n", app, err)
		os.Exit(1)
	}

	err = ioutil.WriteFile(out, b, 0644)
	if err != nil {
		fmt.Fprintf(os.Stderr, "%s: error writing file: %s\n", app, err)
		os.Exit(1)
	}
	os.Exit(0)
}
