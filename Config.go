package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/ulog"
)

var config Config = Config{
	CORS:                 uhelpers.PtrToString("*"),
	CustomLog:            ulog.NewUlog(),
	GzipCompressionLevel: uhelpers.PtrToInt(4),
}

type Config struct {
	CORS                 *string
	CustomLog            ulog.ULogger
	GzipCompressionLevel *int
}

func GetConfig() Config {
	return config
}

// SetConfig set config for all handlers
func SetConfig(_config Config) {
	config = _config

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
}
