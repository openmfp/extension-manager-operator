/*
Package contentconfiguration provides utilities for validating content configurations.

The package supports validating configurations in JSON and YAML formats.

Example:

    import (
        "github.com/username/golang-commons/contentconfiguration"
    )

    func main() {
        jsonConfig := ` + "`" + `{
            "name": "overview",
            "luigiConfigFragment": [
                {
                    "data": {
                        "nodes": [
                            {
                                "entityType": "global",
                                "pathSegment": "home",
                                "label": "Overview",
                                "icon": "home"
                            }
                        ]
                    }
                }
            ]
        }` + "`" + `
		cC := contentconfiguration.NewContentConfiguration()

        valid, errors := cC.ValidateJSON([]byte(jsonConfig))
        if !valid {
            fmt.Println("Validation failed:", errors)
        } else {
            fmt.Println("Validation succeeded")
        }
    }

The package also supports schema validation.
*/

package contentconfiguration
