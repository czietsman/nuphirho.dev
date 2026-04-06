---
source: https://raw.githubusercontent.com/apache/airflow/main/AGENTS.md
commit: 51b7d3b1eb705e48016d1ffeccd7de5c3be7699c
captured: 2026-04-05
lines: 216
---
 <!-- SPDX-License-Identifier: Apache-2.0
      https://www.apache.org/licenses/LICENSE-2.0 -->

# AGENTS instructions

## Environment Setup

- Install prek: `uv tool install prek`
- Enable commit hooks: `prek install`
- **Never run pytest, python, or airflow commands directly on the host** — always use `breeze`.
- Place temporary scripts in `dev/` (mounted as `/opt/airflow/dev/` inside Breeze).

## Commands

`<PROJECT>` is folder where pyproject.toml of the package you want to test is located. For example, `airflow-core` or `providers/amazon`.
`<target_branch>` is the branch the PR will be merged into — usually `main`, but could be `v3-1-test` when creating a PR for the 3.1 branch.

- **Run a single test:** `uv run --project <PROJECT> pytest path/to/test.py::TestClass::test_method -xvs`
- **Run a test file:** `uv run --project <PROJECT> pytest path/to/test.py -xvs`
- **Run all tests in package:** `uv run --project <PROJECT> pytest path/to/package -xvs`
- **If uv tests fail with missing system dependencies, run the tests with breeze**: `breeze run pytest <tests> -xvs`
- **Run a Python script:** `uv run --project <PROJECT> python dev/my_script.py`
- **Run core or provider tests suite in parallel:** `breeze testing <test_group> --run-in-parallel` (test groups: `core-tests`, `providers-tests`)
- **Run core or provider db tests suite in parallel:** `breeze testing <test_group> --run-db-tests-only --run-in-parallel` (test groups: `core-tests`, `providers-tests`)
- **Run core or provider non-db tests suite in parallel:** `breeze testing <test_group> --skip-db-tests --use-xdist` (test groups: `core-tests`, `providers-tests`)
- **Run single provider complete test suite:** `breeze testing providers-tests --test-type "Providers[PROVIDERS_LIST]"` (e.g., `Providers[google]` or `Providers[amazon]` or "Providers[amazon,google]")
- **Run Helm tests in parallel with xdist** `breeze testing helm-tests --use-xdist`
- **Run Helm tests with specific K8s version:** `breeze testing helm-tests --use-xdist --kubernetes-version 1.35.0`
- **Run specific Helm test type:** `breeze testing helm-tests --use-xdist --test-type <type>` (types: `airflow_aux`, `airflow_core`, `apiserver`, `dagprocessor`, `other`, `redis`, `security`, `statsd`, `webserver`)
- **Run other suites of tests** `breeze testing <test_group>` (test groups: `airflow-ctl-tests`, `docker-compose-tests`, `task-sdk-tests`)
- **Run scripts tests:** `uv run --project scripts pytest scripts/tests/ -xvs`
- **Run Airflow CLI:** `breeze run airflow dags list`
- **Type-check:** `breeze run mypy path/to/code`
- **Lint with ruff only:** `prek run ruff --from-ref <target_branch>`
- **Format with ruff only:** `prek run ruff-format --from-ref <target_branch>`
- **Run regular (fast) static checks:** `prek run --from-ref <target_branch> --stage pre-commit`
- **Run manual (slower) checks:** `prek run --from-ref <target_branch> --stage manual`
- **Build docs:** `breeze build-docs`
- **Determine which tests to run based on changed files:** `breeze selective-checks --commit-ref <commit_with_squashed_changes>`

SQLite is the default backend. Use `--backend postgres` or `--backend mysql` for integration tests that need those databases. If Docker networking fails, run `docker network prune`.

## Repository Structure

UV workspace monorepo. Key paths:

- `airflow-core/src/airflow/` — core scheduler, API, CLI, models
  - `models/` — SQLAlchemy models (DagModel, TaskInstance, DagRun, Asset, etc.)
  - `jobs/` — scheduler, triggerer, Dag processor runners
  - `api_fastapi/core_api/` — public REST API v2, UI endpoints
  - `api_fastapi/execution_api/` — task execution communication API
  - `dag_processing/` — Dag parsing and validation
  - `cli/` — command-line interface
  - `ui/` — React/TypeScript web interface (Vite)
- `task-sdk/` — lightweight SDK for Dag authoring and task execution runtime
  - `src/airflow/sdk/execution_time/` — task runner, supervisor
- `providers/` — 100+ provider packages, each with its own `pyproject.toml`
- `airflow-ctl/` — management CLI tool
- `chart/` — Helm chart for Kubernetes deployment
- `dev/` — development utilities and scripts used to bootstrap the environment, releases, breeze dev env
- `scripts/` — utility scripts for CI, Docker, and prek hooks (workspace distribution `apache-airflow-scripts`)
  - `ci/prek/` — prek (pre-commit) hook scripts; shared utilities in `common_prek_utils.py`
  - `tests/` — pytest tests for the scripts; run with `uv run --project scripts pytest scripts/tests/`


## Architecture Boundaries

1. Users author Dags with the Task SDK (`airflow.sdk`).
2. Dag Processor parses Dag files in isolated processes and stores serialized Dags in the metadata DB.
3. Scheduler reads serialized Dags — **never runs user code** — and creates Dag runs / task instances.
4. Workers execute tasks via Task SDK and communicate with the API server through the Execution API — **never access the metadata DB directly**.
5. API Server serves the React UI and handles all client-database interactions.
6. Triggerer evaluates deferred tasks/sensors in isolated processes.
7. Shared libraries that are symbolically linked to different Python distributions are in `shared` folder.
8. Airflow uses `uv workspace` feature to keep all the distributions sharing dependencies and venv
9. Each of the distributions should declare other needed distributions: `uv --project <FOLDER> sync` command acts on the selected project in the monorepo with only dependencies that it has

# Shared libraries

- shared libraries provide implementation of some common utilities like logging, configuration where the code should be reused in different distributions (potentially in different versions)
- we have a number of shared libraries that are separate, small Python distributions located under `shared` folder
- each of the libraries has it's own src, tests, pyproject.toml and dependencies
- sources of those libraries are symbolically linked to the distributions that are using them (`airflow-core`, `task-sdk` for example)
- tests for the libraries (internal) are in the shared distribution's test and can be run from the shared distributions
- tests of the consumers using the shared libraries are present in the distributions that use the libraries and can be run from there

## Coding Standards

- **Always format and check Python files with ruff immediately after writing or editing them:** `uv run ruff format <file_path>` and `uv run ruff check --fix <file_path>`. Do this for every Python file you create or modify, before moving on to the next step.
- No `assert` in production code.
- `time.monotonic()` for durations, not `time.time()`.
- In `airflow-core`, functions with a `session` parameter must not call `session.commit()`. Use keyword-only `session` parameters.
- Imports at top of file. Valid exceptions: circular imports, lazy loading for worker isolation, `TYPE_CHECKING` blocks.
- Guard heavy type-only imports (e.g., `kubernetes.client`) with `TYPE_CHECKING` in multi-process code paths.
- Define dedicated exception classes or use existing exceptions such as `ValueError` instead of raising the broad `AirflowException` directly. Each error case should have a specific exception type that conveys what went wrong.
- Apache License header on all new files (prek enforces this).
- Newsfragments are only added if a major change or breaking change is applied. This is usually coordinate during review. Please do not add newsfragments per default as in most cases this needs a reversion during review.

## Testing Standards

- Add tests for new behavior — cover success, failure, and edge cases.
- Use pytest patterns, not `unittest.TestCase`.
- Use `spec`/`autospec` when mocking.
- Use `time_machine` for time-dependent tests.
- Use `@pytest.mark.parametrize` for multiple similar inputs.
- Use `@pytest.mark.db_test` for tests that require database access.
- Test fixtures: `devel-common/src/tests_common/pytest_plugin.py`.
- Test location mirrors source: `airflow/cli/cli_parser.py` → `tests/cli/test_cli_parser.py`.


## Commits and PRs

Write commit messages focused on user impact, not implementation details.

- **Good:** `Fix airflow dags test command failure without serialized Dags`
- **Good:** `UI: Fix Grid view not refreshing after task actions`
- **Bad:** `Initialize DAG bundles in CLI get_dag function`

Add a newsfragment for user-visible changes:
`echo "Brief description" > airflow-core/newsfragments/{PR_NUMBER}.{bugfix|feature|improvement|doc|misc|significant}.rst`

- NEVER add Co-Authored-By with yourself as co-author of the commit. Agents cannot be authors, humans can be, Agents are assistants.

### Creating Pull Requests

**Always push to the user's fork**, not to the upstream `apache/airflow` repo. Never push
directly to `main`.

Before pushing, determine the fork remote. Check `git remote -v` — if `origin` does **not**
point to `apache/airflow`, use `origin` (it's the user's fork). If `origin` points to
`apache/airflow`, look for another remote that points to the user's fork. If no fork remote
exists, create one:

```bash
gh repo fork apache/airflow --remote --remote-name fork
```

Before pushing, perform a self-review of your changes following the Gen-AI review guidelines
in [`contributing-docs/05_pull_requests.rst`](contributing-docs/05_pull_requests.rst) and the
code review checklist in [`.github/instructions/code-review.instructions.md`](.github/instructions/code-review.instructions.md):

1. Review the full diff (`git diff main...HEAD`) and verify every change is intentional and
   related to the task — remove any unrelated changes.
2. Read `.github/instructions/code-review.instructions.md` and check your diff against every
   rule — architecture boundaries, database correctness, code quality, testing requirements,
   API correctness, and AI-generated code signals. Fix any violations before pushing.
3. Confirm the code follows the project's coding standards and architecture boundaries
   described in this file.
4. Run regular (fast) static checks (`prek run --from-ref <target_branch> --stage pre-commit`)
   and fix any failures.
5. Run manual (slower) checks (`prek run --from-ref <target_branch> --stage manual`) and fix any failures.
6. Run relevant individual tests and confirm they pass.
7. Find which tests to run for the changes with selective-checks and run those tests in parallel to confirm they pass and check for CI-specific issues.
8. Check for security issues — no secrets, no injection vulnerabilities, no unsafe patterns.

Before pushing, always rebase your branch onto the latest target branch (usually `main`)
to avoid merge conflicts and ensure CI runs against up-to-date code:

```bash
git fetch <upstream-remote> <target_branch>
git rebase <upstream-remote>/<target_branch>
```

If there are conflicts, resolve them and continue the rebase. If the rebase is too complex,
ask the user for guidance.

Then push the branch to the fork remote and open the PR creation page in the browser
with the body pre-filled (including the generative AI disclosure already checked):

```bash
git push -u <fork-remote> <branch-name>
gh pr create --web --title "Short title (under 70 chars)" --body "$(cat <<'EOF'
Brief description of the changes.

closes: #ISSUE  (if applicable)

---

##### Was generative AI tooling used to co-author this PR?

- [X] Yes — <Agent Name and Version>

Generated-by: <Agent Name and Version> following [the guidelines](https://github.com/apache/airflow/blob/main/contributing-docs/05_pull_requests.rst#gen-ai-assisted-contributions)

EOF
)"
```

The `--web` flag opens the browser so the user can review and submit. The `--body` flag
pre-fills the PR template with the generative AI disclosure already completed.

Remind the user to:

1. Review the PR title — keep it short (under 70 chars) and focused on user impact.
2. Add a brief description of the changes at the top of the body.
3. Reference related issues when applicable (`closes: #ISSUE` or `related: #ISSUE`).

## Boundaries

- **Ask first**
  - Large cross-package refactors.
  - New dependencies with broad impact.
  - Destructive data or migration changes.
- **Never**
  - Commit secrets, credentials, or tokens.
  - Edit generated files by hand when a generation workflow exists.
  - Use destructive git operations unless explicitly requested.

## References

- [`contributing-docs/03a_contributors_quick_start_beginners.rst`](contributing-docs/03a_contributors_quick_start_beginners.rst)
- [`contributing-docs/05_pull_requests.rst`](contributing-docs/05_pull_requests.rst)
- [`contributing-docs/07_local_virtualenv.rst`](contributing-docs/07_local_virtualenv.rst)
- [`contributing-docs/08_static_code_checks.rst`](contributing-docs/08_static_code_checks.rst)
- [`contributing-docs/12_provider_distributions.rst`](contributing-docs/12_provider_distributions.rst)
- [`contributing-docs/19_execution_api_versioning.rst`](contributing-docs/19_execution_api_versioning.rst)


---
## Appended governance file: .github/instructions/code-review.instructions.md
source: https://raw.githubusercontent.com/apache/airflow/main/.github/instructions/code-review.instructions.md
commit: a06453823a13960387d73210253fc729504396ca
captured: 2026-04-05
lines: 71
reason: Explicit: "Read this file"
---

---
applyTo: "**"
excludeAgent: "coding-agent"
---

# Airflow Code Review Instructions

Use these rules when reviewing pull requests to the Apache Airflow repository.

## Architecture Boundaries

- **Scheduler must never run user code.** It only processes serialized Dags. Flag any scheduler-path code that deserializes or executes Dag/task code.
- **Flag any task execution code that accesses the metadata DB directly** instead of through the Execution API (`/execution` endpoints).
- **Flag any code in Dag Processor or Triggerer that breaks process isolation** — these components run user code in isolated processes.
- **Flag any provider importing core internals** like `SUPERVISOR_COMMS` or task-runner plumbing. Providers interact through the public SDK and execution API only.

## Database and Query Correctness

- **Flag any SQLAlchemy relationship access inside a loop** without `joinedload()` or `selectinload()` — this is an N+1 query.
- **Flag any query on `run_id` without `dag_id`.** `run_id` is only unique per Dag. Queries that filter, group, partition, or join on `run_id` alone will silently collide across Dags.
- **Flag any `session.commit()` call in `airflow-core`** code that receives a `session` parameter. Session lifecycle is managed by the caller, not the callee.
- **Flag any `session` parameter that is not keyword-only** (`*, session`) in `airflow-core`.
- **Flag any database-specific SQL** (e.g., `LATERAL` joins, PostgreSQL-only functions, MySQL-only syntax) without cross-DB handling. SQL must work on PostgreSQL, MySQL, and SQLite.

## Code Quality Rules

- **Flag any `assert` in non-test code.** `assert` is stripped in optimized Python (`python -O`), making it a silent no-op in production.
- **Flag any `time.time()` used for measuring durations.** Use `time.monotonic()` instead — `time.time()` is affected by system clock adjustments.
- **Flag any `from` or `import` statement inside a function or method body.** Imports must be at the top of the file. The only valid exceptions are: (1) circular import avoidance, (2) lazy loading for worker isolation, (3) `TYPE_CHECKING` blocks. If the import is inside a function, ask the author to justify why it cannot be at module level.
- **Flag any `@lru_cache(maxsize=None)`.** This creates an unbounded cache — every unique argument set is cached forever. Note: `@lru_cache()` without arguments defaults to `maxsize=128` and is fine.
- **Flag any heavy import** (e.g., `kubernetes.client`) in multi-process code paths that is not behind a `TYPE_CHECKING` guard.
- **Flag any file, connection, or session opened without a context manager or `try/finally`.**

## Testing Requirements

- **Flag any new public method or behavior without corresponding tests.** Tests must cover success, failure, and edge cases.
- **Flag any `unittest.TestCase` subclass.** Use pytest patterns instead.
- **Flag any `mock.Mock()` or `mock.MagicMock()` without `spec` or `autospec`.** Unspec'd mocks silently accept any attribute access, hiding real bugs.
- **Flag any `time.sleep` or `datetime.now()` in tests.** Use `time_machine` for time-dependent tests.
- **Flag any issue number in test docstrings** (e.g., `"""Fix for #12345"""`) — test names should describe behavior, not track tickets.

## API Correctness

- **Flag any query on mapped task instances that does not filter on `map_index`.** Without it, queries may return arbitrary instances from the mapped set.
- **Flag any Execution API change without a Cadwyn version migration** (CalVer format).

## UI Code (React/TypeScript)

- Avoid `useState + useEffect` to sync derived state. Use nullish coalescing or nullable override patterns instead.
- Extract shared logic into custom hooks rather than copy-pasting across components.

## Generated Files

- **Flag any manual edits to files in `openapi-gen/`** or Task SDK generated models. These must be regenerated, not hand-edited.

## AI-Generated Code Signals

Flag these patterns that indicate low-quality AI-generated contributions:

- **Fabricated diffs**: Changes to files or code paths that don't exist in the repository.
- **Unrelated files included**: Changes to files that have nothing to do with the stated purpose of the PR.
- **Description doesn't match code**: PR description describes something different from what the code actually does.
- **No evidence of testing**: Claims of fixes without test evidence, or author admitting they cannot run the test suite.
- **Over-engineered solutions**: Adding caching layers, complex locking, or benchmark scripts for problems that don't exist or are misunderstood.
- **Narrating comments**: Comments that restate what the next line does (e.g., `# Add the item to the list` before `list.append(item)`).
- **Empty PR descriptions**: PRs with just the template filled in and no actual description of the changes.

## Quality Signals to Check

- **For bug-fix PRs, flag if there is no regression test** — a test that fails without the fix and passes with it.
- **Flag any existing test modified to accommodate new behavior** — this may indicate a behavioral regression rather than a genuine fix.
