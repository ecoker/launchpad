# agent-kit

Opinionated scaffolding for AI agent instructions.

This toolkit creates starter folders for new projects with:

- `.github/copilot-instructions.md` for always-on project standards
- `.github/instructions/*.instructions.md` for language/framework-specific rules
- `AGENTS.md` compatibility for multi-agent workflows

## Why this exists

You want consistent project setup shaped by strong opinions:

- Beauty and UX craft
- Organization and maintainability
- Clean Code and refactoring discipline
- Functional composition and deterministic behavior
- Typed code by default (especially TypeScript)

## Profiles

- `typescript-react`
- `python-data`
- `elixir-phoenix`
- `dotnet-api`
- `laravel`
- `go-service`

## Usage

```bash
npm install
npm run build

node dist/cli.js list
node dist/cli.js init ./my-app --profile typescript-react
```

Or in development mode:

```bash
npm run dev -- init ./my-app --profile typescript-react
```

## Generated structure (example)

```text
my-app/
  .github/
    copilot-instructions.md
    instructions/
      architecture.instructions.md
      typescript-react.instructions.md
  AGENTS.md
  README.md
```

## Extending

Add new profiles under `templates/profiles/<profile-id>/` with one or more `*.instructions.md` files.

The CLI merges:

1. `templates/core`
2. `templates/profiles/<profile-id>`

Use `{{PROJECT_NAME}}` placeholder in template text files when needed.