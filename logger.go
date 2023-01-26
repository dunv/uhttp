package uhttp

import "log"

type Logger interface {
	Infof(template string, args ...interface{})
	Errorf(template string, args ...interface{})
}

type discardLogger struct{}

func (discardLogger) Infof(string, ...interface{})  {}
func (discardLogger) Errorf(string, ...interface{}) {}

func NewDiscardLogger() Logger {
	return discardLogger{}
}

type defaultLogger struct{}

func (defaultLogger) Infof(template string, args ...interface{}) {
	log.Printf(template, args...)
}
func (defaultLogger) Errorf(template string, args ...interface{}) {
	log.Printf(template, args...)
}

func NewDefaultLogger() Logger {
	return defaultLogger{}
}
