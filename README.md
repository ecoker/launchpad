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
brew install ecoker/tap/launchpad

# Or curl
curl -sSfL https://raw.githubusercontent.com/ecoker/launchpad/main/install.sh | sh

# Or build from source
go install github.com/ecoker/launchpad/cmd/launchpad@latest
```

## Usage

```bash
# Set your OpenAI API key
export OPENAI_API_KEY="sk-..."

# Optionally override the model (default: gpt-4.1)
export LAUNCHPAD_MODEL="gpt-4.1-mini"

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
you're building and it recommends a framework based on conceptual integrity,
explicit contracts, and minimal runtime magic — not popularity. Disagree? Pick
a different option in the conversation; every supported stack gets the same
generation quality.

### Canonical stacks (coherence-first philosophy)

These stacks were chosen for conceptual clarity, long-term stability, and
clean standards across every architectural layer. Each one has a principled
reason for being here.

| Stack | Layer | Best for | Scaffold |
|-------|-------|----------|----------|
| Elixir + Phoenix | Coordination | Real-time web, distributed systems, live data | `mix phx.new` |
| TypeScript + SvelteKit | Web UI | JS full-stack web, SSR, content sites | `npm create svelte@latest` |
| Ruby on Rails | Rapid Product | CRUD apps, MVPs, admin panels, SaaS | `rails new` |
| Go Service | Worker | High-perf APIs, CLIs, infrastructure | `go mod init` |
| Rust + Axum | Worker | Performance-critical services, systems | `cargo new` |
| .NET API | Enterprise | Enterprise APIs, C# ecosystem | `dotnet new webapi` |
| Java + Spring Boot | Enterprise | Large-scale enterprise, JVM ecosystem | `spring init` |
| Python + FastAPI | AI Boundary | ML backends, data APIs, LLM integration | `python -m venv .venv` |
| Dart + Flutter | Mobile UI | Cross-platform native apps | `flutter create` |

### Additional supported stacks

| Stack | Layer | Use case | Scaffold |
|-------|-------|----------|----------|
| TypeScript + Next.js | Web UI | React ecosystem, Vercel deployment | `npx create-next-app@latest` |
| TypeScript + Fastify | Worker | Node.js API services | `npm init -y` |
| Python + Django | Rapid Product | Python full-stack, admin-heavy | `django-admin startproject` |
| Laravel | Rapid Product | PHP full-stack, SaaS | `composer create-project` |

### Layer taxonomy

Every stack maps to an architectural role:

| Layer | Role | Canonical stacks |
|-------|------|-----------------|
| Coordination | Distributed orchestration, real-time, supervision | Phoenix |
| Worker | High-performance stateless services | Go, Rust |
| Enterprise | Structured integration, regulated environments | .NET, Spring Boot |
| AI Boundary | LLM integration, schema-driven data APIs | FastAPI |
| Web UI | Browser-based product surfaces | SvelteKit |
| Mobile UI | Cross-platform native experiences | Flutter |
| Rapid Product | Convention-maximalist fast iteration | Rails |

### Add-ons

| Add-on | Coverage |
|--------|----------|
| Data-intensive | Postgres, NATS, Parquet, event-driven |
| Frontend craft | Visual discipline, component composition, accessibility, motion |

**Visual assets are automatic.** Any stack with a UI surface automatically
gets frontend-craft guidance, a default color palette (Obsidian + Indigo),
and font pairing (Inter + JetBrains Mono). No opt-in needed.

**Key opinions:**
- Coherence over popularity — SvelteKit over Next.js, Fastify over Express
- Real-time → Phoenix, not React + server
- Mobile → Flutter, not React Native
- Framework CLIs scaffold the project — AI writes app code, not boilerplate

## Philosophy

Shaped by:

- **Coherence as north star** — conceptual integrity, explicit contracts, minimal magic
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
- Optionally `LAUNCHPAD_MODEL` to use a different OpenAI model
- That's it