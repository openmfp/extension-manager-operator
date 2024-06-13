package validation

import (
	"os"
)

func loadSchemaJSONFromFile(filePath string) ([]byte, error) {
	schemaJSON, err := os.ReadFile(filePath)
	if err != nil {
		return nil, err
	}

	return schemaJSON, nil
}
