package main

import (
	"flag"
	"log"
	"pokic/parser"
)

func main() {
	input := flag.String("input", "", "The input file to process with pokic")
	_ = flag.String("output-basename", "", "The basename of the output files")

	flag.Parse()

	doku, err := parser.ParseFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	println(doku.Output())
}
