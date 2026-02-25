# {{PROJECT_NAME}}

## AI-assisted development

This project includes opinionated AI instruction files that guide code generation
tools (GitHub Copilot, Claude, and other AI agents) toward consistent, high-quality output.

### How it works

| File | Purpose |
|------|---------|
| `.github/copilot-instructions.md` | Always-on project standards — applies to every chat and suggestion |
| `.github/instructions/*.instructions.md` | Scoped rules for specific languages, frameworks, or concerns |
| `AGENTS.md` | Ground rules for AI agent collaboration |

### Customizing

These files are **yours**. They're a starting point, not a constraint.

- Edit `.github/copilot-instructions.md` to reflect your team's actual values.
- Add new `.instructions.md` files for frameworks or patterns specific to your project.
- Delete anything that doesn't fit. Every rule should earn its place.

### Philosophy

These instructions are shaped by a few core beliefs:

- **Typed code is reliable code.** Explicit types at boundaries prevent entire
  categories of bugs.
- **Composition over inheritance.** Small, focused functions compose better than
  deep class hierarchies.
- **Beauty matters.** Software that is pleasant to read is pleasant to maintain.
- **Data is king.** Schema evolution, observability, and reliability are
  first-class concerns — not afterthoughts.

---

*Scaffolded with [launchpad](https://github.com/ecoker/launchpad).*