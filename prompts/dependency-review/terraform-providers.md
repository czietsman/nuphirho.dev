# Terraform Provider Review

## Purpose

Use this brief to record reviewed decisions for Terraform provider
updates.

## Review metadata

- Date: 2026-09-28
- Reviewer: pending owner review of the pull request's `terraform plan` comment
- Scope: `terraform/`, Cloudflare provider v4 to v5 (supersedes dependabot PR #168)

## Approved provider updates

| Provider | Current | Recommended | Reason | Risk |
|---|---|---|---|---|
| `registry.terraform.io/cloudflare/cloudflare` | `~> 4.0` (resolves to 4.52.9) | `~> 5.26` | v4 is superseded by v5; dependabot proposed `~> 5.25` | High: resource renames and schema changes across all 14 managed resources; DNS and email routing are production-critical |

## Security notes

- No new secrets or permissions. The provider still authenticates with `CLOUDFLARE_API_TOKEN`.

## Provenance and trust notes

- Record provider source addresses exactly.
- Note any state, schema, or plan-review implications.
- Provider binaries used for local validation (4.52.9, 5.26.0) were downloaded from the provider's GitHub releases and matched the published SHA256SUMS. Terraform 1.9.8 came from releases.hashicorp.com and matched its SHA256SUMS.
- Configuration changes were generated with Cloudflare's `tf-migrate` (built from github.com/cloudflare/tf-migrate at 2bd5ec4) and then applied by hand. `tf-migrate` has no transformer for `cloudflare_email_routing_rule`, so that change was written against the v5.26.0 provider schema (`terraform providers schema -json`).
- State: six `cloudflare_record` resources move to `cloudflare_dns_record` through `moved` blocks. Cross-type moves require Terraform 1.8 or later, so `required_version` is now `>= 1.8`. The other resources rely on the v5 provider's built-in state upgraders.
- The `blog_record`, `root_record` and `www_record` outputs now return the record name (for example `blog`) instead of the full hostname, because v5 has no `hostname` attribute. Nothing in the repository reads these outputs.
- `terraform validate` passes with v5.26.0. It was not possible to run `terraform plan` locally because the state backend and API credentials are unavailable, so the plan comment on the pull request is the review gate.
- Accept the plan only if every resource is moved or updated in place. Any destroy or replace, especially of a DNS record or email routing rule, blocks the merge, because merging to main runs `terraform apply`.

## Execution instructions

- Update only the approved providers listed above.
- Run `terraform plan` and review the output before any apply.
- Do not apply unreviewed provider changes.
