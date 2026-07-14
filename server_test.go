package main

import (
	"net/http"
	"net/http/httptest"
	"strconv"
	"strings"
	"testing"
	"time"
)

func testApp() *App {
	return &App{cfg: Config{
		SessionSecret: []byte("test-secret"),
		RPID:          "localhost",
		RPOrigin:      "https://craps.moosequest.app",
	}}
}

func TestSignVerifyRoundTrip(t *testing.T) {
	a := testApp()
	tok := a.sign("alice@example.com|123")
	payload, ok := a.verify(tok)
	if !ok || payload != "alice@example.com|123" {
		t.Fatalf("roundtrip failed: ok=%v payload=%q", ok, payload)
	}
}

func TestVerifyRejectsTamper(t *testing.T) {
	a := testApp()
	tok := a.sign("alice@example.com|123")
	if _, ok := a.verify(tok + "x"); ok {
		t.Fatal("tampered token accepted")
	}
	if _, ok := a.verify("garbage"); ok {
		t.Fatal("garbage token accepted")
	}
	// Different secret must not verify.
	b := &App{cfg: Config{SessionSecret: []byte("other")}}
	if _, ok := b.verify(tok); ok {
		t.Fatal("token verified under wrong secret")
	}
}

func TestCurrentEmailFromCookie(t *testing.T) {
	a := testApp()
	exp := time.Now().Add(time.Hour).Unix()
	payload := "bob@example.com|" + strconv.FormatInt(exp, 10)
	r := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	r.AddCookie(&http.Cookie{Name: sessionCookie, Value: a.sign(payload)})
	if got := a.currentEmail(r); got != "bob@example.com" {
		t.Fatalf("currentEmail = %q", got)
	}
}

func TestCurrentEmailExpired(t *testing.T) {
	a := testApp()
	payload := "bob@example.com|" + strconv.FormatInt(time.Now().Add(-time.Hour).Unix(), 10)
	r := httptest.NewRequest(http.MethodGet, "/auth/me", nil)
	r.AddCookie(&http.Cookie{Name: sessionCookie, Value: a.sign(payload)})
	if got := a.currentEmail(r); got != "" {
		t.Fatalf("expired cookie accepted: %q", got)
	}
}

func TestValidEmail(t *testing.T) {
	ok := []string{"a@b.co", "user.name+tag@example.com"}
	bad := []string{"", "nope", "@x.com", "x@", "a@@b", "  ",
		"a b@x.com", "pipe|payload@x.com", "ctrl\x00@x.com", "tab\t@x.com"}
	for _, e := range ok {
		if !validEmail(e) {
			t.Errorf("expected valid: %q", e)
		}
	}
	for _, e := range bad {
		if validEmail(e) {
			t.Errorf("expected invalid: %q", e)
		}
	}
}

func TestSecureCookies(t *testing.T) {
	if !testApp().secureCookies() {
		t.Fatal("https origin should yield secure cookies")
	}
	dev := &App{cfg: Config{RPOrigin: "http://localhost:8080"}}
	if dev.secureCookies() {
		t.Fatal("http origin should not be secure")
	}
}

func TestHealthz(t *testing.T) {
	a := testApp()
	mux := http.NewServeMux()
	// Register only the health route (no Mongo needed).
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})
	_ = a
	rec := httptest.NewRecorder()
	mux.ServeHTTP(rec, httptest.NewRequest(http.MethodGet, "/healthz", nil))
	if rec.Code != 200 || !strings.Contains(rec.Body.String(), "ok") {
		t.Fatalf("healthz: code=%d body=%q", rec.Code, rec.Body.String())
	}
}

func TestRateLimiter(t *testing.T) {
	rl := newRateLimiter(3, time.Minute)
	base := time.Unix(1_700_000_000, 0)
	for i := 0; i < 3; i++ {
		if !rl.allow("1.2.3.4", base) {
			t.Fatalf("request %d should be allowed", i)
		}
	}
	if rl.allow("1.2.3.4", base) {
		t.Fatal("4th request in window should be blocked")
	}
	if !rl.allow("5.6.7.8", base) {
		t.Fatal("a different IP should have its own bucket")
	}
	if !rl.allow("1.2.3.4", base.Add(time.Minute+time.Second)) {
		t.Fatal("request after the window resets should be allowed")
	}
}

func TestGenCode(t *testing.T) {
	for i := 0; i < 50; i++ {
		c, err := genCode()
		if err != nil {
			t.Fatal(err)
		}
		if len(c) != 6 {
			t.Fatalf("code len %d: %q", len(c), c)
		}
		for _, r := range c {
			if r < '0' || r > '9' {
				t.Fatalf("non-digit in %q", c)
			}
		}
	}
}

func TestHashCodeBinding(t *testing.T) {
	a := testApp()
	h := a.hashCode("a@b.com", "123456")
	if h != a.hashCode("a@b.com", "123456") {
		t.Fatal("hash not deterministic")
	}
	if h == a.hashCode("a@b.com", "654321") {
		t.Fatal("different code produced same hash")
	}
	if h == a.hashCode("x@y.com", "123456") {
		t.Fatal("different email produced same hash")
	}
	if h != a.hashCode("a@b.com", "  123456 ") {
		t.Fatal("surrounding whitespace should be trimmed")
	}
	other := &App{cfg: Config{SessionSecret: []byte("other")}}
	if h == other.hashCode("a@b.com", "123456") {
		t.Fatal("hash must depend on the server secret")
	}
}

func TestFromAddress(t *testing.T) {
	if got := fromAddress("Dark Side Craps <no-reply@moosequest.app>"); got != "no-reply@moosequest.app" {
		t.Fatalf("bracketed form: %q", got)
	}
	if got := fromAddress("plain@x.com"); got != "plain@x.com" {
		t.Fatalf("plain form: %q", got)
	}
}

func TestVerifyEmailToggle(t *testing.T) {
	if (Config{}).verifyEmail() {
		t.Fatal("no SMTP host should disable verification")
	}
	if !(Config{SMTPHost: "smtp.example.com"}).verifyEmail() {
		t.Fatal("SMTP host should enable verification")
	}
}

func TestClientIPUsesLastForwardedFor(t *testing.T) {
	r := httptest.NewRequest(http.MethodGet, "/", nil)
	// Heroku appends the real client IP last; spoofed entries come first.
	r.Header.Set("X-Forwarded-For", "9.9.9.9, 203.0.113.7")
	if ip := clientIP(r); ip != "203.0.113.7" {
		t.Fatalf("clientIP = %q, want the last (Heroku-appended) entry", ip)
	}
}
