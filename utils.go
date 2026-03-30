package deferlog

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"sync/atomic"
	"time"
)

// Secret wraps a string to prevent accidental exposure in logs and fmt output.
type Secret string

func (s Secret) Value() string                { return string(s) }
func (s Secret) String() string               { return "***" }
func (s Secret) GoString() string             { return "***" }
func (s Secret) MarshalText() ([]byte, error) { return []byte("***"), nil }

// Password stores a password in plain form internally, accepts plain text or
// base64-encoded input, and keeps fmt output redacted.
//
// Text/JSON marshaling intentionally emits the base64 form for transport or
// config storage, so structured loggers that honor MarshalText/MarshalJSON will
// serialize the encoded value instead of "***".
type Password string

func (p Password) Value() string                { return string(p) }
func (p Password) String() string               { return "***" }
func (p Password) GoString() string             { return "***" }
func (p Password) MarshalText() ([]byte, error) { return []byte(p.Base64()), nil }
func (p Password) MarshalBinary() ([]byte, error) {
	return p.MarshalText()
}
func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(p.Base64())
}
func (p Password) Base64() string       { return base64.StdEncoding.EncodeToString([]byte(p)) }
func NewPassword(value string) Password { return Password(value) }

// ParsePassword is a convenience helper for config-style inputs that may be
// plain text or canonical base64. Because it uses a heuristic, plain values
// that are themselves valid canonical base64 will be decoded; use NewPassword
// when the input must be preserved verbatim.
func ParsePassword(value string) Password { return Password(decodeBase64OrPlain(value)) }
func (p *Password) UnmarshalText(text []byte) error {
	*p = Password(decodeBase64OrPlain(string(text)))
	return nil
}
func (p *Password) UnmarshalBinary(data []byte) error {
	return p.UnmarshalText(data)
}
func (p *Password) UnmarshalJSON(data []byte) error {
	if string(data) == "null" {
		*p = ""
		return nil
	}

	var value string
	if err := json.Unmarshal(data, &value); err != nil {
		return fmt.Errorf("deferlog: password must be a JSON string or null: %w", err)
	}
	return p.UnmarshalText([]byte(value))
}

// decodeBase64OrPlain treats canonical base64 as encoded data and falls back to
// the original string otherwise. This is intended for operator-provided config
// values, not arbitrary user input.
func decodeBase64OrPlain(value string) string {
	decoded, err := base64.StdEncoding.DecodeString(value)
	if err != nil || base64.StdEncoding.EncodeToString(decoded) != value {
		return value
	}
	return string(decoded)
}

// Throttle gates log emissions by time interval and/or occurrence count.
// A zero-value Throttle allows every call. It is safe for concurrent use.
type Throttle struct {
	interval int64 // nanoseconds; 0 = disabled
	every    int64 // count interval; 0 = disabled

	lastTime   atomic.Int64 // UnixNano of last Allow()==true
	count      atomic.Int64 // total Allow() calls
	suppressed atomic.Int64 // calls suppressed since last emission
}

// NewThrottleWith creates a throttle that allows a call when either the time interval
// has elapsed or the count interval is reached (OR logic).
func NewThrottle(interval time.Duration, n int) *Throttle {
	return &Throttle{interval: int64(interval), every: int64(n)}
}

// Allow reports whether the current call should be emitted.
// When ok is true, suppressed returns the number of calls that were
// suppressed since the previous emission.
func (t *Throttle) Allow() (ok bool, suppressed int64) {
	if t.interval == 0 && t.every == 0 {
		return true, 0
	}

	cnt := t.count.Add(1)

	// Check time condition
	if t.interval > 0 {
		now := time.Now().UnixNano()
		last := t.lastTime.Load()
		if now-last >= t.interval {
			if t.lastTime.CompareAndSwap(last, now) {
				return true, t.suppressed.Swap(0)
			}
		}
	}

	// Check count condition
	if t.every > 0 && cnt%t.every == 0 {
		if t.interval > 0 {
			t.lastTime.Store(time.Now().UnixNano())
		}
		return true, t.suppressed.Swap(0)
	}

	t.suppressed.Add(1)
	return false, 0
}
