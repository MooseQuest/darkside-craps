// Command darkside-craps is a single-binary web app: it embeds the client-side
// craps SPA and serves a passkey-authenticated API for saving game history.
package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"strings"
	"time"
)

//go:embed web/index.html web/app.js web/engine.js web/styles.css web/privacy.html web/terms.html
var webFS embed.FS

// Config is resolved from the environment (Heroku config vars / rvault).
type Config struct {
	Port          string
	MongoURI      string
	SessionSecret []byte
	RPID          string // WebAuthn Relying Party ID (the domain), e.g. craps.moosequest.app
	RPOrigin      string // full origin, e.g. https://craps.moosequest.app
	DBName        string

	// SMTP for signup email-verification. When SMTPHost is empty, verification
	// is disabled and registration proceeds without an emailed code (dev).
	SMTPHost string
	SMTPPort string
	SMTPUser string
	SMTPPass string
	MailFrom string
}

// verifyEmail reports whether signup email-verification is enabled.
func (c Config) verifyEmail() bool { return c.SMTPHost != "" }

func loadConfig() Config {
	port := getenv("PORT", "8080")
	return Config{
		Port:          port,
		MongoURI:      os.Getenv("MONGO_URI"),
		SessionSecret: []byte(os.Getenv("SECRET_KEY")),
		RPID:          getenv("RP_ID", "localhost"),
		RPOrigin:      getenv("RP_ORIGIN", "http://localhost:"+port),
		DBName:        getenv("MONGO_DB", "craps_game"),
		SMTPHost:      os.Getenv("SMTP_HOST"),
		SMTPPort:      getenv("SMTP_PORT", "587"),
		SMTPUser:      os.Getenv("SMTP_USERNAME"),
		SMTPPass:      os.Getenv("SMTP_PASSWORD"),
		MailFrom:      getenv("MAIL_FROM", "Dark Side Craps <no-reply@moosequest.app>"),
	}
}

func getenv(k, def string) string {
	if v := os.Getenv(k); v != "" {
		return v
	}
	return def
}

func main() {
	cfg := loadConfig()
	if cfg.MongoURI == "" {
		log.Fatal("MONGO_URI is not set")
	}
	// Fail closed: a real (https) deployment must set SECRET_KEY, otherwise
	// session cookies would be signed with a publicly-known key and forgeable.
	if len(cfg.SessionSecret) == 0 {
		if strings.HasPrefix(cfg.RPOrigin, "https://") {
			log.Fatal("SECRET_KEY is required in production (RP_ORIGIN is https)")
		}
		log.Println("WARNING: SECRET_KEY unset — using an insecure dev key (localhost only)")
		cfg.SessionSecret = []byte("dev-only-not-secret-change-me")
	}

	ctx, cancel := context.WithTimeout(context.Background(), 15*time.Second)
	defer cancel()
	store, err := NewStore(ctx, cfg)
	if err != nil {
		log.Fatalf("mongo: %v", err)
	}
	log.Printf("connected to mongo (db=%s)", cfg.DBName)

	app, err := NewApp(cfg, store)
	if err != nil {
		log.Fatalf("app init: %v", err)
	}

	mux := http.NewServeMux()
	app.routes(mux)

	// 60 auth/api requests per minute per client IP.
	rl := newRateLimiter(60, time.Minute)
	handler := logRequests(securityHeaders(rateLimit(rl, limitBody(mux))))

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           handler,
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("listening on :%s (rp_id=%s origin=%s email_verification=%t)", cfg.Port, cfg.RPID, cfg.RPOrigin, cfg.verifyEmail())
	if err := srv.ListenAndServe(); err != nil {
		log.Fatal(err)
	}
}

func logRequests(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		start := time.Now()
		next.ServeHTTP(w, r)
		log.Printf("%s %s %s", r.Method, r.URL.Path, time.Since(start).Round(time.Millisecond))
	})
}
