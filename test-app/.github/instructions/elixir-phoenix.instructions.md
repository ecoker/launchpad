---
name: Elixir + Phoenix
description: Functional architecture with explicit boundaries, fault tolerance, and joy
applyTo: "**/*.{ex,exs,heex,leex}"
---

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

## Context boundaries (the big one)

Phoenix contexts are your domain boundary. Get them right.

- **Each context owns a concept.** `Accounts`, `Billing`, `Inventory` — not
  `Helpers` or `Utils`.
- **Contexts are the public API.** Other parts of the app call context functions.
  They never reach into a context's internal modules directly.
- **Keep contexts focused.** If a context file exceeds 300 lines, it's probably
  two contexts. Split it.
- **Cross-context communication** goes through explicit function calls or PubSub.
  Never share Ecto schemas across contexts.

```
lib/
  my_app/
    accounts/            # Context boundary
      accounts.ex        # Public API
      user.ex            # Ecto schema (internal)
      user_token.ex
    billing/
      billing.ex
      invoice.ex
    inventory/
      inventory.ex
      product.ex
```

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
