# deferlog
deferlog provide a useful toolset for rs/zerolog.

## default config
1. use console logger while running in a terminal
2. use json logger while log into file (such as running in container)
3. read the environement `LOG_LEVEL` to change default log level


## Defer logger example
```go
func example() (err error) {
	start := time.Now()
	defer func() {
		deferlog.DebugWarn(err).
			Dur("took", time.Since(start)).
			Msg("run the example function")
	}()

    // do something here

	return nil
}
```