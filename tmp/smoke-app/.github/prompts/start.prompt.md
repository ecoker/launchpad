---
description: Bootstrap the smoke-app project using Elixir + Phoenix with strong functional boundaries, layering, and collaborative standards.
agent: "agent"
tools: []
---

You are initializing the **smoke-app** project. This is an internal voting web app for up to 50 simultaneous users, with no authentication, built with Elixir and Phoenix.

**Instructions:**
- Apply all published standards from `.github/copilot-instructions.md` and framework/language instructions in `.github/instructions/`.
- Structure application code with clear context boundaries and functional layering:
  - UI (LiveView or controllers) is thin.
  - Business logic lives in contexts.
  - Ecto models are internal to their contexts.
  - Side effects and I/O are at the boundaries.
- Favor pipelines, pattern matching, composition, and small, pure functions.
- All new code must improve clarity and standards conformance.
- If you see opportunities to extract or refactor for smaller, more explicit modules, do so with minimal, reviewable steps.
- For all changes: leave the codebase better than you found it.
