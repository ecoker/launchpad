package ai

import (
	"testing"
)

func TestValidateSelectionCompatibility(t *testing.T) {
	tests := []struct {
		name       string
		selection  Selection
		wantIssues int
	}{
		{
			name:       "valid tier 1 with compatible addon",
			selection:  Selection{ProfileID: "elixir-phoenix", AddonIDs: []string{"data-intensive", "frontend-craft"}},
			wantIssues: 0,
		},
		{
			name:       "empty profile",
			selection:  Selection{},
			wantIssues: 1,
		},
		{
			name:       "unknown profile",
			selection:  Selection{ProfileID: "cobol-mainframe"},
			wantIssues: 1,
		},
		{
			name:       "frontend-craft incompatible with go-service",
			selection:  Selection{ProfileID: "go-service", AddonIDs: []string{"frontend-craft"}},
			wantIssues: 1,
		},
		{
			name:       "frontend-craft incompatible with typescript-fastify",
			selection:  Selection{ProfileID: "typescript-fastify", AddonIDs: []string{"frontend-craft"}},
			wantIssues: 1,
		},
		{
			name:       "data-intensive compatible with everything",
			selection:  Selection{ProfileID: "rust-axum", AddonIDs: []string{"data-intensive"}},
			wantIssues: 0,
		},
		{
			name:       "duplicate addon",
			selection:  Selection{ProfileID: "ruby-rails", AddonIDs: []string{"data-intensive", "data-intensive"}},
			wantIssues: 1,
		},
		{
			name:       "duplicate asset",
			selection:  Selection{ProfileID: "ruby-rails", AssetIDs: []string{"asset.lint.strict", "asset.lint.strict"}},
			wantIssues: 1,
		},
		{
			name: "multiple palettes rejected",
			selection: Selection{
				ProfileID: "ruby-rails",
				AssetIDs:  []string{"asset.palette.heroui-blue", "asset.palette.obsidian-indigo"},
			},
			wantIssues: 1,
		},
		{
			name: "multiple testing assets rejected",
			selection: Selection{
				ProfileID: "ruby-rails",
				AssetIDs:  []string{"asset.testing.pragmatic", "asset.testing.comprehensive"},
			},
			wantIssues: 1,
		},
		{
			name: "one of each category is fine",
			selection: Selection{
				ProfileID: "typescript-sveltekit",
				AddonIDs:  []string{"frontend-craft"},
				AssetIDs:  []string{"asset.palette.heroui-blue", "asset.fonts.inter-jetbrains", "asset.lint.strict", "asset.testing.pragmatic"},
			},
			wantIssues: 0,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			issues := ValidateSelectionCompatibility(tt.selection)
			if len(issues) != tt.wantIssues {
				t.Errorf("got %d issues, want %d: %v", len(issues), tt.wantIssues, issues)
			}
		})
	}
}
