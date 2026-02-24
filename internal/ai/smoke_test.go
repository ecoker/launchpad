//go:build smoke

// This is a manual smoke test that exercises the full GenerateFiles pipeline
// against the real OpenAI API. It requires a valid API key in .env or
// OPENAI_API_KEY environment variable.
//
// Run:  go test -tags smoke -run TestSmoke -v -timeout 120s ./internal/ai/
package ai

import (
	"context"
	"fmt"
	"os"
	"path/filepath"
	"strings"
	"testing"
	"time"
)

func loadTestAPIKey(t *testing.T) string {
	t.Helper()

	if key := os.Getenv("OPENAI_API_KEY"); key != "" {
		return key
	}

	// Try .env in repo root (test runs from internal/ai/)
	for _, rel := range []string{"../../.env", ".env"} {
		data, err := os.ReadFile(rel)
		if err != nil {
			continue
		}
		for _, line := range strings.Split(string(data), "\n") {
			line = strings.TrimSpace(line)
			if line == "" || strings.HasPrefix(line, "#") {
				continue
			}
			parts := strings.SplitN(line, "=", 2)
			if len(parts) != 2 {
				continue
			}
			key := strings.TrimSpace(parts[0])
			val := strings.TrimSpace(parts[1])
			val = strings.Trim(val, `"'`)
			if key == "OPENAI_API_KEY" || key == "KEY" {
				return val
			}
		}
	}

	t.Skip("No API key found — set OPENAI_API_KEY or create .env with KEY=...")
	return ""
}

// TestSmokeGenerateFiles sends a hardcoded Selection to GenerateFiles and
// writes the output to tmp/smoke-output/ for manual inspection.
func TestSmokeGenerateFiles(t *testing.T) {
	apiKey := loadTestAPIKey(t)

	provider := NewOpenAIProvider(apiKey)
	engine := NewEngine(provider)

	sel := &Selection{
		ProfileID:  "elixir-phoenix",
		AddonIDs:   []string{"frontend-craft"},
		AssetIDs:   []string{"asset.palette.obsidian-indigo", "asset.fonts.inter-jetbrains", "asset.testing.pragmatic"},
		Confidence: 0.95,
		Rationale:  "Smoke test: Phoenix + frontend-craft + obsidian palette + Inter fonts + pragmatic testing",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 90*time.Second)
	defer cancel()

	t.Log("Calling GenerateFiles (this hits the real API)...")
	start := time.Now()
	files, err := engine.GenerateFiles(ctx, "smoke-test-app", sel)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("GenerateFiles failed after %s: %v", elapsed, err)
	}

	t.Logf("Generated %d files in %s", len(files), elapsed)

	// Write to tmp/smoke-output/
	outDir := filepath.Join("..", "..", "tmp", "smoke-output")
	if err := os.RemoveAll(outDir); err != nil {
		t.Fatalf("cleaning output dir: %v", err)
	}

	for _, f := range files {
		dest := filepath.Join(outDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(dest, []byte(f.Content), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		t.Logf("  %s (%d bytes)", f.Path, len(f.Content))
	}

	t.Logf("\nOutput written to %s", outDir)

	// Basic assertions
	if len(files) < 3 {
		t.Errorf("Expected at least 3 files (copilot-instructions, AGENTS, start.prompt), got %d", len(files))
	}

	var hasCopilot, hasAgents, hasStartPrompt bool
	for _, f := range files {
		content := strings.ToLower(f.Content)
		switch {
		case strings.Contains(f.Path, "copilot-instructions.md"):
			hasCopilot = true
			// Should reference Tailwind
			if !strings.Contains(content, "tailwind") {
				t.Error("copilot-instructions.md does not mention Tailwind")
			}
		case strings.Contains(f.Path, "AGENTS.md"):
			hasAgents = true
		case strings.Contains(f.Path, "start.prompt.md"):
			hasStartPrompt = true
			// Should reference mix phx.new
			if !strings.Contains(content, "mix phx.new") {
				t.Error("start.prompt.md does not reference 'mix phx.new' scaffold command")
			}
		}

		// Check for design token synthesis
		if strings.Contains(f.Path, "design-system") || strings.Contains(f.Path, "frontend-craft") {
			// Should have framework-appropriate content — Phoenix/LiveView, not React
			if strings.Contains(content, "shadcn/ui") || strings.Contains(content, "framer motion") {
				t.Errorf("%s contains React-specific references (shadcn/ui or Framer Motion) for a Phoenix project", f.Path)
			}
			// Should reference obsidian palette colors
			if !strings.Contains(content, "#0f0f0f") && !strings.Contains(content, "#6366f1") && !strings.Contains(content, "obsidian") && !strings.Contains(content, "indigo") {
				t.Logf("WARNING: %s doesn't appear to reference obsidian-indigo palette tokens", f.Path)
			}
		}
	}

	if !hasCopilot {
		t.Error("Missing .github/copilot-instructions.md")
	}
	if !hasAgents {
		t.Error("Missing AGENTS.md")
	}
	if !hasStartPrompt {
		t.Error("Missing start.prompt.md")
	}

	// Print a summary of all files for manual review
	fmt.Println("\n=== GENERATED FILES ===")
	for _, f := range files {
		fmt.Printf("\n--- %s ---\n", f.Path)
		// Print first 20 lines
		lines := strings.Split(f.Content, "\n")
		limit := 20
		if len(lines) < limit {
			limit = len(lines)
		}
		for _, line := range lines[:limit] {
			fmt.Println(line)
		}
		if len(lines) > 20 {
			fmt.Printf("... (%d more lines)\n", len(lines)-20)
		}
	}
}

// TestSmokeUIAutoInclude verifies that UI stacks with NO explicit addons/assets
// still get frontend-craft, palette, and font assets auto-included.
func TestSmokeUIAutoInclude(t *testing.T) {
	apiKey := loadTestAPIKey(t)

	provider := NewOpenAIProvider(apiKey)
	engine := NewEngine(provider)

	// Rails with deliberately empty addons/assets — tests auto-include
	sel := &Selection{
		ProfileID:  "ruby-rails",
		AddonIDs:   []string{},
		AssetIDs:   []string{},
		Confidence: 0.95,
		Rationale:  "UI smoke test: Rails with zero addons — testing auto-include of visual assets",
	}

	ctx, cancel := context.WithTimeout(context.Background(), 120*time.Second)
	defer cancel()

	t.Log("Calling GenerateFiles for Rails with no explicit UI addons...")
	start := time.Now()
	files, err := engine.GenerateFiles(ctx, "ui-smoke-test", sel)
	elapsed := time.Since(start)
	if err != nil {
		t.Fatalf("GenerateFiles failed after %s: %v", elapsed, err)
	}

	t.Logf("Generated %d files in %s", len(files), elapsed)

	// Write to tmp/ui-smoke/
	outDir := filepath.Join("..", "..", "tmp", "ui-smoke")
	if err := os.RemoveAll(outDir); err != nil {
		t.Fatalf("cleaning output dir: %v", err)
	}

	for _, f := range files {
		dest := filepath.Join(outDir, f.Path)
		if err := os.MkdirAll(filepath.Dir(dest), 0o755); err != nil {
			t.Fatalf("mkdir: %v", err)
		}
		if err := os.WriteFile(dest, []byte(f.Content), 0o644); err != nil {
			t.Fatalf("write: %v", err)
		}
		t.Logf("  %s (%d bytes)", f.Path, len(f.Content))
	}

	// Check for UI signals
	hasGradient := false
	hasGlow := false
	hasPalette := false
	hasTailwind := false
	hasFrontendCraft := false

	for _, f := range files {
		content := strings.ToLower(f.Content)
		if strings.Contains(content, "gradient") {
			hasGradient = true
		}
		if strings.Contains(content, "glow") {
			hasGlow = true
		}
		if strings.Contains(content, "#6366f1") || strings.Contains(content, "obsidian") || strings.Contains(content, "indigo") {
			hasPalette = true
		}
		if strings.Contains(content, "tailwind") {
			hasTailwind = true
		}
		if strings.Contains(content, "frontend craft") || strings.Contains(content, "visual discipline") ||
			strings.Contains(content, "component composition") || strings.Contains(content, "accessibility") {
			hasFrontendCraft = true
		}
	}

	if !hasTailwind {
		t.Error("FAIL: No Tailwind references found in generated output")
	}
	if !hasPalette {
		t.Error("FAIL: No Obsidian/Indigo palette references found — auto-include may not be working")
	}
	if !hasFrontendCraft {
		t.Error("FAIL: No frontend-craft concepts found — auto-include may not be working")
	}

	// These are nice-to-haves from the enhanced design system
	if !hasGradient {
		t.Log("NOTE: No gradient references found (expected from enhanced design system)")
	}
	if !hasGlow {
		t.Log("NOTE: No glow references found (expected from enhanced design system)")
	}

	// Print design-relevant files for manual inspection
	for _, f := range files {
		if strings.Contains(f.Path, "design") || strings.Contains(f.Path, "copilot-instructions") ||
			strings.Contains(f.Path, "frontend") {
			fmt.Printf("\n--- %s ---\n", f.Path)
			lines := strings.Split(f.Content, "\n")
			limit := 50
			if len(lines) < limit {
				limit = len(lines)
			}
			for _, line := range lines[:limit] {
				fmt.Println(line)
			}
			if len(lines) > limit {
				fmt.Printf("... (%d more lines)\n", len(lines)-limit)
			}
		}
	}

	t.Logf("\nOutput written to %s", outDir)
}
