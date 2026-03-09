---
title: "Threat modeling is not just for enterprise"
slug: "threat-modeling-is-not-just-for-enterprise"
subtitle: "A worked example on a personal blog project — and what GitHub gives you for free"
tags: [security, threat-modeling, devops, github, engineering-process]
draft: true
---

I run a personal blog. It has a public GitHub repository, a live publishing pipeline, a custom domain on Cloudflare, and posts that go to Hashnode and Dev.to automatically when I push markdown to main. It is not a bank. It is not a SaaS product. It does not hold customer data.

I threat modeled it anyway.

Not because I expected to find catastrophic vulnerabilities. Because threat modeling is a thinking tool, and the thinking scales down just as well as it scales up. What I found was instructive — not because the system was insecure, but because the exercise made explicit what I had assumed and surfaced a handful of things I had not thought through.

This post walks through the method, the system, and the findings. The follow-up post covers how AI fits into that process and where the new risks live.

---

## What threat modeling actually is

Threat modeling is the practice of analysing a system from an adversarial perspective before something goes wrong. You model what you have built, identify where it could be exploited, and decide what to do about each finding.

The four questions that anchor every approach, regardless of methodology, are Adam Shostack's:

1. What are we building?
2. What can go wrong?
3. What are we going to do about it?
4. Did we do a good enough job?

That is the entirety of the framework at its most distilled. Everything else — STRIDE, PASTA, LINDDUN, ATT&CK — is methodology layered on top of those four questions to make the process systematic and repeatable.

The reason threat modeling has an enterprise reputation is not that it requires enterprise resources. It is that enterprise organisations have compliance requirements that mandate it, and the tooling and consultancy that grew up around those requirements carries the weight of that context. Strip that away and you have a structured thinking exercise that works on anything with a trust boundary.

---

## Choosing a method

There are more methods than you need to know for most purposes. The ones worth understanding are:

**STRIDE** — the most widely adopted starting point. Developed at Microsoft in the late 1990s, it categorises threats into six types: Spoofing, Tampering, Repudiation, Information Disclosure, Denial of Service, and Elevation of Privilege. It works by mapping each component in a Data Flow Diagram against those six categories. It is system-centric, approachable, and produces concrete findings. The critique is that it identifies threat categories but does not inherently prioritise them or connect them to business impact.

**PASTA** — a seven-stage, risk-centric methodology that works from business objectives down to technical threats. It connects security findings to organisational risk, which makes it useful when findings need to drive business decisions rather than engineering tickets. Appropriate for mature security programmes and complex SaaS platforms. Overkill for most other contexts.

**LINDDUN** — privacy-focused, structured like STRIDE but covering data protection rather than security. Relevant whenever personal data is processed, particularly under GDPR, CCPA, or similar frameworks.

**MITRE ATT&CK** — not a methodology but a knowledge base of adversary tactics and techniques derived from real-world observations. Used as an enrichment layer on top of other methods: once STRIDE identifies a threat category, ATT&CK maps it to concrete techniques, giving specificity about how an attack would actually unfold.

For a small system with no personal data, no complex business logic, and a narrow attack surface, STRIDE is the right choice. PASTA adds overhead that is not warranted. LINDDUN does not apply. ATT&CK is useful for a second pass on the highest-priority findings.

---

## The system

The nuphirho.dev stack is straightforward:

| Concern | Tool |
|---------|------|
| Source control | GitHub (public repo) |
| CI/CD | GitHub Actions |
| IaC | Terraform + Cloudflare provider |
| DNS / CDN / SSL | Cloudflare (free tier) |
| Landing page | GitHub Pages |
| Primary blog | Hashnode |
| Cross-post | Dev.to (automated), Medium (manual) |
| Secret detection | Husky pre-push hook |
| Secret storage | GitHub Secrets |

The key characteristic that shapes the threat model: the repository is public. Source, workflow definitions, Terraform configuration, and pipeline logic are all visible to anyone. The attack surface is not larger than a private repo, but the reconnaissance cost for an adversary is zero.

The pipeline has no inbound API surface — it only makes outbound calls to Hashnode and Dev.to. Secrets used: `CLOUDFLARE_API_TOKEN`, `HASHNODE_TOKEN`, `DEVTO_API_KEY`.

---

## The trust boundaries

Rather than a diagram, the four trust boundaries can be stated plainly:

**Developer machine to GitHub** — code and workflow definitions cross this boundary on every push. The Husky pre-push hook scans for secrets before they leave the machine.

**GitHub Actions to external APIs** — the Actions runner pulls secrets from GitHub Secrets and makes outbound calls to Hashnode and Dev.to. No inbound surface.

**GitHub Actions to Cloudflare via Terraform** — DNS changes are applied using the Cloudflare API token. Changes propagate globally.

**Cloudflare to readers** — nuphirho.dev served via GitHub Pages, blog.nuphirho.dev served via Hashnode, both proxied through Cloudflare with TLS terminated at the edge.

---

## The STRIDE findings

Running the six categories across the system produced thirteen findings. The most instructive ones:

**Supply chain via Actions (Spoofing / Tampering)**
All six third-party GitHub Actions in the workflow files were pinned to mutable tags — `@v4`, `@v3`, and so on. A mutable tag means the action author can push a new commit under that tag at any time. If a dependency in the action's supply chain is compromised, the malicious code runs in the next pipeline execution with access to GitHub Secrets.

The fix is to pin every action to a full commit SHA rather than a tag, with the tag preserved as a comment for human readability. This is a low-effort, high-value change that is easy to overlook and easy to automate.

**GITHUB_TOKEN scope (Elevation of Privilege)**
GitHub Actions automatically provides a `GITHUB_TOKEN` scoped to the repository. If workflow jobs do not declare explicit `permissions:` blocks, the token defaults to broad access. In a publishing pipeline, most jobs need only `contents: read`. The deploy job needs `pages: write` and `id-token: write`. The Terraform job needs `pull-requests: write` for plan comments.

Adding explicit job-level permissions blocks and setting the top-level permissions to `{}` means the token carries only what each job demonstrably requires. Any supply chain compromise that attempts to use the token beyond that scope fails.

**Secret leakage (Information Disclosure)**
The Husky pre-push hook scans for secrets using grep patterns. It can be bypassed with `git push --no-verify`. GitHub's native secret scanning — available at no cost on public repositories — provides a second layer that cannot be bypassed from the developer side. It was not enabled by default.

**Cloudflare token scope (Tampering)**
The Cloudflare API token existed. The question was whether it was scoped to the minimum required — DNS edit on the specific zone — or whether it carried broader account-level permissions. A leaked broad token could modify security settings, add or remove zones, or access billing. This was already handled correctly, but the threat model made it explicit.

---

## What GitHub gives you for free

One of the more useful outputs of the exercise was a clear picture of what GitHub enables by default versus what requires deliberate configuration.

**Enabled by default on public repositories:**
- Secret scanning — GitHub scans commits for known credential patterns and alerts on matches
- Workflow file access restriction — `.github/workflows/*` is not accessible via the unauthenticated API, limiting reconnaissance
- `AGENTS.md` and `.secretscanignore` access restriction — GitHub blocks unauthenticated access to these files
- Actions `GITHUB_TOKEN` scoped to the repository — the token cannot access other repositories

**Not enabled by default — requires deliberate action:**
- Branch protection — nothing prevents force pushes or deletion of `main` without explicit configuration
- Secret scanning push protection — scanning alerts after the fact; push protection blocks the push before it lands
- Job-level `permissions:` blocks — without them, jobs inherit broad defaults
- Actions SHA pinning — the default is mutable tags, which is convenient and insecure
- PR validation workflow — there is no automatic gate on pull requests unless you build one

The defaults are a starting point, not a security posture. The gap between the two is not large, but it is deliberate work.

---

## What the exercise validated

Threat modeling does not only find gaps. It also confirms what is already done correctly, and that confirmation has value.

Pre-existing controls the exercise validated: the Cloudflare API token was already scoped to least privilege. Terraform state was already stored remotely on Cloudflare R2 rather than committed to the repository. All credentials were already in GitHub Secrets with no hardcoded values anywhere in the codebase. The Terraform configuration resolved the Cloudflare zone ID via a data source lookup rather than a hardcoded hex string.

These were not discoveries — they were verifications. The threat model produced a documented, reasoned account of why the system is secure where it is, not just a list of what to fix.

---

## The process is the point

The output of this exercise was not a document. It was a set of changes — Actions pinned to SHAs, permissions blocks added, branch protection configured, secret scanning push protection enabled — each with an audit trail of what was done and why.

More usefully, it was a shift in how the system is understood. The public repository is not just "readable." It is a system where an adversary already knows which APIs the pipeline calls, which token names to look for, and which workflow triggers are in play. The threat model makes that explicit in a way that intuition alone does not.

The next post covers how AI participated in this process — and why that participation introduces its own threat surface that needs the same systematic treatment.

---
*This is part one of a three-part series. Part two covers the AI collaboration layer and the security risks it introduces: [AI as collaborator: research, execution, and the boundary between them](<!-- POST_2_URL -->).*
