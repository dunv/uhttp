package uhttp

import (
	"fmt"
	"strconv"
	"time"
)

func parseString(value string, key string, validatedMap R, errors *[]error) {
	if value != "" {
		validatedMap[key] = value
		return
	}
	*errors = append(*errors, fmt.Errorf("could not validate string. needs to be not empty"))
}

func parseEnum(value string, enum []string, key string, validatedMap R, errors *[]error) {
	for _, enumValue := range enum {
		if enumValue == value {
			validatedMap[key] = value
			return
		}
	}
	*errors = append(*errors, fmt.Errorf("could not validate enum. needs to be one of %s", enum))
}

func parseBool(value string, key string, validatedMap R, errors *[]error) {
	if value == "true" {
		validatedMap[key] = true
		return
	} else if value == "false" {
		validatedMap[key] = false
		return
	}
	*errors = append(*errors, fmt.Errorf("could not validate bool. needs to be true or false, got %s", value))
}

func parseInt(value string, key string, bits int, validatedMap R, errors *[]error) {
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

func parseFloat(value string, key string, bits int, validatedMap R, errors *[]error) {
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

func parseDate(value string, key string, format string, validatedMap R, errors *[]error) {
	parsedDate, err := time.Parse(format, value)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("could not validate date (%s). got %s", format, value))
		return
	}
	validatedMap[key] = parsedDate
}

func parseDuration(value string, key string, validatedMap R, errors *[]error) {
	parsedDuration, err := time.ParseDuration(value)
	if err != nil {
		*errors = append(*errors, fmt.Errorf("could not validate duration. got %s", value))
		return
	}
	validatedMap[key] = parsedDuration
}
