package main

import (
	"net"
	"net/http"
	"strings"

	"github.com/oschwald/geoip2-golang"
)

// geoBlocker resolves a request's country (via an embedded DB-IP Lite database)
// and blocks a configured set of ISO-3166 alpha-2 country codes. When the set is
// empty it is inert.
type geoBlocker struct {
	db      *geoip2.Reader
	blocked map[string]bool
}

func newGeoBlocker(mmdb []byte, blockedCSV string) (*geoBlocker, error) {
	blocked := map[string]bool{}
	for _, c := range strings.Split(blockedCSV, ",") {
		if c = strings.ToUpper(strings.TrimSpace(c)); c != "" {
			blocked[c] = true
		}
	}
	g := &geoBlocker{blocked: blocked}
	if len(blocked) == 0 {
		return g, nil // disabled — don't even open the DB
	}
	db, err := geoip2.FromBytes(mmdb)
	if err != nil {
		return nil, err
	}
	g.db = db
	return g, nil
}

func (g *geoBlocker) enabled() bool { return len(g.blocked) > 0 && g.db != nil }

// countryOf returns the ISO country code for an IP, or "" if unknown/unparseable.
func (g *geoBlocker) countryOf(ip string) string {
	if g.db == nil {
		return ""
	}
	parsed := net.ParseIP(ip)
	if parsed == nil {
		return ""
	}
	rec, err := g.db.Country(parsed)
	if err != nil {
		return ""
	}
	return rec.Country.IsoCode
}

// geoBlock serves a 451 "unavailable for legal reasons" page to visitors in a
// blocked country. It fails open: unknown/unresolvable IPs (localhost, gaps in
// the DB) are allowed through. /styles.css and /healthz are always allowed so
// the notice can render and health checks keep working.
func geoBlock(g *geoBlocker, next http.Handler) http.Handler {
	if !g.enabled() {
		return next
	}
	return http.HandlerFunc(func(w http.ResponseWriter, r *http.Request) {
		switch r.URL.Path {
		case "/styles.css", "/healthz", "/privacy", "/terms":
			next.ServeHTTP(w, r)
			return
		}
		if country := g.countryOf(clientIP(r)); country != "" && g.blocked[country] {
			b, err := webFS.ReadFile("web/blocked.html")
			if err != nil {
				http.Error(w, "Not available in your region.", http.StatusUnavailableForLegalReasons)
				return
			}
			w.Header().Set("Content-Type", "text/html; charset=utf-8")
			w.WriteHeader(http.StatusUnavailableForLegalReasons) // 451
			_, _ = w.Write(b)
			return
		}
		next.ServeHTTP(w, r)
	})
}
