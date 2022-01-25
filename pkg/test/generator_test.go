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
	RunGenerator("blockchain")
	got, want := GetResult("blockchain")
	if !Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}
func TestGeneratorEmbedded(t *testing.T) {
	RunGenerator("embedded")
	got, want := GetResult("embedded")
	if !Equal(got, want) {
		t.Errorf("got %v want %v", got, want)
	}
}

func GetResult(testName string) ([]byte, []byte) {
	got, err := ioutil.ReadFile("generated/" + testName + "_generated.d.ts")
	if err != nil {
		log.Fatal(err)
	}

	want, err := ioutil.ReadFile(testName + ".d.ts")
	if err != nil {
		log.Fatal(err)
	}

	return got, want
}

func RunGenerator(inputSchemaName string) {
	content, err := ioutil.ReadFile("./" + inputSchemaName + ".json")
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
	tsFileName := "./generated/" + inputSchemaName + "_generated.d.ts"
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
