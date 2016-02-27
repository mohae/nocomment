package main

import (
	"flag"
	"fmt"
	"io/ioutil"
	"os"

	"github.com/mohae/nocomment"
)

var in, out string

func init() {
	flag.StringVar(&in, "input", "", "input file: required")
	flag.StringVar(&in, "i", "", "input file: required (short)")
	flag.StringVar(&out, "output", "", "output file: required")
	flag.StringVar(&out, "o", "", "output file: required (short)")
}

func main() {
	flag.Parse()

	if in == "" {
		fmt.Fprintln(os.Stderr, "\"input\" is a required flag")
		flag.Usage()
		os.Exit(1)
	}

	// If the output wasn't specified, save
	if out == "" {
		fmt.Fprintln(os.Stderr, "\"output\" is a required flag")
		flag.Usage()
		os.Exit(1)
	}
	// read the input
	b, err := ioutil.ReadFile(in)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}

	s := nocomment.NewStripper()
	b = s.Clean(b)
	err = ioutil.WriteFile(out, b, 0644)
	if err != nil {
		fmt.Fprintln(os.Stderr, err)
		os.Exit(1)
	}
	os.Exit(0)
}
