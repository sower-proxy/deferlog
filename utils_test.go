package deferlog

import (
	"encoding/json"
	"strings"
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

// TestPasswordBase64LikeValueStaysVerbatim pins the regression that started
// this redesign: "fnhkwfnh" is valid canonical base64, and a heuristic decode
// silently corrupted it into unrelated bytes. Plain input must never be
// reinterpreted.
func TestPasswordBase64LikeValueStaysVerbatim(t *testing.T) {
	for _, input := range []string{"fnhkwfnh", "c2VjcmV0", "aGVsbG8="} {
		p := ParsePassword(input)
		if got := p.Value(); got != input {
			t.Fatalf("ParsePassword(%q) = %q, want verbatim %q", input, got, input)
		}
	}
}

func TestPasswordExplicitBase64Input(t *testing.T) {
	p := ParsePassword("b64:c2VjcmV0")
	if got := p.Value(); got != "secret" {
		t.Fatalf("Value() = %q, want %q", got, "secret")
	}
}

func TestPasswordMalformedPrefixFallsBackToPlain(t *testing.T) {
	p := ParsePassword("b64:!!!")
	if got := p.Value(); got != "b64:!!!" {
		t.Fatalf("Value() = %q, want input kept verbatim", got)
	}
}

func TestPasswordUnmarshalText(t *testing.T) {
	var p Password
	if err := p.UnmarshalText([]byte("plain-password")); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}
	if got := p.Value(); got != "plain-password" {
		t.Fatalf("Value() = %q, want %q", got, "plain-password")
	}

	if err := p.UnmarshalText([]byte("b64:c2VjcmV0")); err != nil {
		t.Fatalf("UnmarshalText() error = %v", err)
	}
	if got := p.Value(); got != "secret" {
		t.Fatalf("Value() = %q, want %q", got, "secret")
	}
}

func TestPasswordUnmarshalTextRejectsMalformedBase64(t *testing.T) {
	var p Password
	if err := p.UnmarshalText([]byte("b64:!!!")); err == nil {
		t.Fatal("UnmarshalText() expected error for malformed b64 input")
	}
}

func TestPasswordMarshalTextUsesTaggedBase64(t *testing.T) {
	p := NewPassword("secret")
	text, err := p.MarshalText()
	if err != nil {
		t.Fatalf("MarshalText() error = %v", err)
	}
	if got := string(text); got != "b64:c2VjcmV0" {
		t.Fatalf("MarshalText() = %q, want %q", got, "b64:c2VjcmV0")
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
	if got := string(encoded); got != `{"password":"b64:c2VjcmV0"}` {
		t.Fatalf("json.Marshal() = %s, want %s", got, `{"password":"b64:c2VjcmV0"}`)
	}

	// Marshal -> Unmarshal must round-trip exactly.
	var decoded config
	if err := json.Unmarshal(encoded, &decoded); err != nil {
		t.Fatalf("json.Unmarshal() round trip error = %v", err)
	}
	if got := decoded.Password.Value(); got != "secret" {
		t.Fatalf("round-tripped password = %q, want %q", got, "secret")
	}

	if err := json.Unmarshal([]byte(`{"password":"plain-password"}`), &decoded); err != nil {
		t.Fatalf("json.Unmarshal() plain error = %v", err)
	}
	if got := decoded.Password.Value(); got != "plain-password" {
		t.Fatalf("decoded plain password = %q, want %q", got, "plain-password")
	}

	if err := json.Unmarshal([]byte(`{"password":"b64:c2VjcmV0"}`), &decoded); err != nil {
		t.Fatalf("json.Unmarshal() b64 error = %v", err)
	}
	if got := decoded.Password.Value(); got != "secret" {
		t.Fatalf("decoded b64 password = %q, want %q", got, "secret")
	}
}

func TestPasswordJSONMarshalOmitsPlaintext(t *testing.T) {
	encoded, err := json.Marshal(NewPassword("topsecret-admin"))
	if err != nil {
		t.Fatalf("json.Marshal() error = %v", err)
	}
	if strings.Contains(string(encoded), "topsecret-admin") {
		t.Fatalf("json output leaks plaintext: %s", encoded)
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
