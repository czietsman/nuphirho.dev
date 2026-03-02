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
