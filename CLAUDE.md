# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## What this is

A guided tracker for the **"Long Dark Side"** craps strategy (disciplined Don't-Pass /
lay-odds play). The player rolls at a real table (or clicks Roll to generate dice) and the
app tells them exactly what to bet next while tracking bankroll, strikes, and points.
Live at **craps.moosequest.app**. Strategy: `docs/KH-Long-Dark-Side.md`.

## Stack (rewritten 2026-07-03; was Python/Flask)

A single **Go** binary that `go:embed`s a dependency-free vanilla-JS SPA. Game logic runs
**client-side**; the server only does passkey auth + history persistence.

## Commands

```bash
# Run (needs a MongoDB)
MONGO_URI=mongodb://localhost:27017 SECRET_KEY=dev \
  RP_ID=localhost RP_ORIGIN=http://localhost:8080 PORT=8080 go run .

go test ./...          # Go: signed-cookie sessions, email validation, helpers
cd web && node --test  # JS engine: every craps transition (run a single test with --test-name-pattern)
go build -o bin/darkside-craps .
```

CI (`.github/workflows/deploy.yml`) is the authoritative check: `go vet` + `go test` +
`node --test` + a boot smoke test against a Mongo service container, then (on `main`)
`git push` to Heroku and a prod `/healthz` smoke test.

## Architecture

**The craps strategy lives in TWO places that must stay in sync conceptually:** the
client engine `web/engine.js` (`CrapsGame`) is a faithful port of the original Python
`CrapsStateMachine`. It is a pure state machine with two states (`comeOutRoll` /
point-established) — `roll()` dispatches to `_comeOutRoll` / `_pointEstablishedRoll`. Dice
are injected via a `roller` callback so `web/engine.test.mjs` can script outcomes
deterministically. **Bet sizing is bankroll-tiered** (`initialBankroll === 1000` → $25
base + smaller lay odds; else $50) and that branch is duplicated across `_baseBet`,
`_placeOddsBet`, and `_getLayBet` — change them together. Strike logic: repeating a working
point is a strike; 3 strikes clears bets until a 7; a 7 with points working wins all lays
and resets.

**The server is stateless about game play.** It never runs the engine. Responsibilities:

- `main.go` — env `Config`, `//go:embed web/...`, HTTP server, `/healthz`. Fatal if
  `MONGO_URI` is unset.
- `app.go` — routing (exact-path static asset handlers + `/` catch-all), JSON helpers,
  and **stateless signed session cookies**: `sign`/`verify` HMAC a `email|expiryUnix`
  payload with `SECRET_KEY`. There is no server-side session table — auth state is the
  cookie. `secureCookies()` keys off an `https://` `RP_ORIGIN`.
- `auth.go` — passkey ceremonies via **go-webauthn**. Flow quirk: `FinishRegistration`/
  `FinishLogin` read the credential from the request body, but the **email comes from a
  separate short-lived signed `craps_wa` cookie** set during `begin` (finish has no email
  in its body). The in-flight `webauthn.SessionData` (challenge) is stashed in Mongo
  `wa_sessions` keyed by email and consumed (find-and-delete) on finish.
- `store.go` — MongoDB. `User` implements `webauthn.User`; a `webauthn.Credential` is
  stored whole as BSON (round-trips through the same struct — no custom tags). Collections:
  `users`, `credentials`, `wa_sessions`, `game_sessions`.
- `api.go` — auth-gated history/summary; anonymous users can play but not persist.

**Client ↔ server WebAuthn contract:** `app.js` base64url-encodes/decodes the ArrayBuffer
fields (`challenge`, `user.id`, credential `rawId`, `attestationObject`, etc.) because the
go-webauthn JSON uses URL-safe base64. If you change the auth JSON shape on either side,
change both `app.js` and the go-webauthn call sites.

## Secrets & deploy

Runtime secrets live in the **mqcraps** vault project (vault.raxx.app) at `prod:/`
(`MONGO_URI`, `SESSION_MONGO_URI`, `SECRET_KEY`); `.vault.toml` is the manifest. On Heroku
they are injected as config vars (plus `RP_ID`, `RP_ORIGIN`). The Heroku app is `mq-craps`;
`craps.moosequest.app` is a Cloudflare CNAME to it (zone `moosequest.app`, account
MooseQuest). **Passkeys are bound to `RP_ID`** — they only work when the page is served
from `craps.moosequest.app`, not the `*.herokuapp.com` URL.

The loose `*.png` files and `.playwright-mcp/` at the repo root are QA/walkthrough
artifacts (gitignored), not app inputs.
