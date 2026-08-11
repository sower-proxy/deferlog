package deferlog

import (
	"encoding/base64"
	"encoding/json"
	"fmt"
	"strings"
	"sync/atomic"
	"time"
)

// Secret wraps a string to prevent accidental exposure in logs and fmt output.
type Secret string

func (s Secret) Value() string                { return string(s) }
func (s Secret) String() string               { return "***" }
func (s Secret) GoString() string             { return "***" }
func (s Secret) MarshalText() ([]byte, error) { return []byte("***"), nil }

// Password stores a password in plain form internally and keeps fmt output
// redacted.
//
// Input is always taken verbatim unless it carries the explicit "b64:"
// prefix, which decodes the following canonical base64. There is no guessing:
// a plain value that happens to be valid canonical base64 (e.g. "fnhkwfnh")
// stays byte-for-byte, so it can never be silently corrupted.
//
// Text/JSON marshaling emits the "b64:" prefixed form for transport or
// config storage, making marshal/unmarshal round-trips exact. Structured
// loggers that honor MarshalText/MarshalJSON serialize the encoded value
// instead of "***" (and never the plaintext).
type Password string

const passwordBase64Prefix = "b64:"

func (p Password) Value() string                { return string(p) }
func (p Password) String() string               { return "***" }
func (p Password) GoString() string             { return "***" }
func (p Password) MarshalText() ([]byte, error) { return []byte(passwordBase64Prefix + p.Base64()), nil }
func (p Password) MarshalBinary() ([]byte, error) {
	return p.MarshalText()
}
func (p Password) MarshalJSON() ([]byte, error) {
	return json.Marshal(passwordBase64Prefix + p.Base64())
}
func (p Password) Base64() string { return base64.StdEncoding.EncodeToString([]byte(p)) }

// NewPassword wraps a value verbatim.
func NewPassword(value string) Password { return Password(value) }

// ParsePassword interprets a config-style input: a "b64:"-prefixed value is
// decoded, anything else is kept verbatim. Malformed base64 after the prefix
// falls back to the input as written, so the value is never silently lost;
// strict validation is available through UnmarshalText, which returns an
// error instead.
func ParsePassword(value string) Password {
	if decoded, err := decodePasswordInput(value); err == nil {
		return decoded
	}
	return Password(value)
}

// decodePasswordInput decodes the explicit "b64:" form. Any other input is
// verbatim. A "b64:" prefix with malformed base64 is an error, never a
// silent fallback, so a typo surfaces instead of becoming a wrong password.
func decodePasswordInput(value string) (Password, error) {
	if !strings.HasPrefix(value, passwordBase64Prefix) {
		return Password(value), nil
	}
	decoded, err := base64.StdEncoding.DecodeString(strings.TrimPrefix(value, passwordBase64Prefix))
	if err != nil {
		return "", fmt.Errorf("deferlog: invalid b64 password: %w", err)
	}
	return Password(decoded), nil
}

func (p *Password) UnmarshalText(text []byte) error {
	decoded, err := decodePasswordInput(string(text))
	if err != nil {
		return err
	}
	*p = decoded
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
