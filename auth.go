package main

import (
	"crypto/hmac"
	"crypto/rand"
	"crypto/sha256"
	"encoding/base64"
	"encoding/hex"
	"fmt"
	"math/big"
	"net/http"
	"strings"
	"time"
)

const verifyCodeTTL = 10 * time.Minute

// genCode returns a cryptographically-random 6-digit code.
func genCode() (string, error) {
	n, err := rand.Int(rand.Reader, big.NewInt(1_000_000))
	if err != nil {
		return "", err
	}
	return fmt.Sprintf("%06d", n.Int64()), nil
}

// hashCode binds the code to the email under the server secret so plaintext
// codes are never stored.
func (a *App) hashCode(email, code string) string {
	m := hmac.New(sha256.New, a.cfg.SessionSecret)
	m.Write([]byte(email + "|" + strings.TrimSpace(code)))
	return hex.EncodeToString(m.Sum(nil))
}

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

// handleRegisterSendCode emails a one-time verification code for signup. When
// email verification is disabled (no SMTP configured), it reports sent=false so
// the client proceeds straight to the passkey step.
func (a *App) handleRegisterSendCode(w http.ResponseWriter, r *http.Request) {
	var req emailReq
	if err := readJSON(r, &req); err != nil || !validEmail(req.Email) {
		writeErr(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	if !a.cfg.verifyEmail() {
		writeJSON(w, http.StatusOK, map[string]any{"sent": false})
		return
	}
	email := normEmail(req.Email)
	ctx, cancel := reqCtx(r)
	defer cancel()

	code, err := genCode()
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not generate code")
		return
	}
	if err := a.store.SaveVerification(ctx, email, a.hashCode(email, code), verifyCodeTTL); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start verification")
		return
	}
	if err := a.sendVerificationEmail(email, code); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not send the verification email")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"sent": true})
}

type registerBeginReq struct {
	Email string `json:"email"`
	Code  string `json:"code"`
}

func (a *App) handleRegisterBegin(w http.ResponseWriter, r *http.Request) {
	var req registerBeginReq
	if err := readJSON(r, &req); err != nil || !validEmail(req.Email) {
		writeErr(w, http.StatusBadRequest, "a valid email is required")
		return
	}
	email := normEmail(req.Email)
	ctx, cancel := reqCtx(r)
	defer cancel()

	// Require a valid, single-use emailed code before issuing a challenge, so a
	// passkey can only be registered against an email the caller controls.
	if a.cfg.verifyEmail() {
		ok, err := a.store.ConsumeVerification(ctx, email, a.hashCode(email, req.Code))
		if err != nil {
			writeErr(w, http.StatusInternalServerError, "could not verify the code")
			return
		}
		if !ok {
			writeErr(w, http.StatusBadRequest, "invalid or expired verification code")
			return
		}
	}

	// Signup is for NEW accounts only. If the email already has a passkey, refuse
	// and steer to sign-in — adding another passkey is an authenticated action
	// (/auth/passkey/add/*). This is reached only after a valid emailed code, so
	// it doesn't leak account existence to an attacker.
	_, err := a.store.GetUser(ctx, email)
	if err == nil {
		writeErr(w, http.StatusConflict, "This email already has a passkey — sign in instead.")
		return
	}
	if err != errNotFound {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	handle, herr := newHandle()
	if herr != nil {
		writeErr(w, http.StatusInternalServerError, "could not start registration")
		return
	}
	user := &User{email: email, handle: handle}

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

// handleAddPasskeyBegin starts adding an additional passkey to the currently
// signed-in account (multi-device). No email code — the session proves identity,
// and BeginRegistration excludes the credentials already on the account.
func (a *App) handleAddPasskeyBegin(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in first")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()
	user, err := a.store.GetUser(ctx, email)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	options, session, err := a.wa.BeginRegistration(user)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start")
		return
	}
	if err := a.store.SaveWASession(ctx, email, session, user.handle); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not persist ceremony")
		return
	}
	a.setCeremony(w, email)
	writeJSON(w, http.StatusOK, options)
}

func (a *App) handleAddPasskeyFinish(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in first")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()
	session, _, err := a.store.TakeWASession(ctx, email)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "session expired — start over")
		return
	}
	user, err := a.store.GetUser(ctx, email)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	cred, err := a.wa.FinishRegistration(user, *session, r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "could not add passkey")
		return
	}
	if err := a.store.AddCredential(ctx, email, cred); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not save passkey")
		return
	}
	a.clearCookie(w, ceremonyCookie)
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
