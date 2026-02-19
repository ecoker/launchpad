# test-app — Engineering Standards

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

```elixir
# ✅ Pure transformation — easy to test, easy to trust
def calculate_discount(price, rate), do: price * (1 - rate)

# ❌ Hidden dependency — harder to test, surprising in production
def calculate_discount(price), do: price * (1 - get_config_rate())
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

# Elixir + Phoenix

Elixir is beautiful. It's functional, concurrent, fault-tolerant, and *fun*.
These instructions help you write Elixir that leans into the language's
strengths rather than fighting them.

## Functional core

- **Pattern match everything.** Multi-clause functions are clearer than
  conditionals. Let the runtime find the right clause.
- **Pipe with purpose.** The `|>` operator makes data flow visible. Use it for
  clear transformation chains, not for one-line calls that don't benefit.
- **Keep functions small and pure.** A function that takes data and returns
  data is trivially testable and composable.

```elixir
# ✅ Clear pipeline — data flows visibly
def process_order(params) do
  params
  |> validate_order()
  |> calculate_totals()
  |> apply_discounts()
  |> persist_order()
  |> notify_customer()
end

# ✅ Multi-clause instead of if/case
def format_status(:active), do: "Active"
def format_status(:suspended), do: "Suspended"
def format_status(:cancelled), do: "Cancelled"
```

## Context boundaries

- **Each context owns a concept.** `Accounts`, `Billing`, `Inventory` — not
  `Helpers` or `Utils`.
- **Contexts are the public API.** Other parts of the app call context functions.
  They never reach into a context's internal modules directly.
- **Keep contexts focused.** If a context file exceeds 300 lines, it's probably
  two contexts. Split it.
- **Cross-context communication** goes through explicit function calls or PubSub.
  Never share Ecto schemas across contexts.

## Ecto discipline

- **Changesets validate at the boundary.** Every user-facing write goes through
  a changeset with explicit validations.
- **Keep queries explicit.** Use `Ecto.Query` composably, but avoid hiding
  complex query logic deep in schemas.
- **Use multi for transactional operations.** `Ecto.Multi` makes multi-step
  operations explicit, composable, and safe.

## OTP and supervision

- **Let it crash.** Design processes to restart cleanly under supervision rather
  than trying to handle every possible failure inline.
- **Supervision trees are architecture.** Think carefully about what supervises
  what and what restart strategy is appropriate.
- **GenServers should be simple.** If a GenServer is doing heavy business logic,
  extract the logic into a pure module and keep the GenServer as a thin wrapper.

## LiveView

- Keep socket assigns minimal. Derive what you can in the template.
- Use `assign_async` and `start_async` for data loading.
- Decompose into function components and live components when complexity grows.
- Handle events explicitly — avoid catch-all `handle_event` clauses.

## What to avoid

- Umbrella apps for small/medium projects — start with a flat context structure.
- Raw SQL instead of Ecto queries (unless genuinely necessary for performance).
- Mutable state outside of GenServers and ETS.
- Ignoring dialyzer/credo warnings — fix them or document why they're acceptable.
