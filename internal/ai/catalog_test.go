package ai

import (
	"strings"
	"testing"

	"github.com/ecoker/launchpad/templates"
)

// TestCatalogAssetTemplatesExist verifies every registered catalog entry
// has a corresponding template file in the embedded filesystem.
func TestCatalogAssetTemplatesExist(t *testing.T) {
	for _, asset := range catalog() {
		t.Run(asset.ID, func(t *testing.T) {
			data, err := templates.FS.ReadFile(asset.TemplatePath)
			if err != nil {
				t.Fatalf("template %q not found in embedded FS: %v", asset.TemplatePath, err)
			}
			if len(data) == 0 {
				t.Errorf("template %q is empty", asset.TemplatePath)
			}
		})
	}
}

// TestNewAssetsRegistered confirms the two new assets exist in the catalog.
func TestNewAssetsRegistered(t *testing.T) {
	byID := catalogMap()

	t.Run("server-patterns exists", func(t *testing.T) {
		asset, ok := byID["asset.server.patterns"]
		if !ok {
			t.Fatal("asset.server.patterns not found in catalog")
		}
		if asset.Category != "server" {
			t.Errorf("category = %q, want %q", asset.Category, "server")
		}
	})

	t.Run("testing asset updated", func(t *testing.T) {
		asset, ok := byID["asset.testing.pragmatic"]
		if !ok {
			t.Fatal("asset.testing.pragmatic not found in catalog")
		}
		// Verify the template content has the new comprehensive content
		data, err := templates.FS.ReadFile(asset.TemplatePath)
		if err != nil {
			t.Fatalf("read template: %v", err)
		}
		content := string(data)
		if !strings.Contains(content, "## Test pyramid") {
			t.Error("testing template missing '## Test pyramid' section")
		}
		if !strings.Contains(content, "## Framework-specific guidance") {
			t.Error("testing template missing '## Framework-specific guidance' section")
		}
	})
}

// TestServerPatternsContent verifies the server-patterns template has key sections.
func TestServerPatternsContent(t *testing.T) {
	data, err := templates.FS.ReadFile("assets/server/server-patterns.instructions.md")
	if err != nil {
		t.Fatalf("read server-patterns: %v", err)
	}
	content := string(data)

	sections := []string{
		"## Validation",
		"## Error handling",
		"## Form actions and mutations",
		"## Data access",
		"## Background jobs",
	}
	for _, s := range sections {
		if !strings.Contains(content, s) {
			t.Errorf("missing section %q", s)
		}
	}
}

// TestResolveContextAssetsWithServerPatterns verifies the new asset
// can be resolved through the standard selection pipeline.
func TestResolveContextAssetsWithServerPatterns(t *testing.T) {
	sel := Selection{
		ProfileID: "elixir-phoenix",
		AssetIDs:  []string{"asset.server.patterns", "asset.testing.pragmatic"},
	}
	assets, err := resolveContextAssets(sel)
	if err != nil {
		t.Fatalf("resolveContextAssets: %v", err)
	}

	found := map[string]bool{}
	for _, a := range assets {
		found[a.ID] = true
	}

	// Core assets always included
	for _, id := range []string{"core.copilot", "core.architecture", "core.agents", "core.design-system"} {
		if !found[id] {
			t.Errorf("expected core asset %q in resolved set", id)
		}
	}

	// Explicitly selected assets
	if !found["asset.server.patterns"] {
		t.Error("asset.server.patterns not in resolved set")
	}
	if !found["asset.testing.pragmatic"] {
		t.Error("asset.testing.pragmatic not in resolved set")
	}

	// Profile
	if !found["profile.elixir-phoenix"] {
		t.Error("profile.elixir-phoenix not in resolved set")
	}

	// UI profile auto-includes frontend-craft
	if !found["addon.frontend-craft"] {
		t.Error("addon.frontend-craft should be auto-included for UI profile")
	}
}
