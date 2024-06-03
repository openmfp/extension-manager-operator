package validation

import (
	//"errors"
	"io/ioutil"
)

func loadSchemaJSONFromFile(filePath string) ([]byte, error) {
	schemaJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return schemaJSON, nil
}
