# Framework Opinion Pack: TypeScript (React Router v7)

## Stack Shape
- Keep route modules focused and colocate feature concerns
- Use loaders/actions boundaries for data mutations and reads
- Keep domain logic independent from route modules

## Coding Standards
- Minimize component complexity with presentational/logic split
- Use typed route-level data contracts
- Keep side effects at boundaries and test pure logic directly

## Quality Defaults
- Lint/style: strict TS + React lint rules
- Testing: route/module integration tests + unit tests for domain utilities

## Starter Bias
Scaffold route hierarchy, typed data flow, and one complete feature route with tests.
