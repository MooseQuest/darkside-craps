# Dark Side Craps

A disciplined **Don't-Pass / lay-odds** craps strategy tracker — the "Long Dark Side."
The goal isn't to win big; it's to play close to the house edge for a long, controlled
session, and to know when to walk away. See `docs/KH-Long-Dark-Side.md` for the strategy.

**Live:** https://craps.moosequest.app

## Stack

The lightest thing that ships: a single **Go** binary that embeds a dependency-free
vanilla-JS single-page app. The craps engine runs entirely in the browser (instant, no
per-roll round-trip); the server exists only to authenticate players with **passkeys**
(email + WebAuthn — no passwords) and to save their session history in MongoDB.

```
web/            client SPA (served embedded)
  engine.js       craps state machine (browser + Node)
  engine.test.mjs unit tests (node --test)
  index.html / app.js / styles.css
main.go         config, embed, HTTP server, /healthz
store.go        MongoDB: users, credentials, ceremony + game sessions
app.go          routing, static assets, signed session cookies
auth.go         passkey register/login/logout (go-webauthn)
api.go          per-account game history + summary
```

## Develop

Requires Go 1.23+, Node 20+, and a MongoDB.

```bash
# start a throwaway mongo
docker run -d --name craps-mongo -p 27017:27017 mongo:7

# run the app
MONGO_URI=mongodb://localhost:27017 SECRET_KEY=dev \
  RP_ID=localhost RP_ORIGIN=http://localhost:8080 PORT=8080 \
  go run .
# → http://localhost:8080
```

### Tests

```bash
go test ./...            # server: signed sessions, email validation, helpers
cd web && node --test    # craps engine: every transition, strikes, seven-out
```

## Configuration (env)

| Var | Purpose |
|-----|---------|
| `MONGO_URI` | MongoDB connection string (required) |
| `SECRET_KEY` | HMAC key for signed session cookies (required in prod) |
| `RP_ID` | WebAuthn Relying Party ID = the domain (`craps.moosequest.app`) |
| `RP_ORIGIN` | Full origin (`https://craps.moosequest.app`) |
| `MONGO_DB` | Database name (default `craps_game`) |
| `PORT` | Listen port (default `8080`; set by Heroku) |
| `SMTP_HOST` | Transactional-email relay host. **When set, signup email-verification is enforced** (an emailed 6-digit code is required before a passkey can be created). Unset ⇒ open registration (dev). |
| `SMTP_PORT` | SMTP submission port (default `587`) |
| `SMTP_USERNAME` / `SMTP_PASSWORD` | SMTP auth (optional; PLAIN over STARTTLS) |
| `MAIL_FROM` | From address (default `Dark Side Craps <no-reply@moosequest.app>`) |

Secrets are sourced from the **mqcraps** project in the vault (vault.raxx.app); see
`.vault.toml`. On Heroku they are injected as config vars.

## Deploy

Pushing to `main` runs CI (build + Go tests + JS engine tests + a boot smoke test against
a Mongo service container) and, when green, deploys to Heroku (Go buildpack) and
smoke-tests `https://craps.moosequest.app/healthz`. Passkeys are bound to the RP ID, so
the app must be reached at `craps.moosequest.app` for registration/login to work.
