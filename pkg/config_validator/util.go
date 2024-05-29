package config_validator

import (
	"github.com/invopop/jsonschema"
)

func createJson(v any) []byte {
	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	s := r.Reflect(v)
	jsonData, err := s.MarshalJSON()
	if err != nil {
		panic(err)
	}

	return jsonData
}
