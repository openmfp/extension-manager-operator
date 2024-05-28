package config_validator

import (
	"github.com/invopop/jsonschema"
)

func createJson(v any) []byte {
	var result []byte

	r := new(jsonschema.Reflector)
	r.ExpandedStruct = true

	s := r.Reflect(v)
	json, err := s.MarshalJSON()
	if err != nil {
		panic(err)
	}

	result = []byte(json)

	return result
}
