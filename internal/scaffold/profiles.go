package scaffold

// Profile represents a language/framework profile that can be scaffolded.
type Profile struct {
	ID          string
	Title       string
	Summary     string
	Dir         string // directory name inside templates/profiles/
	ScaffoldCmd string // CLI command the framework provides to bootstrap a project
	UseCase     string // what kind of projects this is best for
	Layer       string // architectural role: coordination, worker, enterprise, ai-boundary, web-ui, mobile-ui, rapid-product
	HasUI       bool   // whether this profile includes a user interface surface
	Tier        int    // 1 = canonical coherence set, 2 = additional supported stacks
}

// Addon represents an optional add-on instruction set.
type Addon struct {
	ID      string
	Title   string
	Summary string
	Dir     string // directory name inside templates/addons/
}

// Profiles lists every available profile, ordered by recommendation strength.
//
// Tier 1 — Canonical coherence set. These stacks were chosen for conceptual
// integrity, explicit contracts, minimal runtime magic, and long-term
// stability. They span the full spectrum of architectural layers.
//
// Tier 2 — Additional supported stacks for specific domains or ecosystem needs.
var Profiles = []Profile{
	// ── Tier 1: Canonical coherence set ──────────────────────────────

	{
		ID:          "elixir-phoenix",
		Title:       "Elixir + Phoenix",
		Summary:     "Full-stack real-time web — LiveView, Ecto, OTP, no frontend/backend split",
		Dir:         "elixir-phoenix",
		ScaffoldCmd: "mix phx.new {{name}}",
		UseCase:     "Real-time web apps, collaborative tools, dashboards, chat, IoT — anything with live data",
		Layer:       "coordination",
		HasUI:       true,
		Tier:        1,
	},
	{
		ID:          "typescript-sveltekit",
		Title:       "TypeScript + SvelteKit",
		Summary:     "Full-stack JS web — intuitive reactivity, SSR, minimal boilerplate",
		Dir:         "typescript-sveltekit",
		ScaffoldCmd: "npm create svelte@latest",
		UseCase:     "JS-ecosystem full-stack web apps, content sites, SSR apps needing rich interactivity",
		Layer:       "web-ui",
		HasUI:       true,
		Tier:        1,
	},
	{
		ID:          "ruby-rails",
		Title:       "Ruby on Rails",
		Summary:     "Rapid full-stack web — convention over configuration, incredible generators",
		Dir:         "ruby-rails",
		ScaffoldCmd: "rails new {{name}}",
		UseCase:     "CRUD apps, MVPs, admin panels, content platforms, SaaS — fast to production",
		Layer:       "rapid-product",
		HasUI:       true,
		Tier:        1,
	},
	{
		ID:          "go-service",
		Title:       "Go Service",
		Summary:     "Idiomatic Go — stdlib-first, small binaries, excellent concurrency",
		Dir:         "go-service",
		ScaffoldCmd: "go mod init {{module}}",
		UseCase:     "High-performance APIs, CLI tools, infrastructure services, platform tooling",
		Layer:       "worker",
		HasUI:       false,
		Tier:        1,
	},
	{
		ID:          "rust-axum",
		Title:       "Rust + Axum",
		Summary:     "Performance-critical services — type-safe, zero-cost abstractions, Tokio-based",
		Dir:         "rust-axum",
		ScaffoldCmd: "cargo new {{name}}",
		UseCase:     "Performance-critical APIs, systems programming, infrastructure where correctness matters",
		Layer:       "worker",
		HasUI:       false,
		Tier:        1,
	},
	{
		ID:          "dotnet-api",
		Title:       ".NET API",
		Summary:     "C# minimal APIs — Entity Framework, clean architecture, enterprise-grade",
		Dir:         "dotnet-api",
		ScaffoldCmd: "dotnet new webapi -n {{name}}",
		UseCase:     "Enterprise APIs, C# ecosystem services, Azure-native workloads",
		Layer:       "enterprise",
		HasUI:       false,
		Tier:        1,
	},
	{
		ID:          "java-spring",
		Title:       "Java + Spring Boot",
		Summary:     "Enterprise Java — DI, auto-configuration, massive ecosystem, battle-tested at scale",
		Dir:         "java-spring",
		ScaffoldCmd: "spring init --dependencies=web,data-jpa,validation {{name}}",
		UseCase:     "Large-scale enterprise systems, integration-heavy services, JVM ecosystem workloads",
		Layer:       "enterprise",
		HasUI:       false,
		Tier:        1,
	},
	{
		ID:          "python-fastapi",
		Title:       "Python + FastAPI",
		Summary:     "Python APIs — async, typed, Pydantic-centric, ML/data-native",
		Dir:         "python-fastapi",
		ScaffoldCmd: "mkdir {{name}} && cd {{name}} && python -m venv .venv",
		UseCase:     "Python API services, ML model serving, data pipelines, AI agent backends",
		Layer:       "ai-boundary",
		HasUI:       false,
		Tier:        1,
	},
	{
		ID:          "dart-flutter",
		Title:       "Dart + Flutter",
		Summary:     "Cross-platform native apps — single codebase for iOS, Android, web, desktop",
		Dir:         "dart-flutter",
		ScaffoldCmd: "flutter create {{name}}",
		UseCase:     "Mobile apps, cross-platform native experiences — Flutter over React Native",
		Layer:       "mobile-ui",
		HasUI:       true,
		Tier:        1,
	},

	// ── Tier 2: Additional supported stacks ──────────────────────────

	{
		ID:          "typescript-nextjs",
		Title:       "TypeScript + Next.js",
		Summary:     "React ecosystem full-stack — App Router, RSC, Vercel-optimized",
		Dir:         "typescript-nextjs",
		ScaffoldCmd: "npx create-next-app@latest",
		UseCase:     "Apps requiring React ecosystem libraries, Vercel deployment, marketing sites with dynamic sections",
		Layer:       "web-ui",
		HasUI:       true,
		Tier:        2,
	},
	{
		ID:          "typescript-fastify",
		Title:       "TypeScript + Fastify",
		Summary:     "Node.js API — schema-driven, typed routes, plugin architecture",
		Dir:         "typescript-fastify",
		ScaffoldCmd: "npm init -y",
		UseCase:     "Node.js API services, microservices, typed backends — Fastify over Express, always",
		Layer:       "worker",
		HasUI:       false,
		Tier:        2,
	},
	{
		ID:          "python-django",
		Title:       "Python + Django",
		Summary:     "Python full-stack web — admin, ORM, batteries-included",
		Dir:         "python-django",
		ScaffoldCmd: "django-admin startproject {{name}}",
		UseCase:     "Admin-heavy apps, content management, Python full-stack web, rapid prototyping",
		Layer:       "rapid-product",
		HasUI:       true,
		Tier:        2,
	},
	{
		ID:          "laravel",
		Title:       "Laravel",
		Summary:     "PHP full-stack — Eloquent ORM, queues, Inertia, blade templates",
		Dir:         "laravel",
		ScaffoldCmd: "composer create-project laravel/laravel {{name}}",
		UseCase:     "PHP teams, rapid SaaS prototyping, content-driven web apps",
		Layer:       "rapid-product",
		HasUI:       true,
		Tier:        2,
	},
}

// Addons lists every available add-on.
var Addons = []Addon{
	{
		ID:      "data-intensive",
		Title:   "Data-Intensive",
		Summary: "Postgres, NATS, Parquet, event-driven architecture",
		Dir:     "data-intensive",
	},
	{
		ID:      "frontend-craft",
		Title:   "Frontend Craft",
		Summary: "Visual discipline, component composition, accessibility, and motion — framework agnostic",
		Dir:     "frontend-craft",
	},
}

// FindProfile returns the profile with the given ID, or nil if not found.
func FindProfile(id string) *Profile {
	for i := range Profiles {
		if Profiles[i].ID == id {
			return &Profiles[i]
		}
	}
	return nil
}

// FindAddon returns the addon with the given ID, or nil if not found.
func FindAddon(id string) *Addon {
	for i := range Addons {
		if Addons[i].ID == id {
			return &Addons[i]
		}
	}
	return nil
}

// ProfileIDs returns a slice of all profile IDs.
func ProfileIDs() []string {
	ids := make([]string, len(Profiles))
	for i, p := range Profiles {
		ids[i] = p.ID
	}
	return ids
}

// AddonIDs returns a slice of all addon IDs.
func AddonIDs() []string {
	ids := make([]string, len(Addons))
	for i, a := range Addons {
		ids[i] = a.ID
	}
	return ids
}
