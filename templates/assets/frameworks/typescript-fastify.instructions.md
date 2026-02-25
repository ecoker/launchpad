# Framework Opinion Pack: TypeScript (Fastify)

## Stack Shape
- Build around modular plugins and typed route contracts
- Keep request validation and serialization explicit
- Separate transport concerns from core business logic

## Coding Standards
- Prefer small handlers delegating to service functions
- Use schema-first validation with shared types
- Keep side effects behind injectable adapters

## TypeScript Config
- `strict: true` in `tsconfig.json` with `noUncheckedIndexedAccess`
- TypeBox schemas as single source of truth for validation + types
- `unknown` over `any`; narrow explicitly

## Quality Defaults
- Lint/style: strict TS + fail-on-warning lint
- Testing: Vitest for service unit tests, `inject` + `light-my-request` for HTTP integration

## Starter Bias
Scaffold observability, health checks, typed config, and one production-style route module.
