# Agent Collaboration — {{PROJECT_NAME}}

Welcome, agent. This file covers **workflow rules specific to agent collaboration**.
For engineering standards, see `.github/copilot-instructions.md`. For scoped
framework/concern rules, see `.github/instructions/*.instructions.md`.

## Instruction hierarchy

1. `.github/copilot-instructions.md` — always-on project standards.
2. `.github/instructions/*.instructions.md` — scoped rules matched by `applyTo`
   glob patterns. These govern specific languages, frameworks, or concerns.
3. This file — agent workflow rules.

When instructions conflict, the more specific scoped file wins over the general
copilot instructions. If two scoped files conflict on the same file, follow the
one whose `applyTo` pattern is more specific.

## Change discipline

- **One concern per change.** Solve the stated problem. Don't bundle drive-by
  refactors, formatting changes, or dependency upgrades into the same diff.
- **Explain the "why" in commit messages.** The code shows what changed; the
  message should explain why.
- **Run tests before declaring done.** If the project has tests, run them. If
  your change breaks a test, fix it. If the change adds behavior, add a test.
  Check the testing conventions file in `.github/instructions/` for
  framework-specific test runners and patterns.
- **Don't generate boilerplate files.** If the framework has a CLI scaffold
  command (`rails new`, `mix phx.new`, `npm create svelte@latest`, etc.),
  use it. Don't hand-write `package.json`, `mix.exs`, `Cargo.toml`, or
  framework config files from scratch.

## What to do when context is missing

- If a task requires domain knowledge that isn't documented, **look at the
  existing code first** — naming, routing, and model structure reveal intent.
- If a testing framework isn't specified, check `package.json`, `mix.exs`,
  `pyproject.toml`, `go.mod`, or the equivalent dependency file for what's
  already installed.
- If styling approach is unclear, check for `tailwind.config.*` or a
  `ThemeData` definition before introducing a new system.

## Multi-agent coordination

When multiple agents work on the same codebase:

- **Don't touch files another agent owns** unless your task explicitly
  requires it. Check recent diffs for active ownership.
- **Shared types and schemas are high-contention zones.** Change them
  carefully and expect merge conflicts.
- **Prefer additive changes over modifications.** Adding a new function is
  less likely to conflict than rewriting an existing one.
