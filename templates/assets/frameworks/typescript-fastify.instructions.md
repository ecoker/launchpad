# Framework Opinion Pack: TypeScript (Fastify)

## Stack Shape
- Build around modular plugins and typed route contracts
- Keep request validation and serialization explicit
- Separate transport concerns from core business logic

## Coding Standards
- Prefer small handlers delegating to service functions
- Use schema-first validation with shared types
- Keep side effects behind injectable adapters

## Quality Defaults
- Lint/style: strict TS + fail-on-warning lint
- Testing: service-level unit tests + HTTP integration contract tests

## Starter Bias
Scaffold observability, health checks, typed config, and one production-style route module.
