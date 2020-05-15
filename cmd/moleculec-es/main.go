package main

import (
	"encoding/json"
	"flag"
	"io"
	"io/ioutil"
	"log"
	"os"
	"regexp"

	"github.com/xxuejie/moleculec-es/pkg/generator"
)

var inputFile = flag.String("inputFile", "", "Input file to use")
var outputFile = flag.String("outputFile", "-", "Output file to generate, use '-' to print to stdout")
var generateTS = flag.Bool("generateTypeScriptDefinition", false, "True to generate TypeScript definition")
var hasBigInt = flag.Bool("hasBigInt", false, "True to generate BigInt related functions")

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
	var tsWriter io.Writer
	if *outputFile != "-" {
		writer, err = os.OpenFile(*outputFile, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
		if err != nil {
			log.Fatal(err)
		}
		if *generateTS {
			tsFileName := regexp.MustCompile(".js$").ReplaceAllString(*outputFile, ".d.ts")
			tsWriter, err = os.OpenFile(tsFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
			if err != nil {
				log.Fatal(err)
			}
		}
	} else {
		if *generateTS {
			log.Fatal("You must specific outputFile when generating TypeScript definitions!")
		}
	}

	options := generator.Options{
		HasBigInt: *hasBigInt,
	}
	err = generator.Generate(options, schema, writer, tsWriter)
	if err != nil {
		log.Fatal(err)
	}
}
