package uhttp

// CustomLogger <-
type CustomLogger struct {
	Infof  func(string, ...interface{})
	Errorf func(string, ...interface{})
}

// NewCustomLogger <-
func NewCustomLogger(infoFn func(string, ...interface{}), errorFn func(string, ...interface{})) *CustomLogger {
	return &CustomLogger{
		Infof:  infoFn,
		Errorf: errorFn,
	}
}
