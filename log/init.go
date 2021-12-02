package log

import "github.com/wweir/deferlog"

func init() {
	Logger = deferlog.Logger.With().Caller().Logger()
}
