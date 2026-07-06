package main

import (
	"encoding/json"
	"net/http"
)

type saveReq struct {
	Summary gameSummary `json:"summary"`
}

func (a *App) handleSaveSession(w http.ResponseWriter, r *http.Request) {
	email := a.currentEmail(r)
	if email == "" {
		writeErr(w, http.StatusUnauthorized, "sign in to save sessions")
		return
	}
	// Decode into a strictly-typed numeric struct: unknown fields (e.g. the
	// client's per-roll `results` array) are ignored and never stored, and any
	// non-numeric value in a known field fails decoding — so no arbitrary or
	// oversized JSON, and no markup, can reach the database or the history view.
	var req saveReq
	if err := json.NewDecoder(r.Body).Decode(&req); err != nil {
		writeErr(w, http.StatusBadRequest, "invalid summary")
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
	out := make([]map[string]any, 0, len(sessions))
	for _, s := range sessions {
		out = append(out, map[string]any{
			"initial_bankroll": s.Summary.InitialBankroll,
			"final_bankroll":   s.Summary.FinalBankroll,
			"roll_count":       s.Summary.RollCount,
			"points_hit":       s.Summary.PointsHit,
			"strikes":          s.Summary.Strikes,
			"total_wagered":    s.Summary.TotalWagered,
			"net":              s.Summary.Net,
			"created_at":       s.CreatedAt,
		})
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
	var netTotal, rollsTotal int64
	for _, s := range sessions {
		netTotal += int64(s.Summary.Net)
		rollsTotal += int64(s.Summary.RollCount)
	}
	writeJSON(w, http.StatusOK, map[string]any{
		"count":       len(sessions),
		"net_total":   netTotal,
		"rolls_total": rollsTotal,
	})
}
