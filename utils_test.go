package deferlog

import (
	"encoding/json"
	"testing"
	"time"
)

func TestPasswordPlain(t *testing.T) {
	p := ParsePassword("plain-password")
	if got := p.Value(); got != "plain-password" {
		t.Fatalf("Value() = %q, want %q", got, "plain-password")
	}
	if got := p.Base64(); got != "cGxhaW4tcGFzc3dvcmQ=" {
		t.Fatalf("Base64() = %q, want %q", got, "cGxhaW4tcGFzc3dvcmQ=")
	}
}

func TestPasswordBase64Input(t *testing.T) {
	p := ParsePassword("c2VjcmV0")
	if got := p.Value(); got != "secret" {
		t.Fatalf("Value() = %q, want %q", got, "secret")
	}
	if got := p.Base64(); got != "c2VjcmV0" {
		t.Fatalf("Base64() = %q, want %q", got, "c2VjcmV0")
	}
}

func TestPasswordInvalidBase64FallsBackToPlain(t *testing.T) {
	p := ParsePassword("not-base64!")
	if got := p.Value(); got != "not-base64!" {
		t.Fatalf("Value() = %q, want %q", got, "not-base64!")
	}
}

func TestPasswordUnmarshalText(t *testing.T) {
	var p Password
	if err := p.UnmarshalText([]byte("c2VjcmV0")); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}
	if got := p.Value(); got != "secret" {
		t.Fatalf("Value() = %q, want %q", got, "secret")
	}
}

func TestPasswordMarshalTextUsesBase64(t *testing.T) {
	p := NewPassword("secret")
	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}
	if got := string(text); got != "c2VjcmV0" {
		t.Fatalf("MarshalText() = %q, want %q", got, "c2VjcmV0")
	}
}

func TestPasswordJSONRoundTrip(t *testing.T) {
	type config struct {
		Password Password `json:"password"`
	}

	encoded, err := json.Marshal(config{Password: NewPassword("secret")})
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if got := string(encoded); got != `{"password":"c2VjcmV0"}` {
		t.Fatalf("json.Marshal() = %s, want %s", got, `{"password":"c2VjcmV0"}`)
	}

	var decoded config
	if err := json.Unmarshal([]byte(`{"password":"plain-password"}`), &decoded); err != nil {
		t.Fatalf("json.Unmarshal() plain error = %v", err)
	}
	if got := decoded.Password.Value(); got != "plain-password" {
		t.Fatalf("decoded plain password = %q, want %q", got, "plain-password")
	}

	if err := json.Unmarshal([]byte(`{"password":"c2VjcmV0"}`), &decoded); err != nil {
		t.Fatalf("json.Unmarshal() base64 error = %v", err)
	}
	if got := decoded.Password.Value(); got != "secret" {
		t.Fatalf("decoded base64 password = %q, want %q", got, "secret")
	}
}

func TestPasswordRedactedFmtOutput(t *testing.T) {
	p := NewPassword("secret")
	if got := p.String(); got != "***" {
		t.Fatalf("String() = %q, want %q", got, "***")
	}
	if got := p.GoString(); got != "***" {
		t.Fatalf("GoString() = %q, want %q", got, "***")
	}
}

func TestThrottleCountEmissionRefreshesIntervalWindow(t *testing.T) {
	throttle := NewThrottle(75*time.Millisecond, 2)

	if ok, suppressed := throttle.Allow(); !ok || suppressed != 0 {
		t.Fatalf("first Allow() = (%v, %d), want (true, 0)", ok, suppressed)
	}

	time.Sleep(60 * time.Millisecond)

	if ok, suppressed := throttle.Allow(); !ok || suppressed != 0 {
		t.Fatalf("second Allow() = (%v, %d), want (true, 0)", ok, suppressed)
	}

	time.Sleep(20 * time.Millisecond)

	if ok, suppressed := throttle.Allow(); ok || suppressed != 0 {
		t.Fatalf("third Allow() = (%v, %d), want (false, 0)", ok, suppressed)
	}
}
