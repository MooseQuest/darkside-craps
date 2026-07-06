package main

import (
	"crypto/rand"
	"encoding/base64"
	"net/http"
	"strings"
)

// validEmail is intentionally strict: exactly one '@', no surrounding '@', and
// no whitespace, control characters, or '|' (the session-payload delimiter).
func validEmail(e string) bool {
	e = strings.TrimSpace(e)
	if len(e) < 3 || len(e) > 254 {
		return false
	}
	if strings.Count(e, "@") != 1 || strings.HasPrefix(e, "@") || strings.HasSuffix(e, "@") {
		return false
	}
	for _, r := range e {
		if r <= ' ' || r == '|' || r == 0x7f {
			return false
		}
	}
	return true
}

type emailReq struct {
	Email string `json:"email"`
}

func normEmail(s string) string { return strings.ToLower(strings.TrimSpace(s)) }

func (a *App) handleMe(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeJSON(w, http.StatusOK, map[string]any{"authenticated": false})
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"authenticated": true, "email": email})
}

func (a *App) handleRegisterBegin(w http.ResponseWriter, r *http.Request) {
	var req emailReq
	if err := readJSON(r, &req); err != nil || !validEmail(req.Email) {
		writeErr(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	email := normEmail(req.Email)
	ctx, cancel := reqCtx(r)
	defer cancel()

	// Do NOT create the user yet — only persist on a successful finish.
	existing, err := a.store.GetUser(ctx, email)
	var user *User
	var handle []byte
	switch {
	case err == nil:
		user = existing
		handle = existing.handle
	case err == errNotFound:
		h, herr := newHandle()
		if herr != nil {
			writeErr(w, http.StatusInternalServerError, "could not start registration")
			return
		}
		handle = h
		user = &User{email: email, handle: handle}
	default:
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}

	options, session, err := a.wa.BeginRegistration(user)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start registration")
		return
	}
	if err := a.store.SaveWASession(ctx, email, session, handle); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not persist ceremony")
		return
	}
	a.setCeremony(w, email)
	writeJSON(w, http.StatusOK, options)
}

func (a *App) handleRegisterFinish(w http.ResponseWriter, r *http.Request) {
	email, ok := a.ceremonyEmail(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "registration session expired — start over")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()

	session, handle, err := a.store.TakeWASession(ctx, email)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "registration session expired — start over")
		return
	}
	// Reconstruct the ceremony user (existing creds excluded during verify).
	user, gerr := a.store.GetUser(ctx, email)
	if gerr == errNotFound {
		user = &User{email: email, handle: handle}
	} else if gerr != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}

	cred, err := a.wa.FinishRegistration(user, *session, r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "passkey registration failed")
		return
	}
	// Persist the user (idempotent) and the new credential.
	if _, err := a.store.EnsureUser(ctx, email, handle); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not save account")
		return
	}
	if err := a.store.AddCredential(ctx, email, cred); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not save passkey")
		return
	}
	a.clearCookie(w, ceremonyCookie)
	a.setSession(w, email)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "email": email})
}

func (a *App) handleLoginBegin(w http.ResponseWriter, r *http.Request) {
	var req emailReq
	if err := readJSON(r, &req); err != nil || !validEmail(req.Email) {
		writeErr(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	email := normEmail(req.Email)
	ctx, cancel := reqCtx(r)
	defer cancel()

	user, err := a.store.GetUser(ctx, email)
	if err != nil && err != errNotFound {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	// Anti-enumeration: for an unknown email (or one with no passkeys) return a
	// well-formed decoy challenge with the same 200 shape. The ceremony simply
	// fails at finish (no stored session), so existence isn't revealed here.
	if err == errNotFound || len(user.WebAuthnCredentials()) == 0 {
		a.setCeremony(w, email)
		writeJSON(w, http.StatusOK, a.decoyAssertion())
		return
	}

	options, session, err := a.wa.BeginLogin(user)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start login")
		return
	}
	if err := a.store.SaveWASession(ctx, email, session, user.handle); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not persist ceremony")
		return
	}
	a.setCeremony(w, email)
	writeJSON(w, http.StatusOK, options)
}

// decoyAssertion returns a syntactically valid, useless assertion options object
// (random challenge, no allowed credentials) so login/begin looks identical for
// existing and non-existing accounts.
func (a *App) decoyAssertion() map[string]any {
	c := make([]byte, 32)
	_, _ = rand.Read(c)
	fakeID := make([]byte, 32)
	_, _ = rand.Read(fakeID)
	// Include one plausible (random) credential so the response shape matches a
	// real single-passkey account — an empty allowCredentials would itself reveal
	// that the email is unknown.
	return map[string]any{"publicKey": map[string]any{
		"challenge": base64.RawURLEncoding.EncodeToString(c),
		"timeout":   60000,
		"rpId":      a.cfg.RPID,
		"allowCredentials": []any{map[string]any{
			"type": "public-key",
			"id":   base64.RawURLEncoding.EncodeToString(fakeID),
		}},
		"userVerification": "preferred",
	}}
}

func (a *App) handleLoginFinish(w http.ResponseWriter, r *http.Request) {
	email, ok := a.ceremonyEmail(r)
	if !ok {
		writeErr(w, http.StatusBadRequest, "login session expired — start over")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()

	user, err := a.store.GetUser(ctx, email)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "no such account")
		return
	}
	session, _, err := a.store.TakeWASession(ctx, email)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "login session expired — start over")
		return
	}
	cred, err := a.wa.FinishLogin(user, *session, r)
	if err != nil {
		writeErr(w, http.StatusUnauthorized, "passkey verification failed")
		return
	}
	// Persist the updated sign counter (clone-detection hygiene).
	_ = a.store.UpdateCredential(ctx, email, cred)

	a.clearCookie(w, ceremonyCookie)
	a.setSession(w, email)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true, "email": email})
}

func (a *App) handleLogout(w http.ResponseWriter, r *http.Request) {
	a.clearCookie(w, sessionCookie)
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}
