// Command darkside-craps is a single-binary web app: it embeds the client-side
// craps SPA and serves a passkey-authenticated API for saving game history.
package main

import (
	"context"
	"embed"
	"log"
	"net/http"
	"os"
	"time"
)

//go:embed web/index.html web/app.js web/engine.js web/styles.css
var webFS embed.FS

// Config is resolved from the environment (Heroku config vars / rvault).
type Config struct {
	Port          string
	MongoURI      string
	SessionSecret []byte
	RPID          string // WebAuthn Relying Party ID (the domain), e.g. craps.moosequest.app
	RPOrigin      string // full origin, e.g. https://craps.moosequest.app
	DBName        string
}

func loadConfig() Config {
	port := getenv("PORT", "8080")
	secret := os.Getenv("SECRET_KEY")
	if secret == "" {
		secret = "dev-only-not-secret-change-me"
	}
	return Config{
		Port:          port,
		MongoURI:      os.Getenv("MONGO_URI"),
		SessionSecret: []byte(secret),
		RPID:          getenv("RP_ID", "localhost"),
		RPOrigin:      getenv("RP_ORIGIN", "http://localhost:"+port),
		DBName:        getenv("MONGO_DB", "craps_game"),
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

	srv := &http.Server{
		Addr:              ":" + cfg.Port,
		Handler:           logRequests(mux),
		ReadHeaderTimeout: 10 * time.Second,
	}
	log.Printf("listening on :%s (rp_id=%s origin=%s)", cfg.Port, cfg.RPID, cfg.RPOrigin)
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
