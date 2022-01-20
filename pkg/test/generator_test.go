package main

import (
	"encoding/json"
	"io"
	"io/ioutil"
	"log"
	"os"
	"testing"

	"github.com/xxuejie/moleculec-es/pkg/generator"
)

func TestGenerator(t *testing.T) {
	content, err := ioutil.ReadFile("./schema.json")
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
	tsFileName := "./generated/test_generated.d.ts"
	tsWriter, err = os.OpenFile(tsFileName, os.O_CREATE|os.O_WRONLY|os.O_TRUNC, 0644)
	if err != nil {
		log.Fatal(err)
	}

	options := generator.Options{
		HasBigInt: false,
	}
	err = generator.Generate(options, schema, writer, tsWriter)
	if err != nil {
		log.Fatal(err)
	}

	got, err := ioutil.ReadFile("generated/test_generated.d.ts")
	if err != nil {
		log.Fatal(err)
	}

	want, err := ioutil.ReadFile("blockchain.d.ts")
	if err != nil {
		log.Fatal(err)
	}

	if !Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}

}

func Equal(a, b []byte) bool {
	if len(a) != len(b) {
		return false
	}
	for i, v := range a {
		if v != b[i] {
			return false
		}
	}
	return true
}
