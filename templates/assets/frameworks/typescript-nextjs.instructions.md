# Framework Opinion Pack: TypeScript (Next.js)

## Stack Shape
- Prefer App Router and server-first data paths where practical
- Keep domain logic outside route handlers and components
- Use typed contracts at API and data boundaries

## Coding Standards
- Minimize client components to interaction-heavy surfaces
- Avoid business logic in UI components
- Keep data fetching and caching behavior explicit

## Quality Defaults
- Lint/style: strict TypeScript + ESLint gate
- Testing: unit tests for domain logic + route/component smoke coverage

## Starter Bias
Scaffold with clear app shell, route groups, typed env config, and one end-to-end feature slice.
