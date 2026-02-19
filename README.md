# 🚀 launchpad

Set up AI coding instructions for your project through a conversation.

Tell Launchpad what you're building — the language, framework, your style
preferences — and it generates tailored `.github/copilot-instructions.md`,
scoped `.instructions.md` files, and `AGENTS.md` for your repo.

Powered by OpenAI. Your copilot should write code the way *you* would.

## How it works

1. Run `launchpad init ./my-project`
2. Have a brief conversation about what you're building
3. Launchpad generates customized AI instruction files
4. Your AI copilot is now briefed on your standards

No rigid templates. No picking from a list. Just describe your project and get
instructions tailored to your stack, team, and coding philosophy.

## What gets generated

| File | Purpose |
|------|---------|
| `.github/copilot-instructions.md` | Always-on project standards for every chat and suggestion |
| `.github/instructions/*.instructions.md` | Scoped rules by language, framework, or concern |
| `AGENTS.md` | Ground rules for multi-agent collaboration |

## Install

```bash
# Homebrew (macOS / Linux)
brew install ehrencoker/tap/launchpad

# Or curl
curl -sSfL https://raw.githubusercontent.com/ehrencoker/agent-kit/main/install.sh | sh

# Or build from source
go install github.com/ehrencoker/agent-kit/cmd/launchpad@latest
```

## Usage

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="sk-..."

# Start a conversation to generate instructions
launchpad init ./my-app

# Force overwrite in existing directory
launchpad init ./existing-project --force

# See the template knowledge base
launchpad list
```

## Knowledge base

Launchpad carries a curated library of example instructions as seed knowledge.
These inform the tone, depth, and structure of what gets generated — but every
output is customized to your conversation.

| Area | Coverage |
|------|----------|
| TypeScript + React | React Router v7, Tailwind, shadcn/ui, Motion |
| Python | Pydantic, functional style, data pipelines |
| Elixir + Phoenix | LiveView, Ecto, OTP, pattern matching |
| .NET | C# minimal APIs, Clean Architecture, EF Core |
| Laravel | Inertia, Eloquent, actions, queues |
| Go | Stdlib-first, small interfaces, explicit errors |
| Data-intensive | Postgres, NATS, Parquet, event-driven |
| Frontend craft | CSS architecture, animation, accessibility |

## Philosophy

Shaped by:

- **Grokking Simplicity** — actions, calculations, data
- **Clean Code & Clean Architecture** — naming, structure, discipline
- **Refactoring** — continuous, safe, mechanical improvement
- **Designing Data-Intensive Applications** — reliability at scale

## Extending

Add example templates under `templates/profiles/<id>/` or `templates/addons/<id>/`
to expand the knowledge base. Launchpad reads all templates and uses them as
context when generating custom instructions.

## Requirements

- An OpenAI API key (`OPENAI_API_KEY` env var or entered interactively)
- That's it