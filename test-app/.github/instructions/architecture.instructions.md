---
name: Architecture and Refactoring
description: How we structure code and how we improve it safely over time
applyTo: "**"
---

# Architecture and refactoring

## The shape of good code

Good architecture is boring in the best way. You should be able to open any
module and immediately understand what it does, what it depends on, and where
its boundaries are.

### Layering

Organize code into layers with clear dependency direction:

```
  ┌─────────────────────────────┐
  │   Transport / UI / CLI      │  ← Thin. Translates external input.
  ├─────────────────────────────┤
  │   Application / Use Cases   │  ← Orchestrates domain logic + side effects.
  ├─────────────────────────────┤
  │   Domain / Core             │  ← Pure. No framework imports. Testable.
  ├─────────────────────────────┤
  │   Infrastructure / Adapters │  ← DB, HTTP clients, queues, filesystems.
  └─────────────────────────────┘
```

Domain code never imports from transport or infrastructure. Dependencies point inward.

### Naming

- Name functions for **what they do**, not **how** they do it.
- Name modules for the **concept** they own, not the **pattern** they implement.
- `calculateShippingCost` > `shippingHelper`. `OrderPricing` > `OrderManager`.

## Refactoring discipline

Refactoring is not a separate task. It's part of every change.

- **Small steps.** Each commit should leave the code working. If a refactoring
  is bigger than one PR, break it into a chain of safe, reviewable steps.
- **Extract, don't duplicate.** When you see similar logic in two places, extract
  a shared function. Compose it — don't copy it.
- **Rename fearlessly.** If a name no longer matches its purpose, change it now.
  Bad names compound into bad understanding.
- **Delete dead code.** If it's not called, it's not documentation — it's noise.
  Version control remembers. You don't need to.

## Composition patterns

```elixir
# ✅ Compose small, focused functions
def process_order(order) do
  order
  |> validate_items()
  |> calculate_totals()
  |> apply_discounts()
  |> format_receipt()
end
```

- Prefer **function composition** and **pipelines** over deep call hierarchies.
- Prefer **plain data objects** (e.g. maps, structs) over class instances for domain models.
- Prefer **explicit parameters** over ambient context or service locators.
