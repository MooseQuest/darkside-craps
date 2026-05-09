"""
One-shot script that pushes the craps service's secrets into Infisical
under /MooseQuest/craps/ and creates the folder if it does not exist.

Usage:
    INFISICAL_CLIENT_ID=...   \
    INFISICAL_CLIENT_SECRET=...   \
    INFISICAL_PROJECT_ID=...   \
    INFISICAL_HOST=https://vault.raxx.app   \
    INFISICAL_ENV=prod   \
    CF_ACCESS_CLIENT_ID=...   \
    CF_ACCESS_CLIENT_SECRET=...   \
    python scripts/bootstrap_vault.py

The script reads the three secret values from stdin (prompted, no echo)
so they never live on disk or in shell history. Run once after rotating
the leaked Atlas credentials.
"""

from __future__ import annotations

import getpass
import os
import sys
from typing import Optional

import requests

VENDOR_PATH = "/MooseQuest/craps"
TIMEOUT_S = 15
USER_AGENT = "craps-bootstrap-vault/1.0"


def _cf_headers() -> dict:
    cid = os.environ.get("CF_ACCESS_CLIENT_ID", "")
    csec = os.environ.get("CF_ACCESS_CLIENT_SECRET", "")
    if cid and csec:
        return {"CF-Access-Client-Id": cid, "CF-Access-Client-Secret": csec}
    return {}


def _require(name: str) -> str:
    val = os.environ.get(name, "").strip()
    if not val:
        sys.exit(f"ERROR: env var {name} is required")
    return val


def get_bearer_token() -> str:
    """Return a bearer token for the /api/v3 + /api/v1/folders endpoints.

    Prefers INFISICAL_SERVICE_TOKEN (st.xxx format) used directly as the
    bearer; otherwise falls back to Universal Auth machine identity via
    INFISICAL_CLIENT_ID + INFISICAL_CLIENT_SECRET.
    """
    svc = os.environ.get("INFISICAL_SERVICE_TOKEN", "").strip()
    if svc:
        return svc
    host = os.environ.get("INFISICAL_HOST", "https://vault.raxx.app").rstrip("/")
    client_id = _require("INFISICAL_CLIENT_ID")
    client_secret = _require("INFISICAL_CLIENT_SECRET")
    resp = requests.post(
        f"{host}/api/v1/auth/universal-auth/login",
        json={"clientId": client_id, "clientSecret": client_secret},
        headers={"User-Agent": USER_AGENT, **_cf_headers()},
        timeout=TIMEOUT_S,
    )
    resp.raise_for_status()
    token = resp.json().get("accessToken")
    if not token:
        sys.exit("ERROR: Infisical auth returned no accessToken")
    return token


def ensure_folder(host: str, token: str, project_id: str, env: str) -> None:
    """Create /MooseQuest/craps/. The parent /MooseQuest must already exist
    (created when the TradeMasterAPI vendors were set up). 409 = already
    there, treated as success."""
    parent, leaf = VENDOR_PATH.rsplit("/", 1)
    parent = parent or "/"
    resp = requests.post(
        f"{host}/api/v1/folders",
        json={
            "workspaceId": project_id,
            "environment": env,
            "path": parent,
            "name": leaf,
        },
        headers={
            "Authorization": f"Bearer {token}",
            "Content-Type": "application/json",
            "User-Agent": USER_AGENT,
            **_cf_headers(),
        },
        timeout=TIMEOUT_S,
    )
    if resp.status_code == 409:
        print(f"[ok] folder already exists: {VENDOR_PATH}")
        return
    resp.raise_for_status()
    print(f"[ok] created folder: {VENDOR_PATH}")


def write_secret(host: str, token: str, project_id: str, env: str, name: str, value: str) -> None:
    """PATCH (update) → fall back to POST (create) if 404."""
    url = f"{host}/api/v3/secrets/raw/{name}"
    payload = {
        "workspaceId": project_id,
        "environment": env,
        "secretPath": VENDOR_PATH,
        "secretValue": value,
    }
    headers = {
        "Authorization": f"Bearer {token}",
        "Content-Type": "application/json",
        "User-Agent": USER_AGENT,
        **_cf_headers(),
    }

    resp = requests.patch(url, json=payload, headers=headers, timeout=TIMEOUT_S)
    if resp.ok:
        print(f"[ok] updated secret: {VENDOR_PATH}/{name}")
        return
    if resp.status_code != 404:
        sys.exit(f"ERROR: PATCH {name} returned HTTP {resp.status_code}: {resp.text[:200]}")

    resp = requests.post(url, json=payload, headers=headers, timeout=TIMEOUT_S)
    if resp.ok:
        print(f"[ok] created secret: {VENDOR_PATH}/{name}")
        return
    sys.exit(f"ERROR: POST {name} returned HTTP {resp.status_code}: {resp.text[:200]}")


def prompt_secret(label: str, *, allow_empty: bool = False) -> Optional[str]:
    val = getpass.getpass(f"{label} (will not echo, blank to skip): ").strip()
    if not val and not allow_empty:
        return None
    return val or None


def main() -> None:
    host = os.environ.get("INFISICAL_HOST", "https://vault.raxx.app").rstrip("/")
    env = os.environ.get("INFISICAL_ENV", "prod").strip()
    project_id = _require("INFISICAL_PROJECT_ID")

    print(f"Target: {host}  env={env}  path={VENDOR_PATH}")
    print()
    print("Enter the new secret values. Blank input skips that secret.")

    mongo_uri = prompt_secret("MONGO_URI (full mongodb+srv:// connection string)")
    session_mongo_uri = prompt_secret("SESSION_MONGO_URI (blank = same as MONGO_URI)", allow_empty=True)
    secret_key = prompt_secret("SECRET_KEY (Flask session signing key)")

    if not (mongo_uri or session_mongo_uri or secret_key):
        sys.exit("Nothing to write.")

    print()
    print("Obtaining bearer token…")
    token = get_bearer_token()

    print(f"Ensuring folder {VENDOR_PATH} exists…")
    ensure_folder(host, token, project_id, env)

    if mongo_uri:
        write_secret(host, token, project_id, env, "MONGO_URI", mongo_uri)
    if session_mongo_uri:
        write_secret(host, token, project_id, env, "SESSION_MONGO_URI", session_mongo_uri)
    if secret_key:
        write_secret(host, token, project_id, env, "SECRET_KEY", secret_key)

    print()
    print("Done. Verify in Infisical UI: " + host)


if __name__ == "__main__":
    main()
