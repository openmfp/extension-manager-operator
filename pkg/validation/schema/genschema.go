package main

import (
	"encoding/json"
	"fmt"
	"os"
	"path/filepath"
	"runtime"

	"github.com/invopop/jsonschema"
	"github.com/openmfp/extension-content-operator/pkg/validation"
)

func reflectContentConfiguration() {
	r := new(jsonschema.Reflector)
	// r.ExpandedStruct = false
	schemaCategory := r.Reflect(&validation.Category{})
	schemaRoot := r.Reflect(&validation.ContentConfiguration{})
	schemaRoot.Definitions["Category"] = schemaCategory.Definitions["Category"]
	data, err := json.MarshalIndent(schemaRoot, "", "  ")
	if err != nil {
		panic(err.Error())
	}

	// Get the directory of the current file
	_, currentFile, _, ok := runtime.Caller(0)
	if !ok {
		fmt.Println("Unable to get the current file path")
		return
	}
	currentDir := filepath.Dir(currentFile)
	targetPath := filepath.Join(currentDir, "schema_autogen.json")

	// Save the schema to a file
	file, err := os.Create(targetPath)
	if err != nil {
		panic(err.Error())
	}
	defer file.Close()

	_, err = file.Write(data)
	if err != nil {
		panic(err.Error())
	}
	err = file.Sync()
	if err != nil {
		panic(err.Error())
	}
	file.Close()

	// // print data
	// fmt.Println(string(data))

	fmt.Println("Schema saved to schema_autogen.json")

}

func main() {
	reflectContentConfiguration()
}
