# Framework Opinion Pack: TypeScript (React Router v7)

## Stack Shape
- Keep route modules focused and colocate feature concerns
- Use loaders/actions boundaries for data mutations and reads
- Keep domain logic independent from route modules

## Coding Standards
- Minimize component complexity with presentational/logic split
- Use typed route-level data contracts
- Keep side effects at boundaries and test pure logic directly

## TypeScript Config
- `strict: true` in `tsconfig.json` with `noUncheckedIndexedAccess`
- Zod for loader/action validation; infer types with `z.infer<>`
- `type` over `interface`; `as const` objects over `enum`

## Quality Defaults
- Lint/style: strict TS + React lint rules
- Testing: Vitest for unit + route integration tests, Playwright for E2E

## Starter Bias
Scaffold route hierarchy, typed data flow, and one complete feature route with tests.
