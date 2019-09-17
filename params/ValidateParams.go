package params

import (
	"fmt"
	"time"

	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/logging"
)

func ValidateParams(requirement R, actual map[string]string, destination R, required bool) error {
	errors := []error{}
	keys := uhelpers.StringKeysFromMap(requirement)
	for _, key := range keys {
		// Publish an error only in the logs, if a param does already exist in the destination map
		// it obviously points to a bug in the code not an error on the user's side
		if _, ok := destination[key]; ok {
			logging.Logger.Errorf("key %s already present when parsing more params, check the requirements in the handler's definition", key)
		}

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
				destination[key] = actualValue
			case BOOL:
				ParseBool(actual[key], key, destination, &errors)
			case INT:
				ParseInt(actual[key], key, 0, destination, &errors)
			case INT32:
				ParseInt(actual[key], key, 32, destination, &errors)
			case INT64:
				ParseInt(actual[key], key, 64, destination, &errors)
			case FLOAT32:
				ParseFloat(actual[key], key, 32, destination, &errors)
			case FLOAT64:
				ParseFloat(actual[key], key, 64, destination, &errors)
			case SHORT_DATE:
				ParseDate(actual[key], key, "2006-01-02", destination, &errors)
			case RFC3339_DATE:
				ParseDate(actual[key], key, time.RFC3339, destination, &errors)
			default:
				return fmt.Errorf("unknown param requirement")
			}
		case []string:
			ParseEnum(actual[key], requirement[key].([]string), key, destination, &errors)
		default:
			errors = append(errors, fmt.Errorf("don't know what to do with %+v \n", requirement[key]))
		}

	}

	if required && len(errors) != 0 {
		return fmt.Errorf("could not parse %v", errors)
	}

	return nil
}
