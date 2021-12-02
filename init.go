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
	With().Timestamp().CallerWithSkipFrameCount(3).Stack().
	Logger().Level(zerolog.InfoLevel)

var ConsoleLogger = zerolog.New(zerolog.ConsoleWriter{
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
	Logger().Level(zerolog.InfoLevel)

func init() {
	zerolog.ErrorStackMarshaler = func(err error) interface{} {
		if formater, ok := err.(interface {
			Format(fmt.State, rune)
		}); ok {
			fmtState := &fmtState{}
			formater.Format(fmtState, 'v')

			_, _ = fmtState.ReadString(os.PathSeparator)
			line, _ := fmtState.ReadString('\n')
			return ShortCaller(line)
		}
		return nil
	}

	if fi, _ := os.Stdout.Stat(); (fi.Mode() & os.ModeCharDevice) == 0 {
		Logger = StructLogger
	} else {
		Logger = ConsoleLogger
	}
	if ok, _ := strconv.ParseBool(os.Getenv("DEBUG")); ok {
		Logger = Logger.Level(zerolog.DebugLevel)
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
