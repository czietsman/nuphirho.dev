# nuphirho.dev

Enterprise-grade engineering at startup speed.

A technical blog by [Christo Zietsman](https://www.linkedin.com/in/christo-zietsman/), managed like a software project. The infrastructure, pipeline, and process are all in this repository.

## What is this

This is the source repository for [nuphirho.dev](https://nuphirho.dev), a technical blog about AI-assisted software delivery, engineering process, and technology innovation. The root domain is a static landing page served via GitHub Pages. The blog lives at [blog.nuphirho.dev](https://blog.nuphirho.dev), hosted on Hashnode with cross-posting to Dev.to.

The repository contains:

- **site/** -- Static site built from Go templates (deployed to GitHub Pages)
  - **site/templates/** -- HTML templates (base layout, partials, page content)
  - **site/css/** -- Stylesheets (shared + page-specific)
  - **site/js/** -- JavaScript (theme toggle, roadmap calendar)
  - **site/static/** -- Files copied directly to output (CNAME, calendar data)
- **posts/** -- Markdown blog post source files
- **prompts/** -- Reviewed prompt material, including dependency review briefs
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
| Site build | Go html/template (cmd/site-build) |
| Landing page | GitHub Pages (nuphirho.dev) |
| Primary platform | Hashnode (blog.nuphirho.dev) |
| Cross-post | Dev.to (automated), Medium (manual) |
| Amplification | LinkedIn |
| Secret detection | Husky pre-push hook (grep-based) |
| Telegram notifications | Manual or workflow dispatch via Telegram Bot API |

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

`prompts/dependency-review/` contains reviewed research briefs for Go modules, GitHub Actions, npm packages, Terraform providers, and standalone tools.
PR validation also runs mutation testing against `internal/frontmatter`.
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

Push changes to `posts/` on the `main` branch or run the workflow manually. The GitHub Actions pipeline scans the local `posts/` tree and reconciles it against Hashnode and Dev.to so missing posts are published and changed posts are updated.

Posts with `draft: true` in the front matter are skipped by the publishing pipeline.

Medium cross-posting is manual via the "Import a story" feature using the canonical URL.

### Notifications

Telegram notifications can be sent manually from the repository root:

```bash
go run ./cmd/notify "Post 4 is live. Monitor engagement."
```

Or via the `Send Notification` GitHub Actions workflow using the `message` input.

The scheduled publish run also sends a daily Telegram notification when there is something to report.
That notification includes posts queued for tomorrow plus target-level publish changes from the scheduled run. Unchanged posts are not reported.

Required secrets:

- `TELEGRAM_BOT_TOKEN`: Telegram bot token created via BotFather
- `TELEGRAM_CHAT_ID`: Chat ID for the phone or chat that should receive notifications

## Licence

Content (posts, documentation) is copyright Christo Zietsman. Code (Terraform, workflows, scripts) is MIT licensed.
