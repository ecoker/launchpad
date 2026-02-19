# Framework Opinion Pack: Rust

## Stack Shape
- Keep core domain logic in pure modules with explicit trait boundaries
- Isolate IO, networking, and persistence behind adapters
- Prefer explicit error types over opaque panics

## Coding Standards
- Favor composition and small traits
- Keep lifetimes/ownership readable over cleverness
- Use clear module boundaries and minimal public surface

## Quality Defaults
- Lint/style: clippy as gate, rustfmt enforced
- Testing: unit tests at module level + integration tests at crate boundary

## Starter Bias
Scaffold for correctness and observability first: typed errors, logging hooks, and deterministic tests.
