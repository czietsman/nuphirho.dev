# Agent Instructions

These rules apply to every AI agent working on this repository. No exceptions.

## Test-driven development

Follow strict TDD for every code change:

1. Write a failing test or BDD scenario first. Confirm it fails.
2. Write the minimum code to make it pass. Confirm it passes.
3. Refactor. Confirm all tests still pass.

Do not skip step 1. Do not write production code without a failing test.

## BDD specifications

Every change to pipeline logic, validation, or publishing behaviour must have a corresponding Gherkin scenario in `specs/`. If a scenario does not exist for the behaviour you are changing, write one before writing the code.

## No backwards-compatibility code

Delete dead code. Do not keep unused functions, deprecated parameters, re-exports, compatibility shims, or commented-out code. If removing something would break an external contract, raise it with the user before proceeding. Do not silently preserve it.

## Keep README.md current

If a change affects the repository structure, the stack, the setup instructions, or the publishing workflow, update `README.md` in the same commit. The README must always reflect the current state of the project.

## Style guide

All written content, including commit messages and documentation, must follow `docs/STYLE_GUIDE.md`. British English. No em dashes. No emoji.

---

## Security context

This is a public repository. The full directory structure, workflow definitions, Terraform configuration, secret variable names, and pipeline logic are visible to anyone. Treat this as a given constraint, not a problem to solve. Every decision in this repository must be made with that visibility in mind.

## Agent scope

Agents operating in this repository are execution agents only. Their scope is strictly limited to acting on reviewed, human-authored prompts.

Agents must not perform research tasks. This includes but is not limited to:

- Searching for, evaluating, or selecting libraries or dependencies
- Fetching or reading external URLs, documentation, or web content
- Querying package registries or version APIs
- Looking up anything outside the repository to resolve ambiguity

If a task requires external information to complete, stop and raise it with the user. Do not attempt to resolve it independently. Do not make a best-effort guess based on training data as a substitute for research.

Research tasks are handled in a separate, sandboxed context with no write access to the repository or pipeline. That context produces a reviewed prompt. This file is the boundary between that context and this one.

## Secret hygiene

The pre-push hook in `.husky/` scans for secrets before code leaves the development machine. Do not bypass it with `--no-verify` except for documented false positives listed in `.secretscanignore`.

When writing workflow YAML, scripts, or Terraform configuration:

- Do not hardcode values that belong in GitHub Secrets or environment variables
- Do not log, echo, or print secret variable values, even partially
- Do not construct secret values from fragments in ways that would survive grep-based detection

If a change requires a new secret, document the variable name and purpose in `README.md` under the stack table. Do not create the secret itself.

## GitHub Actions permissions

Every workflow job must declare a `permissions:` block. Grant only what the job requires for that specific job. Default to `contents: read`. Do not omit the block and rely on repository defaults.

## Terraform discipline

Run `terraform plan` and review the output before every `terraform apply`. Do not apply without a reviewed plan. Do not commit Terraform state files. If a state file appears in the working directory, add it to `.gitignore` and raise it with the user before proceeding.

## Why this matters

This repository is a deliberate, high-stakes experiment. The author has chosen a public repository with a live publishing pipeline and real infrastructure as the environment for applying and testing engineering and security practices. The low error margin is intentional. It sharpens thinking and produces real signal. The constraints above are not theoretical — they reflect threat modeling applied directly to this system, and they parallel practices under evaluation in a professional context. Treat them accordingly.
