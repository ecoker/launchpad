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
			ID:           "profile.typescript-react",
			Category:     "framework",
			Label:        "TypeScript + React",
			Summary:      "React Router v7, TypeScript-first frontend conventions, and pragmatic component patterns",
			TemplatePath: "profiles/typescript-react/.github/instructions/typescript-react.instructions.md",
		},
		{
			ID:           "profile.python-data",
			Category:     "framework",
			Label:        "Python Data / AI",
			Summary:      "Pydantic-centric Python data and agent engineering conventions",
			TemplatePath: "profiles/python-data/.github/instructions/python-data.instructions.md",
		},
		{
			ID:           "profile.elixir-phoenix",
			Category:     "framework",
			Label:        "Elixir + Phoenix",
			Summary:      "Phoenix + LiveView functional conventions for highly concurrent real-time systems",
			TemplatePath: "profiles/elixir-phoenix/.github/instructions/elixir-phoenix.instructions.md",
		},
		{
			ID:           "profile.dotnet-api",
			Category:     "framework",
			Label:        ".NET API",
			Summary:      "C# API architecture with clear boundaries and maintainable service design",
			TemplatePath: "profiles/dotnet-api/.github/instructions/dotnet-api.instructions.md",
		},
		{
			ID:           "profile.laravel",
			Category:     "framework",
			Label:        "Laravel",
			Summary:      "Laravel + Inertia project conventions for product-focused web apps",
			TemplatePath: "profiles/laravel/.github/instructions/laravel.instructions.md",
		},
		{
			ID:           "profile.go-service",
			Category:     "framework",
			Label:        "Go Service",
			Summary:      "Idiomatic Go service architecture with stdlib-first bias and explicit boundaries",
			TemplatePath: "profiles/go-service/.github/instructions/go-service.instructions.md",
		},
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
			Summary:      "Visual quality, accessibility, and interaction craftsmanship conventions",
			TemplatePath: "addons/frontend-craft/.github/instructions/frontend-craft.instructions.md",
		},
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
		{
			ID:           "asset.framework.csharp-dotnet",
			Category:     "framework-opinion",
			Label:        "C# (.NET) Opinion Pack",
			Summary:      "Opinionated .NET service standards and architecture defaults",
			TemplatePath: "assets/frameworks/csharp-dotnet.instructions.md",
		},
		{
			ID:           "asset.framework.dart-flutter",
			Category:     "framework-opinion",
			Label:        "Dart (Flutter) Opinion Pack",
			Summary:      "Opinionated Flutter app architecture and quality defaults",
			TemplatePath: "assets/frameworks/dart-flutter.instructions.md",
		},
		{
			ID:           "asset.framework.rust",
			Category:     "framework-opinion",
			Label:        "Rust Opinion Pack",
			Summary:      "Opinionated Rust architecture and correctness-first defaults",
			TemplatePath: "assets/frameworks/rust.instructions.md",
		},
		{
			ID:           "asset.framework.typescript-nextjs",
			Category:     "framework-opinion",
			Label:        "TypeScript Next.js Opinion Pack",
			Summary:      "Opinionated Next.js app structure and server/client boundary guidance",
			TemplatePath: "assets/frameworks/typescript-nextjs.instructions.md",
		},
		{
			ID:           "asset.framework.typescript-fastify",
			Category:     "framework-opinion",
			Label:        "TypeScript Fastify Opinion Pack",
			Summary:      "Opinionated Fastify service standards with typed route contracts",
			TemplatePath: "assets/frameworks/typescript-fastify.instructions.md",
		},
		{
			ID:           "asset.framework.typescript-svelte",
			Category:     "framework-opinion",
			Label:        "TypeScript Svelte Opinion Pack",
			Summary:      "Opinionated Svelte project structure and reactive design practices",
			TemplatePath: "assets/frameworks/typescript-svelte.instructions.md",
		},
		{
			ID:           "asset.framework.typescript-react-router-v7",
			Category:     "framework-opinion",
			Label:        "TypeScript React Router v7 Opinion Pack",
			Summary:      "Opinionated React Router v7 route-module architecture guidance",
			TemplatePath: "assets/frameworks/typescript-react-router-v7.instructions.md",
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

	base := []string{"core.copilot", "core.architecture", "core.agents"}
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
