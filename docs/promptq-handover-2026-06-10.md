# PromptQ Handover: nuphirho.dev Governance Documents

**Date:** 2026-06-10
**Framework version applied:** 1.3.1
**Documents evaluated:** `AGENTS.md`, `docs/STYLE_GUIDE.md`
**Repository:** czietsman/nuphirho.dev (public)

---

## Context

nuphirho.dev is a technical blog managed as a software project. The repository is public and contains live infrastructure, a publishing pipeline, and two SvelteKit sites. `AGENTS.md` is the primary governance document for AI agents working in the repository. `docs/STYLE_GUIDE.md` governs written content and design decisions.

PromptQ was applied to both documents in a live working session. Findings were remediated immediately, with the exception of one open question requiring a decision.

---

## AGENTS.md

**Purpose:** Governs all AI agent behaviour in the repository. Covers TDD, BDD, scope, secret hygiene, completion evidence format, deployment architecture constraints, and Terraform discipline.

**Baseline score (pre-session):** 4.5 / 7

| Principle | Score | Finding |
|---|---|---|
| P1 Success Definition | 1 | Four-field completion evidence report with explicit done condition |
| P2 Assessment Rubric | 0.5 | No observability requirement; no active human engagement step |
| P3 Scope Boundary | 1 | Execution-agents-only with explicit prohibition list and escalation format |
| P4 Data Classification | 0.5 | External content classified as untrusted; no formal taxonomy |
| P5 Quality Gate | 0.5 | Self-certification; no confirmatory human action required |
| P6 Internal Consistency | 0.5 | Two gaps: (1) TDD rule applied to "every code change" with no carve-out for YAML/Terraform/docs; (2) blog section required "confirm no secrets exposed" with no verification mechanism |
| P7 Contextual Currency | 0 | `last-modified` date present; no staleness declaration of any kind |

### Remediations applied

**P6 gap 1 (TDD scope):** Added a paragraph after the TDD section clarifying that TDD applies to Go, TypeScript, and Svelte code only. Terraform, YAML, documentation, and configuration are excluded. Per-type verification guidance provided. The "Tests run" field in the evidence report should reflect the applicable verification type.

**P6 gap 2 (blog confirmation mechanism):** Replaced the vague "confirm it does not expose secrets" obligation with a specific post-build check: search `.svelte-kit/cloudflare/` output for KV binding names and environment variable names after build. If the output cannot be inspected, raise with the repo owner before merging.

**P2 + P5 combined (rubric and quality gate):** Added a three-criterion rubric to the Completion evidence section. The repo owner reviews against: (1) Traceability -- files listed in "Files changed" match the PR diff; (2) Reproducibility -- the cited test command can be re-run and produces the stated outcome; (3) Specificity -- "Confirmed" names a specific observable behaviour, not a general claim. The rubric explicitly requires repo owner review before merge, closing the self-certification gap.

**P7 (contextual currency):** Added a "Re-evaluation" section declaring all three trigger types: event-based (model update, new tool/MCP, architecture change), evidence-based (two or more consecutive runs with unexpected output), and time-based (three months). Owner: the repo owner.

**Final score:** 6.5 / 7

**Remaining gap: P4 (Data Classification, 0.5)** -- three implicit data categories exist (external/untrusted, repository/authoritative, secrets/restricted) but no formal taxonomy names them. Closing this to 1 would require a short classification table naming the categories and their handling rules explicitly. Low priority; the current behaviour is correct.

---

## STYLE_GUIDE.md

**Purpose:** Governs voice, tone, formatting, and editorial standards for all written content. Also serves as the brief for AI-assisted drafting. Extended during this session to govern design system compliance.

**Baseline score (pre-session, applicable principles only):** 2.0 / 5

| Principle | Score | Finding |
|---|---|---|
| P1 Success Definition | 1 | "Definition of done for every post" stated; length, structure, and tone criteria present |
| P2 Assessment Rubric | 0 | Rules defined but no rubric for evaluating compliance |
| P3 Scope Boundary | 0.5 | Scope stated as "all written content"; no negative boundary |
| P6 Internal Consistency | 0.5 | Three stale Hashnode references contradicted the current Cloudflare Pages architecture (Security section, Diagrams section, Distribution section) |
| P7 Contextual Currency | 0 | No staleness declaration |

### Remediations applied

**P6 (stale references):** Corrected three Hashnode references to reflect the current platform:

- Security section: "platform (Hashnode)" -- now "platform (Cloudflare Pages)"
- Diagrams section: "supported by Hashnode" -- now "render correctly in the SvelteKit blog"
- Distribution section: "blog.nuphirho.dev (Hashnode, custom subdomain)" -- now "blog.nuphirho.dev (SvelteKit on Cloudflare Pages)"

**P7 (contextual currency):** Added a "Re-evaluation" section with event-based triggers (platform change, new cross-posting target, canonical URL change), an evidence-based trigger (repeated style violations despite following the guide), and a three-month time-based fallback. Owner: the repo owner.

**P2 (assessment rubric) -- applied in two passes:**

*First pass:* Added a three-part compliance rubric covering written content (British English, no em dashes, no emoji, first person, canonical URL, attribution, word count, opening structure) and pipeline/frontmatter (draft flag, publish_date, slug uniqueness, required fields).

*Second pass (design system extension):* Following discussion, the rubric was extended to cover design and code compliance. A new "Design system" section was added documenting the CSS token vocabulary. All hardcoded `border-radius` values in both `blog/src/app.css` and `main-site/src/app.css` were replaced with a new `--radius-xs/sm/md/lg` token scale. Hardcoded confidence badge colours in `main-site/src/app.css` were moved to `--confidence-*` CSS custom properties with proper `[data-theme="dark"]` overrides, removing the duplicate inline overrides. The design rubric checklist covers: colour token compliance, radius token compliance, new token definition process, light/dark testing, alt text, and no secrets in visual content.

**Final score:** 4.5 / 5

**Remaining gap: P3 (Scope Boundary, 0.5)** -- the guide states it covers "all written content" but does not define what it does not cover (design system engineering beyond CSS tokens, accessibility beyond alt text, cross-format rendering). Closing this to 1 would require a brief statement of scope exclusions.

---

## Open question (raised, not resolved)

The rubric currently covers web rendering on the canonical blog. The following surfaces were flagged as under-specified:

- **Referenced files:** whether a formal asset manifest per post is required, or whether constraints on committed versus externally linked assets are needed
- **Cross-format rendering:** whether the rubric should cover Dev.to rendering (which strips CSS) or is scoped to the canonical blog only
- **Print:** whether there is an active or planned print stylesheet use case

These were deferred pending a decision from the repo owner on scope.

---

## Summary of changes

All changes are on branch `claude/add-claude-md` (PR #100):

| File | Changes |
|---|---|
| `AGENTS.md` | P6 TDD scope, P6 blog check mechanism, P2+P5 rubric, P7 re-evaluation section, owner and cadence update |
| `docs/STYLE_GUIDE.md` | P6 Hashnode corrections, P7 re-evaluation section, Design system section, Compliance rubric, owner and cadence update |
| `blog/src/app.css` | `--radius-*` tokens added; all `border-radius` values tokenised |
| `main-site/src/app.css` | `--radius-*` tokens added; `--confidence-*` colour tokens added; all `border-radius` and badge colours tokenised |
| `CLAUDE.md` | New file: `@AGENTS.md` embed |
| `docs/PROJECT_BRIEF.md` | Architecture updated for current platform; completed todo section removed; CPTO approval clause removed |
| `README.md` | Comprehensive rewrite for current stack, structure, and workflows |
