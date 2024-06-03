package validation

import (
	"io/ioutil"

	"github.com/invopop/jsonschema"
)

func createJSON(v any) []byte {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	s := r.Reflect(v)
	jsonData, err := s.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return jsonData
}

func loadSchemaJSONFromFile(filePath string) ([]byte, error) {
	schemaJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return schemaJSON, nil
}
