"""
Infisical-backed secret loader for the craps service.

Reads secrets from /MooseQuest/craps/ in Infisical (vault.raxx.app). When
INFISICAL_CLIENT_ID is unset, falls back to plain environment variables so
local development (e.g. local Mongo on 27017) keeps working without vault
access.

Pattern mirrors velvet/client/infisical.py from the TradeMasterAPI repo.
"""

from __future__ import annotations

import logging
import os
from typing import Optional

import requests

logger = logging.getLogger(__name__)

_TIMEOUT_S = 10
_USER_AGENT = "craps-secrets-loader/1.0"
_VENDOR_PATH = "/MooseQuest/craps"


class SecretError(Exception):
    """Raised when a required secret cannot be resolved."""


def _infisical_base() -> str:
    return os.environ.get("INFISICAL_HOST", "https://app.infisical.com").rstrip("/")


def _project_id() -> Optional[str]:
    return os.environ.get("INFISICAL_PROJECT_ID", "").strip() or None


def _infisical_env() -> str:
    return os.environ.get("INFISICAL_ENV", "prod").strip()


def _cf_access_headers() -> dict:
    cid = os.environ.get("CF_ACCESS_CLIENT_ID", "")
    csec = os.environ.get("CF_ACCESS_CLIENT_SECRET", "")
    if cid and csec:
        return {"CF-Access-Client-Id": cid, "CF-Access-Client-Secret": csec}
    return {}


def _service_token() -> Optional[str]:
    """Direct Infisical service token (st.xxx.yyy.zzz format)."""
    return os.environ.get("INFISICAL_SERVICE_TOKEN", "").strip() or None


def _vault_enabled() -> bool:
    if not _project_id():
        return False
    if _service_token():
        return True
    return bool(
        os.environ.get("INFISICAL_CLIENT_ID", "").strip()
        and os.environ.get("INFISICAL_CLIENT_SECRET", "").strip()
    )


def _authenticate() -> str:
    """Return a bearer token usable for /api/v3/secrets/raw.

    Prefers an Infisical service token when present; otherwise performs a
    Universal Auth (machine identity) login and returns the short-lived
    access token.
    """
    svc = _service_token()
    if svc:
        return svc
    host = _infisical_base()
    client_id = os.environ["INFISICAL_CLIENT_ID"].strip()
    client_secret = os.environ["INFISICAL_CLIENT_SECRET"].strip()
    resp = requests.post(
        f"{host}/api/v1/auth/universal-auth/login",
        json={"clientId": client_id, "clientSecret": client_secret},
        headers={"User-Agent": _USER_AGENT, **_cf_access_headers()},
        timeout=_TIMEOUT_S,
    )
    if not resp.ok:
        raise SecretError(f"Infisical auth failed: HTTP {resp.status_code}")
    token = resp.json().get("accessToken")
    if not token:
        raise SecretError("Infisical auth response missing accessToken")
    return token


def _fetch_from_vault(name: str, token: str) -> str:
    host = _infisical_base()
    pid = _project_id()
    resp = requests.get(
        f"{host}/api/v3/secrets/raw/{name}",
        params={
            "workspaceId": pid,
            "environment": _infisical_env(),
            "secretPath": _VENDOR_PATH,
        },
        headers={
            "Authorization": f"Bearer {token}",
            "User-Agent": _USER_AGENT,
            **_cf_access_headers(),
        },
        timeout=_TIMEOUT_S,
    )
    if resp.status_code == 404:
        raise SecretError(f"Secret not found in vault: {_VENDOR_PATH}/{name}")
    if not resp.ok:
        raise SecretError(
            f"Vault fetch failed for {name}: HTTP {resp.status_code}"
        )
    secret = resp.json().get("secret", {})
    return secret.get("secretValue", "")


def load(name: str, *, env_fallback: Optional[str] = None, default: Optional[str] = None) -> Optional[str]:
    """Resolve a secret by name.

    Resolution order:
        1. If Infisical creds are configured, fetch /MooseQuest/craps/{name}.
           A 404 (not-found) for an optional secret falls through to env/default;
           any other vault failure (auth, network, 5xx) raises SecretError.
        2. Otherwise read os.environ[env_fallback or name]
        3. Otherwise return `default`

    The intent: vault-enabled environments still get clean error messages on
    misconfigured creds or network outages, but a single missing optional
    secret doesn't break startup if a sane default is provided in code.
    """
    env_key = env_fallback or name
    if _vault_enabled():
        try:
            token = _authenticate()
            value = _fetch_from_vault(name, token)
            if value:
                logger.info("secrets_loader: loaded %s from vault", name)
                return value
            logger.warning("secrets_loader: vault returned empty value for %s", name)
        except SecretError as exc:
            # 404 on an optional secret → fall through; any other vault error
            # is a real problem (bad auth, CF gate, 5xx) and must surface.
            if "not found" not in str(exc).lower():
                raise
            logger.info("secrets_loader: %s not in vault, falling back", name)
        except requests.exceptions.RequestException as exc:
            raise SecretError(f"Vault unreachable: {exc}") from exc
    value = os.environ.get(env_key)
    if value is not None:
        logger.info("secrets_loader: loaded %s from env", name)
        return value
    return default
