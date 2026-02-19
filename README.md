# agent-kit

Opinionated scaffolding for AI agent instructions.

Spin up a new project with deeply considered coding standards baked in from the start.
Your AI copilot should write code the way *you* would — with type safety,
functional composition, clean architecture, and visual craft.

## What you get

When you run `agent-kit init`, your new project gets:

| File | What it does |
|------|-------------|
| `.github/copilot-instructions.md` | Always-on project standards for every chat and suggestion |
| `.github/instructions/*.instructions.md` | Scoped rules by language, framework, or concern |
| `AGENTS.md` | Ground rules for multi-agent collaboration |
| `README.md` | Project documentation explaining the setup |

## Philosophy

These instructions are shaped by years of reading, building, and refactoring:

- **Grokking Simplicity** — separate actions, calculations, and data
- **Clean Code & Clean Architecture** — structure, naming, and discipline
- **Refactoring** — continuous, safe, mechanical improvement
- **Designing Data-Intensive Applications** — reliable systems at scale

And some strong opinions:

- **TypeScript** is the default. Type everything.
- **PostgreSQL** is the database. Schema changes are first-class.
- **Composition over inheritance.** Small functions. Pure cores.
- **Tailwind + shadcn/ui + Motion** = beautiful, accessible interfaces.
- **Elixir** is wonderful. **Go** is great when you need it. **C#/.NET** is better than people think.
- **NATS** for messaging. **Parquet** for archival data.

## Quick start

```bash
npm install
npm run build

# Interactive mode — asks you questions
node dist/cli.js init ./my-app

# Or go fast with flags
node dist/cli.js init ./my-app -p typescript-react -y
```

## CLI reference

```bash
# List all available profiles and add-ons
agent-kit list

# Interactive scaffold (default)
agent-kit init ./my-app

# Non-interactive with specific profile
agent-kit init ./my-app --profile python-data --yes

# With add-ons
agent-kit init ./api --profile go-service --addon data-intensive

# Force overwrite in existing directory
agent-kit init ./existing-project --profile typescript-react --force
```

## Profiles

| Profile | Stack |
|---------|-------|
| `typescript-react` | TypeScript, React, React Router v7, discriminated unions, feature folders |
| `python-data` | Python, Pydantic, typed pipelines, composable transforms |
| `elixir-phoenix` | Elixir, Phoenix, contexts, pattern matching, supervision |
| `dotnet-api` | C#, .NET, Clean Architecture, records, EF Core |
| `laravel` | PHP, Laravel, thin controllers, actions, form requests |
| `go-service` | Go, small interfaces, explicit errors, domain packages |

## Add-ons

| Add-on | What it covers |
|--------|---------------|
| `data-intensive` | Postgres best practices, NATS messaging, Parquet, idempotency, observability |
| `frontend-craft` | Tailwind discipline, shadcn/ui composition, Motion with purpose, accessibility |

## Extending

Add new profiles under `templates/profiles/<profile-id>/` or new add-ons under
`templates/addons/<addon-id>/`.

The CLI merges templates in order:
1. `templates/core` (always)
2. `templates/profiles/<profile>` (chosen profile)
3. `templates/addons/<addon>` (each selected add-on)

Use `{{PROJECT_NAME}}` as a placeholder in any template text file — it's replaced
with the target directory name at scaffold time.