package deferlog

import (
	"sync/atomic"
	"time"
)

// Secret wraps a string to prevent accidental exposure in logs and fmt output.
type Secret string

func (s Secret) Value() string                { return string(s) }
func (s Secret) String() string               { return "***" }
func (s Secret) GoString() string             { return "***" }
func (s Secret) MarshalText() ([]byte, error) { return []byte("***"), nil }

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
		return true, t.suppressed.Swap(0)
	}

	t.suppressed.Add(1)
	return false, 0
}
