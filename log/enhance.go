package log

import (
	"github.com/rs/zerolog"
	"github.com/sower-proxy/deferlog"
)

func init() {
	Logger = deferlog.StdLogger
}

func InfoWarn(err error) *zerolog.Event {
	if err != nil {
		return Logger.Warn().Err(err)
	}
	return Logger.Info()
}
func InfoFatal(err error) *zerolog.Event {
	if err != nil {
		return Logger.Fatal().Err(err)
	}
	return Logger.Info()
}

func DebugWarn(err error) *zerolog.Event {
	if err != nil {
		return Logger.Warn().Err(err)
	}
	return Logger.Debug()
}
func DebugFatal(err error) *zerolog.Event {
	if err != nil {
		return Logger.Fatal().Err(err)
	}
	return Logger.Debug()
}
