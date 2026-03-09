# nuphirho.dev

Enterprise-grade engineering at startup speed.

A technical blog by [Christo Zietsman](https://www.linkedin.com/in/christo-zietsman/), managed like a software project. The infrastructure, pipeline, and process are all in this repository.

## What is this

This is the source repository for [nuphirho.dev](https://nuphirho.dev), a technical blog about AI-assisted software delivery, engineering process, and technology innovation. The root domain is a static landing page served via GitHub Pages. The blog lives at [blog.nuphirho.dev](https://blog.nuphirho.dev), hosted on Hashnode with cross-posting to Dev.to.

The repository contains:

- **site/** -- Static landing page (deployed to GitHub Pages)
- **posts/** -- Markdown blog post source files
- **specs/** -- BDD feature files describing pipeline behaviour
- **terraform/** -- Cloudflare DNS configuration as code
- **.github/workflows/** -- GitHub Actions CI/CD pipelines
- **docs/** -- Project brief and style guide
- **tests/** -- Shell-based test scripts

## Philosophy

Process matters more than technology. AI has changed the economics of rigorous engineering practices. This project demonstrates that enterprise-grade tooling and practices are available at near-zero cost.

AI assists in research, drafting, and refinement. The thinking, decisions, direction, and accountability are the author's. The thinking is mine. The clarity is a collaboration.

## Stack

| Concern | Tool |
|---|---|
| Source control | GitHub (public) |
| CI/CD | GitHub Actions |
| IaC | Terraform + Cloudflare provider |
| DNS/CDN/SSL | Cloudflare (free tier) |
| Landing page | GitHub Pages (nuphirho.dev) |
| Primary platform | Hashnode (blog.nuphirho.dev) |
| Cross-post | Dev.to (automated), Medium (manual) |
| Amplification | LinkedIn |
| Secret detection | Husky pre-push hook (grep-based) |

The domain is the only cost.

## Development setup

```bash
npm install
```

This installs [husky](https://typicode.github.io/husky/) and configures a `pre-push` git hook that scans for secrets before code leaves your machine. The hook catches:

- AWS / R2 access key IDs (`AKIA...`)
- GitHub token variants (`ghp_`, `github_pat_`, `gho_`, `ghu_`, `ghr_`, `ghs_`)
- PEM private key headers
- Assignments to known secret variables (`CLOUDFLARE_API_TOKEN`, `HASHNODE_TOKEN`, `DEVTO_API_KEY`, and others)
- Generic secret patterns (`api_key`, `token`, `password`, etc. followed by long values)

Paths listed in `.secretscanignore` are excluded from scanning. To bypass the hook for a known false positive, use `git push --no-verify`.

To run the pattern tests:

```bash
bash tests/test-secret-patterns.sh
```

## Getting started

### Terraform

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Requires a `CLOUDFLARE_API_TOKEN` environment variable or equivalent configuration in GitHub Secrets.

### Publishing

Push a markdown file to `posts/` on the `main` branch. The GitHub Actions pipeline handles publishing to Hashnode and cross-posting to Dev.to.

Posts with `draft: true` in the front matter are pushed as unpublished drafts to both Hashnode and Dev.to, allowing end-to-end pipeline validation without going live. Removing the `draft` flag and pushing again publishes the post.

Medium cross-posting is manual via the "Import a story" feature using the canonical URL.

## Licence

Content (posts, documentation) is copyright Christo Zietsman. Code (Terraform, workflows, scripts) is MIT licensed.

