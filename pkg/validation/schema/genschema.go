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
	s := jsonschema.Reflect(&validation.ContentConfiguration{})
	data, err := json.MarshalIndent(s, "", "  ")
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
	file.Sync()
	file.Close()

	// // print data
	// fmt.Println(string(data))

	fmt.Println("Schema saved to schema_autogen.json")

}

func main() {
	reflectContentConfiguration()
}
