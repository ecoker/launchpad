# Framework Opinion Pack: TypeScript (Next.js)

## Stack Shape
- Prefer App Router and server-first data paths where practical
- Keep domain logic outside route handlers and components
- Use typed contracts at API and data boundaries

## Coding Standards
- Minimize client components to interaction-heavy surfaces
- Avoid business logic in UI components
- Keep data fetching and caching behavior explicit

## TypeScript Config
- `strict: true` in `tsconfig.json` — non-negotiable
- Enable `noUncheckedIndexedAccess` for safer array/object access
- Zod for runtime validation; infer types with `z.infer<>` to avoid duplication
- `type` over `interface`; `as const` objects over `enum`

## Quality Defaults
- Lint/style: strict TypeScript + ESLint gate
- Testing: Vitest for unit tests, Playwright for E2E, Testing Library for components

## Starter Bias
Scaffold with clear app shell, route groups, typed env config, and one end-to-end feature slice.
