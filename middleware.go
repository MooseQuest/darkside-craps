package main

import (
	"net"
	"net/http"
	"strings"
	"sync"
	"time"
)

// securityHeaders sets defensive response headers on every request. The CSP is
// strict: only same-origin scripts/styles/assets, no inline script, no framing.
func securityHeaders(next http.Handler) http.Handler {
	const csp = "default-src 'self'; script-src 'self'; style-src 'self'; " +
		"img-src 'self' data:; connect-src 'self'; font-src 'self'; object-src 'none'; " +
		"base-uri 'none'; form-action 'self'; frame-ancestors 'none'"
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		h := w.Header()
		h.Set("Content-Security-Policy", csp)
		h.Set("X-Content-Type-Options", "nosniff")
		h.Set("X-Frame-Options", "DENY")
		h.Set("Referrer-Policy", "no-referrer")
		h.Set("Cross-Origin-Opener-Policy", "same-origin")
		next.ServeHTTP(w, r)
	})
}

// limitBody caps request bodies to guard against memory-exhaustion.
func limitBody(next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		if r.Body != nil {
			r.Body = http.MaxBytesReader(w, r.Body, maxBodyBytes)
		}
		next.ServeHTTP(w, r)
	})
}

// ---- per-IP fixed-window rate limiting ------------------------------------

type rateLimiter struct {
	mu     sync.Mutex
	hits   map[string]*rlWindow
	limit  int
	window time.Duration
}

type rlWindow struct {
	count int
	reset time.Time
}

func newRateLimiter(limit int, window time.Duration) *rateLimiter {
	return &rateLimiter{hits: make(map[string]*rlWindow), limit: limit, window: window}
}

func (rl *rateLimiter) allow(ip string, now time.Time) bool {
	rl.mu.Lock()
	defer rl.mu.Unlock()
	w := rl.hits[ip]
	if w == nil || now.After(w.reset) {
		if len(rl.hits) > 10000 { // opportunistic cleanup of expired entries
			for k, v := range rl.hits {
				if now.After(v.reset) {
					delete(rl.hits, k)
				}
			}
		}
		rl.hits[ip] = &rlWindow{count: 1, reset: now.Add(rl.window)}
		return true
	}
	if w.count >= rl.limit {
		return false
	}
	w.count++
	return true
}

// rateLimit throttles the auth + api surfaces per client IP.
func rateLimit(rl *rateLimiter, next http.Handler) http.Handler {
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		p := r.URL.Path
		if strings.HasPrefix(p, "/auth/") || strings.HasPrefix(p, "/api/") {
			if !rl.allow(clientIP(r), time.Now()) {
				w.Header().Set("Retry-After", "60")
				writeErr(w, http.StatusTooManyRequests, "too many requests — slow down")
				return
			}
		}
		next.ServeHTTP(w, r)
	})
}

// clientIP returns the real caller IP. On Heroku the router appends the true
// client IP as the LAST entry of X-Forwarded-For; earlier entries are
// client-supplied and spoofable, so we deliberately take the last one. (DNS is
// grey-cloud, so CF-Connecting-IP is untrusted and ignored.)
func clientIP(r *http.Request) string {
	if xff := r.Header.Get("X-Forwarded-For"); xff != "" {
		parts := strings.Split(xff, ",")
		return strings.TrimSpace(parts[len(parts)-1])
	}
	if host, _, err := net.SplitHostPort(r.RemoteAddr); err == nil {
		return host
	}
	return r.RemoteAddr
}
