package ai

import (
	"fmt"
	"sort"
	"strings"
)

// ContextAsset is a selectable instruction source defined in this repository.
type ContextAsset struct {
	ID           string
	Category     string
	Label        string
	Summary      string
	TemplatePath string
}

func catalog() []ContextAsset {
	return []ContextAsset{
		// ── Core (always included) ───────────────────────────────────
		{
			ID:           "core.copilot",
			Category:     "core",
			Label:        "Core Copilot Standards",
			Summary:      "Always-on engineering standards for architecture, naming, and implementation quality",
			TemplatePath: "core/.github/copilot-instructions.md",
		},
		{
			ID:           "core.architecture",
			Category:     "practices",
			Label:        "Architecture Practices",
			Summary:      "Functional-first decomposition, pure core / imperative edge boundaries, and layered composition",
			TemplatePath: "core/.github/instructions/architecture.instructions.md",
		},
		{
			ID:           "core.agents",
			Category:     "collaboration",
			Label:        "Agent Collaboration Rules",
			Summary:      "Ground rules for multi-agent workflow, ownership boundaries, and quality checks",
			TemplatePath: "core/AGENTS.md",
		},
		{
			ID:           "core.design-system",
			Category:     "design",
			Label:        "Design System Baseline",
			Summary:      "Dark-first visual identity, typography, spacing, and component DNA — the visual foundation that all generated apps share",
			TemplatePath: "core/.github/instructions/design-system.instructions.md",
		},

		// ── Tier 1 Profiles (author's opinionated picks) ────────────
		{
			ID:           "profile.elixir-phoenix",
			Category:     "framework",
			Label:        "Elixir + Phoenix",
			Summary:      "Full-stack real-time web — LiveView, Ecto, OTP. Best AI context: entire app in one framework",
			TemplatePath: "profiles/elixir-phoenix/.github/instructions/elixir-phoenix.instructions.md",
		},
		{
			ID:           "profile.typescript-sveltekit",
			Category:     "framework",
			Label:        "TypeScript + SvelteKit",
			Summary:      "Full-stack JS web — intuitive reactivity, SSR, minimal boilerplate. Best JS framework for AI",
			TemplatePath: "profiles/typescript-sveltekit/.github/instructions/typescript-sveltekit.instructions.md",
		},
		{
			ID:           "profile.ruby-rails",
			Category:     "framework",
			Label:        "Ruby on Rails",
			Summary:      "Rapid full-stack web — generators, convention over configuration, fast to production",
			TemplatePath: "profiles/ruby-rails/.github/instructions/ruby-rails.instructions.md",
		},

		// ── Tier 2 Profiles (domain-specific) ────────────────────────
		{
			ID:           "profile.typescript-nextjs",
			Category:     "framework",
			Label:        "TypeScript + Next.js",
			Summary:      "React ecosystem full-stack — App Router, RSC, Vercel-optimized",
			TemplatePath: "profiles/typescript-nextjs/.github/instructions/typescript-nextjs.instructions.md",
		},
		{
			ID:           "profile.typescript-fastify",
			Category:     "framework",
			Label:        "TypeScript + Fastify",
			Summary:      "Node.js API service — schema-driven routes, typed contracts, plugin architecture",
			TemplatePath: "profiles/typescript-fastify/.github/instructions/typescript-fastify.instructions.md",
		},
		{
			ID:           "profile.go-service",
			Category:     "framework",
			Label:        "Go Service",
			Summary:      "Idiomatic Go service architecture with stdlib-first bias and explicit boundaries",
			TemplatePath: "profiles/go-service/.github/instructions/go-service.instructions.md",
		},
		{
			ID:           "profile.dotnet-api",
			Category:     "framework",
			Label:        ".NET API",
			Summary:      "C# API architecture with clear boundaries and maintainable service design",
			TemplatePath: "profiles/dotnet-api/.github/instructions/dotnet-api.instructions.md",
		},
		{
			ID:           "profile.python-fastapi",
			Category:     "framework",
			Label:        "Python + FastAPI",
			Summary:      "Async Python APIs with Pydantic types, ideal for ML/data service backends",
			TemplatePath: "profiles/python-fastapi/.github/instructions/python-fastapi.instructions.md",
		},
		{
			ID:           "profile.python-django",
			Category:     "framework",
			Label:        "Python + Django",
			Summary:      "Batteries-included Python web — admin, ORM, auth, content management",
			TemplatePath: "profiles/python-django/.github/instructions/python-django.instructions.md",
		},
		{
			ID:           "profile.dart-flutter",
			Category:     "framework",
			Label:        "Dart + Flutter",
			Summary:      "Cross-platform native apps — single codebase, widget composition, platform channels",
			TemplatePath: "profiles/dart-flutter/.github/instructions/dart-flutter.instructions.md",
		},
		{
			ID:           "profile.rust-axum",
			Category:     "framework",
			Label:        "Rust + Axum",
			Summary:      "Performance-critical services — Tokio-based, type-safe, zero-cost abstractions",
			TemplatePath: "profiles/rust-axum/.github/instructions/rust-axum.instructions.md",
		},
		{
			ID:           "profile.laravel",
			Category:     "framework",
			Label:        "Laravel",
			Summary:      "Laravel + Inertia project conventions for product-focused web apps",
			TemplatePath: "profiles/laravel/.github/instructions/laravel.instructions.md",
		},

		// ── Add-ons ──────────────────────────────────────────────────
		{
			ID:           "addon.data-intensive",
			Category:     "architecture",
			Label:        "Data-Intensive Add-on",
			Summary:      "Patterns for event streams, durable storage, and resilient data processing",
			TemplatePath: "addons/data-intensive/.github/instructions/data-intensive.instructions.md",
		},
		{
			ID:           "addon.frontend-craft",
			Category:     "ui",
			Label:        "Frontend Craft Add-on",
			Summary:      "Framework-agnostic visual discipline, component composition, accessibility, motion, and styling system guidance",
			TemplatePath: "addons/frontend-craft/.github/instructions/frontend-craft.instructions.md",
		},

		// ── Design Assets ────────────────────────────────────────────
		{
			ID:           "asset.palette.heroui-blue",
			Category:     "palette",
			Label:        "HeroUI Blue Scale Palette",
			Summary:      "Blue-centered semantic scale inspired by your attached `colors.ts` palette structure",
			TemplatePath: "assets/palettes/heroui-blue.instructions.md",
		},
		{
			ID:           "asset.palette.obsidian-indigo",
			Category:     "palette",
			Label:        "Obsidian + Indigo Palette",
			Summary:      "Dark Phoenix-style UI palette inspired by your attached LiveView layout styling",
			TemplatePath: "assets/palettes/obsidian-indigo.instructions.md",
		},
		{
			ID:           "asset.fonts.inter-jetbrains",
			Category:     "fonts",
			Label:        "Inter + JetBrains Mono",
			Summary:      "Sans + monospace pairing for product UI and dev-facing surfaces",
			TemplatePath: "assets/fonts/inter-jetbrains.instructions.md",
		},

		// ── Quality Assets ───────────────────────────────────────────
		{
			ID:           "asset.lint.strict",
			Category:     "linting",
			Label:        "Strict Linting",
			Summary:      "Fail-on-warning lint posture and formatting consistency expectations",
			TemplatePath: "assets/linting/strict.instructions.md",
		},
		{
			ID:           "asset.testing.pragmatic",
			Category:     "testing",
			Label:        "Pragmatic Testing",
			Summary:      "Fast feedback testing pyramid with contract and integration confidence",
			TemplatePath: "assets/testing/pragmatic.instructions.md",
		},
	}
}

func catalogMap() map[string]ContextAsset {
	byID := make(map[string]ContextAsset)
	for _, item := range catalog() {
		byID[item.ID] = item
	}
	return byID
}

func catalogSummaryLines() []string {
	items := catalog()
	sort.Slice(items, func(i, j int) bool { return items[i].ID < items[j].ID })
	lines := make([]string, 0, len(items))
	for _, item := range items {
		lines = append(lines, fmt.Sprintf("- %s | %s | %s", item.ID, item.Category, item.Summary))
	}
	return lines
}

func resolveContextAssets(selection Selection) ([]ContextAsset, error) {
	byID := catalogMap()

	base := []string{"core.copilot", "core.architecture", "core.agents", "core.design-system"}
	resolvedIDs := make([]string, 0, len(base)+len(selection.AddonIDs)+len(selection.AssetIDs)+2)
	resolvedIDs = append(resolvedIDs, base...)

	if selection.ProfileID != "" {
		profileID := selection.ProfileID
		if !strings.HasPrefix(profileID, "profile.") {
			profileID = "profile." + profileID
		}
		resolvedIDs = append(resolvedIDs, profileID)
	}
	for _, addonID := range selection.AddonIDs {
		id := addonID
		if !strings.HasPrefix(id, "addon.") {
			id = "addon." + id
		}
		resolvedIDs = append(resolvedIDs, id)
	}
	resolvedIDs = append(resolvedIDs, selection.AssetIDs...)

	seen := make(map[string]bool)
	resolved := make([]ContextAsset, 0, len(resolvedIDs))
	for _, id := range resolvedIDs {
		if id == "" || seen[id] {
			continue
		}
		asset, ok := byID[id]
		if !ok {
			return nil, fmt.Errorf("unknown context asset %q", id)
		}
		seen[id] = true
		resolved = append(resolved, asset)
	}

	return resolved, nil
}
