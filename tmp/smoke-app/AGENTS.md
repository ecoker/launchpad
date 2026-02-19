# Agent Collaboration — smoke-app

Welcome, agent. Here are our ground rules for working together.

1. **Read `.github/copilot-instructions.md` first.** This file defines the non-negotiable engineering standards for the entire project, including architectural guidelines, style, layers, and boundaries.
2. **Check `.github/instructions/` for focused rules.** Instructions in this directory are scoped by language, framework, or project concern. The `applyTo` patterns in each file tell you where and when to apply them.
3. **Make minimal, focused changes.** Address the root cause and keep each change tightly scoped. Avoid wide-ranging refactors unless they're intentionally planned and reviewed.
4. **Keep tests aligned.** Your changes must maintain or improve test coverage. Always run the tests before marking work as complete.
5. **Respect existing patterns.** If the codebase follows a style or convention, continue with it—even if you personally would do things differently.

**What great output looks like:**
- Explicit types at every boundary.
- Pure functions wherever possible, with side effects isolated at the edges.
- Names that communicate intent, not implementation.
- Small, cohesive, independently understandable modules.
- Code that's cleaner after your change than it was before.

**When in doubt:**
- Prefer simple over clever.
- Prefer explicit over implicit.
- Prefer composition over inheritance.
- Ask: "Would someone enjoy reading this code six months from now?"
