# Privacy Policy

**Effective date: July 20, 2026**

_This is a plain-language privacy policy for a small, free hobby app. It is not legal advice._

---

## Who we are

Dark Side Craps (available at https://craps.moosequest.app) is operated by **MooseQuest LLC**. We can be reached at **privacy@moosequest.app** for any privacy-related questions.

---

## What this app is

Dark Side Craps is a free, entertainment-only web app for learning and tracking a disciplined craps strategy using simulated dice rolls. **No real money is involved.** No wagering, no deposits, no payouts, no prizes.

---

## What we collect and why

We collect the minimum necessary to run an account-based app.

| Data | Why we collect it |
|---|---|
| **Email address** | Used as your account identifier and to send a one-time verification code at sign-in. |
| **Passkey (WebAuthn public-key credential)** | Lets you sign in without a password. We store only the public-key portion — your private key never leaves your device. No passwords are ever collected or stored. |
| **Game history** | Per-session summaries of your simulated play (starting bankroll, roll count, wins/losses, net result). All figures are pretend money. Stored so you can review your practice history. |
| **Session cookie** | A signed, HttpOnly cookie that keeps you signed in during and across sessions. |
| **Verification cookie** | A short-lived cookie issued during the sign-up/sign-in flow, discarded once verified. |

We do **not** collect payment information, location data, advertising identifiers, or any sensitive personal data.

---

## Cookies

We use exactly two types of cookies:

- **Session cookie** — keeps you authenticated. Signed and HttpOnly (not accessible to JavaScript). Expires when you sign out or after a reasonable inactivity period.
- **Verification cookie** — short-lived; used only during the email-verification step at sign-up or sign-in.

We do **not** use advertising cookies, tracking pixels, analytics cookies, or any third-party cookies.

---

## How your data is stored and secured

- All data is transmitted over HTTPS (TLS-encrypted in transit).
- Account data (email, passkey credentials, game history) is stored in **MongoDB** hosted by a managed MongoDB database host.
- Passkey credentials are stored as public-key material only; no passwords exist in our system.
- Session cookies are signed and HttpOnly to resist tampering and cross-site script access.

We are a small hobby operation and implement reasonable security measures, but we cannot guarantee absolute security.

---

## Third-party processors

We share data only with the vendors needed to operate the service. We do **not** sell your data.

| Vendor | Role | Data involved |
|---|---|---|
| **Postmark** | Transactional email (verification codes) | Your email address |
| **Heroku** | Application hosting | All app data in transit and at rest on Heroku infrastructure |
| **Cloudflare** | DNS and edge/CDN layer | IP addresses and request metadata (standard CDN handling) |
| **MongoDB host** | Database storage | Email, passkey credentials, game history |

Each vendor has its own privacy practices. We encourage you to review them if you have concerns.

---

## Data retention and deletion

- We retain your account data for as long as your account is active.
- **You can delete your account at any time** using the "Delete account" button inside the app. This permanently and immediately removes your email address, passkey credentials, and all game history. Deletion is irreversible.
- We may retain anonymized, non-identifiable aggregate statistics (e.g., "N sessions were played this month") that cannot be linked back to you.
- If we receive a verifiable deletion request by email, we will process it within 30 days.

---

## Children

This app is not directed to children under **13**, and in line with the U.S. Children's Online Privacy Protection Act (**COPPA**) we do not knowingly collect personal information from children under 13. If you are in the EU/EEA, a higher minimum age (typically 16, or your country's digital-consent age) may apply. If you believe a child under the applicable age has created an account, contact us at privacy@moosequest.app and we will delete it promptly.

---

## Your rights

Regardless of where you live, you can:

- **Access** your data — contact us and we will tell you what we hold about you.
- **Delete** your data — use the in-app "Delete account" button (instant) or contact privacy@moosequest.app.
- **Correct** your data — contact us if your email address needs to change.

**If you are in the European Union / European Economic Area / United Kingdom:** You may have additional rights under the GDPR (or UK GDPR), including the right to object to processing, the right to data portability, and the right to lodge a complaint with your local supervisory authority. Our legal basis for processing your data is the performance of the service you signed up for (Art. 6(1)(b) GDPR) and, where required, your consent.

**If you are in California:** You may have rights under the CCPA/CPRA. We do not sell or share your personal information. We do not use it for cross-context behavioral advertising.

---

## Changes to this policy

We may update this policy. If we make material changes, we will update the effective date at the top of this page. For a hobby app, we will make reasonable efforts to note significant changes; continued use after a change means you accept the updated policy.

---

## Contact

Privacy questions or requests: **privacy@moosequest.app**

---
