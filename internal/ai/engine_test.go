package ai

import (
	"testing"
)

func TestParseSelection_ValidJSON(t *testing.T) {
	tests := []struct {
		name    string
		input   string
		wantID  string
		wantErr bool
	}{
		{
			name:   "clean JSON",
			input:  `{"profile_id":"elixir-phoenix","addon_ids":["data-intensive"],"asset_ids":["asset.palette.obsidian-indigo"],"confidence":0.9,"rationale":"good fit"}`,
			wantID: "elixir-phoenix",
		},
		{
			name:   "strips profile. prefix",
			input:  `{"profile_id":"profile.typescript-nextjs","confidence":0.88,"rationale":"react"}`,
			wantID: "typescript-nextjs",
		},
		{
			name:    "garbage input",
			input:   "I don't know what to pick",
			wantErr: true,
		},
		{
			name:    "empty input",
			input:   "",
			wantErr: true,
		},
	}

	for _, tt := range tests {
		t.Run(tt.name, func(t *testing.T) {
			sel, err := ParseSelection(tt.input)
			if tt.wantErr {
				if err == nil {
					t.Fatal("expected error, got nil")
				}
				return
			}
			if err != nil {
				t.Fatalf("unexpected error: %v", err)
			}
			if sel.ProfileID != tt.wantID {
				t.Errorf("profile_id = %q, want %q", sel.ProfileID, tt.wantID)
			}
		})
	}
}

func TestParseSelection_NormalizesAddons(t *testing.T) {
	input := `{"profile_id":"elixir-phoenix","addon_ids":["addon.data-intensive","addon.data-intensive","frontend-craft"],"confidence":0.9,"rationale":"test"}`
	sel, err := ParseSelection(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sel.AddonIDs) != 2 {
		t.Fatalf("expected 2 addons, got %d: %v", len(sel.AddonIDs), sel.AddonIDs)
	}
	if sel.AddonIDs[0] != "data-intensive" {
		t.Errorf("addon[0] = %q, want %q", sel.AddonIDs[0], "data-intensive")
	}
	if sel.AddonIDs[1] != "frontend-craft" {
		t.Errorf("addon[1] = %q, want %q", sel.AddonIDs[1], "frontend-craft")
	}
}

func TestParseSelection_FiltersProfileAndAddonFromAssets(t *testing.T) {
	input := `{"profile_id":"ruby-rails","addon_ids":[],"asset_ids":["profile.ruby-rails","addon.data-intensive","asset.lint.strict"],"confidence":0.9,"rationale":"test"}`
	sel, err := ParseSelection(input)
	if err != nil {
		t.Fatalf("unexpected error: %v", err)
	}
	if len(sel.AssetIDs) != 1 {
		t.Fatalf("expected 1 asset, got %d: %v", len(sel.AssetIDs), sel.AssetIDs)
	}
	if sel.AssetIDs[0] != "asset.lint.strict" {
		t.Errorf("asset[0] = %q, want %q", sel.AssetIDs[0], "asset.lint.strict")
	}
}

func TestParseFileOutput(t *testing.T) {
	input := "===FILE: .github/copilot-instructions.md===\n# Project Standards\n\nSome content here.\n===END_FILE===\n\n===FILE: AGENTS.md===\n# Agent Rules\n\nMore content.\n===END_FILE===\n"
	files := ParseFileOutput(input)
	if len(files) != 2 {
		t.Fatalf("expected 2 files, got %d", len(files))
	}
	if files[0].Path != ".github/copilot-instructions.md" {
		t.Errorf("file[0].Path = %q", files[0].Path)
	}
	if files[1].Path != "AGENTS.md" {
		t.Errorf("file[1].Path = %q", files[1].Path)
	}
}

func TestParseFileOutput_Empty(t *testing.T) {
	files := ParseFileOutput("No file blocks here at all.")
	if len(files) != 0 {
		t.Fatalf("expected 0 files, got %d", len(files))
	}
}

func TestIsReady(t *testing.T) {
	tests := []struct {
		input string
		want  bool
	}{
		{"Great choice! READY_TO_GENERATE", true},
		{"ready_to_generate", true},
		{"Still thinking about options...", false},
		{"", false},
	}
	for _, tt := range tests {
		if got := IsReady(tt.input); got != tt.want {
			t.Errorf("IsReady(%q) = %v, want %v", tt.input, got, tt.want)
		}
	}
}
