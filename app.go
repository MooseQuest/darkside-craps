package main

import (
	"context"
	"crypto/hmac"
	"crypto/sha256"
	"encoding/base64"
	"encoding/json"
	"net/http"
	"strconv"
	"strings"
	"time"

	"github.com/go-webauthn/webauthn/webauthn"
)

const sessionCookie = "craps_session"
const ceremonyCookie = "craps_wa"
const sessionTTL = 7 * 24 * time.Hour
const maxBodyBytes = 1 << 20 // 1 MiB request-body cap

// App holds shared dependencies for the HTTP handlers.
type App struct {
	cfg   Config
	store *Store
	wa    *webauthn.WebAuthn
}

func NewApp(cfg Config, store *Store) (*App, error) {
	wa, err := webauthn.New(&webauthn.Config{
		RPID:          cfg.RPID,
		RPDisplayName: "Dark Side Craps",
		RPOrigins:     []string{cfg.RPOrigin},
	})
	if err != nil {
		return nil, err
	}
	return &App{cfg: cfg, store: store, wa: wa}, nil
}

func (a *App) routes(mux *http.ServeMux) {
	mux.HandleFunc("/healthz", func(w http.ResponseWriter, r *http.Request) {
		w.WriteHeader(http.StatusOK)
		_, _ = w.Write([]byte("ok"))
	})

	// Static SPA.
	mux.HandleFunc("/", a.serveAsset("web/index.html", "text/html; charset=utf-8"))
	mux.HandleFunc("/app.js", a.serveAsset("web/app.js", "text/javascript; charset=utf-8"))
	mux.HandleFunc("/engine.js", a.serveAsset("web/engine.js", "text/javascript; charset=utf-8"))
	mux.HandleFunc("/styles.css", a.serveAsset("web/styles.css", "text/css; charset=utf-8"))

	// Auth.
	mux.HandleFunc("/auth/me", a.handleMe)
	mux.HandleFunc("/auth/register/send-code", a.postOnly(a.handleRegisterSendCode))
	mux.HandleFunc("/auth/register/begin", a.postOnly(a.handleRegisterBegin))
	mux.HandleFunc("/auth/register/finish", a.postOnly(a.handleRegisterFinish))
	mux.HandleFunc("/auth/login/begin", a.postOnly(a.handleLoginBegin))
	mux.HandleFunc("/auth/login/finish", a.postOnly(a.handleLoginFinish))
	mux.HandleFunc("/auth/logout", a.postOnly(a.handleLogout))

	// Game history (auth required).
	mux.HandleFunc("/api/sessions", a.postOnly(a.handleSaveSession))
	mux.HandleFunc("/api/history", a.handleHistory)
	mux.HandleFunc("/api/summary", a.handleSummary)
}

func (a *App) serveAsset(path, ctype string) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		b, err := webFS.ReadFile(path)
		if err != nil {
			http.NotFound(w, r)
			return
		}
		w.Header().Set("Content-Type", ctype)
		_, _ = w.Write(b)
	}
}

func (a *App) postOnly(h http.HandlerFunc) http.HandlerFunc {
	return func(w http.ResponseWriter, r *http.Request) {
		if r.Method != http.MethodPost {
			writeErr(w, http.StatusMethodNotAllowed, "method not allowed")
			return
		}
		h(w, r)
	}
}

// ---- JSON helpers ----------------------------------------------------------

func writeJSON(w http.ResponseWriter, code int, v any) {
	w.Header().Set("Content-Type", "application/json")
	w.WriteHeader(code)
	_ = json.NewEncoder(w).Encode(v)
}

func writeErr(w http.ResponseWriter, code int, msg string) {
	writeJSON(w, code, map[string]string{"error": msg})
}

func readJSON(r *http.Request, v any) error {
	return json.NewDecoder(r.Body).Decode(v)
}

// ---- Signed cookies (stateless sessions) ----------------------------------

// sign returns base64(payload).base64(hmac) for a value.
func (a *App) sign(payload string) string {
	mac := hmac.New(sha256.New, a.cfg.SessionSecret)
	mac.Write([]byte(payload))
	sig := mac.Sum(nil)
	return b64(payload) + "." + base64.RawURLEncoding.EncodeToString(sig)
}

func (a *App) verify(token string) (string, bool) {
	parts := strings.SplitN(token, ".", 2)
	if len(parts) != 2 {
		return "", false
	}
	payloadBytes, err := base64.RawURLEncoding.DecodeString(parts[0])
	if err != nil {
		return "", false
	}
	payload := string(payloadBytes)
	expected := a.sign(payload)
	if !hmac.Equal([]byte(expected), []byte(token)) {
		return "", false
	}
	return payload, true
}

func b64(s string) string { return base64.RawURLEncoding.EncodeToString([]byte(s)) }

// setSession issues a signed session cookie carrying email|expiryUnix.
func (a *App) setSession(w http.ResponseWriter, email string) {
	exp := time.Now().Add(sessionTTL)
	payload := email + "|" + strconv.FormatInt(exp.Unix(), 10)
	http.SetCookie(w, &http.Cookie{
		Name:     sessionCookie,
		Value:    a.sign(payload),
		Path:     "/",
		Expires:  exp,
		HttpOnly: true,
		Secure:   a.secureCookies(),
		SameSite: http.SameSiteLaxMode,
	})
}

func (a *App) clearCookie(w http.ResponseWriter, name string) {
	http.SetCookie(w, &http.Cookie{
		Name: name, Value: "", Path: "/", MaxAge: -1,
		HttpOnly: true, Secure: a.secureCookies(), SameSite: http.SameSiteLaxMode,
	})
}

func (a *App) secureCookies() bool { return strings.HasPrefix(a.cfg.RPOrigin, "https://") }

// currentUser returns the authenticated email, or "" if none.
func (a *App) currentEmail(r *http.Request) string {
	c, err := r.Cookie(sessionCookie)
	if err != nil {
		return ""
	}
	payload, ok := a.verify(c.Value)
	if !ok {
		return ""
	}
	parts := strings.SplitN(payload, "|", 2)
	if len(parts) != 2 {
		return ""
	}
	exp, err := strconv.ParseInt(parts[1], 10, 64)
	if err != nil || time.Now().Unix() > exp {
		return ""
	}
	return parts[0]
}

// ceremony cookie carries the signed email between begin/finish.
func (a *App) setCeremony(w http.ResponseWriter, email string) {
	http.SetCookie(w, &http.Cookie{
		Name:     ceremonyCookie,
		Value:    a.sign(email),
		Path:     "/",
		Expires:  time.Now().Add(5 * time.Minute),
		HttpOnly: true,
		Secure:   a.secureCookies(),
		SameSite: http.SameSiteLaxMode,
	})
}

func (a *App) ceremonyEmail(r *http.Request) (string, bool) {
	c, err := r.Cookie(ceremonyCookie)
	if err != nil {
		return "", false
	}
	return a.verify(c.Value)
}

// ctx is a short helper for request-scoped Mongo calls.
func reqCtx(r *http.Request) (context.Context, context.CancelFunc) {
	return context.WithTimeout(r.Context(), 12*time.Second)
}
