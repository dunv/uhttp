package params

import (
	"fmt"
	"strconv"
	"time"
)

func ParseEnum(value string, enum []string, key string, validatedMap map[string]interface{}, errors *[]error) {
	for _, enumValue := range enum {
		if enumValue == value {
			validatedMap[key] = value
			return
		}
	}
	*errors = append(*errors, fmt.Errorf("could not validate enum. needs to be one of %s", enum))
}

func ParseBool(value string, key string, validatedMap map[string]interface{}, errors *[]error) {
	if value == "true" {
		validatedMap[key] = true
		return
	} else if value == "false" {
		validatedMap[key] = false
		return
	}
	*errors = append(*errors, fmt.Errorf("could not validate bool. needs to be true or false, got %s", value))
}

func ParseInt(value string, key string, bits int, validatedMap map[string]interface{}, errors *[]error) {
	intValue, err := strconv.ParseInt(value, 10, 64)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("could not validate int%d. got %s", bits, value))
		return
	}

	switch bits {
	case 0:
		validatedMap[key] = int(intValue)
	case 32:
		validatedMap[key] = int32(intValue)
	case 64:
		validatedMap[key] = int64(intValue)
	}
}

func ParseFloat(value string, key string, bits int, validatedMap map[string]interface{}, errors *[]error) {
	floatValue, err := strconv.ParseFloat(value, 64)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("could not validate float%d. got %s", bits, value))
		return
	}

	switch bits {
	case 32:
		validatedMap[key] = float32(floatValue)
	case 64:
		validatedMap[key] = float64(floatValue)
	}
}

func ParseDate(value string, key string, format string, validatedMap map[string]interface{}, errors *[]error) {
	parsedDate, err := time.Parse(format, value)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("could not validate date (%s). got %s", format, value))
		return
	}
	validatedMap[key] = parsedDate
}
