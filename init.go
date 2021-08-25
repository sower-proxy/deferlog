package deferlog

import (
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
	"github.com/rs/zerolog/pkgerrors"
	logEnhance "github.com/wweir/deferlog/log"
)

var StructLogger = zerolog.New(os.Stdout).
	With().Timestamp().Logger()

var ConsoleLogger = zerolog.New(zerolog.ConsoleWriter{
	Out:        os.Stdout,
	TimeFormat: time.StampMilli,
	FormatCaller: func(i interface{}) string {
		caller := i.(string)
		if idx := strings.Index(caller, "/pkg/mod/"); idx > 0 {
			return caller[idx+9:]
		}
		if idx := strings.LastIndexByte(caller, '/'); idx > 0 {
			return caller[idx+1:]
		}
		return caller
	},
}).With().Timestamp().Logger()

func init() {
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		return pkgerrors.MarshalStack(err)
	}

	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		SetDefaultLogger(ConsoleLogger, 1, zerolog.DebugLevel)
	} else {
		SetDefaultLogger(ConsoleLogger, 1, zerolog.InfoLevel)
	}
}

func SetDefaultLogger(logger zerolog.Logger, deferSkip int, logLevel zerolog.Level) {
	logEnhance.Logger = logger.With().Caller().Logger().Level(logLevel)
	Logger = logger.With().CallerWithSkipFrameCount(deferSkip + 2).Logger().Level(logLevel)
}
