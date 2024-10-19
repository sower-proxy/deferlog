# deferlog

deferlog provide a useful toolset for log/slog.

## Defer logger example

```go
func example() (err error) {
	start := time.Now()
	defer func() {
		deferlog.DebugWarn(err, "run the example function", "took", time.Since(start))
	}()

    // do something

	return nil
}
```
