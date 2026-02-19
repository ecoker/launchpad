---
name: Architecture and Refactoring
description: Baseline architecture, organization, and refactoring standards
applyTo: "**"
---

# Architecture and refactoring rules

- Keep business logic framework-agnostic where practical.
- Push side effects to boundaries and keep core logic pure.
- Prefer explicit data flow over hidden mutable state.
- Refactor in small, behavior-preserving steps.
- Remove duplication by extracting composable functions.
- Name modules and functions for intent, not implementation details.
- Avoid deep inheritance trees; prefer composition.