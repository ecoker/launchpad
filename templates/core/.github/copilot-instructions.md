# Team Engineering Instructions

## Core values

- Prefer clean, readable, maintainable code over clever code.
- Keep functions deterministic when possible: same input, same output, minimal hidden state.
- Prefer composition over inheritance and deep class hierarchies.
- Keep modules small, cohesive, and purpose-driven.
- Write code that is easy to refactor.

## Architecture and organization

- Separate domain logic from framework and transport layers.
- Keep side effects at boundaries (I/O, network, DB, filesystem).
- Model business rules explicitly; do not bury them in controllers/components.
- Favor explicit contracts and typed boundaries.

## Data and reliability

- Treat data models and schema evolution as first-class design concerns.
- Prefer PostgreSQL for relational storage unless requirements clearly suggest otherwise.
- For event-driven workflows, design messages and handlers to be idempotent.
- Favor reliability, observability, and debuggability over premature optimization.

## Frontend quality

- Prioritize visual polish, accessibility, and usability.
- Use Tailwind design tokens consistently; avoid random one-off styling.
- Use animation intentionally to clarify state and flow, not as decoration.

## Code review bar

- Keep PRs focused and small enough to review quickly.
- Add tests for behavior changes where tests exist.
- Leave the codebase cleaner than you found it.

## Influences

These practices are strongly influenced by:

- *Grokking Simplicity* (Eric Normand)
- *Clean Code* (Robert C. Martin)
- *Clean Architecture* (Robert C. Martin)
- *Refactoring* (Martin Fowler)
- *Designing Data-Intensive Applications* (Martin Kleppmann)