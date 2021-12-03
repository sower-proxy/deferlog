package deferlog

import (
	"bytes"
	"fmt"
	"os"
	"strconv"
	"strings"
	"time"

	"github.com/rs/zerolog"
)

var StructLogger = zerolog.New(os.Stdout).
	With().Timestamp().Stack().Logger()
var ConsoleLogger = zerolog.New(
	zerolog.ConsoleWriter{
		Out:        os.Stdout,
		TimeFormat: time.StampMilli,
	}).
	With().Timestamp().Stack().Logger()
var StdLogger zerolog.Logger

func init() {
	zerolog.CallerMarshalFunc = func(file string, line int) string {
		return ShortCaller(file) + ":" + strconv.Itoa(line)
	}

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
		Logger = Logger.Level(zerolog.TraceLevel)
	case "DEBUG", "Debug", "debug":
		Logger = Logger.Level(zerolog.DebugLevel)
	case "INFO", "Info", "info":
		Logger = Logger.Level(zerolog.InfoLevel)
	case "WARN", "Warn", "warn":
		Logger = Logger.Level(zerolog.WarnLevel)
	case "ERROR", "Error", "error":
		Logger = Logger.Level(zerolog.ErrorLevel)
	default:
		Logger = Logger.Level(zerolog.InfoLevel)
	}

	StdLogger = Logger.With().Caller().Logger()
	Logger = Logger.With().CallerWithSkipFrameCount(3).Logger()
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
