# Agent Collaboration — smoke-app

Welcome, agent. Here's how we work together.

## Ground rules

1. **Read `.github/copilot-instructions.md` first.** Those are the project-wide
   standards. They apply to every file, every language, every change.
2. **Check `.github/instructions/` for focused rules.** These files are scoped
   to specific languages, frameworks, or concerns. The `applyTo` patterns tell
   you which files they govern.
3. **Make minimal, focused changes.** Solve the root cause. Don't refactor the
   world in a single diff.
4. **Keep tests aligned.** If the project has tests, your changes should maintain
   or improve coverage. Run them before declaring done.
5. **Respect existing patterns.** If the codebase uses a particular style or
   convention, follow it — even if you'd prefer something else.

## What great output looks like

- Types are explicit at every boundary.
- Functions are pure when possible, with side effects at the edges.
- Names describe intent, not implementation.
- Modules are small, cohesive, and independently understandable.
- The code is cleaner after your change than it was before.

## When in doubt

- Prefer simple over clever.
- Prefer explicit over implicit.
- Prefer composition over inheritance.
- Ask yourself: "Would someone enjoy reading this code six months from now?"
