---
name: Elixir Phoenix Standards
description: Functional Elixir architecture with explicit boundaries
applyTo: "**/*.{ex,exs,heex}"
---

# Elixir/Phoenix conventions

- Keep domain logic in pure modules and context boundaries.
- Prefer pattern matching and pipelines for readability.
- Use supervision trees intentionally; design for failure recovery.
- Keep Ecto schemas and changesets explicit and validated.
- Avoid leaking web-layer concerns into domain modules.
- Favor small, composable modules over large god-contexts.