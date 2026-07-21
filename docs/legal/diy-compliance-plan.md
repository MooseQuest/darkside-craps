# Dark Side Craps — DIY Compliance Plan (No-Attorney Path)

> **Status: research-only. This document does NOT constitute legal or tax advice.**
> The operator has made a business decision not to obtain a paid attorney sign-off on this
> hobby app. This document is written for that reality: it gives the facts, the honest
> risk assessment, and a concrete self-executed record-keeping path. It does not replace
> a licensed attorney's judgment — it documents a reasonable, good-faith DIY posture in
> lieu of one. If circumstances change (real money, ads/monetization, material EU/Brazil/
> Australia traffic, an actual complaint or regulator inquiry), the calculus in this
> document changes too — reopen it with a privacy/COPPA attorney at that point.
> Last updated: 2026-07-21. Sources as of that date — verify freshness before relying on them.
>
> Builds on: `docs/legal/compliance-research.md` (original briefing), `docs/legal/terms.md`,
> `docs/legal/privacy.md`. Already SHIPPED and not re-litigated here: 13+ base age with
> AU 18+ / EU-EEA 16+ tiering, first-visit age-gate self-attestation, GeoIP geo-block of the
> 20-country list, Terms §3a geographic-restrictions clause, COPPA cited by name in the
> Privacy Policy.

---

## TL;DR (3 sentences)

The 2025 COPPA amendments add real new obligations (verifiable-parental-consent mechanics,
data-minimization/retention limits, a formal "mixed audience" category) — but every one of
them is triggered by being **child-directed or having actual knowledge of under-13 users**,
neither of which applies to a self-attested-13+, adult-casino-subject-matter, no-ads,
no-behavioral-tracking hobby app; the "audit deadline" reported as overdue is therefore a
best-practice recommendation for this app's risk profile, not a filing obligation the app
has missed. The one governing-law gap (Terms §12) is resolvable without a lawyer: MooseQuest
LLC's own formation records show it was organized in **Pennsylvania** (PA Certificate of
Organization approved 2026-05-22 via Northwest Registered Agent; EIN issued same day) — name
Pennsylvania in §12. The rest of this doc is a self-audit checklist Kris can fill out and
date himself this week, plus a final, decisive call on the remaining borderline items (UK
AADC, Brazil LGPD/Digital ECA, the 9 borderline geo-block countries) so nothing sits in limbo.

---

## 1. COPPA 2025 amendment reality check

### 1a. What the 2025 amendments actually changed

The FTC's final rule (published in the Federal Register April 22, 2025; effective June 23,
2025; compliance deadline **April 22, 2026**) made these substantive changes, all of which
regulate operators who are **already covered** by COPPA (child-directed, or actual knowledge
of under-13 collection):

- A new, separate regulatory definition of **"mixed audience" service** (general-audience
  site with a distinguishable child-directed section or feature).
- Expanded definitions of **personal information** (adds biometric identifiers and
  government-issued identifiers to the list of regulated data types).
- New, more specific methods and content requirements for **verifiable parental consent**
  and parent notices.
- Tighter **data retention and deletion** requirements — data may be retained only "as long
  as reasonably necessary" for the specific purpose collected, with a written data-retention
  policy required of covered operators.
- Expanded **Safe Harbor program** oversight (public membership disclosure, annual
  disciplinary-action reporting, triennial technology reviews) — this only affects operators
  who are *members* of an FTC-approved Safe Harbor program.
- Source: FTC, [FTC Finalizes Changes to Children's Privacy Rule (Jan 2025)](https://www.ftc.gov/news-events/news/press-releases/2025/01/ftc-finalizes-changes-childrens-privacy-rule-limiting-companies-ability-monetize-kids-data)
- Source: Federal Register, [COPPA Final Rule Amendments](https://www.federalregister.gov/documents/2025/04/22/2025-05904/childrens-online-privacy-protection-rule)
- Source: Securiti, [FTC's 2025 COPPA Final Rule Amendments — plain-language summary](https://securiti.ai/ftc-coppa-final-rule-amendments/)
- Source: Latham & Watkins, [FTC Publishes Updates to COPPA Rule](https://www.lw.com/en/insights/ftc-publishes-updates-to-coppa-rule)

**None of these new obligations apply to an operator outside COPPA's threshold trigger.**
COPPA has always applied only to services that are (a) directed to children under 13, or
(b) have actual knowledge they're collecting personal information from a child under 13.
The 2025 amendments did not lower that threshold or add a new "general audience site must
prove it's not covered" filing requirement. Source: 16 CFR § 312.3, [via eCFR](https://www.ecfr.gov/current/title-16/chapter-I/subchapter-C/part-312).

### 1b. Does any of it bind Dark Side Craps?

Applying the threshold test directly:

| Trigger | Dark Side Craps status |
|---|---|
| "Directed to children under 13" (16 CFR § 312.2 totality test — subject matter, visual content, animated characters, music, model age, celebrities, advertising, empirical audience data, marketing, third-party representations, reviews, comparable-site demographics) | No factor points toward children. Adult casino-strategy subject matter is one of the strongest "not directed to children" fact patterns available — already assessed factor-by-factor in `compliance-research.md` §1a. |
| "Actual knowledge" of under-13 collection | None. The app collects only an email + WebAuthn public key at signup; there is no birthdate field, no school field, no other signal that would create actual knowledge. The self-attestation age gate is the FTC-recommended mechanism for exactly this purpose (see 1c below). |
| Ads / monetization of children's data (the specific harm the 2025 amendments target — "limiting companies' ability to monetize kids' data") | Not applicable at all — the app carries **no advertising, no ad network SDKs, no behavioral tracking, no data sale**, to anyone, child or adult. The entire monetization-of-minors problem the 2025 amendments were written to solve does not exist in this product. |

**Plain answer: the substantive posture is already compliant, and the "compliance deadline"
that reads as overdue is a deadline for operators who were already covered by COPPA to adopt
new mechanics (parental consent flows, retention schedules, Safe Harbor reporting) — not a
universal audit-and-file requirement that this app missed.** The prior briefing's language
("audit overdue... consult a privacy/COPPA attorney promptly") was conservative advice
appropriate for a paid-attorney-review posture; under the DIY constraint, the more precise
and equally honest framing is: *there is no COPPA filing or deadline this app has missed,
because this app was never inside COPPA's coverage trigger to begin with.* What's actually
useful is a documented, dated self-assessment proving that conclusion was reached in good
faith — see the checklist in Section 2.

### 1c. The self-attestation age gate is the FTC's own recommended mechanism

The 2025 amendments require age-screening mechanisms to be "neutral" (not defaulting to a
passing age, not encouraging false entry) and "reasonably calculated, in light of available
technology, to determine whether the site is directed to children." Separately, in
**February 2026** the FTC issued an enforcement-policy statement specifically encouraging
general-audience and mixed-audience operators to use **age-verification/self-attestation
technologies to determine user age**, and committing that the FTC will not bring a COPPA
enforcement action against operators who collect/use/disclose personal information **solely**
to determine age, provided the operator gives adequate notice, uses reasonable safeguards,
and doesn't repurpose that data. Source: FTC, [Enforcement Policy Statement Promoting the
Adoption of Age-Verification Technology (Feb 25, 2026)](https://www.ftc.gov/legal-library/browse/enforcement-policy-statement-promoting-adoption-age-verification-technology);
FTC press release, [FTC Issues COPPA Policy Statement to Incentivize the Use of Age
Verification Technologies (Feb 2026)](https://www.ftc.gov/news-events/news/press-releases/2026/02/ftc-issues-coppa-policy-statement-incentivize-use-age-verification-technologies-protect-children).

**Implication:** the shipped first-visit self-attestation age gate (AU 18+ / EU-EEA 16+ /
else 13+) is not just adequate, it is the specific mechanism the FTC is currently promoting
as the model general-audience operators should use. This is a point *in favor* of the
current implementation, not a gap.

### 1d. Realistic enforcement exposure — who enforces, and against whom

- **Who enforces COPPA:** the FTC (primary), and separately, **state attorneys general** have
  independent statutory authority to sue on behalf of state residents (parens patriae) for
  COPPA rule violations, with a duty to notify the FTC when they do. Source: 15 U.S.C.
  § 6504, [via Cornell LII](https://www.law.cornell.edu/uscode/text/15/6504). There is **no
  private right of action** under COPPA — individuals cannot sue directly; only the FTC or a
  state AG can bring an enforcement action.
- **What actual COPPA enforcement in 2025–2026 looks like:** every publicized action is
  against operators with a genuine child-directed product or documented actual knowledge —
  e.g., Cognosphere/HoYoverse ($20M, Jan 2025, game with child-directed elements and known
  child users), Apitor Technology (Sept 2025, child-directed app + third-party SDK data
  collection), Disney ($10M, Sept 2025, child-directed YouTube videos with ad tracking
  enabled). None resemble a free, no-ads, no-tracking, adult-subject-matter hobby app with a
  neutral age gate. Source: [Davis Polk, FTC Prioritizes COPPA Enforcement](https://www.davispolk.com/insights/client-update/ftc-prioritizes-coppa-enforcement-new-compliance-obligations-take-effect);
  general FTC enforcement action index at [ftc.gov/news-events/news/press-releases](https://www.ftc.gov/news-events/news/press-releases).
- **Statutory penalty exposure, in the abstract:** COPPA violations are treated as unfair-or-
  deceptive-practice violations under the FTC Act, currently carrying civil penalties of
  roughly $50,000+ per violation (inflation-adjusted periodically). This number is real but
  is calculated **per violation** in cases where a court finds the operator was, in fact,
  covered by COPPA and violated it — it is not a standing fine for merely operating a
  general-audience site.
- **Realistic exposure for Dark Side Craps specifically:** low. The FTC's own enforcement
  pattern targets scale (millions of child users), a genuine child-directed product, or
  documented actual knowledge — none of which describes this app. The residual risk is not
  zero (an aggressive plaintiff's-bar-adjacent state AG office could theoretically open an
  inquiry if a complaint were filed, e.g., a parent alleging their child under 13 used the
  app), but the app's own self-attestation gate + adult subject matter + "we do not knowingly
  collect from under-13 users, and will delete promptly if discovered" Privacy Policy
  language is exactly the record that would close such an inquiry quickly and without
  penalty. That is the record Section 2 below builds and dates.

---

## 2. The DIY self-audit — fill this in, date it, keep it

**Purpose:** this is the actual audit. Complete it yourself, save a signed/dated copy (PDF
or dated commit in this repo), and treat it as your compliance file. If anyone ever asks
"did you assess your COPPA posture," this is the answer — a real, dated, good-faith
self-assessment, done by the operator, for free.

```
DARK SIDE CRAPS — COPPA SELF-ASSESSMENT
Operator: MooseQuest LLC
Completed by: _______________________  Date: __________________
Product: Dark Side Craps (https://craps.moosequest.app)

SECTION A — DATA MAP (what we collect, and why)
[ ] Email address — account identifier + one-time verification code at sign-in
[ ] WebAuthn public-key credential — passkey authentication (private key never leaves device)
[ ] Game history — session summaries of simulated (fictional-money) play
[ ] Session cookie — signed, HttpOnly, authentication only
[ ] Verification cookie — short-lived, sign-up/sign-in flow only
[ ] Confirmed: NO birthdate, NO school, NO real name required, NO physical address,
    NO payment data, NO advertising identifiers, NO location data collected
[ ] Confirmed: NO data sold to third parties (see Privacy Policy vendor table:
    Postmark, Heroku, Cloudflare, MongoDB host — all service providers, none receive
    data for their own marketing/ad purposes)

SECTION B — "DIRECTED TO CHILDREN" TOTALITY TEST (16 CFR § 312.2 factors)
For each factor, state the fact and your conclusion:
[ ] Subject matter — adult casino gambling strategy (craps). Conclusion: NOT child-directed.
[ ] Visual content — no cartoon characters, no bright toy-like palette aimed at children.
    Conclusion: __________
[ ] Animated characters / child-oriented activities — none. Conclusion: __________
[ ] Music / audio content aimed at children — none. Conclusion: __________
[ ] Age of models/actors depicted — none depicted. Conclusion: __________
[ ] Child celebrities or celebrities who appeal to children — none. Conclusion: __________
[ ] Advertising on/promoting the service — no ads run on the service at all; confirm any
    external marketing (social posts, etc.) is not child-targeted. Conclusion: __________
[ ] Empirical evidence of audience composition — no analytics/audience data collected;
    if this ever changes (e.g., adding analytics), redo this line item. Conclusion: __________
[ ] Marketing or promotional materials/plans (2025 addition) — describe current marketing
    (if any) and confirm it targets adults. Conclusion: __________
[ ] Representations made to consumers/third parties (2025 addition) — ToS/Privacy Policy
    state 13+/16+/18+ tiered minimum age and "not directed to children under 13."
    Conclusion: __________
[ ] User/third-party reviews (2025 addition) — if reviews exist anywhere (app store,
    social, forums), note whether any reviewer identifies as a minor. Conclusion: __________
[ ] Age of users on comparable sites/services (2025 addition) — comparable products
    (social-casino apps, craps strategy content) skew adult. Conclusion: __________
OVERALL CONCLUSION: [ ] Not directed to children under 13   [ ] Reassess — flag raised

SECTION C — ACTUAL KNOWLEDGE CHECK
[ ] Confirmed: no signup field requests birthdate/age directly (only the neutral age-gate
    attestation, which does not store a specific birthdate)
[ ] Confirmed: no support ticket, email, or user report has ever identified a user as under
    the applicable minimum age. If one ever does: [ ] deletion completed within
    __________ (date) and logged here.
[ ] Confirmed: the age gate is neutral — it does not pre-select a passing age, does not use
    dark patterns to encourage false entry, and blocks/warns rather than silently proceeding
    on a failing answer.

SECTION D — NO BEHAVIORAL ADVERTISING / NO DATA MONETIZATION (the core harm the 2025
amendments target)
[ ] Confirmed: zero ad networks, zero ad SDKs, zero pixel trackers, zero analytics-for-ads
    integrations present in the codebase as of the date of this assessment.
[ ] Confirmed: game history and account data are used only to render the user's own history
    back to them — never aggregated for ad targeting, never sold, never shared with a data
    broker.
[ ] Confirmed: Privacy Policy vendor table is current and accurate as of this assessment date.

SECTION E — DELETION MECHANISM
[ ] Confirmed: in-app "Delete account" button provides immediate, permanent deletion of
    email, passkey credentials, and game history (verify by testing on this date: ______).
[ ] Confirmed: email-request deletion path (privacy@moosequest.net, 30-day SLA) is staffed
    and monitored.

SECTION F — SIGN-OFF
This self-assessment was completed in good faith by the operator on the date above, based
on the product's actual code, data flows, and policies as they existed on that date. It
should be redone whenever: (a) the product adds ads, analytics, or any new data collection;
(b) the product's subject matter or visual design changes materially; (c) any user report
suggests an under-13 (or jurisdiction-applicable-age) user; or (d) at minimum, annually.

Signed: _______________________   Date: __________________
Next scheduled re-assessment date: __________________
```

Save the completed checklist as a dated file (e.g. `docs/legal/coppa-self-assessment-YYYY-MM-DD.md`)
each time it's redone, so there's a running, dated compliance history rather than one static
document that silently goes stale.

---

## 3. Other open items — DIY-resolved

### 3a. Governing-law state (Terms §12)

**Resolved. MooseQuest LLC is organized in Pennsylvania.**

- PA Certificate of Organization filed 2026-05-21, state approval confirmed 2026-05-22, via
  Northwest Registered Agent (order/invoice `2WB4E2F6`). Internal source: the operator's own
  business-formation records (`llc-post-formation-questions-2026-05-25.md` and
  `2026-05-13-northwest-order-record.md` in the MooseQuest business-docs history for the
  companion Raxx/TradeMasterAPI project — same MooseQuest LLC entity).
- EIN issued 2026-05-22 (same records).
- Free public verification available anytime at the PA Department of State's PENN File /
  Business One-Stop Hub entity search — search "MooseQuest LLC" to confirm active status
  before finalizing: [pa.gov Business One-Stop Hub](https://www.pa.gov/agencies/dos/programs/business/).
- **Recommended §12 language:** replace "the state in which MooseQuest LLC is organized"
  with **"the Commonwealth of Pennsylvania"**. No attorney is needed to make this factual
  substitution — it is simply naming the state where the LLC's own Certificate of
  Organization was filed. (This document does not edit `terms.md` directly; the operator or
  a future task can paste this in.)
- One nuance worth a single sentence, not a blocker: Pennsylvania courts generally respect a
  contractually chosen governing-law clause that matches the entity's actual state of
  organization — this is the least-contestable version of a governing-law clause (compare to
  picking an unrelated state with no connection to the parties, which is more easily
  challenged). Naming PA here is the conservative, defensible choice for a PA LLC.

### 3b. GDPR Article 8 — is the 16+ EU/EEA gate a sufficient good-faith simplification?

**Yes, and here's why in one paragraph:** GDPR Art. 8 sets the digital age of consent at 16
by default, letting member states lower it (not below 13); 9 EU/EEA member states set a
threshold above 13 (14, 15, or 16). Rather than build 27 different jurisdiction-specific
gates or parental-consent flows, setting the EU/EEA gate at the **ceiling** (16) satisfies
every member state's threshold simultaneously — a 16-year-old clears the bar everywhere in
the EU/EEA, full stop, regardless of which of the 9 higher-threshold states applies or which
16 states use the 13-year floor. This is the standard "set to the strictest applicable
number" simplification pattern, and it eliminates the nuanced (and genuinely unresolved
without an attorney) question of whether Art. 6(1)(b) "performance of contract" as the
stated legal basis sidesteps Art. 8's consent-specific age threshold — because at 16+, that
question never has to be answered. Source: [Art. 8 GDPR full text](https://gdpr-info.eu/art-8-gdpr/);
[EuConsent.eu member-state digital-consent-age table](https://euconsent.eu/digital-age-of-consent-under-the-gdpr/).
**Residual, genuinely un-DIY-able nuance:** if the operator ever wants to *lower* the EU gate
below 16 to capture the 13-15 cohort in the 17 lower-threshold states, that requires an actual
attorney opinion on the Art. 6(1)(b)-vs-Art. 8 question. Don't do that without counsel — but
there's no product reason to; the 16+ gate is already shipped and costs nothing to keep.

### 3c. UK Age Appropriate Design Code (AADC / Children's Code) — real obligation or monitor-only?

**Monitor-only for this app, with a documented reason why.** The AADC applies to services
"likely to be accessed by children" in the UK — a broad standard that has caught general-
audience gaming and social products in ICO guidance and enforcement history. But:

- The AADC's own published guidance and enforcement focus (per ICO's game-design guidance and
  its "top tips for game designers") targets products with mechanics that *specifically*
  interact with likely-child audiences — loot boxes, in-game purchases marketed to minors,
  behavioral nudges, streaks/rewards aimed at engagement-maximizing for young users. Source:
  [ICO — Age appropriate design: a code of practice for online services](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/childrens-information/childrens-code-guidance-and-resources/age-appropriate-design-a-code-of-practice-for-online-services/);
  [DLA Piper — New ICO guidance on Children's Code for games companies](https://mse.dlapiper.com/post/102i83a/new-ico-guidance-on-childrens-code-for-games-companies-identify-under-18s-and-p).
  Dark Side Craps has: no loot boxes, no purchases of any kind (nothing is for sale), no
  behavioral-engagement mechanics beyond a strategy tracker, no ads, adult-only subject
  matter, and the app's own 13+ base / 16+ EU / 18+ AU age gate applied globally (the UK gate
  is 13+, matching UK GDPR's own floor).
  **Because the UK gate already sits at 13+ and the product has zero monetization or
  engagement-nudge mechanics, the "likely to be accessed by children" argument is weak, and
  the product design already happens to satisfy the AADC's substantive asks (no ads to
  minors, no manipulative engagement design, no purchases) even without a formal AADC
  compliance program.**
- **DIY action:** note this reasoning in the self-assessment file (Section 2 above can be
  extended with a UK-specific line) and monitor the ICO's enforcement priorities annually —
  [ICO priorities announcement](https://ico.org.uk/about-the-ico/media-centre/news-and-blogs/2024/04/ico-sets-out-priorities-to-protect-childrens-privacy-online/).
  Revisit if the ICO issues gambling/social-casino-specific AADC guidance, or if the app
  adds any monetization, ads, or engagement-nudge features.
- **Residual risk:** low but not theoretically zero — the "likely accessed" standard is
  broader than actual-knowledge, so a UK regulator could in principle open an inquiry. No
  documented ICO enforcement action exists against a free, no-money, no-ads casino-strategy
  tool. This is a defensible monitor-only posture for a hobby app at this traffic scale, not
  a certainty; an attorney opinion would remove the residual ambiguity but is not proportionate
  to the risk at this scale.

### 3d. Brazil LGPD / Digital ECA — real obligation, or monitor-only?

**Real, but low-priority given current scale — monitor, and consider blocking if Brazilian
traffic ever becomes material.** Update since the original briefing: Brazil's Digital ECA
(Law 15,211/2025) was signed and published September 17, 2025, and **took effect March 17,
2026** — it is now in force, not a pending bill. It layers child-specific obligations (age
verification, default privacy settings, restrictions on profiling/harmful content targeting
minors) on top of the existing LGPD (which itself requires parental consent for processing
data of children under 12, LGPD Art. 14). A staged enforcement rollout continues through
2026, with a second enforcement phase beginning August 2026. Penalties are steep in the
abstract (up to R$50M or 10% of Brazil revenue) but are calibrated to revenue — a zero-
revenue hobby app's realistic exposure under the revenue-percentage prong is nominal, though
the flat-fee prong is not entirely dollar-linked to revenue and deserves monitoring.
Source: [Baker McKenzie — Brazil Regulates the Digital ECA](https://www.bakermckenzie.com/en/insight/publications/2026/03/brazil-regulates-the-children-and-adolescents-online-safety-act);
[INPLP — The Digital ECA: Brazil's New Age Verification Framework and Enforcement Timeline](https://inplp.com/latest-news/article/the-digital-eca-brazils-new-age-verification-framework-and-enforcement-timeline/);
[LGPD Art. 14](https://lgpd-brazil.info/chapter_02/article_14).

**DIY action:** Brazil is not currently on the geo-block or borderline list (gambling-themed
content is not the trigger here — child-protection law is). No action needed today given the
app's global 13+/16+/18+ age-gate already functions as a basic age-assurance mechanism
consistent with the Digital ECA's general direction. Add a monitoring line to the annual
self-assessment: "check whether Brazilian regulators have issued Digital ECA guidance
specific to casino-themed/simulated-gambling content aimed at minors" and reassess if
Brazilian user volume becomes non-trivial (e.g., if analytics are ever added).

### 3e. The 9 borderline geo-block countries — final DIY call

The original briefing left these nine for the operator to decide. Final, decisive DIY
recommendation for each, given the "minimize unresolved items" goal:

| Code | Country | Final call | One-line reason |
|---|---|---|---|
| SG | Singapore | **Do NOT block.** | Remote Gambling Act 2014 explicitly excludes no-stakes games; IMDA's own clarification note confirms this. This is the one country on the list with an affirmative statutory carve-out, not just "limited enforcement." Blocking it would give up real, law-supported access for no benefit. Source: [IMDA clarification note](https://www.imda.gov.sg/-/media/imda/files/regulation-licensing-and-consultations/codes-of-practice-and-guidelines/acts-codes/14-remote-gambling-act-clarification-note-on-the-scope-of-games.pdf). |
| AE | UAE | **Block.** | Islamic-law gambling ban with no clear carve-out for simulated content; GCGRA (2024) covers only a single licensed physical resort, not online/simulated play. Conservative default holds. |
| BH | Bahrain | **Block.** | Same Sharia-based rationale as AE/OM; no simulated-content carve-out documented. |
| CU | Cuba | **Block.** | State-controlled internet means practical traffic is near zero either way; blocking costs nothing and removes the ambiguous-law question entirely. |
| JO | Jordan | **Block.** | All gambling illegal; no simulated-content carve-out found; low traffic impact either way. |
| OM | Oman | **Block.** | Same rationale as AE/BH. |
| TM | Turkmenistan | **Block.** | Gambling banned, highly restricted internet; blocking is low-cost and removes ambiguity. |
| TN | Tunisia | **Block.** | Lower enforcement risk than the others, but no affirmative carve-out (unlike Singapore) — the conservative default is more defensible than carving out an exception on ambiguous law. |
| UZ | Uzbekistan | **Block.** | Gambling broadly restricted; no simulated-content carve-out found. |

**Net effect: add AE, BH, CU, JO, OM, TM, TN, UZ to `BLOCKED_COUNTRIES`; explicitly do NOT
add SG.** This resolves all 9 borderline items with a documented one-line rationale each —
the only one with an affirmative "don't block" case is Singapore, and that case rests on
statutory text, not just low enforcement odds.

---

## 4. Bottom line — ranked DIY action list

**This week (do these, no professional needed):**

1. **Complete the COPPA self-assessment checklist in Section 2**, dated and saved as
   `docs/legal/coppa-self-assessment-YYYY-MM-DD.md`. This is the actual "audit" — done, for
   free, by the operator, on the record. ~20 minutes.
2. **Update Terms §12** to name **"the Commonwealth of Pennsylvania"** as the governing law
   (text provided in Section 3a). Factual substitution, no legal judgment required — PA is
   MooseQuest LLC's actual state of formation.
3. **Add the 8 remaining borderline countries to `BLOCKED_COUNTRIES`** (AE, BH, CU, JO, OM,
   TM, TN, UZ) and confirm SG is explicitly excluded from the block list. Engineering task,
   file as a PM card.
4. **Add one line to the Privacy Policy's UK section (or a new short UK-specific line)**
   documenting the AADC monitor-only reasoning from Section 3c, so the good-faith record
   exists in the public-facing document, not just this internal research file.

**This quarter (low effort, keeps the record current):**

5. **Re-run the self-assessment** if any of these happen: ads/analytics added, a user report
   claims an under-13 (or jurisdiction-applicable) user, or the product's visual design/
   subject matter changes materially.
6. **Check the ICO's annual children's-privacy enforcement priorities** and Brazil's Digital
   ECA rulemaking progress (second enforcement phase begins August 2026) — both take under
   30 minutes to scan once a year and directly inform whether the monitor-only calls in
   Sections 3c/3d still hold.

**What genuinely can't be de-risked without a lawyer (kept deliberately short):**

- **The Art. 6(1)(b)-vs-Art. 8 legal-basis question**, if the operator ever wants to *lower*
  the EU gate below 16 — not needed today since the 16+ gate already moots this question, but
  flagged so it's not silently forgotten if EU expansion is ever considered.
- **A formal opinion on whether Australia's Classification Act R18+ requirement is fully
  satisfied by a self-attestation gate alone** (versus a more robust age-verification
  method) — the current 18+ self-attestation gate is a reasonable, industry-consistent
  position, but no attorney has opined on whether the Australian Classification Board would
  consider self-attestation sufficient versus a stronger verification method. Given zero
  known enforcement actions against comparable free hobby apps, this is a low-priority,
  accept-the-residual-risk item, not a blocker.
- **Any actual regulator inquiry or user complaint**, should one ever arrive — that is the
  one scenario where a DIY posture stops being sufficient and retaining counsel (even a
  single paid consult) becomes the right call, specifically because at that point the
  question shifts from "good-faith prospective posture" to "defending a specific claim,"
  which is a different and higher-stakes exercise.

---

## Jurisdiction flags

| Jurisdiction | Status after this doc | Residual risk |
|---|---|---|
| US federal (COPPA) | Not covered by COPPA's threshold trigger; self-assessment documents why | Low — no enforcement pattern against comparable apps |
| Pennsylvania (governing law) | Resolved — PA is MooseQuest LLC's actual formation state | None — factual, not judgment-based |
| EU/EEA (GDPR Art. 8) | 16+ gate already shipped; satisfies all 27 member states | Low — only relevant if EU gate is ever lowered |
| UK (AADC) | Monitor-only; documented reasoning added to Privacy Policy | Low-moderate — "likely accessed" standard is broad but product design already satisfies substance |
| Brazil (LGPD / Digital ECA) | Monitor-only; now-effective law but low realistic exposure at current scale | Low at current traffic; revisit if Brazilian traffic grows |
| Australia (Classification Act) | Already shipped 18+ gate; self-attestation adequacy is the one unresolved nuance | Low — no attorney opinion obtained, but no enforcement pattern exists either |
| 9 borderline geo-block countries | Resolved — 8 blocked, Singapore explicitly not blocked | None remaining — decisive call made |

## Timing / deadlines

| Deadline | Item | Status |
|---|---|---|
| Already passed (April 22, 2026) | COPPA 2025 amendment compliance | Not applicable to this app — see Section 1 |
| Ongoing | Annual self-assessment re-run | Schedule the first one this week; recur annually or on any material product change |
| August 2026 | Brazil Digital ECA — second enforcement phase begins | Monitor only; no action required at current traffic |
| Ongoing | ICO children's-privacy enforcement priorities | Check annually |

## Sources

- FTC, [FTC Finalizes Changes to Children's Privacy Rule (Jan 2025)](https://www.ftc.gov/news-events/news/press-releases/2025/01/ftc-finalizes-changes-childrens-privacy-rule-limiting-companies-ability-monetize-kids-data)
- Federal Register, [COPPA Final Rule Amendments (Apr 22, 2025)](https://www.federalregister.gov/documents/2025/04/22/2025-05904/childrens-online-privacy-protection-rule)
- eCFR, [16 CFR Part 312](https://www.ecfr.gov/current/title-16/chapter-I/subchapter-C/part-312)
- Securiti, [FTC's 2025 COPPA Final Rule Amendments](https://securiti.ai/ftc-coppa-final-rule-amendments/)
- Latham & Watkins, [FTC Publishes Updates to COPPA Rule](https://www.lw.com/en/insights/ftc-publishes-updates-to-coppa-rule)
- FTC, [Enforcement Policy Statement Promoting the Adoption of Age-Verification Technology (Feb 2026)](https://www.ftc.gov/legal-library/browse/enforcement-policy-statement-promoting-adoption-age-verification-technology)
- FTC, [FTC Issues COPPA Policy Statement to Incentivize the Use of Age Verification Technologies (Feb 2026)](https://www.ftc.gov/news-events/news/press-releases/2026/02/ftc-issues-coppa-policy-statement-incentivize-use-age-verification-technologies-protect-children)
- 15 U.S.C. § 6504, [Actions by States (Cornell LII)](https://www.law.cornell.edu/uscode/text/15/6504)
- Davis Polk, [FTC Prioritizes COPPA Enforcement](https://www.davispolk.com/insights/client-update/ftc-prioritizes-coppa-enforcement-new-compliance-obligations-take-effect)
- GDPR, [Article 8 full text](https://gdpr-info.eu/art-8-gdpr/)
- EuConsent.eu, [Digital Age of Consent Under the GDPR](https://euconsent.eu/digital-age-of-consent-under-the-gdpr/)
- ICO, [Age appropriate design: a code of practice for online services](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/childrens-information/childrens-code-guidance-and-resources/age-appropriate-design-a-code-of-practice-for-online-services/)
- ICO, [ICO Sets Out Priorities to Protect Children's Privacy Online (2024)](https://ico.org.uk/about-the-ico/media-centre/news-and-blogs/2024/04/ico-sets-out-priorities-to-protect-childrens-privacy-online/)
- DLA Piper, [New ICO guidance on Children's Code for games companies](https://mse.dlapiper.com/post/102i83a/new-ico-guidance-on-childrens-code-for-games-companies-identify-under-18s-and-p)
- Baker McKenzie, [Brazil Regulates the Digital ECA (Mar 2026)](https://www.bakermckenzie.com/en/insight/publications/2026/03/brazil-regulates-the-children-and-adolescents-online-safety-act)
- INPLP, [The Digital ECA: Brazil's New Age Verification Framework and Enforcement Timeline](https://inplp.com/latest-news/article/the-digital-eca-brazils-new-age-verification-framework-and-enforcement-timeline/)
- [LGPD Article 14](https://lgpd-brazil.info/chapter_02/article_14)
- Singapore IMDA, [Remote Gambling Act Clarification Note on Scope of Games (Jan 2015)](https://www.imda.gov.sg/-/media/imda/files/regulation-licensing-and-consultations/codes-of-practice-and-guidelines/acts-codes/14-remote-gambling-act-clarification-note-on-the-scope-of-games.pdf)
- Internal: MooseQuest LLC formation records (`llc-post-formation-questions-2026-05-25.md`,
  `2026-05-13-northwest-order-record.md` — companion-project business-docs history;
  Northwest Registered Agent order confirmation, PA state approval, EIN issuance)
- PA Department of State, [Business One-Stop Hub — entity search](https://www.pa.gov/agencies/dos/programs/business/) (use to independently re-verify MooseQuest LLC's active PA status at any time)
- Prior briefing: `docs/legal/compliance-research.md` (full geo-block country-by-country sourcing, carried forward by reference, not repeated here)

---

> **This document is a good-faith, self-executed compliance record prepared without a
> licensed attorney, per the operator's explicit direction. It is research and preparation
> material, not legal advice, and creates no attorney-client relationship or legal opinion.
> If a regulator inquiry, user complaint, or material business change (real money, ads,
> material EU/Brazil/Australia scale) ever occurs, that is the point to engage a licensed
> attorney (privacy/COPPA specialist for the US items; a UK/EU data-protection solicitor for
> AADC/GDPR items) — the DIY posture in this document is calibrated for a free, no-monetization
> hobby app at low traffic, not for a scaled or monetized product.**
