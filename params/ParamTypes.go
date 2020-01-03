package params

import (
	"fmt"
	"strconv"
	"time"
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
