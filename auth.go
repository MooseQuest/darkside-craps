package main

import (
	"net/http"
	"strings"
)

func validEmail(e string) bool {
	e = strings.TrimSpace(e)
	return len(e) >= 3 && len(e) <= 254 && strings.Count(e, "@") == 1 &&
		!strings.HasPrefix(e, "@") && !strings.HasSuffix(e, "@")
}

type emailReq struct {
	Email string `json:"email"`
}

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
	email := strings.ToLower(strings.TrimSpace(req.Email))
	ctx, cancel := reqCtx(r)
	defer cancel()

	user, err := a.store.GetOrCreateUser(ctx, email)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	options, session, err := a.wa.BeginRegistration(user)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start registration")
		return
	}
	if err := a.store.SaveWASession(ctx, email, session); err != nil {
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

	user, err := a.store.GetOrCreateUser(ctx, email)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	session, err := a.store.TakeWASession(ctx, email)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "registration session expired — start over")
		return
	}
	cred, err := a.wa.FinishRegistration(user, *session, r)
	if err != nil {
		writeErr(w, http.StatusBadRequest, "passkey registration failed")
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
	email := strings.ToLower(strings.TrimSpace(req.Email))
	ctx, cancel := reqCtx(r)
	defer cancel()

	user, err := a.store.GetUser(ctx, email)
	if err == errNotFound || (user != nil && len(user.WebAuthnCredentials()) == 0) {
		writeErr(w, http.StatusNotFound, "no passkey found for that email")
		return
	}
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load account")
		return
	}
	options, session, err := a.wa.BeginLogin(user)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not start login")
		return
	}
	if err := a.store.SaveWASession(ctx, email, session); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not persist ceremony")
		return
	}
	a.setCeremony(w, email)
	writeJSON(w, http.StatusOK, options)
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
	session, err := a.store.TakeWASession(ctx, email)
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
