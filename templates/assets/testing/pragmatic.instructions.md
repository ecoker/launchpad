---
name: Testing Conventions
description: How we write, organize, and run tests — framework-aware, behavior-focused
applyTo: "**"
---

# Testing conventions

> "A test that doesn't tell you what broke is just a slow build step."

Tests exist to give fast, trustworthy feedback about behavior. Not to hit a
coverage number. Not to prove implementation details haven't changed. To catch
real breakage before users do.

## Philosophy

- **Test behavior, not implementation.** Assert on outputs and effects, not on
  internal method calls or data shapes. If the implementation changes but the
  behavior doesn't, tests should still pass.
- **Prefer integration over mocks.** Tests that exercise real boundaries
  (database, HTTP, file system) catch more real bugs. Mock only at boundaries
  you don't own (third-party APIs, payment providers).
- **Fast by default.** Unit tests should run in milliseconds. Integration tests
  should run in low seconds. If your test suite takes minutes, something is
  wrong — split it or parallelize.
- **Deterministic always.** No flaky tests. No test ordering dependencies. No
  reliance on wall-clock time. Use dependency injection to control clocks,
  randomness, and external state.

## Test pyramid

```
          ┌───────────────┐
          │   E2E / Smoke  │  Few. Critical user paths only.
          ├───────────────┤
          │  Integration   │  Moderate. Real DB, real HTTP, real boundaries.
          ├───────────────┤
          │     Unit       │  Many. Pure logic, fast, isolated.
          └───────────────┘
```

- **Unit tests** — pure domain logic, transformations, calculations. No I/O.
  No framework. These are your fastest feedback loop.
- **Integration tests** — exercise a full slice: HTTP request → handler →
  database → response. Use real infrastructure (test databases, in-memory
  queues) over mocks wherever practical.
- **E2E / smoke tests** — one or two critical user flows to verify the system
  assembles correctly. Not your primary safety net.

## File conventions

Co-locate tests with the code they test when the framework supports it.

| Stack | Test file pattern | Runner |
|-------|-------------------|--------|
| TypeScript (Vitest) | `*.test.ts` / `*.test.tsx` next to source | `vitest` |
| TypeScript (Jest) | `*.test.ts` / `*.spec.ts` next to source | `jest` |
| Elixir | `test/**/*_test.exs` mirroring `lib/` | `mix test` |
| Python (pytest) | `tests/test_*.py` or `*_test.py` co-located | `pytest` |
| Go | `*_test.go` same package, same directory | `go test ./...` |
| Rust | `#[cfg(test)] mod tests` inline + `tests/` for integration | `cargo test` |
| Ruby (RSpec) | `spec/**/*_spec.rb` mirroring `app/` | `bundle exec rspec` |
| Ruby (Minitest) | `test/**/*_test.rb` mirroring `app/` | `rails test` |
| C# (xUnit) | `*.Tests` project alongside source project | `dotnet test` |
| Java (JUnit 5) | `src/test/java/` mirroring `src/main/java/` | `./gradlew test` |
| Dart (Flutter) | `test/` mirroring `lib/`, `*_test.dart` | `flutter test` |
| PHP (Pest/PHPUnit) | `tests/Feature/`, `tests/Unit/` | `php artisan test` |

## What to test

- **Domain logic** — every pure function with meaningful branching gets unit
  tests. If a function has three code paths, write three tests minimum.
- **API boundaries** — every public endpoint gets at least one happy-path and
  one error-path integration test.
- **Validation rules** — if you validate input (and you should), test that
  invalid input is rejected with the expected error shape.
- **State transitions** — if your domain models have lifecycle states (draft →
  published → archived), test the transitions and the rejections.
- **Edge cases** — empty collections, null/nil inputs at boundaries,
  concurrent access where relevant.

## What NOT to test

- Private implementation details that will change during refactoring.
- Framework internals (don't test that your ORM generates SQL correctly).
- Exact snapshot matches of rendered HTML/UI unless you have a specific
  regression reason.
- Logging output (test the behavior that triggers the log, not the log itself).

## Test structure

Every test follows the same shape:

```
Arrange → Act → Assert
(Given  → When → Then)
```

Name tests for the behavior they verify:
```
✅ "returns shipping cost of zero for orders over $100"
✅ "rejects login with expired token"
❌ "test1"
❌ "it works"
```

## Coverage expectations

- **Don't set a global coverage percentage target.** Coverage is a tool for
  finding untested code, not a goal to game.
- **Do expect full coverage of domain logic.** Pure functions in the core layer
  should be thoroughly tested.
- **Do expect at least one happy + one error test per endpoint.**
- **Do flag untested public functions** in code review — they should be tested
  or explicitly marked as trivial.

## Framework-specific guidance

### TypeScript (Vitest / Jest)
- Vitest is preferred for new projects (faster, ESM-native, Vite-aligned).
- Use `describe` blocks to group by function/behavior, not by file.
- Use `beforeEach` for shared setup, `afterEach` for cleanup.
- For component tests: Testing Library (`@testing-library/react`,
  `@testing-library/svelte`) — query by role and text, not by CSS selector.
- For E2E: Playwright over Cypress. Fewer flaky tests, better parallelism.

### Elixir
- ExUnit is the only testing framework. Use it.
- `async: true` on every `describe` block unless tests share state.
- Use `Ecto.Adapters.SQL.Sandbox` for concurrent database tests.
- Test contexts through their public API; don't reach into internal modules.
- Use ExMachina or plain factory functions for test data.

### Python
- pytest over unittest. Always.
- Use fixtures for setup/teardown. Keep them in `conftest.py`.
- Use `pytest-asyncio` for async code.
- Use `httpx.AsyncClient` (FastAPI) or Django's `TestClient` for integration.
- Type-check test files too — they're code, not a special zone.

### Go
- Table-driven tests with `t.Run` subtests.
- `testify/assert` or `testify/require` for readable assertions.
- `httptest.NewServer` for HTTP integration.
- `t.Parallel()` for independent tests.
- `testcontainers-go` for database integration when needed.

### Rust
- `#[test]` in same-file `mod tests` for unit tests.
- `tests/` directory for integration tests.
- `#[tokio::test]` for async tests.
- `assert_eq!`, `assert_matches!` for assertions.
- `mockall` crate only at external boundaries.

### C# / .NET
- xUnit over NUnit or MSTest.
- `WebApplicationFactory<T>` for integration tests.
- FluentAssertions for readable assertions.
- `Testcontainers` for database integration.
- Separate `*.Tests` and `*.IntegrationTests` projects.

### Java / Spring
- JUnit 5 with `@Nested` classes for grouping.
- `@WebMvcTest` for controller slices, `@DataJpaTest` for repository slices.
- `@SpringBootTest` only for full integration tests.
- AssertJ over Hamcrest.
- Testcontainers for database and messaging.

### Ruby / Rails
- RSpec or Minitest — pick one for the project and commit.
- Factory Bot (RSpec) or fixtures (Minitest) for test data.
- System tests with Capybara for critical user flows.
- `rails test` or `bundle exec rspec` — both are fine.

### Dart / Flutter
- `flutter test` for widget and unit tests.
- `testWidgets` for widget behavior tests — pump, tap, verify.
- `mockito` or `mocktail` for dependency mocking.
- `integration_test/` for device/emulator integration tests.

### PHP / Laravel
- Pest is preferred over raw PHPUnit for readability.
- `RefreshDatabase` trait for test isolation.
- Feature tests hit full HTTP cycle; Unit tests hit classes directly.
- Factories and seeders for test data.
- Assert against response status, JSON structure, and database state.
