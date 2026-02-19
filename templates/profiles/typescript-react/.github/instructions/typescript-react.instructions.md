---
name: TypeScript React Standards
description: Strongly typed React code with composition and clean architecture
applyTo: "**/*.{ts,tsx,mts,cts,js,jsx}"
---

# TypeScript + React conventions

- Type everything at boundaries: API payloads, component props, and public functions.
- Prefer `type` and discriminated unions for domain modeling.
- Avoid `any`; use `unknown` plus narrowing when necessary.
- Keep React components focused; extract logic into pure utility hooks/functions.
- Favor React Router v7 patterns for route modules and data loading.
- Use Tailwind classes consistently with design tokens and spacing rhythm.
- Prefer `shadcn/ui` style component composition for UI primitives.
- Use Motion intentionally for transitions that improve clarity.
- Keep client state minimal and derived whenever possible.

# Preferred patterns

- Composition over class-based abstractions.
- Small pure functions for transformations.
- Feature folders with clear separation of UI, domain, and data access.