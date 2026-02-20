# 🚀 launchpad

Set up AI coding instructions for your project through a conversation.

Tell Launchpad what you're building — not what language you want, but what
problem you're solving — and it recommends the best framework for AI-assisted
development, then generates tailored `.github/copilot-instructions.md`,
scoped `.instructions.md` files, and `AGENTS.md` for your repo.

Powered by OpenAI. Opinionated by design.

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

Launchpad carries a curated, opinionated library of instruction templates.
These inform the tone, depth, and structure of what gets generated — but every
output is customized to your conversation.

The system is **use-case driven**, not language-driven. Tell Launchpad what
you're building and it recommends the best framework for AI-assisted
development — not what's popular, but what gives you the best chance of
shipping.

### Recommended stacks (best for AI)

| Stack | Best for | Scaffold |
|-------|----------|----------|
| Elixir + Phoenix | Real-time web, dashboards, chat, collaboration | `mix phx.new` |
| TypeScript + SvelteKit | JS full-stack web, SSR, content sites | `npm create svelte@latest` |
| Ruby on Rails | CRUD apps, MVPs, admin panels | `rails new` |

### All supported stacks

| Stack | Use case | Scaffold |
|-------|----------|----------|
| TypeScript + Next.js | React ecosystem, Vercel deployment | `npx create-next-app@latest` |
| TypeScript + Fastify | Node.js API services | `npm init -y` |
| Go Service | High-perf APIs, CLIs, infrastructure | `go mod init` |
| .NET API | Enterprise APIs, C# ecosystem | `dotnet new webapi` |
| Python + FastAPI | ML backends, data APIs | `python -m venv .venv` |
| Python + Django | Python full-stack, admin-heavy | `django-admin startproject` |
| Dart + Flutter | Mobile, cross-platform native | `flutter create` |
| Rust + Axum | Performance-critical services | `cargo new` |
| Laravel | PHP full-stack, SaaS | `composer create-project` |

### Add-ons

| Add-on | Coverage |
|--------|----------|
| Data-intensive | Postgres, NATS, Parquet, event-driven |
| Frontend craft | CSS architecture, animation, accessibility |

**Key opinions:**
- Real-time → Phoenix, not React + server
- Node.js API → Fastify, not Express
- Mobile → Flutter, not React Native
- Framework CLIs scaffold the project — AI writes app code, not boilerplate

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