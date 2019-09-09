package params

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
