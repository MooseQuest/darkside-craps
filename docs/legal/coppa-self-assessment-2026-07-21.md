# Dark Side Craps — COPPA Self-Assessment

> **Good-faith operator self-assessment. Not legal advice, and not a filing.**
> This is the dated, self-executed record of Dark Side Craps' COPPA posture, prepared per
> the reasoning in `docs/legal/diy-compliance-plan.md` §1–2. The factual answers below were
> pre-filled from the product's actual code, data flows, and policies as of the date shown;
> the operator should review each line, correct anything that has changed, then sign and date
> Section F to make it a complete record. Re-run and save a new dated copy whenever the
> product materially changes (see Section F triggers), or at minimum annually.

**Operator:** MooseQuest LLC
**Product:** Dark Side Craps (https://craps.moosequest.app)
**Assessment date:** 2026-07-21
**Completed by:** _______________________  (operator to sign in Section F)

---

## Section A — Data map (what we collect, and why)

- [x] **Email address** — account identifier + one-time verification code at sign-in.
- [x] **WebAuthn public-key credential** — passkey authentication; the private key never leaves the user's device.
- [x] **Game history** — session summaries of simulated (fictional-money) play.
- [x] **Session cookie** — signed, HttpOnly, authentication only.
- [x] **Verification cookie** — short-lived, sign-up/sign-in flow only.
- [x] Confirmed: **NO** birthdate, **NO** school, **NO** real name required, **NO** physical address, **NO** payment data, **NO** advertising identifiers, **NO** location data collected. *(The age gate records only a boolean "meets age requirement" in local storage — no specific birthdate is collected or stored.)*
- [x] Confirmed: **NO** data sold to third parties. The vendors in the Privacy Policy table (Postmark, Heroku, Cloudflare, MongoDB host) are service providers only; none receive data for their own marketing/advertising purposes.

## Section B — "Directed to children" totality test (16 CFR § 312.2 factors)

| Factor | Fact | Conclusion |
|---|---|---|
| Subject matter | Adult casino gambling strategy (Don't-Pass / lay-odds craps). | **Not child-directed** — one of the strongest "not directed to children" fact patterns available. |
| Visual content | Retro/"dark side" aesthetic; no cartoon characters, no toy-like palette aimed at children. | Not child-directed. |
| Animated characters / child-oriented activities | None. | Not child-directed. |
| Music / audio aimed at children | None (no audio content). | Not child-directed. |
| Age of models/actors depicted | None depicted. | Not child-directed. |
| Child celebrities / celebrities appealing to children | None. | Not child-directed. |
| Advertising on / promoting the service | No ads run on the service at all. Confirm any external marketing (social posts, etc.) is not child-targeted. | Not child-directed. |
| Empirical evidence of audience composition | No analytics/audience data collected. Redo this line if analytics are ever added. | Not child-directed (no contrary evidence). |
| Marketing / promotional materials (2025 addition) | No paid marketing at present; word-of-mouth / personal channels, adult-oriented. | Not child-directed. |
| Representations to consumers / third parties (2025 addition) | Terms & Privacy state a 13+/16+/18+ tiered minimum age and "not directed to children under 13." | Not child-directed. |
| User / third-party reviews (2025 addition) | No known reviews; no reviewer has identified as a minor. | Not child-directed (no contrary evidence). |
| Age of users on comparable sites/services (2025 addition) | Comparable products (social-casino apps, craps strategy content) skew adult. | Not child-directed. |

**Overall conclusion:** ☑ **Not directed to children under 13.**  ☐ Reassess — flag raised.

## Section C — Actual-knowledge check

- [x] Confirmed: no signup field requests birthdate/age directly (only the neutral age-gate attestation, which stores no specific birthdate).
- [x] Confirmed: no support ticket, email, or user report has ever identified a user as under the applicable minimum age. *If one ever does:* deletion completed within __________ (date) and logged here.
- [x] Confirmed: the age gate is **neutral** — it does not pre-select a passing age, does not use dark patterns to encourage false entry, and blocks/warns rather than silently proceeding on a failing answer.

## Section D — No behavioral advertising / no data monetization

*(This is the core harm the 2025 COPPA amendments target — "limiting companies' ability to monetize kids' data." It does not exist in this product.)*

- [x] Confirmed: zero ad networks, zero ad SDKs, zero pixel trackers, zero analytics-for-ads integrations present in the codebase as of this assessment date.
- [x] Confirmed: game history and account data are used only to render the user's own history back to them — never aggregated for ad targeting, never sold, never shared with a data broker.
- [x] Confirmed: the Privacy Policy vendor table is current and accurate as of this assessment date.

## Section E — Deletion mechanism

- [x] Confirmed: in-app **"Delete account"** button provides immediate, permanent deletion of email, passkey credentials, and game history. *Verify by testing on this date:* __________.
- [ ] Confirmed: email-request deletion path (**privacy@moosequest.net**, 30-day SLA) is staffed and monitored. *(Operator to confirm the inbox is actively monitored — see diy-compliance-plan.md action list.)*

## Section F — Sign-off

This self-assessment was completed in good faith by the operator on the date below, based on
the product's actual code, data flows, and policies as they existed on that date.

**Redo this assessment whenever:**
- (a) the product adds ads, analytics, or any new data collection;
- (b) the product's subject matter or visual design changes materially;
- (c) any user report suggests an under-13 (or jurisdiction-applicable-age) user; or
- (d) at minimum, annually.

**Signed:** _______________________   **Date:** __________________

**Next scheduled re-assessment date:** __________________ *(suggest: 2027-07-21)*

---

*Basis and sources for the conclusions above: see `docs/legal/diy-compliance-plan.md` §1
(COPPA threshold analysis, FTC Feb-2026 age-verification policy statement, enforcement
pattern) and `docs/legal/compliance-research.md` §1a (factor-by-factor "not directed to
children" analysis).*
