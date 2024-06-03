package validation

import (
	"fmt"
	"os"
	"testing"

	"github.com/stretchr/testify/assert"
)

func Test_loadSchemaJSONFromFile_ValidFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "mock_schema.json.out")

	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Println("Failed to remove temp file:", err)
		}
	}()

	schemaJSONContent := getJSONSchemaFixture()
	if _, err := tmpFile.Write(schemaJSONContent); err != nil {
		t.Fatal(err)
	}

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	schemaJSON, err := loadSchemaJSONFromFile(tmpFile.Name())

	assert.NoError(t, err)
	assert.NotNil(t, schemaJSON)
	assert.Equal(t, schemaJSONContent, schemaJSON)
}

func Test_loadSchemaJSONFromFile_InvalidFile(t *testing.T) {
	_, err := loadSchemaJSONFromFile("invalid_file_path")

	assert.Error(t, err)
	assert.Contains(t, err.Error(), "no such file or directory")
	assert.True(t, os.IsNotExist(err))
}

func Test_loadSchemaJSONFromFile_EmptyFile(t *testing.T) {
	tmpFile, err := os.CreateTemp("", "mock_schema.json.out")
	if err != nil {
		t.Fatal(err)
	}
	defer func() {
		if err := os.Remove(tmpFile.Name()); err != nil {
			fmt.Println("Failed to remove temp file:", err)
		}
	}()

	if err := tmpFile.Close(); err != nil {
		t.Fatal(err)
	}

	schemaJSON, err := loadSchemaJSONFromFile(tmpFile.Name())

	assert.NoError(t, err)
	assert.Equal(t, "", string(schemaJSON))
}
