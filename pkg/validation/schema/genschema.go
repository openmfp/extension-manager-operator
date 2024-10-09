package main

import (
	"encoding/json"
	"fmt"
	"os"

	"github.com/invopop/jsonschema"
	"github.com/openmfp/extension-content-operator/pkg/validation"
)

func ReflectContentConfiguration() {
	s := jsonschema.Reflect(&validation.ContentConfiguration{})
	data, err := json.MarshalIndent(s, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	file, err := os.Create("schema_autogen.json")
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err.Error())
	}

	fmt.Println("Schema saved to schema.json")

}

func main() {
	ReflectContentConfiguration()
}
