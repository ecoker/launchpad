# Testing: Pragmatic Confidence

## Guidance
- Keep the test pyramid balanced: many unit tests, focused integration tests, minimal e2e smoke tests
- Co-locate tests near features where possible
- Prefer deterministic tests over snapshot-heavy or brittle tests
- Include contract tests around external boundaries (APIs, queues, DB adapters)

## Application Rule
Starter prompt should scaffold test commands, one smoke test path, and one core unit-test example for the selected stack.
