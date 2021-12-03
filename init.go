package deferlog

import (
	"bytes"
	"fmt"
	"os"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var StructLogger = zerolog.New(os.Stdout).
	With().Timestamp().CallerWithSkipFrameCount(3).Stack().
	Logger()

var ConsoleLogger = zerolog.New(
	zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.StampMilli,
		FormatCaller: func(i interface{}) string {
			if caller, ok := i.(string); ok {
				return ShortCaller(caller)
			}
			return ""
		},
	}).
	With().Timestamp().CallerWithSkipFrameCount(3).Stack().
	Logger()

func init() {
	pool := make(chan *fmtState, 1)
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		if formater, ok := err.(interface {
			Format(fmt.State, rune)
		}); ok {
			var state *fmtState
			select {
			case state = <-pool:
				state.Reset()
			default:
				state = new(fmtState)
			}

			formater.Format(state, 'v')
			_, _ = state.ReadString(os.PathSeparator)
			caller, _ := state.ReadString('\n')

			select {
			case pool <- state:
			default:
			}
			return ShortCaller(caller)
		}
		return nil
	}

	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		Logger = StructLogger
	} else {
		Logger = ConsoleLogger
	}

	switch os.Getenv("LOG_LEVEL") {
	case "TRACE", "Trace", "trace":
		Logger.Level(zerolog.TraceLevel)
	case "DEBUG", "Debug", "debug":
		Logger.Level(zerolog.DebugLevel)
	case "INFO", "Info", "info":
		Logger.Level(zerolog.InfoLevel)
	case "WARN", "Warn", "warn":
		Logger.Level(zerolog.WarnLevel)
	case "ERROR", "Error", "error":
		Logger.Level(zerolog.ErrorLevel)
	default:
		Logger.Level(zerolog.InfoLevel)
	}
}

type fmtState struct {
	bytes.Buffer
}

func (*fmtState) Width() (wid int, ok bool)      { return 0, false }
func (*fmtState) Precision() (prec int, ok bool) { return 0, false }
func (*fmtState) Flag(c int) bool                { return c == '+' }

func ShortCaller(line string) string {
	containSep := false
	for i := len(line) - 1; i >= 0; i-- {
		if line[i] == os.PathSeparator {
			if containSep {
				return strings.TrimSpace(line[i+1:])
			}
			containSep = true
		}
	}
	return strings.TrimSpace(line)
}
