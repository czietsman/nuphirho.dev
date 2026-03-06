---
title: "Process Over Technology: Starting With the Blog Itself"
slug: process-over-technology-starting-with-the-blog-itself
subtitle: "Building a blog with BDD specs, Terraform, and a CI/CD pipeline"
tags: [process, engineering, ai, devops, infrastructure-as-code]
---

# Process Over Technology: Starting With the Blog Itself

I have been thinking about starting a blog for a while. Not because the world needs another tech blog, but because I needed a place to think out loud about something I keep coming back to: process matters more than technology.

AI has changed the economics of rigorous engineering. Practices that used to be too expensive or too slow for most teams, things like executable specifications, mutation testing, and formal verification layers, are now economically viable. The tooling is free. The compute is cheap. The only thing standing in the way is how we think about building software.

So when I finally sat down to build this blog, I decided to treat it the way I believe all software should be built. Not as a weekend side project where I pick a static site generator and start writing, but as a properly managed effort with a defined process, infrastructure as code, a publishing pipeline, and security defaults baked in from the start.

This post is the story of that build. Not a tutorial. Just the decisions, the reasoning, and what it cost.

## Starting with why, not how

The first decision was deliberate: no code until the process was defined.

Before choosing a platform, a theme, or a static site generator, I wrote a style guide. It captures the voice, tone, formatting rules, and editorial standards for every post. British English. No em dashes. No emoji. Paragraphs over bullet points. Target length of 1,200 to 1,800 words. Always attribute other people's work.

Then I wrote BDD specifications in Gherkin describing how the publishing pipeline should behave. What happens when I push a markdown file to the main branch? What if the post is marked as a draft? What if it already exists on the target platform?

The pipeline behaviour was fully specified before a single line of workflow YAML existed. This is the approach I advocate for in software delivery, and it felt wrong to skip it for my own project.

## Choosing the platform

The platform decision came down to a simple question: where does the content live, and who owns it?

I chose Hashnode as the primary platform with a custom domain. It is free, supports custom domains at no cost, has a GraphQL API for automated publishing, handles light and dark themes natively, and supports Mermaid diagrams (which are version-controllable as code). The built-in developer community provides discoverability without me having to build an audience from scratch.

Cross-posting goes to Dev.to via their REST API, automated through the same pipeline. Medium is a manual step using their URL import feature. Their API no longer issues new integration tokens, so automation is not an option. The pipeline handles this gracefully: it is documented, repeatable, and takes about 30 seconds.

Every cross-posted article sets its canonical URL back to blog.nuphirho.dev. This is non-negotiable. The custom domain builds SEO authority over time. The platforms provide reach. Both matter, but ownership comes first.

## The domain

The root domain, nuphirho.dev, is kept deliberately separate from the blog. The blog lives at blog.nuphirho.dev. The root hosts a simple static landing page on GitHub Pages, leaving it flexible for whatever comes next.

The .dev TLD was a conscious choice. It sits on the HSTS preload list, which means browsers refuse to load it over plain HTTP. You do not get the option to be insecure. Cloudflare adds a second layer of SSL/TLS and CDN. Hashnode auto-provisions a Let's Encrypt certificate for the custom domain. That is HTTPS enforced at three independent layers, before a single word of content is published.

In a world where AI is generating code and people are shipping software they do not fully understand, security defaults matter more, not less. This blog enforces that principle from the infrastructure up.

## The name

Nu, phi, rho. Three Greek letters I picked up studying physics and mathematics at Stellenbosch University. They stuck as a username during university and never left. The name reflects where I started: grounded in rigour, pattern recognition, and first principles thinking. It is sentimental, not a brand exercise.

## Infrastructure as code

Everything is managed with Terraform. The Cloudflare DNS configuration, including the blog CNAME pointing to Hashnode, the root A records pointing to GitHub Pages, and the www redirect, is all declared in version-controlled HCL files.

Terraform state is stored in Cloudflare R2, which is S3-compatible and sits within the free tier. This means the entire infrastructure layer, DNS, CDN, SSL, state management, is declarative, reproducible, and costs nothing.

Some might call this overkill for a blog. I call it the point. If we are going to argue that enterprise-grade practices are accessible, we need to demonstrate it. Terraform for a blog is not about complexity. It is about showing that the barrier to doing things properly has collapsed.

## The publishing pipeline

The pipeline runs on GitHub Actions, which is free and unlimited for public repositories. The workflow is straightforward:

When a markdown file in the posts directory is pushed to the main branch, the pipeline reads the frontmatter, determines whether the post is new or an update, and publishes it. If the post is marked as a draft, it gets pushed as an unpublished draft to both Hashnode and Dev.to rather than being skipped entirely. This mirrors how code deployments work: you can push to staging without going to production.

Hashnode uses a GraphQL API. Dev.to uses a REST API. The pipeline handles both, sets canonical URLs, manages tags, and reports a summary of what was published and where. Draft posts are checked for duplicates to avoid pushing the same draft twice.

All API tokens and credentials live in GitHub Secrets. Nothing sensitive touches the repository.

## Secrets at the boundary

Speaking of security: there is a Husky pre-push hook that scans for AWS keys, GitHub tokens, PEM headers, and generic secret patterns before any code leaves the machine. It is a simple check, but it catches the most common mistakes at the earliest possible point.

This is defence in depth applied to a blog. Three layers of HTTPS. Secrets in a vault, not in code. Scanning at the git boundary. None of this is complex. All of it is free. The only cost is deciding to do it.

## On AI assistance

This blog is AI-assisted. I want to be upfront about that because it connects directly to the thesis.

I think in systems and architecture. I do not always communicate those ideas clearly on the first pass. AI helps me bridge that gap. The thinking is mine. The clarity is a collaboration.

This post was drafted with AI assistance. The decisions, the architecture, the reasoning, those are mine. The process of turning those thoughts into clear prose was a collaboration. I believe this is how software will increasingly be built: human judgement and accountability, with AI handling the parts that benefit from scale and speed.

Being transparent about this is not a disclaimer. It is a demonstration.

## What it cost

Let me be specific.

| Concern | Tool | Cost |
|---|---|---|
| Source control | GitHub (public) | Free |
| CI/CD | GitHub Actions | Free |
| IaC | Terraform + Cloudflare provider | Free |
| Terraform state | Cloudflare R2 | Free |
| DNS/CDN/SSL | Cloudflare free tier | Free |
| Landing page | GitHub Pages | Free |
| Blog platform | Hashnode | Free |
| Cross-post | Dev.to (automated), Medium (manual) | Free |
| Secret detection | Husky pre-push hook | Free |
| Domain | Cloudflare Registrar | ~$12/year |

The domain registration is the only line item. Everything else, the infrastructure, the pipeline, the platform, the security layers, is enterprise-grade tooling at zero cost.

This is the argument made tangible. The economics have changed. The practices that used to require dedicated teams and significant budgets are available to anyone willing to apply the process.

## What comes next

This post is the first in what I hope becomes a regular practice: writing about what I am working on, what I am learning, and what I get wrong. The topics will span AI-assisted software delivery, engineering process, organisational transformation, and the practical challenges of introducing new tools and practices to existing teams.

The entire source for this blog, including the infrastructure, pipeline, specifications, and style guide, is public at [github.com/czietsman/nuphirho.dev](https://github.com/czietsman/nuphirho.dev). If the process interests you, the receipts are there.

The technology was the easy part. The process is what made it work.



