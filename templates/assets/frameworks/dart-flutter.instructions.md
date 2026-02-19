# Framework Opinion Pack: Dart (Flutter)

## Stack Shape
- Organize by feature modules rather than widget type
- Separate presentation state from side-effect boundaries
- Keep platform integration adapters isolated

## Coding Standards
- Prefer immutable data and explicit state transitions
- Keep widgets small and focused
- Treat design tokens (color, spacing, typography) as first-class config

## Quality Defaults
- Lint/style: strict analyzer rules and formatted code
- Testing: widget tests for behavior + integration smoke path

## Starter Bias
Bootstrap a clear app shell, routing, theme tokens, and one vertical feature slice.
