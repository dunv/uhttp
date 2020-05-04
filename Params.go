package uhttp

import (
	"fmt"
	"strconv"
	"time"

	"github.com/dunv/uhelpers"
)

const (
	STRING       string = "string"
	BOOL         string = "bool"
	INT          string = "int"
	INT32        string = "int32"
	INT64        string = "int64"
	FLOAT32      string = "float32"
	FLOAT64      string = "float64"
	SHORT_DATE   string = "shortDate" // 2006-01-02
	RFC3339_DATE string = "rfc3339Date"
)

func ENUM(values ...string) []string {
	return values
}

type R map[string]interface{}

func (r R) Printable() (map[string]string, error) {
	printable := map[string]string{}
	for key, value := range r {
		switch typed := value.(type) {
		case string:
			printable[key] = typed
		case bool:
			printable[key] = strconv.FormatBool(typed)
		case int:
			printable[key] = strconv.FormatInt(int64(typed), 10)
		case int32:
			printable[key] = strconv.FormatInt(int64(typed), 10)
		case int64:
			printable[key] = strconv.FormatInt(typed, 10)
		case float32:
			printable[key] = strconv.FormatFloat(float64(typed), 'f', 2, 32)
		case float64:
			printable[key] = strconv.FormatFloat(typed, 'f', 2, 64)
		case time.Time:
			printable[key] = typed.Format(time.RFC3339)
		default:
			return nil, fmt.Errorf("could not print type %T", typed)
		}
	}
	return printable, nil
}

func (u *UHTTP) validateParams(requirement R, actual map[string]string, destination R, required bool) error {
	errors := []error{}
	keys := uhelpers.StringKeysFromMap(requirement)
	for _, key := range keys {
		// Publish an error only in the logs, if a param does already exist in the destination map
		// it obviously points to a bug in the code not an error on the user's side
		if _, ok := destination[key]; ok {
			u.opts.log.Errorf("key %s already present when parsing more params, check the requirements in the handler's definition", key)
		}

		// var actualValue string
		// var ok bool
		if _, ok := actual[key]; !ok && required {
			errors = append(errors, fmt.Errorf("required param %s (%s) not present", key, requirement[key]))
			continue
		}

		switch requirement[key].(type) {
		case string:
			switch requirement[key] {
			case STRING:
				parseString(actual[key], key, destination, &errors)
			case BOOL:
				parseBool(actual[key], key, destination, &errors)
			case INT:
				parseInt(actual[key], key, 0, destination, &errors)
			case INT32:
				parseInt(actual[key], key, 32, destination, &errors)
			case INT64:
				parseInt(actual[key], key, 64, destination, &errors)
			case FLOAT32:
				parseFloat(actual[key], key, 32, destination, &errors)
			case FLOAT64:
				parseFloat(actual[key], key, 64, destination, &errors)
			case SHORT_DATE:
				parseDate(actual[key], key, "2006-01-02", destination, &errors)
			case RFC3339_DATE:
				parseDate(actual[key], key, time.RFC3339, destination, &errors)
			default:
				return fmt.Errorf("unknown param requirement")
			}
		case []string:
			parseEnum(actual[key], requirement[key].([]string), key, destination, &errors)
		default:
			errors = append(errors, fmt.Errorf("don't know what to do with %+v \n", requirement[key]))
		}

	}

	if required && len(errors) != 0 {
		return fmt.Errorf("could not parse %v", errors)
	}

	return nil
}
