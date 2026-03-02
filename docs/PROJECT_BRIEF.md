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

Currently Director of Technology Innovation at CyberSentriq. Content referencing CyberSentriq requires CPTO approval before publishing. General approaches and patterns are shareable; proprietary details are not. Approval has not yet been obtained.

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

- Blog: nuphirho on Hashnode (nuphirho.dev)
- Cross-post: nuphirho on Dev.to
- GitHub: czietsman (existing account, repo is czietsman/nuphirho.dev)
- LinkedIn: christo-zietsman (existing profile)
- The about page on Hashnode connects nuphirho to Christo Zietsman by name.

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
| Primary platform | Hashnode (custom domain) | Free (includes custom domain, SSL, CDN) |
| Cross-post | Dev.to (API, automated) | Free |
| Cross-post | Medium (manual, URL import) | Free |
| Amplification | LinkedIn | Free |
| Future consideration | X | Not yet |

The domain is the only cost. Everything else is enterprise-grade tooling at zero cost.

### Why these choices

- **GitHub (public):** The repo itself is a portfolio piece. Demonstrates the process openly. GitHub Actions are free and unlimited for public repos.
- **Hashnode:** Free custom domain mapping, built-in developer community for discoverability, GraphQL API for automated publishing, supports light/dark themes, native Mermaid diagram support.
- **Dev.to:** REST API with simple API key auth. Built-in developer audience. Supports canonical URLs.
- **Medium:** API is no longer issuing new integration tokens. Medium is a manual cross-post target using the "Import a story" URL feature. The pipeline handles this gracefully.
- **Cloudflare:** Free tier includes DNS, CDN, SSL. Mature Terraform provider. Domain registered here so DNS, CDN, and registrar are in one place.
- **Terraform:** Infrastructure as code. Cloudflare DNS configuration managed declaratively.

### Security

- HTTPS enforced at three layers: .dev TLD (HSTS preload list), Cloudflare (free SSL/TLS), Hashnode (Let's Encrypt auto-provisioned).
- No credentials, tokens, or secrets in the repository or post content.
- All secrets stored in GitHub Secrets.
- Security is a core tenet of the blog, both in content and in practice.

### Canonical URLs

Always set to nuphirho.dev. Every cross-posted article must reference the canonical URL on the primary domain. This protects SEO and ensures the custom domain builds authority over time.

---

## Content Strategy

### First blog post

About setting up the blog itself: the platform choices, the pipeline, the process, and the domain. A meta-post that demonstrates the philosophy in action.

### Planned content areas

- AI-assisted software delivery (complementing the LinkedIn series).
- Engineering process and practice (BDD, mutation testing, executable specifications).
- Introducing tools and practices to existing platforms and teams (e.g. split.io, Auth0).
- Organisational transformation and technology innovation.

### LinkedIn series

11 posts on AI-assisted software development covering trust barriers, security, architecture, and organisational transformation. The blog and the LinkedIn series are complementary and cross-reference each other.

---

## Repository Structure

```
czietsman/nuphirho.dev
├── docs/
│   ├── STYLE_GUIDE.md
│   └── PROJECT_BRIEF.md
├── specs/
│   └── (BDD feature files)
├── terraform/
│   └── (Cloudflare DNS configuration)
├── .github/
│   └── workflows/
│       └── (GitHub Actions pipeline)
├── posts/
│   └── (markdown blog post files)
├── .gitignore
└── README.md
```

---

## What Needs to Happen Next

### Immediate

1. Create `.gitignore` and `README.md` for the repo.
2. Commit `docs/STYLE_GUIDE.md` and `docs/PROJECT_BRIEF.md`.
3. Write BDD specs (Gherkin) describing the content pipeline behaviour.
4. Set up Hashnode account, connect nuphirho.dev custom domain.
5. Generate Dev.to API key.
6. Store API keys and tokens in GitHub Secrets.

### Pipeline

7. Terraform configuration for Cloudflare DNS (pointing to Hashnode).
8. GitHub Actions workflow: push to main triggers publish to Hashnode and Dev.to.
9. Handle Medium as a manual step with documented process.

### First post

10. Draft the first blog post about setting up the blog.
11. Review against the style guide.
12. Publish through the pipeline.

---

## Key Principles

- The platform and domain are more important than the tech.
- Everything in the open. Public repo, public process.
- Enterprise-grade practices at zero cost.
- Security is a default, not an afterthought.
- Process over technology.
- Honest, transparent, no fluff.
