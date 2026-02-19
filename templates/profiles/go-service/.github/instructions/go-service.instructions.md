---
name: Go Service Standards
description: Composable Go service design with explicit contracts and observability
applyTo: "**/*.go"
---

# Go conventions

- Keep packages small and domain-oriented.
- Prefer simple structs/functions and explicit interfaces at boundaries.
- Return errors with context; avoid swallowing failures.
- Keep goroutine usage controlled with clear ownership and cancellation.
- Write straightforward, readable code over abstract indirection.
- Keep handlers, transport, and core business logic separated.