# Framework Opinion Pack: C# (.NET)

## Stack Shape
- Prefer ASP.NET Core minimal APIs for focused services
- Use strong DTO boundaries and explicit validation
- Keep domain logic framework-agnostic where practical

## Coding Standards
- Favor small, composable services over inheritance-heavy hierarchies
- Keep async flows explicit and cancellation-aware
- Use dependency injection with clear lifetimes

## Quality Defaults
- Lint/style: enforce analyzers and formatting in CI
- Testing: fast unit tests + integration tests around endpoints and persistence

## Starter Bias
Scaffold for maintainability first: clear projects, predictable naming, and strict compile/lint gates.
