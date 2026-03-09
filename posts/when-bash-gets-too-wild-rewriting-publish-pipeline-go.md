---
title: "When bash gets too wild: rewriting my publish pipeline in Go"
slug: "when-bash-gets-too-wild-rewriting-publish-pipeline-go"
draft: true
tags: ["go", "bdd", "devops", "software-engineering", "blogging"]
---

# When bash gets too wild: rewriting my publish pipeline in Go

There is a particular kind of technical debt that only grows in shell scripts. It starts innocently: a `sed` command to strip some quotes, a `grep` to pull out a frontmatter field, a `jq` one-liner to map tags. Each fix is small and defensible in isolation. But give it a few months and you have something that works in exactly the conditions it was written for and fails in ways that are almost impossible to reason about.

That is where I ended up with the publish pipeline for this blog -- and it did not even take months. We built it fresh, but the pace of iteration compressed what normally takes a year into days. The workarounds accumulated just as fast.

The pipeline cross-posts markdown files to Hashnode and Dev.to via their APIs. It lives in a GitHub Actions workflow, and it started failing almost immediately: quoted YAML values, tags with hyphens, multi-line frontmatter. Each failure led to another workaround: quote-stripping logic, a tag glossary loaded via `jq`, conditionals that were more incantation than logic. The script became the kind of thing you are afraid to touch because you cannot be certain what it is actually doing.

The honest diagnosis: I had reached the natural ceiling of bash for this kind of work. Bash has no good answer to "how do I test this?" And if you cannot test it, you cannot trust it.

## The cost of untestable pipelines

This is not unique to personal blog tooling. Infrastructure scripts, CI workflows, deployment pipelines -- these are often written in bash because bash is always available and the task seems small at the time. The problem is that "small at the time" is doing a lot of work in that sentence. Pipelines accumulate logic. They handle edge cases. They become load-bearing parts of a development workflow. And because they are in bash, they remain untestable by default.

The failure mode is not usually a catastrophic crash. It is subtler: a post that silently publishes with the wrong tags, a frontmatter field that gets parsed incorrectly and corrupts a slug, a workaround that fixes one case while breaking another you did not think to check. You discover these failures after the fact, in production, with no good way to reproduce them locally.

I have spent a lot of time over the past year thinking about what rigorous software engineering actually looks like when AI is in the picture. One of the core arguments I keep coming back to is that AI has collapsed the cost of doing things properly -- not just the glamorous parts like architecture and feature development, but the unglamorous parts too: tests, specs, validation, structured error handling. If the tooling for rigour is now cheaper to write, there is less justification for skipping it. That includes pipelines.

## The rewrite

The goal was straightforward: replace the inline bash in `.github/workflows/publish.yml` with a Go CLI tool. The workflow calls the binary; all logic lives in testable Go code with BDD specifications written before any implementation.

The package structure reflects the actual problem decomposition: `internal/frontmatter` for YAML parsing, `internal/tags` for glossary loading and per-platform mapping, `internal/hashnode` and `internal/devto` for the API clients, `internal/probe` for contract validation, `internal/pipeline` for the orchestrator, and `cmd/publish` as the thin CLI wrapper. The full plan is in the [nuphirho.dev repository](https://github.com/czietsman/nuphirho.dev).

The end result: 98 BDD scenarios, 488 steps across 7 packages, all passing. Around 150 lines of bash collapse into two workflow lines:

```yaml
- name: Build publish tool
  run: go build -o publish ./cmd/publish/

- name: Publish posts
  env:
    HASHNODE_TOKEN: ${{ secrets.HASHNODE_TOKEN }}
    HASHNODE_PUBLICATION_ID: ${{ secrets.HASHNODE_PUBLICATION_ID }}
    DEVTO_API_KEY: ${{ secrets.DEVTO_API_KEY }}
  run: ./publish --tags-file tags.json ${{ steps.changed.outputs.posts }}
```

But the line count is not the interesting part. What is interesting is what the BDD scenarios caught.

## What the specs caught

The first package -- `internal/frontmatter` -- is where most of the bash failures had originated. Writing the scenarios first forced precision about what "correct" actually means. That precision immediately exposed bugs that had been hiding in plain sight.

The secret detection regex in the bash pipeline was `[A-Za-z0-9+/]{20,}`. It did not match GitHub personal access tokens (`ghp_...`) because it excluded underscores and hyphens. The BDD scenario caught this on the first run. The fix was adding `_-` to the character class. One line. But without a scenario that explicitly tested a GitHub token pattern, that bug would have kept silently passing anything with a `ghp_` prefix as clean.

The `gopkg.in/yaml.v3` library eliminated a whole category of problems. The bash pipeline needed separate handling for three tag notation formats: `[a, b]`, `["a", "b"]`, and YAML list. The Go library handles all three identically, returning `[]string` in all cases. That alone justified the rewrite.

Moving to `internal/hashnode`, the GraphQL client needed a helper for navigating untyped JSON responses. The `gqlNode.path()` helper was written to chain safely on nil -- missing fields return empty strings rather than panics. What it did not handle correctly was a null JSON value: `publication: null`. The initial check was `if pub == nil`, which passed because the node existed, just with a nil raw value. The "publication not found" error was never raised. The BDD scenario caught it. Fixed by checking `pub.isNull()` instead. This is the same class of nil-vs-empty bug that bites Go developers with nil interfaces. Without the scenario, it would have manifested as a confusing API error message at publish time.

The `internal/probe` package had the most satisfying catch. The contract probe output format was specified exactly in the plan: `[hashnode] credentials        OK` with precise column padding. The initial implementation used `fmt.Sprintf("[%-8s]", platform)`, which pads inside the brackets: `[devto   ]`. The spec required padding after the closing bracket: `[devto]   `. Five of eleven scenarios failed on that single formatting difference. One line fix. All eleven passing. The spec was doing its job.

## Patterns, not just passes

Something else happened during the implementation that was worth noting. The first package -- `internal/hashnode` -- required iteration. The nil-vs-empty bug needed fixing. The fake HTTP client routing needed thought. Two scenarios failed on the first run and required code changes.

Every package after that passed on the first run. Tags, Dev.to, probe, pipeline, CLI -- five consecutive packages with no post-implementation fixes. That is not coincidence. The patterns established in the hashnode package -- injected HTTP client interface, fake routing by method and path, result structs carrying action and IDs, `io.Writer` for testable output -- transferred cleanly to every subsequent package. By the time the orchestrator was being written, the shape of the solution was so well established that the implementation matched the spec without iteration.

This is one of the underappreciated benefits of BDD when it is done in order. The first package is where you pay the design tax. Everything after inherits the patterns.

## Two decisions worth explaining

Two design decisions in `internal/pipeline` are worth being explicit about, because both have non-obvious implications.

The first is two-phase execution. Phase 1 validates all post files. If any fail validation, the pipeline exits with code 1 before making any API calls. Phase 2 processes each file independently. If Hashnode fails for one file, Dev.to is skipped for that file but other files continue. This means validation is a hard gate -- you never partially publish a batch with a broken post in it -- but publish failures are isolated, not catastrophic.

The second is no rollback. If file 1 publishes successfully to both platforms and file 2 then fails on Hashnode, file 1's results are preserved. The summary reports both the success and the failure. Exit code 2. There is no mechanism to undo file 1's publish. This matches the bash pipeline's behaviour, and it is the right default: a post that has been successfully published should stay published. The failure is in the tooling, not the content.

Both decisions are explicitly tested in BDD scenarios. Writing them as scenarios forces you to state the expected behaviour unambiguously before you implement it.

## The migration strategy

The switch is not a flag day. The Go binary is running in `--dry-run` mode in a parallel workflow job alongside the existing bash pipeline, uploading structured JSON output as a GitHub Actions artifact. Once the dry-run output matches the bash pipeline's behaviour across several real pushes, Hashnode cuts over first -- keeping Dev.to on bash isolates the risk to one platform. Once that is stable, Dev.to follows, the bash pipeline is removed, and the multi-job workflow collapses into a single `./publish` invocation.

The `--dry-run` flag produces structured JSON rather than human-readable logs specifically for this phase. Diffing JSON is straightforward. Diffing log output is not.

## The broader point

This is a small project. It publishes markdown files to two APIs. No reasonable person would call it critical infrastructure.

But the same thinking applies at every scale. Pipelines accumulate logic. Untestable logic accumulates risk. And the cost of writing testable, well-specified code has come down dramatically. This rewrite -- 7 packages, 98 scenarios, 488 steps, full BDD coverage -- happened in days, not months. That is the part that changes the calculus. There is less and less justification for treating pipelines as second-class citizens that do not deserve the same engineering rigour as application code.

The secret regex bug, the nil-vs-empty GraphQL node, the five failing probe scenarios from a single misplaced padding character -- none of these are dramatic failures. They are exactly the kind of subtle, silent, hard-to-reproduce problems that bash pipelines accumulate over time and that a proper test suite catches before they reach production.

The bash script served its purpose. It got the pipeline working. Now it has earned the right to be replaced by something that can actually be trusted.

I will write a follow-up once the migration is complete -- specifically on running go-mutesting against the frontmatter parser, which is the highest-consequence code in the tool and the place where mutation coverage matters most.
