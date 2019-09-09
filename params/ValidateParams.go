package params

import (
	"fmt"
	"time"

	"github.com/dunv/uhelpers"
)

func ValidateParams(requirement R, actual map[string]string, required bool) (R, error) {
	errors := []error{}
	keys := uhelpers.StringKeysFromMap(requirement)
	validatedMap := map[string]interface{}{}
	for _, key := range keys {
		var actualValue string
		var ok bool
		if actualValue, ok = actual[key]; !ok {
			errors = append(errors, fmt.Errorf("required param %s not present", key))
			continue
		}

		switch requirement[key].(type) {
		case string:
			switch requirement[key] {
			case STRING:
				validatedMap[key] = actualValue
			case BOOL:
				ParseBool(actual[key], key, validatedMap, &errors)
			case INT:
				ParseInt(actual[key], key, 0, validatedMap, &errors)
			case INT32:
				ParseInt(actual[key], key, 32, validatedMap, &errors)
			case INT64:
				ParseInt(actual[key], key, 64, validatedMap, &errors)
			case FLOAT32:
				ParseFloat(actual[key], key, 32, validatedMap, &errors)
			case FLOAT64:
				ParseFloat(actual[key], key, 64, validatedMap, &errors)
			case SHORT_DATE:
				ParseDate(actual[key], key, "2006-01-02", validatedMap, &errors)
			case RFC3339_DATE:
				ParseDate(actual[key], key, time.RFC3339, validatedMap, &errors)
			default:
				return nil, fmt.Errorf("unknown param requirement")
			}
		case []string:
			ParseEnum(actual[key], requirement[key].([]string), key, validatedMap, &errors)
		default:
			// fmt.Printf("don't know %+v \n", wished_required["userId"])
		}

	}

	if required && len(errors) != 0 {
		return nil, fmt.Errorf("could not parse %v", errors)
	}

	return validatedMap, nil
}
