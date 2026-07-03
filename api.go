package main

import (
	"net/http"
)

type saveReq struct {
	Summary map[string]any `json:"summary"`
}

func (a *App) handleSaveSession(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in to save sessions")
		return
	}
	var req saveReq
	if err := readJSON(r, &req); err != nil || req.Summary == nil {
		writeErr(w, http.StatusBadRequest, "summary is required")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()
	if err := a.store.SaveGameSession(ctx, email, req.Summary); err != nil {
		writeErr(w, http.StatusInternalServerError, "could not save session")
		return
	}
	writeJSON(w, http.StatusOK, map[string]any{"ok": true})
}

func (a *App) handleHistory(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in to view history")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()
	sessions, err := a.store.History(ctx, email, 100)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load history")
		return
	}
	// Flatten each session's summary + created_at for the client.
	out := make([]map[string]any, 0, len(sessions))
	for _, s := range sessions {
		row := map[string]any{}
		for k, v := range s.Summary {
			row[k] = v
		}
		row["created_at"] = s.CreatedAt
		out = append(out, row)
	}
	writeJSON(w, http.StatusOK, map[string]any{"sessions": out})
}

func (a *App) handleSummary(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in to view summary")
		return
	}
	ctx, cancel := reqCtx(r)
	defer cancel()
	sessions, err := a.store.History(ctx, email, 1000)
	if err != nil {
		writeErr(w, http.StatusInternalServerError, "could not load summary")
		return
	}
	var netTotal, rollsTotal float64
	for _, s := range sessions {
		netTotal += toFloat(s.Summary["net"])
		rollsTotal += toFloat(s.Summary["roll_count"])
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"count":       len(sessions),
		"net_total":   int64(netTotal),
		"rolls_total": int64(rollsTotal),
	})
}

func toFloat(v any) float64 {
	switch n := v.(type) {
	case float64:
		return n
	case float32:
		return float64(n)
	case int:
		return float64(n)
	case int32:
		return float64(n)
	case int64:
		return float64(n)
	default:
		return 0
	}
}
