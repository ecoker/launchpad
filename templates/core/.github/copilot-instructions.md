# {{PROJECT_NAME}} — Engineering Standards

> "Simplicity is prerequisite for reliability." — Edsger Dijkstra

These are the ground rules. They apply to every file, every language, every PR.

## Philosophy

We write software that is **clean**, **typed**, **functional-first**, and **beautiful**.
We care about the humans who read code as much as the machines that run it.
We build data-intensive systems that are reliable, observable, and maintainable.

## Core principles

### Composition over inheritance
Build behavior by composing small, focused functions — not by extending class trees.
If you're reaching for `extends`, pause and ask whether a plain function or a
module-level composition would be clearer.

### Deterministic by default
Functions should be predictable: same inputs → same outputs. Push randomness,
clocks, network calls, and file I/O to the edges. The core should be pure.

```
// ✅ Pure transformation — easy to test, easy to trust
const applyDiscount = (price: number, rate: number): number => price * (1 - rate);

// ❌ Hidden dependency — harder to test, surprising in production
const applyDiscount = (price: number): number => price * (1 - getConfigRate());
```

### Explicit over implicit
Type your boundaries. Name things for intent. Make data flow visible.
When someone reads your code six months from now, they shouldn't need archaeology.

### Small, cohesive modules
Each file should do one thing well. If a module needs a paragraph to explain
what it does, it's doing too much. Split it.

## Architecture

- **Separate domain from infrastructure.** Business rules live in pure modules.
  HTTP handlers, database queries, and queue consumers are plumbing — keep them thin.
- **Push side effects to boundaries.** I/O happens at the edges; transformations
  happen in the middle. This is the functional core, imperative shell pattern.
- **Model business rules explicitly.** Use types and discriminated unions to make
  illegal states unrepresentable instead of relying on runtime validation alone.
- **Design for change.** Code will be refactored. Write it so refactoring is safe
  and mechanical, not scary and manual.

## Data

- **PostgreSQL is the default** for relational/transactional workloads.
- **Schema changes are first-class.** Migrations are versioned, reviewed, and
  reversible. Never hand-edit production schemas.
- **Idempotency matters.** Design message handlers and data pipelines so they
  can be safely retried without side effects.
- **Observability is not optional.** Structured logs, correlation IDs, health
  checks, and metrics from day one — not bolted on after an incident.

## Quality bar

- Leave every file cleaner than you found it.
- Keep PRs focused: one concern per PR.
- Write tests for behavior, not implementation details.
- Prefer integration tests that exercise real boundaries over mocks of everything.

## Influences

These opinions didn't come from nowhere:

- *Grokking Simplicity* (Eric Normand) — actions, calculations, and data
- *Clean Code* & *Clean Architecture* (Robert C. Martin) — structure and discipline
- *Refactoring* (Martin Fowler) — continuous, safe improvement
- *Designing Data-Intensive Applications* (Martin Kleppmann) — building reliable systems at scale