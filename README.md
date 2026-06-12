# nuphirho.dev

Enterprise-grade engineering at startup speed.

A technical blog by [Christo Zietsman](https://www.linkedin.com/in/christo-zietsman/), managed like a software project. The infrastructure, pipeline, and process are all in this repository.

## What is this

This is the source repository for [nuphirho.dev](https://nuphirho.dev), a technical blog about AI-assisted software delivery, engineering process, and technology innovation. The root domain is a SvelteKit static site served via Cloudflare Pages. The blog lives at [blog.nuphirho.dev](https://blog.nuphirho.dev), also a SvelteKit app on Cloudflare Pages, with posts cross-posted to Dev.to.

The repository contains:

- **blog/** -- SvelteKit app for blog.nuphirho.dev (Cloudflare Pages + Workers)
  - **blog/src/routes/** -- SvelteKit routes (post list, post detail, stats API)
  - **blog/src/lib/posts.ts** -- Markdown post loading and parsing
- **main-site/** -- SvelteKit app for nuphirho.dev (Cloudflare Pages)
  - **main-site/src/routes/** -- Landing page, about, words-of-meaning, novel-findings, roadmap, privacy, cookies
  - **main-site/src/lib/Roadmap.svelte** -- Publishing calendar component
- **site/** -- Go template static site builder (legacy; replaced by main-site/)
  - **site/templates/** -- HTML templates (base layout, partials, page content)
  - **site/css/** -- Stylesheets (shared + page-specific)
  - **site/js/** -- JavaScript (theme toggle, roadmap calendar)
  - **site/static/** -- Files copied directly to output (CNAME, calendar data)
- **posts/** -- Markdown blog post source files with YAML frontmatter
- **cmd/** -- Go CLI tools
  - **cmd/publish** -- Reconcile posts/ with Hashnode and Dev.to
  - **cmd/notify** -- Send Telegram notifications
  - **cmd/notify-summary** -- Scheduled notification digest
  - **cmd/site-build** -- Build static site from Go templates
  - **cmd/validate-tags** -- Validate post tag values
- **internal/** -- Shared Go packages
  - **internal/hashnode** -- Hashnode GraphQL API client
  - **internal/devto** -- Dev.to REST API client
  - **internal/frontmatter** -- Post metadata schema and parsing
  - **internal/pipeline** -- Publishing orchestration logic
  - **internal/tags** -- Tag validation
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
| Cross-post | Dev.to (automated) |
| Amplification | LinkedIn |
| Secret detection | Husky pre-push hook (grep-based) |
| Telegram notifications | Manual or workflow dispatch via Telegram Bot API |

The domain is the only cost.

## GitHub Actions workflows

| Workflow | Trigger | Purpose |
|---|---|---|
| `blog.yml` | Push to main, paths: `posts/**` `blog/**` | Build and deploy blog to Cloudflare Pages |
| `main-site.yml` | Push to main, paths: `main-site/**` | Build and deploy main site to Cloudflare Pages |
| `publish.yml` | Push to main, paths: `posts/**`; daily cron | Reconcile posts/ with Dev.to |
| `pages.yml` | Push to main, paths: `site/**` | Build Go template site to GitHub Pages (legacy) |
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
- Assignments to known secret variables (`CLOUDFLARE_API_TOKEN`, `HASHNODE_TOKEN`, `DEVTO_API_KEY`, and others)
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

### Go pipeline

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

Push changes to `posts/` on the `main` branch or run the workflow manually. The GitHub Actions pipeline scans the local `posts/` tree and reconciles it against Dev.to so missing posts are published and changed posts are updated.

Posts with `draft: true` in the front matter are skipped by the publishing pipeline. Posts with a future `publish_date` are skipped until that date arrives. The publish cron runs at 05:00 UTC daily, so time-of-day scheduling is not available.

### Notifications

Telegram notifications can be sent manually from the repository root:

```bash
go run ./cmd/notify "Post 4 is live. Monitor engagement."
```

Or via the `Send Notification` GitHub Actions workflow using the `message` input.

The scheduled publish run also sends a daily Telegram notification when there is something to report. That notification includes posts queued for tomorrow plus target-level publish changes from the scheduled run. Unchanged posts are not reported, and scheduled publish failures are summarised there instead of failing the workflow outright.

Required secrets:

- `TELEGRAM_BOT_TOKEN`: Telegram bot token created via BotFather
- `TELEGRAM_CHAT_ID`: Chat ID for the phone or chat that should receive notifications

## Licence

Content (posts, documentation) is copyright Christo Zietsman. Code (Terraform, workflows, scripts) is MIT licensed.
