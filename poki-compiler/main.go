package main

import (
	"flag"
	"io/ioutil"
	"log"
	"os"
	"pokic/parser"
)

func main() {
	input := flag.String("input", "", "The input file to process with pokic")
	output := flag.String("output", "", "The output file to write to")

	flag.Parse()

	doku, err := parser.ParseFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*output, []byte(doku.Output()), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
