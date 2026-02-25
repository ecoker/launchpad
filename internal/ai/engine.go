package ai

import (
	"context"
	"encoding/json"
	"fmt"
	"sort"
	"strings"

	"github.com/ecoker/launchpad/internal/scaffold"
	"github.com/ecoker/launchpad/templates"
)

// FileOutput represents a single file the AI wants to create.
type FileOutput struct {
	Path    string
	Content string
}

// Selection is the resolved setup used to load context assets.
type Selection struct {
	ProfileID  string   `json:"profile_id"`
	AddonIDs   []string `json:"addon_ids,omitempty"`
	AssetIDs   []string `json:"asset_ids,omitempty"`
	Confidence float64  `json:"confidence"`
	Rationale  string   `json:"rationale"`
}

// confidenceThreshold is the minimum self-reported confidence the model must
// return for us to proceed with generation. This is a soft heuristic — LLM
// confidence scores are uncalibrated — but in practice it catches cases where
// the conversation was too vague to produce a useful selection. Tuned through
// manual testing; not a statistical guarantee.
const confidenceThreshold = 0.72

// ReadyToken is the phrase the model appends to signal readiness.
const ReadyToken = "READY_TO_GENERATE"

// Engine orchestrates the multi-turn conversation and generation workflow.
// It delegates all LLM communication to a Provider implementation.
type Engine struct {
	provider Provider
}

// NewEngine creates a new Engine backed by the given Provider.
func NewEngine(provider Provider) *Engine {
	return &Engine{provider: provider}
}

// Chat sends a user message and returns the assistant's reply.
func (e *Engine) Chat(ctx context.Context, message string) (string, error) {
	if strings.TrimSpace(message) == "" {
		return "", fmt.Errorf("empty message")
	}
	// Always send instructions — the Responses API does NOT carry them
	// across previous_response_id chains.
	return e.provider.Send(ctx, message, conversationSystemPrompt())
}

// IsReady reports whether the assistant reply contains the readiness token.
func IsReady(reply string) bool {
	normalized := strings.ToUpper(strings.TrimSpace(reply))
	if strings.Contains(normalized, ReadyToken) {
		return true
	}
	return strings.Contains(normalized, "READY TO GENERATE")
}

// ExtractDecision silently reads the current thread and returns a structured Selection.
// This call is never shown to the user.
func (e *Engine) ExtractDecision(ctx context.Context) (*Selection, error) {
	extractPrompt := "Based on our conversation, extract the final stack decision.\n\n" +
		"Return ONLY valid JSON — no markdown, no prose:\n" +
		"{\n" +
		"  \"profile_id\": \"<elixir-phoenix|typescript-sveltekit|ruby-rails|typescript-nextjs|typescript-fastify|go-service|dotnet-api|java-spring|python-fastapi|python-django|dart-flutter|rust-axum|laravel>\",\n" +
		"  \"addon_ids\": [],\n" +
		"  \"asset_ids\": [],\n" +
		"  \"confidence\": 0.0,\n" +
		"  \"rationale\": \"one sentence\"\n" +
		"}\n\n" +
		"Asset IDs available:\n" + catalogIDLines()

	raw, err := e.provider.Send(ctx, extractPrompt, "")
	if err != nil {
		return nil, err
	}
	return parseSelection(raw)
}

// GenerateFiles loads the selected context assets and generates instruction files.
func (e *Engine) GenerateFiles(ctx context.Context, projectName string, sel *Selection) ([]FileOutput, error) {
	if sel == nil || sel.ProfileID == "" {
		return nil, fmt.Errorf("no stack selected")
	}
	if sel.Confidence < confidenceThreshold {
		return nil, fmt.Errorf(
			"confidence %.2f is below minimum %.2f — try describing your project in more detail",
			sel.Confidence, confidenceThreshold,
		)
	}
	if issues := ValidateSelectionCompatibility(*sel); len(issues) > 0 {
		return nil, fmt.Errorf("incompatible selection: %s", strings.Join(issues, "; "))
	}

	assets, err := resolveContextAssets(*sel)
	if err != nil {
		return nil, fmt.Errorf("resolving assets: %w", err)
	}

	var contextBlocks strings.Builder
	for _, asset := range assets {
		data, readErr := templates.FS.ReadFile(asset.TemplatePath)
		if readErr != nil {
			return nil, fmt.Errorf("reading asset %s: %w", asset.ID, readErr)
		}
		fmt.Fprintf(&contextBlocks, "===ASSET: %s===\n%s\n===END_ASSET===\n\n", asset.ID, string(data))
	}

	summary := make([]string, 0, len(assets))
	for _, a := range assets {
		summary = append(summary, fmt.Sprintf("%s (%s)", a.ID, a.Category))
	}
	sort.Strings(summary)

	// Pull scaffold command from the profile registry
	scaffoldInfo := scaffoldCommandForProfile(sel.ProfileID)

	// Check which assets are in the selection so we can
	// give the model explicit synthesis instructions.
	hasDesignSystem := false
	hasPalette := false
	hasFonts := false
	hasFrontendCraft := false
	hasServerPatterns := false
	hasTesting := false
	for _, a := range assets {
		switch {
		case a.ID == "core.design-system":
			hasDesignSystem = true
		case strings.HasPrefix(a.ID, "asset.palette."):
			hasPalette = true
		case strings.HasPrefix(a.ID, "asset.fonts."):
			hasFonts = true
		case a.ID == "addon.frontend-craft":
			hasFrontendCraft = true
		case a.ID == "asset.server.patterns":
			hasServerPatterns = true
		case a.ID == "asset.testing.pragmatic":
			hasTesting = true
		}
	}

	// Detect whether the selected profile has a UI surface.
	isUIStack := false
	if profile := scaffold.FindProfile(sel.ProfileID); profile != nil {
		isUIStack = profile.HasUI
	}

	var designGuidance strings.Builder
	if hasDesignSystem || hasPalette || hasFonts || hasFrontendCraft {
		designGuidance.WriteString("DESIGN SYSTEM SYNTHESIS:\n")
		designGuidance.WriteString("The assets below include visual identity guidance. When generating output files:\n")
		designGuidance.WriteString("- Merge the design-system baseline with any selected palette/font assets into\n")
		designGuidance.WriteString("  a single cohesive visual language. Don't repeat conflicting defaults.\n")
		if hasPalette {
			designGuidance.WriteString("- A palette asset is included. Use its specific color tokens as the concrete\n")
			designGuidance.WriteString("  values for the design-system's color guidance. The palette overrides generic\n")
			designGuidance.WriteString("  color suggestions in the baseline.\n")
		}
		if hasFonts {
			designGuidance.WriteString("- A font pairing asset is included. Use its specific fonts as the concrete\n")
			designGuidance.WriteString("  values for the design-system's typography guidance.\n")
		}
		if hasFrontendCraft {
			designGuidance.WriteString("- The frontend-craft addon is included. Its principles are framework-agnostic.\n")
			designGuidance.WriteString("  When generating instruction files, adapt ALL examples, component patterns,\n")
			designGuidance.WriteString("  animation techniques, and styling approaches to the selected framework's\n")
			designGuidance.WriteString("  idioms (e.g. LiveView function components for Phoenix, Svelte components\n")
			designGuidance.WriteString("  for SvelteKit, ViewComponent for Rails, Blade for Laravel, widgets for\n")
			designGuidance.WriteString("  Flutter). Do NOT emit React/JSX examples for non-React stacks.\n")
			designGuidance.WriteString("- IMPORTANT: The frontend-craft file MUST preserve guidance on ALL of these:\n")
			designGuidance.WriteString("  loading/empty/error state patterns, state management, motion/animation,\n")
			designGuidance.WriteString("  accessibility, and performance. These are the most actionable parts —\n")
			designGuidance.WriteString("  do NOT compress them away. Adapt examples to the selected framework.\n")
		}
		designGuidance.WriteString("- Generate a dedicated design-system.instructions.md that synthesizes the\n")
		designGuidance.WriteString("  baseline + palette + fonts into framework-appropriate tokens and setup.\n")
		designGuidance.WriteString("  The applyTo glob MUST match the selected framework's template/style files.\n\n")
	}

	// Build conditional asset instructions.
	var assetGuidance strings.Builder
	if hasServerPatterns {
		assetGuidance.WriteString("SERVER PATTERNS:\n")
		assetGuidance.WriteString("A server-patterns asset is included. Generate a dedicated\n")
		assetGuidance.WriteString("server-patterns.instructions.md file with validation, error handling,\n")
		assetGuidance.WriteString("data access, and form/action conventions adapted to the selected framework.\n")
		assetGuidance.WriteString("The applyTo glob MUST target server-side source files for the framework.\n\n")
	}
	if hasTesting {
		assetGuidance.WriteString("TESTING:\n")
		assetGuidance.WriteString("A testing asset is included. Generate a dedicated testing.instructions.md\n")
		assetGuidance.WriteString("with ONLY the framework-specific testing guidance (runner, file conventions,\n")
		assetGuidance.WriteString("setup/teardown, assertion style). Drop guidance for other frameworks.\n\n")
	}

	// Resolve the actual scaffold command with project name substituted.
	scaffoldResolved := strings.ReplaceAll(scaffoldInfo, "{{name}}", projectName)
	scaffoldResolved = strings.ReplaceAll(scaffoldResolved, "{{module}}", projectName)

	// Build profile file guidance.
	profileFileGlob := "**"
	switch sel.ProfileID {
	case "elixir-phoenix":
		profileFileGlob = "**/*.{ex,exs,heex,leex}"
	case "typescript-sveltekit", "typescript-nextjs", "typescript-fastify":
		profileFileGlob = "**/*.{ts,tsx,svelte,js,jsx}"
	case "ruby-rails":
		profileFileGlob = "**/*.{rb,erb,haml}"
	case "go-service":
		profileFileGlob = "**/*.go"
	case "rust-axum":
		profileFileGlob = "**/*.rs"
	case "dotnet-api":
		profileFileGlob = "**/*.{cs,csproj}"
	case "java-spring":
		profileFileGlob = "**/*.{java,kt}"
	case "python-fastapi", "python-django":
		profileFileGlob = "**/*.py"
	case "dart-flutter":
		profileFileGlob = "**/*.dart"
	case "laravel":
		profileFileGlob = "**/*.{php,blade.php}"
	}

	var uiGuidance string
	if isUIStack {
		uiGuidance = "UI STACK NOTE:\n" +
			"This is a UI framework. The copilot-instructions.md MUST mention the\n" +
			"styling system (e.g. Tailwind CSS) as part of the always-on standards.\n" +
			"A brief reference is sufficient — detailed tokens belong in design-system.instructions.md.\n\n"
	}

	prompt := fmt.Sprintf(
		"Generate AI instruction files for the project %q.\n\n"+
			"Selected: profile=%s | addons=%s | assets=%s\n\n"+
			"IMPORTANT — SCAFFOLD COMMAND:\n"+
			"The framework provides its own CLI scaffold command. The start.prompt.md MUST\n"+
			"use this command as step 1 instead of manually creating project boilerplate:\n"+
			"%s\n\n"+
			"The AI agent should NEVER generate framework boilerplate files (package.json,\n"+
			"mix.exs, Gemfile, etc.). The scaffold command handles all of that. The agent's\n"+
			"job is to write application code AFTER the scaffold is complete.\n\n"+
			"PROJECT NAME SUBSTITUTION:\n"+
			"The project name is %q. In all generated files, use the actual project name —\n"+
			"NEVER output template variables like {{name}} or {{module}}. For example,\n"+
			"write %q not {{name}} in scaffold commands and file references.\n\n"+
			"%s"+
			"%s"+
			"%s"+
			"ADAPTATION RULE:\n"+
			"All generated instruction files MUST use the selected framework's idioms.\n"+
			"Code examples, component patterns, styling approaches, and file globs must\n"+
			"match the framework. Do NOT emit patterns from a different ecosystem.\n\n"+
			"Use ONLY the asset content below as your source. Do not invent conventions.\n\n"+
			"%s\n"+
			"Output ONLY file blocks — no prose before or after:\n"+
			"===FILE: relative/path===\n(content)\n===END_FILE===\n\n"+
			"Required:\n"+
			"1. .github/copilot-instructions.md — always-on standards from core + profile assets\n"+
			"2. .github/instructions/<profile>.instructions.md — framework-specific conventions from the\n"+
			"   profile asset. YAML frontmatter with applyTo: %q to scope to framework source files.\n"+
			"   This MUST be a SEPARATE file from copilot-instructions.md.\n"+
			"3. .github/instructions/*.instructions.md — one per additional concern (architecture,\n"+
			"   design-system, frontend-craft, testing, server-patterns, etc.) with YAML frontmatter applyTo glob\n"+
			"4. AGENTS.md — multi-agent ground rules\n"+
			"5. .github/prompts/start.prompt.md — YAML frontmatter MUST be exactly:\n"+
			"   ---\n"+
			"   description: \"<one-sentence description>\"\n"+
			"   mode: agent\n"+
			"   tools: [\"terminal\", \"editFiles\", \"codebase\"]\n"+
			"   ---\n"+
			"   Do NOT invent tool names. The only valid tools are: terminal, editFiles,\n"+
			"   codebase, fetch. Use exactly these identifiers.\n"+
			"   Body MUST:\n"+
			"   a) Run the framework scaffold command first: %s\n"+
			"   b) Then proceed with application-specific implementation\n"+
			"   c) Never manually create files the scaffold already provides\n",
		projectName,
		sel.ProfileID,
		strings.Join(sel.AddonIDs, ", "),
		strings.Join(summary, ", "),
		scaffoldResolved,
		projectName,
		projectName,
		uiGuidance,
		designGuidance.String(),
		assetGuidance.String(),
		contextBlocks.String(),
		profileFileGlob,
		scaffoldResolved,
	)

	raw, err := e.provider.Send(ctx, prompt, "")
	if err != nil {
		return nil, err
	}
	files := parseFileOutput(raw)
	if len(files) == 0 {
		return nil, fmt.Errorf("model returned no file blocks")
	}
	return files, nil
}

// ParseSelection parses raw LLM JSON output into a normalized Selection.
// Exported for testing.
func ParseSelection(raw string) (*Selection, error) {
	return parseSelection(raw)
}

func parseSelection(raw string) (*Selection, error) {
	clean := strings.TrimSpace(raw)
	clean = strings.TrimPrefix(clean, "```json")
	clean = strings.TrimPrefix(clean, "```")
	clean = strings.TrimSuffix(clean, "```")
	clean = strings.TrimSpace(clean)
	if i := strings.Index(clean, "{"); i != -1 {
		if j := strings.LastIndex(clean, "}"); j > i {
			clean = clean[i : j+1]
		}
	}
	var sel Selection
	if err := json.Unmarshal([]byte(clean), &sel); err != nil {
		return nil, fmt.Errorf("parse selection: %w\nraw output: %s", err, raw)
	}
	sel.ProfileID = strings.TrimPrefix(strings.TrimSpace(sel.ProfileID), "profile.")

	normalizedAddons := make([]string, 0, len(sel.AddonIDs))
	seenAddons := make(map[string]bool)
	for _, addonID := range sel.AddonIDs {
		id := strings.TrimPrefix(strings.TrimSpace(addonID), "addon.")
		if id == "" || seenAddons[id] {
			continue
		}
		seenAddons[id] = true
		normalizedAddons = append(normalizedAddons, id)
	}
	sel.AddonIDs = normalizedAddons

	normalizedAssets := make([]string, 0, len(sel.AssetIDs))
	seenAssets := make(map[string]bool)
	for _, assetID := range sel.AssetIDs {
		id := strings.TrimSpace(assetID)
		if id == "" || strings.HasPrefix(id, "profile.") || strings.HasPrefix(id, "addon.") || seenAssets[id] {
			continue
		}
		seenAssets[id] = true
		normalizedAssets = append(normalizedAssets, id)
	}
	sel.AssetIDs = normalizedAssets

	return &sel, nil
}

// ParseFileOutput parses raw LLM output containing ===FILE: blocks.
// Exported for testing.
func ParseFileOutput(raw string) []FileOutput {
	return parseFileOutput(raw)
}

func parseFileOutput(raw string) []FileOutput {
	var files []FileOutput
	remaining := raw
	for {
		const startMark = "===FILE: "
		si := strings.Index(remaining, startMark)
		if si == -1 {
			break
		}
		after := remaining[si+len(startMark):]
		ep := strings.Index(after, "===")
		if ep == -1 {
			break
		}
		path := strings.TrimSpace(after[:ep])
		cs := si + len(startMark) + ep + 3
		const endMark = "===END_FILE==="
		ei := strings.Index(remaining[cs:], endMark)
		if ei == -1 {
			break
		}
		content := strings.TrimSpace(remaining[cs : cs+ei])
		files = append(files, FileOutput{Path: path, Content: content})
		remaining = remaining[cs+ei+len(endMark):]
	}
	return files
}

func catalogIDLines() string {
	return strings.Join(catalogSummaryLines(), "\n")
}

func conversationSystemPrompt() string {
	var sb strings.Builder

	// CONSTRAINTS FIRST — these override everything
	sb.WriteString("CONSTRAINTS — violating any of these is a failure:\n")
	sb.WriteString("1. NEVER write code, code blocks, folder structures, data models, or architecture.\n")
	sb.WriteString("2. NEVER use markdown headers (###) in replies.\n")
	sb.WriteString("3. ONLY recommend stacks from the catalog below. Express, Socket.IO, and React Native do not exist in the catalog.\n")
	sb.WriteString("4. NEVER skip Phase 1. Your first reply MUST be scope questions, not a recommendation.\n")
	sb.WriteString("5. ONE phase per reply. Never combine phases.\n")
	sb.WriteString("6. Maximum 6 sentences per reply.\n\n")

	sb.WriteString("WRONG OUTPUT (this is what failure looks like — never do this):\n")
	sb.WriteString("User: 'I want a real-time voting app'\n")
	sb.WriteString("BAD: '### Core Features\n1. Room creation...\n### Suggested Tech Stack\nReact + Express + Socket.IO...\n### Starter Template\n```/backend/index.js```'\n")
	sb.WriteString("This is wrong because it skips Phase 1, uses headers, writes code, and recommends stacks not in the catalog.\n\n")

	sb.WriteString("You are Launchpad, a stack advisor. You follow three phases in strict order.\n\n")

	// PHASE 1
	sb.WriteString("PHASE 1 — SCOPE (1-3 rounds, start here ALWAYS):\n")
	sb.WriteString("Ask 2-4 questions about features and behavior the user hasn't mentioned yet. Be specific to their project. Examples: Would you want a leaderboard? Multiple rounds? Countdown timer? Share via link or code? Should results persist after the session?\n")
	sb.WriteString("Do NOT mention any technology. Do NOT recommend a stack. Just explore the idea.\n")
	sb.WriteString("After enough rounds, summarize the features you've captured as a short numbered list and ask if anything is missing. Only move to Phase 2 after confirmation.\n\n")

	// PHASE 2
	sb.WriteString("PHASE 2 — OPTIONS (exactly 1 turn):\n")
	sb.WriteString("Present 2-3 stack options from the catalog. For each: name, one sentence why it fits, and the scaffold command. Mark your top pick with ★.\n")
	sb.WriteString("After presenting stacks, briefly mention relevant add-ons and design assets.\n")
	sb.WriteString("Note: for any stack with a UI surface, frontend-craft visual guidance and default palette/font assets are included automatically — no need for the user to opt in. You can mention this as a bonus.\n")
	sb.WriteString("For data-heavy projects, suggest the data-intensive add-on.\n")
	sb.WriteString("Ask which stack (and optionally which add-ons/assets) they want.\n\n")

	// PHASE 3
	sb.WriteString("PHASE 3 — COMMIT (exactly 1 turn):\n")
	sb.WriteString("Confirm their choice in one sentence. Emit READY_TO_GENERATE on its own line.\n\n")

	// DECISION MAP — derived from profile metadata
	sb.WriteString("DECISION MAP (★ = your top pick for that use case):\n")
	sb.WriteString("real-time/live/presence/chat/voting/collaborative -> ★ elixir-phoenix | typescript-sveltekit\n")
	sb.WriteString("full-stack JS web/SSR/content -> ★ typescript-sveltekit | typescript-nextjs\n")
	sb.WriteString("CRUD/MVP/admin/content platform -> ★ ruby-rails | python-django\n")
	sb.WriteString("React required/Vercel -> typescript-nextjs\n")
	sb.WriteString("Node.js API/microservice -> typescript-fastify\n")
	sb.WriteString("high-perf API/CLI/infra -> ★ go-service | rust-axum\n")
	sb.WriteString("enterprise API/C# -> dotnet-api\n")
	sb.WriteString("enterprise API/Java/JVM -> java-spring\n")
	sb.WriteString("Python API/ML/data -> python-fastapi\n")
	sb.WriteString("Python full-stack/admin/CMS -> python-django\n")
	sb.WriteString("native mobile -> dart-flutter\n")
	sb.WriteString("perf-critical systems -> ★ rust-axum | go-service\n")
	sb.WriteString("PHP -> laravel\n\n")

	// LAYER TAXONOMY — helps the model understand architectural roles
	sb.WriteString("LAYER TAXONOMY (how stacks map to architectural roles):\n")
	for _, p := range scaffold.Profiles {
		sb.WriteString(fmt.Sprintf("- %s: layer=%s", p.ID, p.Layer))
		if p.HasUI {
			sb.WriteString(" (has UI)")
		}
		sb.WriteByte('\n')
	}
	sb.WriteByte('\n')

	sb.WriteString("Catalog IDs (for extraction step):\n")
	for _, line := range catalogSummaryLines() {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// scaffoldCommandForProfile returns the CLI scaffold command for a given profile ID.
func scaffoldCommandForProfile(profileID string) string {
	if p := scaffold.FindProfile(profileID); p != nil && p.ScaffoldCmd != "" {
		return p.ScaffoldCmd
	}
	return "(no scaffold command defined)"
}
