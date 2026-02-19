# Agent Collaboration — test-app

Welcome, agent! These are the ground rules for working smoothly together.

1. **Read `.github/copilot-instructions.md` first.** It contains always-on project standards and applies to every file, language, and change.
2. **Check `.github/instructions/` for concern-specific rules.** Files here are scoped to languages, frameworks, or focus areas. Respect their `applyTo` patterns.
3. **Make minimal, focused changes.** Address the root cause. Avoid sprawling PRs.
4. **Keep or improve test coverage.** All changes should maintain coverage. Run tests before finishing.
5. **Respect existing patterns.** Follow the codebase’s established style and conventions, even if they're different from your preferences.

## Great output

- Explicit types at boundaries.
- Pure functions where possible, with side effects at the edges.
- Intentional, descriptive naming.
- Small, focused, independent modules.
- The codebase is cleaner after your change than before.

## When uncertain

- Prefer simplicity to cleverness.
- Prefer explicitness over implicitness.
- Prefer composition over inheritance.
- Ask: "Will someone enjoy reading this code six months from now?"
