package uhttp

import (
	"github.com/dunv/uhelpers"
	"github.com/dunv/uhttp/logging"
	"github.com/dunv/ulog"
)

var config Config = Config{
	DisableCORS: uhelpers.PtrToBool(false),
	CustomLog:   ulog.NewUlog(),
}

type Config struct {
	DisableCORS *bool
	CustomLog   ulog.ULogger
}

func GetConfig() Config {
	return config
}

// SetConfig set config for all handlers
func SetConfig(_config Config) {
	config = _config

	if _config.CustomLog != nil {
		config.CustomLog = _config.CustomLog
		logging.Logger = _config.CustomLog
	}

	if _config.DisableCORS != nil {
		config.DisableCORS = _config.DisableCORS
	}
}
