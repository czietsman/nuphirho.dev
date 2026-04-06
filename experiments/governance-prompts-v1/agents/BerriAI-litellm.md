---
source: https://raw.githubusercontent.com/BerriAI/litellm/main/AGENTS.md
commit: cfd0e2cf9924c7cf0c64cd37298710f7b125a8d1
captured: 2026-04-05
lines: 274
---
# INSTRUCTIONS FOR LITELLM

This document provides comprehensive instructions for AI agents working in the LiteLLM repository.

## OVERVIEW

LiteLLM is a unified interface for 100+ LLMs that:
- Translates inputs to provider-specific completion, embedding, and image generation endpoints
- Provides consistent OpenAI-format output across all providers
- Includes retry/fallback logic across multiple deployments (Router)
- Offers a proxy server (LLM Gateway) with budgets, rate limits, and authentication
- Supports advanced features like function calling, streaming, caching, and observability

## REPOSITORY STRUCTURE

### Core Components
- `litellm/` - Main library code
  - `llms/` - Provider-specific implementations (OpenAI, Anthropic, Azure, etc.)
  - `proxy/` - Proxy server implementation (LLM Gateway)
  - `router_utils/` - Load balancing and fallback logic
  - `types/` - Type definitions and schemas
  - `integrations/` - Third-party integrations (observability, caching, etc.)

### Key Directories
- `tests/` - Comprehensive test suites
- `docs/my-website/` - Documentation website
- `ui/litellm-dashboard/` - Admin dashboard UI
- `enterprise/` - Enterprise-specific features

## DEVELOPMENT GUIDELINES

### MAKING CODE CHANGES

1. **Provider Implementations**: When adding/modifying LLM providers:
   - Follow existing patterns in `litellm/llms/{provider}/`
   - Implement proper transformation classes that inherit from `BaseConfig`
   - Support both sync and async operations
   - Handle streaming responses appropriately
   - Include proper error handling with provider-specific exceptions

2. **Type Safety**: 
   - Use proper type hints throughout
   - Update type definitions in `litellm/types/`
   - Ensure compatibility with both Pydantic v1 and v2

3. **Testing**:
   - Add tests in appropriate `tests/` subdirectories
   - Include both unit tests and integration tests
   - Test provider-specific functionality thoroughly
   - Consider adding load tests for performance-critical changes

### MAKING CODE CHANGES FOR THE UI (IGNORE FOR BACKEND)

1. **Tremor is DEPRECATED, do not use Tremor components in new features/changes**
   - The only exception is the Tremor Table component and its required Tremor Table sub components.

2. **Use Common Components as much as possible**:
   - These are usually defined in the `common_components` directory
   - Use these components as much as possible and avoid building new components unless needed

3. **Testing**:
   - The codebase uses **Vitest** and **React Testing Library**
   - **Query Priority Order**: Use query methods in this order: `getByRole`, `getByLabelText`, `getByPlaceholderText`, `getByText`, `getByTestId`
   - **Always use `screen`** instead of destructuring from `render()` (e.g., use `screen.getByText()` not `getByText`)
   - **Wrap user interactions in `act()`**: Always wrap `fireEvent` calls with `act()` to ensure React state updates are properly handled
   - **Use `query` methods for absence checks**: Use `queryBy*` methods (not `getBy*`) when expecting an element to NOT be present
   - **Test names must start with "should"**: All test names should follow the pattern `it("should ...")`
   - **Mock external dependencies**: Check `setupTests.ts` for global mocks and mock child components/networking calls as needed
   - **Structure tests properly**:
     - First test should verify the component renders successfully
     - Subsequent tests should focus on functionality and user interactions
     - Use `waitFor` for async operations that aren't already awaited
   - **Avoid using `querySelector`**: Prefer React Testing Library queries over direct DOM manipulation

### IMPORTANT PATTERNS

1. **Function/Tool Calling**:
   - LiteLLM standardizes tool calling across providers
   - OpenAI format is the standard, with transformations for other providers
   - See `litellm/llms/anthropic/chat/transformation.py` for complex tool handling

2. **Streaming**:
   - All providers should support streaming where possible
   - Use consistent chunk formatting across providers
   - Handle both sync and async streaming

3. **Error Handling**:
   - Use provider-specific exception classes
   - Maintain consistent error formats across providers
   - Include proper retry logic and fallback mechanisms

4. **Configuration**:
   - Support both environment variables and programmatic configuration
   - Use `BaseConfig` classes for provider configurations
   - Allow dynamic parameter passing

## PROXY SERVER (LLM GATEWAY)

The proxy server is a critical component that provides:
- Authentication and authorization
- Rate limiting and budget management
- Load balancing across multiple models/deployments
- Observability and logging
- Admin dashboard UI
- Enterprise features

Key files:
- `litellm/proxy/proxy_server.py` - Main server implementation
- `litellm/proxy/auth/` - Authentication logic
- `litellm/proxy/management_endpoints/` - Admin API endpoints

**Database (proxy)**: Use Prisma model methods (`prisma_client.db.<model>.upsert`, `.find_many`, `.find_unique`, etc.), not raw SQL (`execute_raw`/`query_raw`). See COMMON PITFALLS for details.

## MCP (MODEL CONTEXT PROTOCOL) SUPPORT

LiteLLM supports MCP for agent workflows:
- MCP server integration for tool calling
- Transformation between OpenAI and MCP tool formats
- Support for external MCP servers (Zapier, Jira, Linear, etc.)
- See `litellm/experimental_mcp_client/` and `litellm/proxy/_experimental/mcp_server/`

## RUNNING SCRIPTS

Use `poetry run python script.py` to run Python scripts in the project environment (for non-test files).

## GITHUB TEMPLATES

When opening issues or pull requests, follow these templates:

### Bug Reports (`.github/ISSUE_TEMPLATE/bug_report.yml`)
- Describe what happened vs. expected behavior
- Include relevant log output
- Specify LiteLLM version
- Indicate if you're part of an ML Ops team (helps with prioritization)

### Feature Requests (`.github/ISSUE_TEMPLATE/feature_request.yml`)
- Clearly describe the feature
- Explain motivation and use case with concrete examples

### Pull Requests (`.github/pull_request_template.md`)
- Add at least 1 test in `tests/litellm/`
- Ensure `make test-unit` passes


## TESTING CONSIDERATIONS

1. **Provider Tests**: Test against real provider APIs when possible
2. **Proxy Tests**: Include authentication, rate limiting, and routing tests
3. **Performance Tests**: Load testing for high-throughput scenarios
4. **Integration Tests**: End-to-end workflows including tool calling

## DOCUMENTATION

- Keep documentation in sync with code changes
- Update provider documentation when adding new providers
- Include code examples for new features
- Update changelog and release notes

## SECURITY CONSIDERATIONS

- Handle API keys securely
- Validate all inputs, especially for proxy endpoints
- Consider rate limiting and abuse prevention
- Follow security best practices for authentication

## ENTERPRISE FEATURES

- Some features are enterprise-only
- Check `enterprise/` directory for enterprise-specific code
- Maintain compatibility between open-source and enterprise versions

## COMMON PITFALLS TO AVOID

1. **Breaking Changes**: LiteLLM has many users - avoid breaking existing APIs
2. **Provider Specifics**: Each provider has unique quirks - handle them properly
3. **Rate Limits**: Respect provider rate limits in tests
4. **Memory Usage**: Be mindful of memory usage in streaming scenarios
5. **Dependencies**: Keep dependencies minimal and well-justified
6. **UI/Backend Contract Mismatch**: When adding a new entity type to the UI, always check whether the backend endpoint accepts a single value or an array. Match the UI control accordingly (single-select vs. multi-select) to avoid silently dropping user selections
7. **Missing Tests for New Entity Types**: When adding a new entity type (e.g., in `EntityUsage`, `UsageViewSelect`), always add corresponding tests in the existing test files and update any icon/component mocks
8. **Raw SQL in proxy DB code**: Do not use `execute_raw` or `query_raw` for proxy database access. Use Prisma model methods (e.g. `prisma_client.db.litellm_tooltable.upsert()`, `.find_many()`, `.find_unique()`) so behavior stays consistent with the schema, the client stays mockable in tests, and you avoid the pitfalls of hand-written SQL (parameter ordering, type casting, schema drift)

8. **Do not hardcode model-specific flags**: Put model-specific capability flags in `model_prices_and_context_window.json` and read them via `get_model_info` (or existing helpers like `supports_reasoning`). This prevents users from needing to upgrade LiteLLM each time a new model supports a feature.

   **Example of BAD** (hardcoded model checks):

   ```python
   @staticmethod
   def _is_effort_supported_model(model: str) -> bool:
       """Check if the model supports the output_config.effort parameter..."""
       model_lower = model.lower()
       if AnthropicConfig._is_claude_4_6_model(model):
           return True
       return any(
           v in model_lower for v in ("opus-4-5", "opus_4_5", "opus-4.5", "opus_4.5")
       )
   ```

   **Example of GOOD** (config-driven or helper that reads from config):

   ```python
   if (
       "claude-3-7-sonnet" in model
       or AnthropicConfig._is_claude_4_6_model(model)
       or supports_reasoning(
           model=model,
           custom_llm_provider=self.custom_llm_provider,
       )
   ):
       ...
   ```

   Using helpers like `supports_reasoning` (which read from `model_prices_and_context_window.json` / `get_model_info`) allows future model updates to "just work" without code changes.

9. **Never close HTTP/SDK clients on cache eviction**: Do not add `close()`, `aclose()`, or `create_task(close_fn())` inside `LLMClientCache._remove_key()` or any cache eviction path. Evicted clients may still be held by in-flight requests; closing them causes `RuntimeError: Cannot send a request, as the client has been closed.` in production after the cache TTL (1 hour) expires. Connection cleanup is handled at shutdown by `close_litellm_async_clients()`. See PR #22247 for the full incident history.

## HELPFUL RESOURCES

- Main documentation: https://docs.litellm.ai/
- Provider-specific docs in `docs/my-website/docs/providers/`
- Admin UI for testing proxy features

## WHEN IN DOUBT

- Follow existing patterns in the codebase
- Check similar provider implementations
- Ensure comprehensive test coverage
- Update documentation appropriately
- Consider backward compatibility impact

## Cursor Cloud specific instructions

### Environment

- Poetry is installed in `~/.local/bin`; the update script ensures it is on `PATH`.
- Python 3.12, Node 22 are pre-installed.
- The virtual environment lives under `~/.cache/pypoetry/virtualenvs/`.

### Running the proxy server

Start the proxy with a config file:

```bash
poetry run litellm --config dev_config.yaml --port 4000
```

The proxy takes ~15-20 seconds to fully start (it runs Prisma migrations on boot). Wait for `/health` to return before sending requests. Without a PostgreSQL `DATABASE_URL`, the proxy connects to a default Neon dev database embedded in the `litellm-proxy-extras` package.

### Running tests

See `CLAUDE.md` and the `Makefile` for standard commands. Key notes:

- `psycopg-binary` must be installed (`poetry run pip install psycopg-binary`) because the pytest-postgresql plugin requires it and the lock file only includes `psycopg` (no binary).
- `openapi-core` must be installed (`poetry run pip install openapi-core`) for the OpenAPI compliance tests in `tests/test_litellm/interactions/`.
- The `--timeout` pytest flag is NOT available; don't pass it.
- Unit tests: `poetry run pytest tests/test_litellm/ -x -vv -n 4`
- Black `--check` may report pre-existing formatting issues; this does not block test runs.
- If `poetry install` fails with "pyproject.toml changed significantly since poetry.lock was last generated", run `poetry lock` first to regenerate the lock file.

### Lint

```bash
cd litellm && poetry run ruff check .
```

Ruff is the primary fast linter. For the full lint suite (including mypy, black, circular imports), run `make lint` per `CLAUDE.md`.

### UI Dashboard development

- The UI is at `ui/litellm-dashboard/`. Run `npm run dev` from that directory for the Next.js dev server on port 3000.
- The proxy at port 4000 serves a **pre-built** static UI from `litellm/proxy/_experimental/out/`. After making UI code changes, you must run `npm run build` in the dashboard directory and copy the output: `cp -r ui/litellm-dashboard/out/* litellm/proxy/_experimental/out/` for the proxy to serve the updated UI.
- SVGs used as provider logos (loaded via `<img>` tags) must NOT use `fill="currentColor"` — replace with an explicit color like `#000000` or use the `-color` variant from lobehub icons, since CSS color inheritance does not work inside `<img>` elements.
- Provider logos live in `ui/litellm-dashboard/public/assets/logos/` (source) and `litellm/proxy/_experimental/out/assets/logos/` (pre-built). Both locations must have the file for it to work in dev and proxy-served modes.
- UI Vitest tests: `cd ui/litellm-dashboard && npx vitest run`

---
## Appended governance file: CLAUDE.md
source: https://raw.githubusercontent.com/BerriAI/litellm/main/CLAUDE.md
commit: 8e61b32b8e76b3318cd74532fa4dcf1af51803c4
captured: 2026-04-05
lines: 164
reason: Explicit: "See CLAUDE.md for standard commands"
---

# CLAUDE.md

This file provides guidance to Claude Code (claude.ai/code) when working with code in this repository.

## Development Commands

### Installation
- `make install-dev` - Install core development dependencies
- `make install-proxy-dev` - Install proxy development dependencies with full feature set
- `make install-test-deps` - Install all test dependencies

### Testing
- `make test` - Run all tests
- `make test-unit` - Run unit tests (tests/test_litellm) with 4 parallel workers
- `make test-integration` - Run integration tests (excludes unit tests)
- `pytest tests/` - Direct pytest execution

### Code Quality
- `make lint` - Run all linting (Ruff, MyPy, Black, circular imports, import safety)
- `make format` - Apply Black code formatting
- `make lint-ruff` - Run Ruff linting only
- `make lint-mypy` - Run MyPy type checking only

### Single Test Files
- `poetry run pytest tests/path/to/test_file.py -v` - Run specific test file
- `poetry run pytest tests/path/to/test_file.py::test_function -v` - Run specific test

### Running Scripts
- `poetry run python script.py` - Run Python scripts (use for non-test files)

### GitHub Issue & PR Templates
When contributing to the project, use the appropriate templates:

**Bug Reports** (`.github/ISSUE_TEMPLATE/bug_report.yml`):
- Describe what happened vs. what you expected
- Include relevant log output
- Specify your LiteLLM version

**Feature Requests** (`.github/ISSUE_TEMPLATE/feature_request.yml`):
- Describe the feature clearly
- Explain the motivation and use case

**Pull Requests** (`.github/pull_request_template.md`):
- Add at least 1 test in `tests/litellm/`
- Ensure `make test-unit` passes

## Architecture Overview

LiteLLM is a unified interface for 100+ LLM providers with two main components:

### Core Library (`litellm/`)
- **Main entry point**: `litellm/main.py` - Contains core completion() function
- **Provider implementations**: `litellm/llms/` - Each provider has its own subdirectory
- **Router system**: `litellm/router.py` + `litellm/router_utils/` - Load balancing and fallback logic
- **Type definitions**: `litellm/types/` - Pydantic models and type hints
- **Integrations**: `litellm/integrations/` - Third-party observability, caching, logging
- **Caching**: `litellm/caching/` - Multiple cache backends (Redis, in-memory, S3, etc.)

### Proxy Server (`litellm/proxy/`)
- **Main server**: `proxy_server.py` - FastAPI application
- **Authentication**: `auth/` - API key management, JWT, OAuth2
- **Database**: `db/` - Prisma ORM with PostgreSQL/SQLite support
- **Management endpoints**: `management_endpoints/` - Admin APIs for keys, teams, models
- **Pass-through endpoints**: `pass_through_endpoints/` - Provider-specific API forwarding
- **Guardrails**: `guardrails/` - Safety and content filtering hooks
- **UI Dashboard**: Served from `_experimental/out/` (Next.js build)

## Key Patterns

### Provider Implementation
- Providers inherit from base classes in `litellm/llms/base.py`
- Each provider has transformation functions for input/output formatting
- Support both sync and async operations
- Handle streaming responses and function calling

### Error Handling
- Provider-specific exceptions mapped to OpenAI-compatible errors
- Fallback logic handled by Router system
- Comprehensive logging through `litellm/_logging.py`

### Configuration
- YAML config files for proxy server (see `proxy/example_config_yaml/`)
- Environment variables for API keys and settings
- Database schema managed via Prisma (`proxy/schema.prisma`)

## Development Notes

### Code Style
- Uses Black formatter, Ruff linter, MyPy type checker
- Pydantic v2 for data validation
- Async/await patterns throughout
- Type hints required for all public APIs
- **Avoid imports within methods** — place all imports at the top of the file (module-level). Inline imports inside functions/methods make dependencies harder to trace and hurt readability. The only exception is avoiding circular imports where absolutely necessary.
- **Use dict spread for immutable copies** — prefer `{**original, "key": new_value}` over `dict(obj)` + mutation. The spread produces the final dict in one step and makes intent clear.
- **Guard at resolution time** — when resolving an optional value through a fallback chain (`a or b or ""`), raise immediately if the resolved result being empty is an error. Don't pass empty strings or sentinel values downstream for the callee to deal with.
- **Extract complex comprehensions to named helpers** — a set/dict comprehension that calls into the DB or manager (e.g. "which of these server IDs are OAuth2?") belongs in a named helper function, not inline in the caller.
- **FastAPI parameter declarations** — mark required query/form params with `= Query(...)` / `= Form(...)` explicitly when other params in the same handler are optional. Mixing `str` (required) with `Optional[str] = None` in the same signature causes silent 422s when the required param is missing.

### Testing Strategy
- Unit tests in `tests/test_litellm/`
- Integration tests for each provider in `tests/llm_translation/`
- Proxy tests in `tests/proxy_unit_tests/`
- Load tests in `tests/load_tests/`
- **Always add tests when adding new entity types or features** — if the existing test file covers other entity types, add corresponding tests for the new one
- **Keep monkeypatch stubs in sync with real signatures** — when a function gains a new optional parameter, update every `fake_*` / `stub_*` in tests that patch it to also accept that kwarg (even as `**kwargs`). Stale stubs fail with `unexpected keyword argument` and mask real bugs.
- **Test all branches of name→ID resolution** — when adding server/resource lookup that resolves names to UUIDs, test: (1) name resolves and UUID is allowed, (2) name resolves but UUID is not allowed, (3) name does not resolve at all. The silent-fallback path is where access-control bugs hide.

### UI / Backend Consistency
- When wiring a new UI entity type to an existing backend endpoint, verify the backend API contract (single value vs. array, required vs. optional params) and ensure the UI controls match — e.g., use a single-select dropdown when the backend accepts a single value, not a multi-select

### MCP OAuth / OpenAPI Transport Mapping
- `TRANSPORT.OPENAPI` is a UI-only concept. The backend only accepts `"http"`, `"sse"`, or `"stdio"`. Always map it to `"http"` before any API call (including pre-OAuth temp-session calls).
- FastAPI validation errors return `detail` as an array of `{loc, msg, type}` objects. Error extractors must handle: array (map `.msg`), string, nested `{error: string}`, and fallback.
- When an MCP server already has `authorization_url` stored, skip OAuth discovery (`_discovery_metadata`) — the server URL for OpenAPI MCPs is the spec file, not the API base, and fetching it causes timeouts.
- `client_id` should be optional in the `/authorize` endpoint — if the server has a stored `client_id` in credentials, use that. Never require callers to re-supply it.

### MCP Credential Storage
- OAuth credentials and BYOK credentials share the `litellm_mcpusercredentials` table, distinguished by a `"type"` field in the JSON payload (`"oauth2"` vs plain string).
- When deleting OAuth credentials, check type before deleting to avoid accidentally deleting a BYOK credential for the same `(user_id, server_id)` pair.
- Always pass the raw `expires_at` timestamp to the client — never set it to `None` for expired credentials. Let the frontend compute the "Expired" display state from the timestamp.
- Use `RecordNotFoundError` (not bare `except Exception`) when catching "already deleted" in credential delete endpoints.

### Browser Storage Safety (UI)
- Never write LiteLLM access tokens or API keys to `localStorage` — use `sessionStorage` only. `localStorage` survives browser close and is readable by any injected script (XSS).
- Shared utility functions (e.g. `extractErrorMessage`) belong in `src/utils/` — never define them inline in hooks or duplicate them across files.

### Database Migrations
- Prisma handles schema migrations
- Migration files auto-generated with `prisma migrate dev`
- Always test migrations against both PostgreSQL and SQLite

### Proxy database access
- **Do not write raw SQL** for proxy DB operations. Use Prisma model methods instead of `execute_raw` / `query_raw`.
- Use the generated client: `prisma_client.db.<model>` (e.g. `litellm_tooltable`, `litellm_usertable`) with `.upsert()`, `.find_many()`, `.find_unique()`, `.update()`, `.update_many()` as appropriate. This avoids schema/client drift, keeps code testable with simple mocks, and matches patterns used in spend logs and other proxy code.
- **No N+1 queries.** Never query the DB inside a loop. Batch-fetch with `{"in": ids}` and distribute in-memory.
- **Batch writes.** Use `create_many`/`update_many`/`delete_many` instead of individual calls (these return counts only; `update_many`/`delete_many` no-op silently on missing rows). When multiple separate writes target the same table (e.g. in `batch_()`), order by primary key to avoid deadlocks.
- **Push work to the DB.** Filter, sort, group, and aggregate in SQL, not Python. Verify Prisma generates the expected SQL — e.g. prefer `group_by` over `find_many(distinct=...)` which does client-side processing.
- **Bound large result sets.** Prisma materializes full results in memory. For results over ~10 MB, paginate with `take`/`skip` or `cursor`/`take`, always with an explicit `order`. Prefer cursor-based pagination (`skip` is O(n)). Don't paginate naturally small result sets.
- **Limit fetched columns on wide tables.** Use `select` to fetch only needed fields — returns a partial object, so downstream code must not access unselected fields.
- **Check index coverage.** For new or modified queries, check `schema.prisma` for a supporting index. Prefer extending an existing index (e.g. `@@index([a])` → `@@index([a, b])`) over adding a new one, unless it's a `@@unique`. Only add indexes for large/frequent queries.
- **Keep schema files in sync.** Apply schema changes to all `schema.prisma` copies (`schema.prisma`, `litellm/proxy/`, `litellm-proxy-extras/`, `litellm-js/spend-logs/` for SpendLogs) with a migration under `litellm-proxy-extras/litellm_proxy_extras/migrations/`.

### Setup Wizard (`litellm/setup_wizard.py`)
- The wizard is implemented as a single `SetupWizard` class with `@staticmethod` methods — keep it that way. No module-level functions except `run_setup_wizard()` (the public entrypoint) and pure helpers (color, ANSI).
- Use `litellm.utils.check_valid_key(model, api_key)` for credential validation — never roll a custom completion call.
- Do not hardcode provider env-key names or model lists that already exist in the codebase. Add a `test_model` field to each provider entry to drive `check_valid_key`; set it to `None` for providers that can't be validated with a single API key (Azure, Bedrock, Ollama).

### Enterprise Features
- Enterprise-specific code in `enterprise/` directory
- Optional features enabled via environment variables
- Separate licensing and authentication for enterprise features

### HTTP Client Cache Safety
- **Never close HTTP/SDK clients on cache eviction.** `LLMClientCache._remove_key()` must not call `close()`/`aclose()` on evicted clients — they may still be used by in-flight requests. Doing so causes `RuntimeError: Cannot send a request, as the client has been closed.` after the 1-hour TTL expires. Cleanup happens at shutdown via `close_litellm_async_clients()`.

### Troubleshooting: DB schema out of sync after proxy restart
`litellm-proxy-extras` runs `prisma migrate deploy` on startup using **its own** bundled migration files, which may lag behind schema changes in the current worktree. Symptoms: `Unknown column`, `Invalid prisma invocation`, or missing data on new fields.

**Diagnose:** Run `\d "TableName"` in psql and compare against `schema.prisma` — missing columns confirm the issue.

**Fix options:**
1. **Create a Prisma migration** (permanent) — run `prisma migrate dev --name <description>` in the worktree. The generated file will be picked up by `prisma migrate deploy` on next startup.
2. **Apply manually for local dev** — `psql -d litellm -c "ALTER TABLE ... ADD COLUMN IF NOT EXISTS ..."` after each proxy start. Fine for dev, not for production.
3. **Update litellm-proxy-extras** — if the package is installed from PyPI, its migration directory must include the new file. Either update the package or run the migration manually until the next release ships it.
