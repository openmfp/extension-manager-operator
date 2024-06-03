package validation

import (
	//"errors"
	"io/ioutil"
)

//func createJSON(v interface{}) ([]byte, error) {
//	//if _, ok := v.(struct{}); !ok {
//	//	return nil, errors.New("input is not a struct")
//	//}
//
//	r := new(jsonschema.Reflector)
//	r.ExpandedStruct = true
//
//	s := r.Reflect(v)
//	jsonData, err := s.MarshalJSON()
//	if err != nil {
//		return nil, err
//	}
//
//	return jsonData, nil
//}

func loadSchemaJSONFromFile(filePath string) ([]byte, error) {
	schemaJSON, err := ioutil.ReadFile(filePath)
	if err != nil {
		return nil, err
	}
	return schemaJSON, nil
}
