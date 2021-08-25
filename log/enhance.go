package log

import (
	"github.com/rs/zerolog"
)

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
func DebugError(err error) *zerolog.Event {
	if err != nil {
		return Logger.Error().Err(err)
	}

	return Logger.Debug()
}
func DebugFatal(err error) *zerolog.Event {
	if err != nil {
		return Logger.Fatal().Err(err)
	}

	return Logger.Debug()
}
