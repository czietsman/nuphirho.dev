# nuphirho.dev

Enterprise-grade engineering at startup speed.

A technical blog by [Christo Zietsman](https://www.linkedin.com/in/christo-zietsman/), managed like a software project. The infrastructure, pipeline, and process are all in this repository.

## What is this

This is the source repository for [nuphirho.dev](https://nuphirho.dev), a technical blog about AI-assisted software delivery, engineering process, and technology innovation. The root domain and [blog.nuphirho.dev](https://blog.nuphirho.dev) are SvelteKit applications served via Cloudflare Pages. Posts are published only to the first-party blog.

The repository contains:

- **blog/** -- SvelteKit app for blog.nuphirho.dev (Cloudflare Pages + Workers)
  - **blog/src/routes/** -- SvelteKit routes (post list, post detail, stats API)
  - **blog/src/lib/posts.ts** -- Markdown post loading and parsing
- **main-site/** -- SvelteKit app for nuphirho.dev (Cloudflare Pages)
  - **main-site/src/routes/** -- Landing page, about, words-of-meaning, novel-findings, roadmap, privacy, cookies
  - **main-site/src/lib/Roadmap.svelte** -- Publishing calendar component
- **posts/** -- Markdown blog post source files with YAML frontmatter
- **cmd/** -- Go CLI tools
  - **cmd/notify** -- Send Telegram notifications
- **internal/** -- Shared Go packages
  - **internal/frontmatter** -- Post metadata schema and parsing
- **terraform/** -- Cloudflare infrastructure as code (DNS, Pages projects, KV)
- **.github/workflows/** -- GitHub Actions CI/CD pipelines
- **prompts/** -- Reviewed prompt material, including dependency review briefs
- **specs/** -- BDD feature files describing pipeline behaviour
- **tests/** -- Shell-based test scripts
- **docs/** -- Project brief and style guide
- **papers/** -- Academic paper build infrastructure
- **experiments/** -- Research projects

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
| Frontend framework | SvelteKit + adapter-cloudflare |
| Landing page | Cloudflare Pages (nuphirho.dev) |
| Blog hosting | Cloudflare Pages + Workers (blog.nuphirho.dev) |
| Blog visitor counter | Cloudflare KV |
| Amplification | LinkedIn |
| Secret detection | Husky pre-push hook (grep-based) |
| Telegram notifications | Manual or workflow dispatch via Telegram Bot API |

The domain is the only cost.

## GitHub Actions workflows

| Workflow | Trigger | Purpose |
|---|---|---|
| `blog.yml` | Post or blog changes; daily at 05:00 UTC; manual dispatch | Build and deploy the blog, then report the result to Telegram |
| `main-site.yml` | Push to main, paths: `main-site/**` | Build and deploy main site to Cloudflare Pages |
| `terraform.yml` | Push/PR to main, paths: `terraform/**` | Plan and apply Cloudflare infrastructure |
| `validate-pr.yml` | Pull requests | Run Go tests, linters, mutation testing |
| `notify.yml` | Workflow dispatch | Send manual Telegram notification |

## Development setup

```bash
npm install
```

This installs [husky](https://typicode.github.io/husky/) and configures a `pre-push` git hook that scans for secrets before code leaves your machine. The hook catches:

- AWS / R2 access key IDs (`AKIA...`)
- GitHub token variants (`ghp_`, `github_pat_`, `gho_`, `ghu_`, `ghr_`, `ghs_`)
- PEM private key headers
- Assignments to known secret variables (`CLOUDFLARE_API_TOKEN` and others)
- Generic secret patterns (`api_key`, `token`, `password`, etc. followed by long values)

Paths listed in `.secretscanignore` are excluded from scanning. To bypass the hook for a known false positive, use `git push --no-verify`.

To run the pattern tests:

```bash
bash tests/test-secret-patterns.sh
```

### Blog and main site (SvelteKit)

```bash
cd blog && npm install && npm run dev
cd main-site && npm install && npm run dev
```

Both apps use `@sveltejs/adapter-cloudflare` and prerender all pages at build time.

### Go tools

```bash
go test ./...
```

PR validation also runs mutation testing against `internal/frontmatter`.

`prompts/dependency-review/` contains reviewed research briefs for Go modules, GitHub Actions, npm packages, Terraform providers, and standalone tools.

## Getting started

### Terraform

```bash
cd terraform
terraform init
terraform plan
terraform apply
```

Requires `CLOUDFLARE_API_TOKEN`, `CLOUDFLARE_ACCOUNT_ID`, and R2 backend credentials in GitHub Secrets or as environment variables.

### Publishing

Push changes to `posts/` on the `main` branch or run the blog workflow manually. GitHub Actions builds the SvelteKit blog from the Markdown source and deploys it to the `nuphirho-blog` Cloudflare Pages project.

Posts with `draft: true` in the front matter are excluded from the blog build. Posts with a future `publish_date` remain excluded until that date arrives. The blog is rebuilt at 05:00 UTC daily so scheduled posts become visible without an additional commit. Time-of-day scheduling is not available.

### Notifications

Telegram notifications can be sent manually from the repository root:

```bash
go run ./cmd/notify "Post 4 is live. Monitor engagement."
```

Or via the `Send Notification` GitHub Actions workflow using the `message` input.

Production blog workflow runs send the final deployment status and GitHub Actions run URL to Telegram. Pull-request preview deployments do not send notifications. Telegram delivery is non-blocking, so a notification failure does not change the blog deployment result.

Required secrets:

- `TELEGRAM_BOT_TOKEN`: Telegram bot token created via BotFather
- `TELEGRAM_CHAT_ID`: Chat ID for the phone or chat that should receive notifications

## Licence

Content (posts, documentation) is copyright Christo Zietsman. Code (Terraform, workflows, scripts) is MIT licensed.
