package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

// TODO: migrate to GRPC-Style config

var errorLevel = ulog.LEVEL_ERROR
var config Config = Config{
	CORS:                    uhelpers.PtrToString("*"),
	CustomLog:               ulog.NewUlog(),
	GzipCompressionLevel:    uhelpers.PtrToInt(4),
	EncodingErrorLogLevel:   &errorLevel,
	ParseModelErrorLogLevel: &errorLevel,
}

type Config struct {
	CORS                    *string
	CustomLog               ulog.ULogger
	GzipCompressionLevel    *int
	EncodingErrorLogLevel   *ulog.LogLevel
	ParseModelErrorLogLevel *ulog.LogLevel
}

func GetConfig() Config {
	return config
}

// SetConfig set config for all handlers
func SetConfig(_config Config) {
	if _config.CustomLog != nil {
		config.CustomLog = _config.CustomLog
		Logger = _config.CustomLog
	}

	if _config.CORS != nil {
		config.CORS = _config.CORS
	}

	if _config.GzipCompressionLevel != nil {
		config.GzipCompressionLevel = _config.GzipCompressionLevel
	}

	if _config.EncodingErrorLogLevel != nil {
		config.EncodingErrorLogLevel = _config.EncodingErrorLogLevel
	}

	if _config.ParseModelErrorLogLevel != nil {
		config.ParseModelErrorLogLevel = _config.ParseModelErrorLogLevel
	}
}
