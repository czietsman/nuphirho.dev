# nuphirho.dev Project Brief

This document captures the decisions, context, and rationale for nuphirho.dev. It serves as the single source of truth for any collaborator picking up work on this project.

---

## Project Overview

nuphirho.dev is a public-facing technical blog created by Christo Zietsman. The blog is managed like a software project, from ideation through definition, prototyping, implementation, and release. The project itself, including its infrastructure, pipeline, and process, is the subject of the first blog post.

### Goals

- Increase Christo's renown in the developer and engineering leadership community.
- Complement an existing 11-post LinkedIn thought leadership series on AI-assisted software development (covering trust barriers, security, architecture, and organisational transformation).
- Demonstrate that enterprise-grade engineering practices are available today at near-zero cost.

### Core philosophy

Process matters more than technology. AI has changed the economics of rigorous engineering practices, making things like BDD specifications, mutation testing, and executable verification layers economically viable in ways they were not before. Security is a non-negotiable default.

---

## The Author

Christo Zietsman. Based in Somerset West, Western Cape, South Africa. 20+ years of experience spanning mine seismology (ISSI/IMS), expert systems for antenna design (Antenna Magus), enterprise backup (Attix5/Redstor), and cybersecurity (CyberSentriq). Masters in Continuum Mechanics, Cum Laude, Stellenbosch University.

Currently Director of Technology Innovation at CyberSentriq. General approaches and patterns are shareable; proprietary details are not.

### AI transparency

AI assists in research, drafting, and refinement for this blog. The thinking, decisions, direction, and accountability are the author's. This is stated openly on the about page.

The framing: "I think in systems and architecture. AI helps me express those ideas clearly. The thinking is mine. The clarity is a collaboration."

This reflects a broader philosophy: AI is a tool for amplifying human capability, not replacing human judgement. The same principle applies to every technology choice in this project.

---

## Domain and Brand

### Domain

nuphirho.dev, registered on Cloudflare.

### Origin of the name

Nu, phi, rho: Greek letters from physics and mathematics. Used as a gaming profile since university. Sentimental value. The story is part of the brand.

### Tagline

"Enterprise-grade engineering at startup speed" (consistent with LinkedIn profile).

### Identity across platforms

- Main site: nuphirho.dev (SvelteKit on Cloudflare Pages)
- Blog: blog.nuphirho.dev (SvelteKit on Cloudflare Pages)
- Cross-post: nuphirho on Dev.to
- GitHub: czietsman (existing account, repo is czietsman/nuphirho.dev)
- LinkedIn: christo-zietsman (existing profile)

---

## Architecture Decisions

### Platform and distribution

| Concern | Tool | Cost |
|---|---|---|
| Source control | GitHub (public repo) | Free |
| CI/CD | GitHub Actions | Free (public repo, unlimited minutes) |
| Secrets | GitHub Secrets | Free |
| IaC | Terraform + Cloudflare provider | Free |
| DNS/CDN/SSL | Cloudflare (free tier) | Free |
| Domain | Cloudflare Registrar | ~$12-15/year |
| Frontend framework | SvelteKit + adapter-cloudflare | Free |
| Landing page hosting | Cloudflare Pages (nuphirho.dev) | Free |
| Blog hosting | Cloudflare Pages + Workers (blog.nuphirho.dev) | Free |
| Blog visitor counter | Cloudflare KV | Free (within limits) |
| Cross-post | Dev.to (API, automated) | Free |
| Cross-post | Medium (manual, URL import) | Free |
| Amplification | LinkedIn | Free |

The domain is the only cost. Everything else is enterprise-grade tooling at zero cost.

### Why these choices

- **GitHub (public):** The repo itself is a portfolio piece. Demonstrates the process openly. GitHub Actions are free and unlimited for public repos.
- **SvelteKit + Cloudflare Pages:** Full prerender at build time. Zero cold-start latency. The Worker only handles the `/api/stats` endpoint. No tracking cookies, no third-party analytics.
- **Dev.to:** REST API with simple API key auth. Built-in developer audience. Supports canonical URLs.
- **Medium:** API is no longer issuing new integration tokens. Medium is a manual cross-post target using the "Import a story" URL feature. The pipeline handles this gracefully.
- **Cloudflare:** Free tier includes DNS, CDN, SSL, Pages, Workers, and KV. Mature Terraform provider. Domain registered here so DNS, CDN, and registrar are in one place.
- **Terraform:** Infrastructure as code. Cloudflare DNS and Pages configuration managed declaratively.

### Security

- HTTPS enforced at multiple layers: .dev TLD (HSTS preload list), Cloudflare (free SSL/TLS).
- No credentials, tokens, or secrets in the repository or post content.
- All secrets stored in GitHub Secrets.
- No third-party tracking scripts or analytics on either site. Visitor counter is path-only, privacy-respecting, and runs on Cloudflare KV.
- Security is a core tenet of the blog, both in content and in practice.

### Canonical URLs

Always set to blog.nuphirho.dev. Every cross-posted article must reference the canonical URL on the blog subdomain. This protects SEO and ensures the blog builds authority over time.

---

## Content Strategy

### Planned content areas

- AI-assisted software delivery (complementing the LinkedIn series).
- Engineering process and practice (BDD, mutation testing, executable specifications).
- Introducing tools and practices to existing platforms and teams.
- Organisational transformation and technology innovation.

### LinkedIn series

11 posts on AI-assisted software development covering trust barriers, security, architecture, and organisational transformation. The blog and the LinkedIn series are complementary and cross-reference each other.

---

## Repository Structure

```
czietsman/nuphirho.dev
├── blog/                         # SvelteKit app — blog.nuphirho.dev
│   └── src/
│       ├── lib/posts.ts          # Post loading and markdown parsing
│       └── routes/               # Post list, post detail, /api/stats
├── main-site/                    # SvelteKit app — nuphirho.dev
│   └── src/
│       ├── lib/Roadmap.svelte    # Publishing calendar component
│       └── routes/               # Landing, about, words-of-meaning, etc.
├── site/                         # Go template static site (legacy)
├── posts/                        # Markdown blog post source files
├── cmd/
│   ├── publish/                  # Reconcile posts/ with Dev.to
│   ├── notify/                   # Send Telegram notifications
│   ├── notify-summary/           # Scheduled notification digest
│   ├── site-build/               # Build Go template site
│   └── validate-tags/            # Validate post tag values
├── internal/
│   ├── devto/                    # Dev.to REST API client
│   ├── frontmatter/              # Post metadata schema and parsing
│   ├── hashnode/                 # Hashnode GraphQL API client
│   ├── pipeline/                 # Publishing orchestration
│   └── tags/                     # Tag validation
├── terraform/                    # Cloudflare infrastructure as code
├── .github/workflows/            # CI/CD pipelines
├── docs/                         # Project brief and style guide
├── specs/                        # BDD feature files
├── tests/                        # Shell-based test scripts
├── prompts/                      # Reviewed prompt material
├── papers/                       # Academic paper builds
├── experiments/                  # Research projects
├── AGENTS.md                     # Agent instructions (authoritative)
├── CLAUDE.md                     # Points to AGENTS.md
└── README.md
```

---

## Key Principles

- The platform and domain are more important than the tech.
- Everything in the open. Public repo, public process.
- Enterprise-grade practices at zero cost.
- Security is a default, not an afterthought.
- Process over technology.
- Honest, transparent, no fluff.
