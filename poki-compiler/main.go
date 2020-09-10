package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"
	"pokic/parser"
)

func main() {
	input := flag.String("input", "", "The input file to process with pokic")
	output := flag.String("output", "", "The output file to write to")
	parseTree := flag.Bool("parseTree", false, "If set to true, pokic will print the parsed document instead of writing to output")

	flag.Parse()

	if *parseTree {
		doku, err := parser.ParseFile(*input)
		if err != nil {
			log.Fatal(err)
		}
		doku.Verify()
		jsonified, _ := json.MarshalIndent(doku, "", "  ")
		println(string(jsonified))
		return
	}

	doku, err := parser.ParseFile(*input)
	if err != nil {
		log.Fatal(err)
	}

	err = ioutil.WriteFile(*output, []byte(doku.Output()), os.ModePerm)
	if err != nil {
		log.Fatal(err)
	}
}
