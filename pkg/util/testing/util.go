package testing

import (
	"encoding/json"
	"fmt"

	"github.com/google/go-cmp/cmp"
)

// CompareJSON This function is used to compare two JSON strings in unit tests
func CompareJSON(json1, json2 string) (bool, error) { // coverage-ignore
	var obj1, obj2 map[string]interface{}

	err := json.Unmarshal([]byte(json1), &obj1)
	if err != nil {
		return false, err
	}

	err = json.Unmarshal([]byte(json2), &obj2)
	if err != nil {
		return false, err
	}

	equal := cmp.Equal(obj1, obj2)
	if !equal {
		diff := cmp.Diff(obj1, obj2)
		if diff != "" {
			fmt.Printf("Differences:\n%s", diff)
		}
	}
	return equal, nil
}
