# Dark Side Craps — Compliance Research Briefing

> **Status: research-only. This document does NOT constitute legal or tax advice.**
> Before filing or acting on anything in this document, consult a licensed attorney in the relevant
> jurisdictions. Sources were verified as of 2026-07-20. Confirm freshness before relying on them.
>
> Prepared for: MooseQuest LLC. Product: Dark Side Craps (https://craps.moosequest.app).
> Product description: free, browser-only, no-money, no-prize craps strategy tracker/simulator.
> Data collected: email address (account identifier + one-time code), WebAuthn public key, game history.

---

## Machine-readable summary

```
RECOMMENDED_MIN_AGE: 13
AGE_RATIONALE: No US statute mandates 18+ for free, no-money, no-prize web apps; COPPA
  "not-directed-to-children" posture is defensible given adult casino subject matter; app-store
  17+ rules are platform policy, not law, and apply only to store-distributed apps; required
  overrides: 18+ gate for Australian users (Classification Act R18+), and 16+ gate (or parental
  consent workflow) for EU users in the 9 member states whose digital consent age exceeds 13.

GEO_BLOCK_ISO: AF,DZ,BD,BN,CN,ID,IQ,IR,KP,KW,LY,MY,PK,QA,SA,SD,SY,TH,VN,YE
GEO_BLOCK_RATIONALE: These jurisdictions either (a) explicitly ban simulated/casino-themed
  content regardless of real-money involvement (VN, TH, KW's advertising prohibition), or (b)
  have comprehensive Islamic-law gambling bans with active website-blocking infrastructure and
  documented regulatory positions that even simulated gambling violates the prohibition (SA, AE
  legal commentary, ID, PK, BD, MY). CN is included because all gambling—including simulated—is
  banned and enforcement is aggressive. KP is included for completeness (no civilian internet
  access). Only jurisdictions whose bans clearly extend beyond real-money gambling appear here.

BORDERLINE_ISO: AE,BH,CU,JO,OM,SG,TM,TN,UZ
BORDERLINE_RATIONALE: AE has a new GCGRA regulatory framework (2024) suggesting possible
  future licensing paths; BH/OM have Islamic bans but limited online enforcement; SG's Remote
  Gambling Act 2014 explicitly excludes no-stakes games (key quote below); JO is not
  aggressively policed online; TN doesn't actively block; TM/UZ have bans but unclear reach;
  CU has state-controlled internet with ambiguous gambling rules. Operator decides these nine.
```

---

## TL;DR (3 sentences)

A "not directed to children under 13" posture is defensible for COPPA purposes given the adult casino-strategy subject matter, but implementing 13+ globally requires an 18+ gate for Australian users (hard legal requirement since September 2024) and either a 16+ gate or a parental-consent flow for users in 9 EU member states with higher digital consent ages. Approximately 20 countries warrant GeoIP blocking because their laws ban gambling-themed content even without real money; Singapore is the notable exception, as its Remote Gambling Act explicitly excludes no-stakes/no-prize games. The existing Terms and Privacy Policy have strong "not gambling" and liability language, but contain four must-fix gaps: no age gate mechanism on the site, an internal age inconsistency between the two docs, no COPPA citation by name, and no prohibited-territory representation clause.

---

## Question 1: Minimum age — is 13+ defensible?

### 1a. US COPPA analysis

**What COPPA governs.** The Children's Online Privacy Protection Act (15 U.S.C. §§ 6501–6506)
and its implementing rule (16 CFR Part 312, amended effective June 23, 2025; operator compliance
deadline April 22, 2026) apply to operators of websites or online services that are either
(a) "directed to children under 13," or (b) have "actual knowledge" they are collecting personal
information from a child under 13.
- Source: FTC, [Children's Online Privacy Protection Rule](https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa)
- Source: Federal Register, [Final COPPA Rule Amendments (Apr 22, 2025)](https://www.federalregister.gov/documents/2025/04/22/2025-05904/childrens-online-privacy-protection-rule)
- Source: FTC press release, [FTC Finalizes Changes to Children's Privacy Rule (Jan 2025)](https://www.ftc.gov/news-events/news/press-releases/2025/01/ftc-finalizes-changes-childrens-privacy-rule-limiting-companies-ability-monetize-kids-data)

**"Directed to children" test.** Under 16 CFR § 312.2, the Commission evaluates a multi-factor
totality-of-circumstances test. The original eight factors are:
1. Subject matter
2. Visual content
3. Use of animated characters or child-oriented activities
4. Music or other audio content
5. Age of models
6. Presence of child celebrities or celebrities who appeal to children
7. Advertising on or promoting the service
8. Empirical evidence regarding audience composition

The 2025 amendments added four more factors:
9. Marketing or promotional materials and plans
10. Representations made to consumers or third parties
11. User and third-party reviews
12. Age of users on similar websites or services

No single factor is dispositive; it is a totality test.
- Source: 16 CFR § 312.2 via [eCFR](https://www.ecfr.gov/current/title-16/chapter-I/subchapter-C/part-312)
- Source: Hintze Law analysis, [Final COPPA Rule Amendments — Definitional Changes](https://hintzelaw.com/blog/2025/2/6/final-coppa-rule-amendments-definitional-changes)
- Source: FTC, [Complying with COPPA: Frequently Asked Questions](https://www.ftc.gov/business-guidance/resources/complying-coppa-frequently-asked-questions)

**Application to Dark Side Craps.** Applying the test to this app:

| Factor | Dark Side Craps assessment |
|---|---|
| Subject matter | Adult casino gambling strategy — strongly NOT directed at children |
| Visual content | No cartoon characters, bright toy colors, or child-appealing imagery |
| Animated characters | None |
| Music / audio | None noted |
| Age of models | None |
| Child celebrities | None |
| Advertising | None (no ads on the service) |
| Audience composition empirical | Casino gaming enthusiasts — strongly adult demographic |
| Marketing / promotional | Not targeted to minors |
| User reviews | Not yet known; monitor if reviews accumulate from minors |

Assessment: A craps strategy tracker whose entire premise is adult casino gambling is among the
strongest possible fact patterns for "not directed to children." Contrast with FTC enforcement
actions against Musical.ly (video lip-sync platform) and Kidoz (child-targeted ad network), where
the subject matter and visual design clearly targeted children. Casino strategy is an adult
activity; the FTC has not brought an action against a gambling-strategy educational web app.

**"Not directed to children" safe harbor mechanism.** An operator not directed to children under
13 who also has no "actual knowledge" of collecting data from under-13 users is not subject to
COPPA's verifiable-parental-consent requirement. A neutral age screen (that does not default to a
passing age or encourage false entry) plus "not directed to children under 13" language in the
Privacy Policy creates the operational record supporting this posture. The 2025 amendments require
age mechanisms be "reasonably calculated, in light of available technology" and "neutral."
- Source: Hintze Law, [Final COPPA Rule Amendments](https://hintzelaw.com/blog/2025/2/6/final-coppa-rule-amendments-definitional-changes)

**Is 13+ defensible under COPPA?** Yes. The "not directed to children under 13" posture is
defensible given the adult subject matter. A 13+ minimum age in the ToS, combined with a neutral
age attestation gate and a "not directed to children under 13" Privacy Policy statement, satisfies
COPPA compliance obligations without requiring verifiable parental consent. Lowering the ToS floor
from 18 to 13 does not create COPPA risk, because COPPA's threshold is 13 (not 18); operators
above 13 are outside COPPA's direct consent requirement entirely.

**COPPA compliance deadline notice.** The April 22, 2026 compliance deadline has now passed
(today is 2026-07-20). If MooseQuest LLC has not audited the site against the 2025 rule
amendments, that audit is overdue. Key new obligations include: new parental consent mechanisms,
stricter data retention limits, and expanded "directed to children" definition. Consult a
privacy/COPPA attorney promptly.
- Source: Hunton Andrews Kurth, [COPPA Rule Amendment Compliance Deadline](https://www.hunton.com/privacy-and-cybersecurity-law-blog/coppa-rule-amendment-compliance-deadline-approaches)
- Source: Davis Polk, [FTC Prioritizes COPPA Enforcement](https://www.davispolk.com/insights/client-update/ftc-prioritizes-coppa-enforcement-new-compliance-obligations-take-effect)

---

### 1b. GDPR / UK GDPR — age of digital consent

**GDPR Article 8 baseline.** Under Article 8 GDPR, where consent is the legal basis, processing a
child's personal data for information society services (ISS) is lawful only when the child is at
least 16 years old — unless a member state lowers this floor, which may not go below 13.
- Source: [Art. 8 GDPR full text via gdpr-info.eu](https://gdpr-info.eu/art-8-gdpr/)

**Member state digital age of consent (complete table as of 2026-07-20):**

| Member State | Digital Consent Age |
|---|---|
| Belgium | 13 |
| Denmark | 13 |
| Estonia | 13 |
| Finland | 13 |
| Latvia | 13 |
| Malta | 13 |
| Portugal | 13 |
| Sweden | 13 |
| Austria | 14 |
| Bulgaria | 14 |
| Cyprus | 14 |
| Italy | 14 |
| Lithuania | 14 |
| Spain | 14 |
| Czech Republic | 15 |
| France | 15 |
| Greece | 15 |
| Slovenia | 15 (proposed; currently 16) |
| Croatia | 16 |
| Germany | 16 |
| Hungary | 16 |
| Ireland | 16 |
| Luxembourg | 16 |
| Netherlands | 16 |
| Poland | 16 |
| Romania | 16 |
| Slovakia | 16 |

**UK GDPR.** Post-Brexit, the UK GDPR sets the digital consent age at 13. The ICO expects
operators to respect this threshold for UK-based users. The ICO's Children's Code (Age Appropriate
Design Code) applies to services "likely to be accessed by children under 18" — this is a broader
code and its applicability to a gambling-strategy simulator marketed to adults is uncertain; flag
for attorney review.
- Source: ICO, [Children and the UK GDPR](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/childrens-information/children-and-the-uk-gdpr/)

**Application to "13+ globally" posture.** Dark Side Craps's privacy.md states that GDPR legal
basis is Art. 6(1)(b) (performance of contract), not consent. If the legal basis is contractual
performance — not consent — then Art. 8's age thresholds may not apply directly (Art. 8 triggers
specifically on consent as the lawful basis). However, this is a nuanced legal argument; the
safer position is to treat the age thresholds as applicable since an account creation / email
collection for a minor is the exact scenario Art. 8 was designed to address. Confirm with an EU
data-protection attorney.

**Practical implication for 13+ global gate.** Setting 13+ globally means:
- For users in the 8 member states with a 13-year threshold: fully compliant.
- For users in Austria, Bulgaria, Cyprus, Italy, Lithuania, Spain (14-year threshold): processing
  personal data of 13-year-olds requires parental consent.
- For users in Czech Republic, France, Greece (15-year threshold): processing for 13–14-year-olds
  requires parental consent.
- For users in Croatia, Germany, Hungary, Ireland, Luxembourg, Netherlands, Poland, Romania,
  Slovakia (16-year threshold): processing for 13–15-year-olds requires parental consent.

**Simplifying options:**
- Option A: Set the EU/EEA gate to **16+** and the rest-of-world gate to 13+. Covers all EU
  member states without parental consent mechanics.
- Option B: Set a single **16+** global gate. Covers all EU states; UK is fine at 13+ (16 is
  more restrictive than required). No parental consent workflows needed anywhere. Loses the 13-15
  cohort globally.
- Option C: Implement jurisdiction-specific parental consent flows for EU states above 13.
  Most complex engineering but allows true 13+ in high-consent-age states.

- Source: [EuConsent.eu — Digital Age of Consent Under the GDPR (full table)](https://euconsent.eu/digital-age-of-consent-under-the-gdpr/)

---

### 1c. Casino-themed content: does it legally mandate 18+ for a free web app?

**App-store policies vs. law — key distinction.** This is the most common source of confusion.

| Policy/Law | Applies to | Threshold | Binding? |
|---|---|---|---|
| Apple App Store guideline 5.3 | Store-distributed iOS/macOS apps | 17+ | Contractual with Apple; does NOT apply to open web |
| Google Play policy | Store-distributed Android apps | 18+ in many markets | Contractual with Google; does NOT apply to open web |
| ESRB / PEGI / IARC ratings | Retail/store-distributed games | Varies by rating | Voluntary ratings system; not US law |
| Australia Classification Act | Computer games (including web-based interactive games) | R18+ for interactive simulated gambling | Statutory; applies to open web in Australia |

Apple changed its App Store age rating for "Frequent/Intense Simulated Gambling" to 17+ in
August 2019 — but this is Apple's platform policy, not a law.
- Source: AppleInsider, [App Store shakeup limits simulated gambling to users aged 17+](https://appleinsider.com/articles/19/08/20/app-store-shakeup-limits-simulated-gambling-to-users-aged-17)

**US law.** No US federal statute and no US state statute specifically mandates a minimum age of
18+ for a free, no-prize, no-money, browser-based simulated gambling or casino-themed web
application. Social casino games operate legally in 49 US states (Washington is the main
exception for sweepstakes mechanics; Idaho has restrictions). Social casino age minimums are
industry practice (typically 18 or 21), not statutory mandates for free educational apps.
- Source: Richt Law, [Proliferation of Social Gaming Casinos: Legal Compliance Considerations](https://richtfirm.com/the-proliferation-of-social-gaming-casinos-legal-compliance-considerations/)
- Source: gambling.com, [How Social Casinos Navigate Legal Constraints in the US](https://www.gambling.com/us/online-casinos/strategy/how-social-casinos-navigate-legal-constraints-in-the-us)

**Australia: the one hard 18+ legal requirement for open web.** Australia's Classification
(Publications, Films and Computer Games) Act, as amended by new Games Guidelines effective
September 22, 2024, mandates a minimum R18+ classification for "video games containing simulated
gambling, such as casino games." The classification authority has confirmed this applies to "all
devices including phones, tablets, consoles and PCs," which legal analysis indicates extends to
web-based interactive games. Critically, the guidelines distinguish: games that "feature casino
settings, imagery or themes but do NOT allow players to engage or interact with gambling
activities" are not captured — but Dark Side Craps, which has an interactive simulated betting
interface and dice-roll engine, almost certainly IS captured. For Australian users, an 18+ age
gate is effectively a legal requirement.
- Source: Australian Classification Board, [New Classifications for Gambling-Like Content in Video Games](https://www.classification.gov.au/about-us/media-and-news/news/new-classifications-for-gambling-content-video-games)
- Source: Australian Classification Board, [Games Guidelines Fact Sheet](https://www.classification.gov.au/sites/default/files/documents/game-guidelines-for-gambling-like-content-fact-sheet.pdf)
- Source: Corrs Chambers Westgarth, [Mandatory Minimum Classifications for Computer Games with Gambling-Like Content](https://www.corrs.com.au/insights/mandatory-minimum-classifications-for-computer-games-with-gambling-like-content)

---

### 1d. Net recommendation on minimum age

**RECOMMENDED_MIN_AGE: 13**

Rationale (one line): No US statute or global law (outside Australia) mandates 18+ for a free,
no-prize, no-money, open-web casino strategy simulator; COPPA is satisfied by the "not directed
to children" posture given adult casino subject matter; 13+ is the COPPA floor and the UK GDPR
floor; app-store 17+ policies are platform contracts, not law.

**Required jurisdiction overrides — these are not optional:**

| Jurisdiction | Required minimum | Legal basis |
|---|---|---|
| Australia (AU) | 18+ | Classification Act + 2024 Games Guidelines (R18+ for interactive simulated gambling) |
| EU member states with 14–16 digital consent age (see table above) | 16+ gate OR parental consent for 13–15-year-olds | GDPR Art. 8 |
| UK | 13+ (no change) | UK GDPR Art. 8; ICO 13 threshold |

**Implementation options ranked by simplicity:**
1. **Keep current 18+ globally.** No additional engineering. Covers Australia, all EU states.
   Loses 13–17 user cohort globally. Simplest legal posture.
2. **Set 16+ globally.** Covers all EU states and Australia's R18+ requirement would still need
   AU-specific 18+ gate. Simple for EU; still needs AU gate.
3. **13+ globally with AU=18+ gate and EU=16+ gate.** Achieves operator's preference. Requires
   GeoIP-driven age gate logic: AU → 18+, EU → 16+, everywhere else → 13+. Medium engineering.
4. **13+ globally with parental consent flows for EU high-threshold states.** Most complex;
   probably not worth the engineering for a hobby app.

**Bottom line for the operator's attorney and CPA:** 13+ is legally defensible for US/UK. But a
mechanism to enforce 18+ for Australian users is a statutory obligation under the Classification
Act, and GDPR parental consent obligations for 13-15 year olds in 9 EU states must be addressed
either by a 16+ EU gate or actual consent flows. The current ToS says 18+, which is legally
conservative and covers all these issues without jurisdiction-specific complexity.

---

## Question 2: Geo-block list

### 2a. The critical analytical distinction

**Block only if the ban covers simulated/themed content with no money.**
A jurisdiction that bans only real-money online gambling does NOT require blocking of a free,
no-money, no-prize educational tool. The UK, EU, Canada, Australia (Classification Act applies
to age, not access), Japan, South Korea, Brazil, and the US regulate real-money gambling but
do not prohibit free simulated gambling web apps. Over-blocking these jurisdictions would be
wrong and harmful to the product.

Block only when: (a) the jurisdiction's law or regulatory interpretation explicitly covers
simulated or casino-themed content, OR (b) the jurisdiction's enforcement apparatus (ISP blocking)
targets casino-themed content regardless of money, AND the statute language is broad enough to
capture it.

---

### 2b. Definite block list — ISO 3166-1 alpha-2

`AF, DZ, BD, BN, CN, ID, IQ, IR, KP, KW, LY, MY, PK, QA, SA, SD, SY, TH, VN, YE`

**Country-by-country rationale:**

| Code | Country | Basis for block | Key source |
|---|---|---|---|
| AF | Afghanistan | Taliban rule; all gambling and gambling-themed entertainment prohibited; virtually no civilian internet access to foreign sites | [playtoday.co — Countries Where Gambling Is Illegal](https://playtoday.co/blog/guides/countries-where-gambling-is-illegal/) |
| DZ | Algeria | All gambling illegal with no exceptions; active ISP-level blocking of gambling-related websites; 2025 law (No. 25-10) criminalized VPN use to circumvent blocks | [CMS Expert Guide — Algeria](https://cms.law/en/int/expert-guides/cms-expert-guide-to-gambling-laws-in-africa/algeria); [lcb.org — Algeria](https://lcb.org/restrictions/algeria) |
| BD | Bangladesh | Prevention of Gambling Act 1867 (broadly worded); all gambling forms banned; authorities block gambling content | [legalpilot.com — Bangladesh](https://legalpilot.com/country/bangladesh/) |
| BN | Brunei | Sharia law; comprehensive gambling prohibition including any simulated or themed gambling activity | [playtoday.co — Countries Where Gambling Is Illegal](https://playtoday.co/blog/guides/countries-where-gambling-is-illegal/) |
| CN | China | All gambling including simulated prohibited; 4,500+ platforms destroyed in 2024; GFW blocks essentially all foreign casino-related content regardless of real-money involvement; Consumer Council has called for legislation targeting simulated gambling games | [GGRAsia — China 2024 crackdown](https://www.ggrasia.com/china-destroyed-in-2024-more-than-4500-illegal-online-gambling-platforms-ministry-of-security); [Wikipedia — Gambling in China](https://en.wikipedia.org/wiki/Gambling_in_China) |
| ID | Indonesia | All gambling prohibited; government actively blocks gambling websites; Islam is predominant; no licensing framework for any casino-related content | [legalpilot.com — Indonesia](https://legalpilot.com/country/indonesia/) |
| IQ | Iraq | Islamic law; all gambling prohibited; active internet censorship including gambling content | [arabianbetting.com — Country Guides](https://arabianbetting.com/en/country/) |
| IR | Iran | All gambling prohibited under Islamic law; government actively blocks gambling and gaming content including casino-themed material | [playtoday.co — Countries Where Gambling Is Illegal](https://playtoday.co/blog/guides/countries-where-gambling-is-illegal/) |
| KP | North Korea | Gambling banned; civilian internet access to foreign sites essentially non-existent (Kwangmyong domestic network only); included for completeness | [Wikipedia — Internet in North Korea](https://en.wikipedia.org/wiki/Internet_in_North_Korea) |
| KW | Kuwait | All gambling illegal; critically, "displaying gambling ads and publishing material related to gambling is a punishable offense" — language broad enough to cover casino-themed educational content | [arabiancasino.org — Kuwait](https://arabiancasino.org/countries/) |
| LY | Libya | Penal Code bans all gambling and lottery games; active ISP-level blocking of gambling-related websites | [lcb.org — Libya](https://lcb.org/restrictions/libya) |
| MY | Malaysia | Betting Act 1953; Communications and Multimedia Act s.263 empowers MCMC to direct ISPs to block unlawful content including gambling; removed 457,562 pieces of gambling content and blocked 1,778 sites Jan 2025–May 2026 | [AGB Brief — Malaysia](https://agbrief.com/news/malaysia/13/08/2024/majority-of-blocked-websites-in-malaysia-linked-to-online-gambling/); [Bright Side of News — Malaysia 2026](https://brightsideofnews.com/gambling/regulation/malaysia-2026-anti-online-gambling-reform/) |
| PK | Pakistan | Prevention of Gambling Act 1977 (all gambling forms prohibited); online gambling actively blocked; broadly interpreted to include digital content | [legalpilot.com — Pakistan](https://legalpilot.com/country/pakistan/) |
| QA | Qatar | All gambling banned under Islamic law; active enforcement; no exceptions for simulated play | [lcb.org — Qatar](https://lcb.org/restrictions/qatar) |
| SA | Saudi Arabia | All gambling prohibited under Sharia; ISPs required to block gambling-related websites; legal commentary explicitly states even simulated gambling violates Islamic principles; GCGRA established 2024 but only for licensed physical resort (Ras Al Khaimah) — no online licensing | [sumsub.com — UAE Gaming Regulations](https://sumsub.com/blog/uae-gaming-regulations-all-you-need-to-know/); [aladllegal.com — UAE Gambling Laws](https://aladllegal.com/news/uaes-gambling-laws-prohibition-consequences-and-ex/) |
| SD | Sudan | Islamic law; all gambling banned | [playtoday.co — Countries Where Gambling Is Illegal](https://playtoday.co/blog/guides/countries-where-gambling-is-illegal/) |
| SY | Syria | Gambling banned; active conflict; internet infrastructure severely restricted | [arabianbetting.com — Country Guides](https://arabianbetting.com/en/country/) |
| TH | Thailand | Gambling Act 1935 (comprehensive prohibition including casino-style games); legal commentary confirms simulated/casino-style games "without real money remain prohibited"; 717K+ URLs blocked in 2024 crackdown | [Wikipedia — Gambling in Thailand](https://en.wikipedia.org/wiki/Gambling_in_Thailand); [Yogonet — Thailand 717K blocks](https://www.yogonet.com/international/news/2024/04/03/71535-thailand-steps-up-efforts-against-online-gambling-blocks-25-000-websites) |
| VN | Vietnam | Decree 147/2024/ND-CP (effective December 25, 2024) explicitly bans casino-style and card-based games regardless of whether real money is involved; ISPs mandated to block non-compliant cross-border services | [Yogonet — Vietnam Decree 147](https://www.yogonet.com/international/news/2024/12/05/87586-vietnam-39s-new-decree-147-enforces-stricter-regulations-on-online-gaming-industry); [vietnamnet.vn — Casino-style games banned](https://vietnamnet.vn/en/vietnam-tightens-online-gaming-rules-bans-casino-style-and-card-based-games-2346909.html) |
| YE | Yemen | Islamic law; all gambling prohibited; civil conflict; limited internet infrastructure | [arabianbetting.com — Country Guides](https://arabianbetting.com/en/country/) |

---

### 2c. Borderline / operator decides

`AE, BH, CU, JO, OM, SG, TM, TN, UZ`

| Code | Country | Why borderline | Recommendation |
|---|---|---|---|
| AE | UAE | All gambling prohibited; legal commentary says simulated gambling violates Islamic principles; BUT the GCGRA (2024) introduced a new commercial gaming licensing framework for specific venues, suggesting regulatory modernization; unclear if a free educational web tool would be prosecuted | Conservative: block. If AE revenue is material: seek UAE legal counsel. |
| BH | Bahrain | Gambling illegal under Sharia; limited exceptions for tourism casinos; online gambling not legalized; enforcement of foreign websites is inconsistent | Conservative: block. |
| CU | Cuba | State-controlled internet; gambling law unclear for foreign educational tools; most Cubans cannot access international internet | Practical effect of blocking is minimal; conservative: block. |
| JO | Jordan | All gambling illegal; Jordan allows some controlled casino gambling for tourists; online gambling not aggressively policed; unclear for simulated educational content | Borderline; conservative: block. |
| OM | Oman | Gambling prohibited under Islamic law; limited online enforcement; no documented ISP-level blocking of educational tools | Conservative: block. |
| SG | Singapore | Remote Gambling Act 2014 (No. 34 of 2014) **explicitly excludes** games where players cannot "receive money or money's worth consequent to the outcome of that game" — free craps simulator is textually outside RGA scope; IMDA clarification note (Jan 2015) confirms social games without money payouts are not regulated; HOWEVER, government is conservative and this is an area of "regulatory concern" | Do NOT block based on current law; monitor for regulatory changes. SG is a notable safe harbor from this list. |
| TM | Turkmenistan | Gambling banned; highly restricted internet; limited foreign web access generally | Practical effect minimal; conservative: block. |
| TN | Tunisia | Online gambling illegal; casino gambling legal for foreigners with passport; government does not actively block foreign gambling websites for Tunisian users; ambiguous for simulated content | Lower risk than other Islamic-law states; operator decides. |
| UZ | Uzbekistan | Gambling broadly restricted; unclear scope for free no-money educational tools; limited enforcement of foreign web content | Unclear; conservative: block. |

**Jurisdictions specifically NOT on the block or borderline lists (real-money only regulation):**

The following regulate real-money gambling only and do not prohibit free, no-prize, no-money
simulated gambling web apps: United States (all states), United Kingdom, all EU/EEA member states,
Canada, Australia (Classification Act creates an AGE requirement, not a service block), Japan,
South Korea, Brazil, Mexico, New Zealand, South Africa, and most of Latin America and Southeast
Asia not already listed.

- Singapore Remote Gambling Act source: [IMDA Clarification Note — Scope of Games (Jan 2015)](https://www.imda.gov.sg/-/media/imda/files/regulation-licensing-and-consultations/codes-of-practice-and-guidelines/acts-codes/14-remote-gambling-act-clarification-note-on-the-scope-of-games.pdf)
- Source: [sso.agc.gov.sg — Remote Gambling Act 2014](https://sso.agc.gov.sg/Act-Rev/RGA2014/Published/20211231?DocDate=20211231)

**Australia — age gate, not geo-block.** Australia requires 18+ for interactive simulated gambling
under the Classification Act (September 2024), but there is no statutory prohibition on serving
the content at all. Australian users should receive an 18+ age gate, not a geo-block. If the
operator does not implement this, they risk classification non-compliance and potential ACMA
attention.
- Source: [Australian Classification Board — New Classifications for Gambling-Like Content](https://www.classification.gov.au/about-us/media-and-news/news/new-classifications-for-gambling-content-video-games)

---

## Question 3: Current Terms/Privacy coverage assessment

### 3a. What the current documents do well

Reviewing the live Terms of Use (effective July 16, 2026) and Privacy Policy (effective July 16,
2026) at https://craps.moosequest.app:

**Strengths:**
- "This is NOT gambling. There is no real money involved — no wagering, no deposits, no payouts,
  no prizes, no financial transactions of any kind." — explicit and prominent; strong.
- "Do not use this app at a casino, during real craps play, or to make real-money betting
  decisions." — addresses the core misuse risk.
- AS-IS / AS-AVAILABLE warranty disclaimer in caps — standard and adequate.
- Limitation of liability covering gambling losses, data loss, and reliance on guidance — good.
- Privacy Policy enumerates data collected (email, WebAuthn key, game history, cookies only).
- GDPR Art. 6(1)(b) legal basis stated.
- GDPR and UK GDPR rights enumerated; CCPA non-sale representation.
- Account deletion mechanism described (instant via in-app button; 30-day email request process).
- Data processor table (Postmark, Heroku, Cloudflare, MongoDB host).

---

### 3b. Must-fix gaps

**Gap 1 — No age gate mechanism on the site.**
The ToS says users must be 18+ (or 13+ if changed), but the live site has no age attestation
gate before account creation. Without a gate, the "by using the app you represent" language is
the only enforcement, which is weak in the face of a COPPA audit or GDPR inquiry. Add a
click-through age affirmation (or an age input + confirmation) at the account-creation step.
The 2025 COPPA rule amendments require age mechanisms to be "reasonably calculated, in light of
available technology" to determine whether a visitor is a child.

**Gap 2 — Internal age inconsistency between the two documents.**
- Terms of Use § 3 says: "You must be at least **18 years old**."
- Privacy Policy / Children section says: "This app is not directed to children under **13**."
These two numbers must be consistent. Either both documents say 13+ (with the operator accepting
the Australia/EU caveats), or both say 18+ (or 16+, or whatever the chosen minimum is). The
current inconsistency could be read as the operator simultaneously barring under-18s (Terms) and
knowingly not collecting from under-13s only (Privacy) — creating ambiguity about the 13-17
cohort.

**Gap 3 — No COPPA citation by name in Privacy Policy.**
The Privacy Policy says "not directed to children under 13" and describes the deletion-request
process for minors, but does not cite COPPA by name or describe the compliance framework. Many
privacy attorneys recommend citing COPPA explicitly and describing the operator's obligations under
it (specifically that the operator does not knowingly collect personal information from children
under 13 and will delete promptly if discovered). This is a best-practice gap, not a legal void,
but it helps establish good-faith compliance posture.

**Gap 4 — No prohibited-territory representation in ToS.**
Once geo-blocking is implemented, the ToS should include a clause requiring users to represent
that they are not accessing from a prohibited jurisdiction, and stating that service in those
jurisdictions is blocked and prohibited. Without this, a user who bypasses GeoIP with a VPN has
no explicit ToS violation to cite.

Recommended add-on to § 3 (Eligibility) or a new § 3a (Geographic Restrictions):
> This service is not available in certain jurisdictions where local law prohibits simulated
> casino content or gambling-themed entertainment, regardless of whether real money is involved.
> If you access this service from a prohibited jurisdiction, you do so in violation of these
> Terms and at your own risk. We reserve the right to block access by region at any time.

**Gap 5 — Governing law clause does not name the state.**
Terms § 12 says "laws of the state in which MooseQuest LLC is organized" without naming the state.
This is deliberately vague (likely intentional until the formation state is confirmed), but it
is a weakness in any dispute: the other party could contest which state's law applies. Once
MooseQuest LLC's state of formation is confirmed (or if it has already been formed), name the
state explicitly in § 12.

**Gap 6 — UK Children's Code (AADC) — assess applicability.**
The UK ICO's Children's Code (Age Appropriate Design Code) applies to online services "likely to
be accessed by children under 18" in the UK. For a gambling-strategy simulator, the argument
that children would not access it is plausible given subject matter, but the ICO takes an expansive
view of "likely to be accessed." If any UK-resident under-18 has an account, this flag becomes
material. Recommend confirming with a UK data-protection attorney whether the Children's Code
applies, and if so, what additional design standards are required (high-privacy defaults, no
profiling for children, no nudge techniques, etc.).
- Source: ICO, [ICO Sets Out Priorities to Protect Children's Privacy Online (2024)](https://ico.org.uk/about-the-ico/media-centre/news-and-blogs/2024/04/ico-sets-out-priorities-to-protect-childrens-privacy-online/)

**Gap 7 — Brazil LGPD / Digital ECA (2025).**
Brazil's LGPD (Lei Geral de Proteção de Dados) defines a child as under 12 and requires parental
consent for data processing of children under 12. Brazil's Digital ECA (enacted September 2025)
applies to technology products "aimed at or likely accessed by minors" regardless of where the
product is developed. If the operator anticipates significant Brazilian traffic, a Brazil-specific
attorney review of these obligations is warranted. For a casino-strategy simulator, the "likely
accessed by minors" threshold may be lower risk, but the data collection (email) of a 12-year-old
in Brazil without parental consent would be unlawful under LGPD Art. 14.
- Source: [lgpd-brazil.info — Article 14](https://lgpd-brazil.info/chapter_02/article_14)
- Source: Captain Compliance, [Brazil's Digital ECA](https://captaincompliance.com/education/brazils-digital-eca-is-almost-law-what-the-child-online-safety-bill-demands-and-how-it-interlocks-with-lgpd/)

---

### 3c. Recommended age-gate display wording

**For most-of-world (13+ base):**
> Before you continue, confirm you meet the age requirement for your location:
> - In Australia: you must be 18 or older.
> - In the EU / EEA: you must be 16 or older (or have parental consent if your country's law allows
>   younger with consent).
> - Everywhere else: you must be 13 or older.
>
> By creating an account, you confirm that you meet the age requirement applicable in your
> jurisdiction. This app is not for use by anyone below the applicable minimum age.

**For simple 18+ global gate (current behavior):**
> By creating an account, you confirm that you are at least 18 years old.
> This service is not directed to and may not be used by anyone under 18.

---

### 3d. Recommended geo-block display wording

> **Service not available in your region**
>
> Dark Side Craps is not available in your current location. Access from this jurisdiction is
> restricted because local law prohibits or limits simulated casino content, gaming-themed
> entertainment, or similar services, regardless of whether real money is involved.
>
> We apologize for the inconvenience. If you believe this is an error, please contact
> privacy@moosequest.app.

---

## Jurisdiction flags

| Jurisdiction | Issue | Who to call |
|---|---|---|
| United States (federal) | COPPA 2025 amendment audit overdue (deadline April 22, 2026 has passed); age gate required | Privacy/COPPA attorney |
| Australia | Classification Act R18+ for interactive simulated gambling (effective Sept 2024); no 18+ gate currently implemented | Australian media law attorney OR US attorney with Australian law referral |
| EU/EEA (9 member states) | GDPR Art. 8 parental consent required for 13-15 year olds in states with 14-16 digital consent age | EU data-protection attorney (or DPO if needed) |
| United Kingdom | ICO Children's Code applicability uncertain; UK GDPR age 13 met; ICO enforcement active | UK data-protection solicitor |
| Vietnam | Decree 147/2024 bans casino-style games; geo-block recommended | No action required if blocking |
| Thailand | Gambling Act 1935 captures simulated casino content; geo-block recommended | No action required if blocking |
| Brazil | LGPD Art. 14 (under-12 children); Digital ECA (2025) potential reach | Brazilian data-protection attorney if material traffic |
| All 20 block-list countries | GeoIP implementation | Engineering (no attorney needed for implementation; attorney for ToS geo-restriction clause) |

---

## Timing / deadlines

| Deadline | What | Status |
|---|---|---|
| **April 22, 2026** (PASSED — today is July 20, 2026) | COPPA 2025 rule amendment compliance | OVERDUE. Audit COPPA posture now. |
| **September 22, 2024** (PASSED) | Australia Classification Act R18+ for interactive simulated gambling | Requirement is in effect. 18+ AU gate is currently missing. |
| **December 25, 2024** (PASSED) | Vietnam Decree 147/2024 bans casino-style games | VN block should be in place. |
| **Ongoing** | ICO Children's Code enforcement cycle | Monitor ICO enforcement guidance annually. |
| **Ongoing** | GDPR Art. 8 parental consent for 13-15 in 9 EU states | Required before any 13+ EU launch. |

---

## Questions for your attorney (privacy/COPPA + international)

1. Has MooseQuest LLC completed a COPPA 2025 rule amendment compliance audit? The April 22, 2026
   deadline has passed — what is the current exposure?
2. Does the "performance of contract" (Art. 6(1)(b)) legal basis in the Privacy Policy
   effectively sidestep GDPR Art. 8 age thresholds, or do those thresholds apply regardless of
   legal basis when the data subject is a minor?
3. Is Dark Side Craps, as an interactive browser-based craps simulator, captured by Australia's
   Classification Act R18+ requirement for "video games containing simulated gambling"? Does the
   lack of an App Store distribution channel affect the classification obligation?
4. Does the UK ICO Children's Code (AADC) apply to Dark Side Craps given its adult casino-strategy
   subject matter? What design changes would be required if it applies?
5. What state is MooseQuest LLC organized in? The ToS governing-law clause should name it. If the
   state is not yet determined, this is a formation question for a business formation attorney.
6. Is the current "not directed to children under 13" Privacy Policy language sufficient as a
   COPPA posture document, or should it include more explicit COPPA compliance language?
7. Does Brazil's Digital ECA (September 2025) impose obligations on MooseQuest LLC for Dark Side
   Craps, given Brazilian traffic? What are the parental-consent requirements under LGPD Art. 14?
8. For the 9 borderline-ISO countries (AE, BH, CU, JO, OM, SG, TM, TN, UZ): confirm whether
   Singapore's clear statutory carve-out for no-stakes games is safe to rely on without blocking.
9. Does Kuwait's prohibition on "publishing material related to gambling" cover a no-money
   educational craps strategy simulator? (Already on block list as conservative call, but attorney
   should confirm whether the operator has any exposure for prior access from KW.)
10. Is the "AS IS / AS AVAILABLE" warranty disclaimer, as currently drafted, enforceable in
    California (where MooseQuest LLC may operate) and in the UK / EU consumer contexts?

---

## Sources

### US — Federal

- FTC, [Children's Online Privacy Protection Rule (COPPA)](https://www.ftc.gov/legal-library/browse/rules/childrens-online-privacy-protection-rule-coppa)
- FTC, [Complying with COPPA: Frequently Asked Questions](https://www.ftc.gov/business-guidance/resources/complying-coppa-frequently-asked-questions)
- FTC, [FTC Finalizes Changes to Children's Privacy Rule (Jan 2025)](https://www.ftc.gov/news-events/news/press-releases/2025/01/ftc-finalizes-changes-childrens-privacy-rule-limiting-companies-ability-monetize-kids-data)
- Federal Register, [COPPA Final Rule Amendments (Apr 22, 2025)](https://www.federalregister.gov/documents/2025/04/22/2025-05904/childrens-online-privacy-protection-rule)
- eCFR, [16 CFR Part 312 — Children's Online Privacy Protection Rule](https://www.ecfr.gov/current/title-16/chapter-I/subchapter-C/part-312)

### US — COPPA Analysis (law firm commentary, primary statute cited above)

- Hunton Andrews Kurth, [COPPA Rule Amendment Compliance Deadline Approaches](https://www.hunton.com/privacy-and-cybersecurity-law-blog/coppa-rule-amendment-compliance-deadline-approaches)
- Davis Polk, [FTC Prioritizes COPPA Enforcement](https://www.davispolk.com/insights/client-update/ftc-prioritizes-coppa-enforcement-new-compliance-obligations-take-effect)
- Hintze Law, [Final COPPA Rule Amendments — Definitional Changes](https://hintzelaw.com/blog/2025/2/6/final-coppa-rule-amendments-definitional-changes)
- White & Case, [Unpacking the FTC's COPPA Amendments](https://www.whitecase.com/insight-alert/unpacking-ftcs-coppa-amendments-what-you-need-know)
- Securiti, [FTC's 2025 COPPA Final Rule Amendments](https://securiti.ai/ftc-coppa-final-rule-amendments/)

### GDPR / UK GDPR

- GDPR, [Article 8 — Conditions Applicable to Child's Consent (gdpr-info.eu)](https://gdpr-info.eu/art-8-gdpr/)
- EuConsent.eu, [Digital Age of Consent Under the GDPR — full member-state table](https://euconsent.eu/digital-age-of-consent-under-the-gdpr/)
- ICO, [Children and the UK GDPR](https://ico.org.uk/for-organisations/uk-gdpr-guidance-and-resources/childrens-information/children-and-the-uk-gdpr/)
- ICO, [ICO Sets Out Priorities to Protect Children's Privacy Online (Apr 2024)](https://ico.org.uk/about-the-ico/media-centre/news-and-blogs/2024/04/ico-sets-out-priorities-to-protect-childrens-privacy-online/)

### Australia

- Australian Classification Board, [New Classifications for Gambling-Like Content in Video Games](https://www.classification.gov.au/about-us/media-and-news/news/new-classifications-for-gambling-content-video-games)
- Australian Classification Board, [Games Guidelines Fact Sheet (PDF)](https://www.classification.gov.au/sites/default/files/documents/game-guidelines-for-gambling-like-content-fact-sheet.pdf)
- Australian Government Minister, [New Mandatory Minimum Classifications for Gambling-Like Games Content](https://minister.infrastructure.gov.au/rowland/media-release/new-mandatory-minimum-classifications-gambling-games-content)
- Corrs Chambers Westgarth, [Mandatory Minimum Classifications for Computer Games with Gambling-Like Content](https://www.corrs.com.au/insights/mandatory-minimum-classifications-for-computer-games-with-gambling-like-content)
- ACMA, [Online Gambling Services](https://www.acma.gov.au/online-gambling-services)

### App store policy vs. law distinction

- AppleInsider, [App Store Shakeup Limits Simulated Gambling to Users Aged 17+](https://appleinsider.com/articles/19/08/20/app-store-shakeup-limits-simulated-gambling-to-users-aged-17)
- Apple Developer, [Age Ratings Values and Definitions](https://developer.apple.com/help/app-store-connect/reference/app-information/age-ratings-values-and-definitions/)

### US social casino law

- Richt Law, [Proliferation of Social Gaming Casinos: Legal Compliance Considerations](https://richtfirm.com/the-proliferation-of-social-gaming-casinos-legal-compliance-considerations/)
- gambling.com, [How Social Casinos Navigate Legal Constraints in the US](https://www.gambling.com/us/online-casinos/strategy/how-social-casinos-navigate-legal-constraints-in-the-us)

### Geo-block — country law sources

- Singapore: [IMDA — Remote Gambling Act Clarification Note on Scope of Games (Jan 2015)](https://www.imda.gov.sg/-/media/imda/files/regulation-licensing-and-consultations/codes-of-practice-and-guidelines/acts-codes/14-remote-gambling-act-clarification-note-on-the-scope-of-games.pdf)
- Singapore: [sso.agc.gov.sg — Remote Gambling Act 2014](https://sso.agc.gov.sg/Act-Rev/RGA2014/Published/20211231?DocDate=20211231)
- Vietnam: [Yogonet — Vietnam Decree 147/2024](https://www.yogonet.com/international/news/2024/12/05/87586-vietnam-39s-new-decree-147-enforces-stricter-regulations-on-online-gaming-industry)
- Vietnam: [vietnamnet.vn — Casino-style games banned under Decree 147](https://vietnamnet.vn/en/vietnam-tightens-online-gaming-rules-bans-casino-style-and-card-based-games-2346909.html)
- Thailand: [Wikipedia — Gambling in Thailand](https://en.wikipedia.org/wiki/Gambling_in_Thailand)
- Thailand: [Yogonet — Thailand 25,000+ sites blocked](https://www.yogonet.com/international/news/2024/04/03/71535-thailand-steps-up-efforts-against-online-gambling-blocks-25-000-websites)
- China: [GGRAsia — China 4,500+ platforms destroyed 2024](https://www.ggrasia.com/china-destroyed-in-2024-more-than-4500-illegal-online-gambling-platforms-ministry-of-security)
- Malaysia: [AGB Brief — Malaysia blocking](https://agbrief.com/news/malaysia/13/08/2024/majority-of-blocked-websites-in-malaysia-linked-to-online-gambling/)
- Malaysia: [Bright Side of News — Malaysia 2026 reform](https://brightsideofnews.com/gambling/regulation/malaysia-2026-anti-online-gambling-reform/)
- Saudi Arabia: [sumsub.com — UAE/Saudi gaming regulations](https://sumsub.com/blog/uae-gaming-regulations-all-you-need-to-know/)
- UAE: [aladllegal.com — UAE Gambling Laws](https://aladllegal.com/news/uaes-gambling-laws-prohibition-consequences-and-ex/)
- UAE: [GCGRA](https://www.gcgra.gov.ae/)
- Algeria: [CMS Expert Guide — Algeria gambling law](https://cms.law/en/int/expert-guides/cms-expert-guide-to-gambling-laws-in-africa/algeria)
- Kuwait: [arabiancasino.org — Arab countries](https://arabiancasino.org/countries/)
- Pakistan: [legalpilot.com — Pakistan](https://legalpilot.com/country/pakistan/)
- General multi-country: [playtoday.co — Countries Where Gambling Is Illegal](https://playtoday.co/blog/guides/countries-where-gambling-is-illegal/)
- General multi-country: [legalpilot.com — Gambling Regulations Map](https://legalpilot.com/gambling-regulations-map/)

### Brazil

- [lgpd-brazil.info — Article 14 (Children's Data)](https://lgpd-brazil.info/chapter_02/article_14)
- Captain Compliance, [Brazil's Digital ECA](https://captaincompliance.com/education/brazils-digital-eca-is-almost-law-what-the-child-online-safety-bill-demands-and-how-it-interlocks-with-lgpd/)
- IAPP, [Children's Privacy Under Brazil's LGPD](https://iapp.org/news/a/can-mandatory-consent-be-optional-processing-childrens-personal-data-under-brazils-lgpd)

---

> **Before acting on any item in this document, consult a licensed attorney with expertise in
> privacy law (COPPA, GDPR), international internet regulation, and/or Australian media
> classification law, as applicable to each item. This document is research preparation material,
> not legal advice. MooseQuest LLC is responsible for confirming the currency of all sources and
> obtaining qualified legal review before implementing any policy, gate, or block.**
