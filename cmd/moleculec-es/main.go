package main

import (
	"encoding/json"
	"flag"
	"io/ioutil"
	"log"
	"os"

	"github.com/xxuejie/moleculec-es/pkg/generator"
)

var inputFile = flag.String("inputFile", "", "Input file to use")
var outputFile = flag.String("outputFile", "-", "Output file to generate, use '-' to print to stdout")

func main() {
	flag.Parse()

	content, err := ioutil.ReadFile(*inputFile)
	if err != nil {
		log.Fatal(err)
	}

	var schema generator.Schema
	err = json.Unmarshal(content, &schema)
	if err != nil {
		log.Fatal(err)
	}

	writer := os.Stdout
	if *outputFile != "-" {
		writer, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
	}

	err = generator.Generate(schema, writer)
	if err != nil {
		log.Fatal(err)
	}
}
