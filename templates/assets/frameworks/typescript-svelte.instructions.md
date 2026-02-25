# Framework Opinion Pack: TypeScript (Svelte)

## Stack Shape
- Keep stores focused and domain-oriented
- Isolate side effects from component rendering paths
- Prefer simple component composition and explicit props/contracts

## Coding Standards
- Keep reactive statements minimal and readable
- Avoid hidden coupling across global stores
- Use typed boundaries for API/data interactions

## TypeScript Config
- `strict: true` in `tsconfig.json` — SvelteKit default, never relax
- Zod or Valibot for runtime validation in server load/actions
- `type` over `interface`; `as const` objects over `enum`

## Quality Defaults
- Lint/style: strict TS checks + Svelte linting
- Testing: Vitest for unit tests, Playwright for E2E, Testing Library for component behavior

## Starter Bias
Scaffold for velocity with maintainability: theme tokens, route structure, and one real feature slice.
