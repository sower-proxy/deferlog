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

var StdLogger zerolog.Logger

func init() {
	zerolog.CallerMarshalFunc = func(pc uintptr, file string, line int) string {
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

	level := zerolog.InfoLevel
	switch os.Getenv("LOG_LEVEL") {
	case "TRACE", "Trace", "trace":
		level = zerolog.TraceLevel
	case "DEBUG", "Debug", "debug":
		level = zerolog.DebugLevel
	case "INFO", "Info", "info":
		level = zerolog.InfoLevel
	case "WARN", "Warn", "warn":
		level = zerolog.WarnLevel
	case "ERROR", "Error", "error":
		level = zerolog.ErrorLevel
	}

	writer := zerolog.ConsoleWriter{
		Out:        os.Stderr,
		TimeFormat: time.StampMilli,
	}
	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		writer.NoColor = true
		writer.Out = os.Stdout
	}

	StdLogger = zerolog.New(writer).Level(level).
		With().Timestamp().Stack().Caller().Logger()
	Logger = zerolog.New(writer).Level(level).
		With().Timestamp().Stack().CallerWithSkipFrameCount(3).Logger()
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
