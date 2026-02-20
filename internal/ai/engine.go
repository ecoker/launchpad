package ai

import (
"bytes"
"context"
"encoding/json"
"fmt"
"io"
"net/http"
"sort"
"strings"
"time"

"github.com/ehrencoker/agent-kit/templates"
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

const confidenceThreshold = 0.72

// ReadyToken is the phrase the model appends to signal readiness.
const ReadyToken = "READY_TO_GENERATE"

// Engine manages the OpenAI Responses API thread.
type Engine struct {
	apiKey             string
	httpClient         *http.Client
	model              string
	previousResponseID string
	onChunk            func(string)
}

// NewEngine creates a new Engine.
func NewEngine(apiKey string, onChunk func(string)) *Engine {
	return &Engine{
		apiKey:     strings.TrimSpace(apiKey),
		httpClient: &http.Client{Timeout: 120 * time.Second},
		model:      "gpt-4.1",
		onChunk:    onChunk,
	}
}

// Chat sends a user message on the persistent Responses API thread.
// The system prompt is injected only on the first call.
func (e *Engine) Chat(ctx context.Context, message string) (string, error) {
	if strings.TrimSpace(message) == "" {
		return "", fmt.Errorf("empty message")
	}
	// Always send instructions — the Responses API does NOT carry them
	// across previous_response_id chains.
	systemPrompt := conversationSystemPrompt()
	raw, id, err := e.call(ctx, message, systemPrompt, true)
	if err != nil {
		return "", err
	}
	e.previousResponseID = id
	return raw, nil
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
		"  \"profile_id\": \"<elixir-phoenix|typescript-sveltekit|ruby-rails|typescript-nextjs|typescript-fastify|go-service|dotnet-api|python-fastapi|python-django|dart-flutter|rust-axum|laravel>\",\n" +
		"  \"addon_ids\": [],\n" +
		"  \"asset_ids\": [],\n" +
		"  \"confidence\": 0.0,\n" +
		"  \"rationale\": \"one sentence\"\n" +
		"}\n\n" +
		"Asset IDs available:\n" + catalogIDLines()

	raw, id, err := e.call(ctx, extractPrompt, "", false)
	if err != nil {
		return nil, err
	}
	e.previousResponseID = id
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
"Use ONLY the asset content below as your source. Do not invent conventions.\n\n"+
"%s\n"+
"Output ONLY file blocks — no prose before or after:\n"+
"===FILE: relative/path===\n(content)\n===END_FILE===\n\n"+
"Required:\n"+
"1. .github/copilot-instructions.md — always-on standards from core + profile assets\n"+
"2. .github/instructions/*.instructions.md — one per concern, YAML frontmatter applyTo glob required\n"+
"3. AGENTS.md — multi-agent ground rules\n"+
"4. .github/prompts/start.prompt.md — YAML frontmatter with description, agent: \"agent\", tools list.\n"+
"   Body MUST:\n"+
"   a) Run the framework scaffold command first: %s\n"+
"   b) Then proceed with application-specific implementation\n"+
"   c) Never manually create files the scaffold already provides\n",
projectName,
sel.ProfileID,
strings.Join(sel.AddonIDs, ", "),
strings.Join(summary, ", "),
scaffoldInfo,
contextBlocks.String(),
scaffoldInfo,
	)

	raw, _, err := e.call(ctx, prompt, "", false)
	if err != nil {
		return nil, err
	}
	files := parseFileOutput(raw)
	if len(files) == 0 {
		return nil, fmt.Errorf("model returned no file blocks")
	}
	return files, nil
}

// call makes a single Responses API request. systemPrompt is sent only when non-empty.
func (e *Engine) call(ctx context.Context, input, systemPrompt string, stream bool) (string, string, error) {
	type reqBody struct {
		Model              string `json:"model"`
		Instructions       string `json:"instructions,omitempty"`
		PreviousResponseID string `json:"previous_response_id,omitempty"`
		Input              string `json:"input"`
	}
	body := reqBody{
		Model:              e.model,
		Input:              input,
		PreviousResponseID: e.previousResponseID,
		Instructions:       systemPrompt,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", "", fmt.Errorf("marshal: %w", err)
	}

	for attempt := 1; attempt <= 3; attempt++ {
		req, reqErr := http.NewRequestWithContext(
ctx, http.MethodPost,
"https://api.openai.com/v1/responses",
bytes.NewReader(payload),
)
		if reqErr != nil {
			return "", "", fmt.Errorf("build request: %w", reqErr)
		}
		req.Header.Set("Authorization", "Bearer "+e.apiKey)
		req.Header.Set("Content-Type", "application/json")

		res, doErr := e.httpClient.Do(req)
		if doErr != nil {
			return "", "", fmt.Errorf("http: %w", doErr)
		}
		respBytes, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", "", fmt.Errorf("read body: %w", readErr)
		}

		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return "", "", fmt.Errorf("api %d: %s", res.StatusCode, strings.TrimSpace(string(respBytes)))
		}

		var out responsesAPIResponse
		if jsonErr := json.Unmarshal(respBytes, &out); jsonErr != nil {
			return "", "", fmt.Errorf("decode: %w", jsonErr)
		}
		text := out.text()
		if text == "" {
			return "", "", fmt.Errorf("empty response from api")
		}
		if stream && e.onChunk != nil {
			e.onChunk(text)
		}
		return text, out.ID, nil
	}
	return "", "", fmt.Errorf("rate limit exceeded after retries")
}

type responsesAPIResponse struct {
	ID     string `json:"id"`
	Output []struct {
		Content []struct {
			Type string `json:"type"`
			Text string `json:"text"`
		} `json:"content"`
	} `json:"output"`
	OutputText string `json:"output_text"`
}

func (r responsesAPIResponse) text() string {
	if t := strings.TrimSpace(r.OutputText); t != "" {
		return t
	}
	var sb strings.Builder
	for _, o := range r.Output {
		for _, c := range o.Content {
			if t := strings.TrimSpace(c.Text); t != "" {
				if sb.Len() > 0 {
					sb.WriteByte('\n')
				}
				sb.WriteString(t)
			}
		}
	}
	return strings.TrimSpace(sb.String())
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
	sb.WriteString("Present 2-3 stack options from the catalog. For each: name, one sentence why it fits, and the scaffold command. Mark your top pick with ★. Ask which they want.\n\n")

	// PHASE 3
	sb.WriteString("PHASE 3 — COMMIT (exactly 1 turn):\n")
	sb.WriteString("Confirm their choice in one sentence. Emit READY_TO_GENERATE on its own line.\n\n")

	// DECISION MAP
	sb.WriteString("DECISION MAP (★ = your top pick for that use case):\n")
	sb.WriteString("real-time/live/presence/chat/voting/collaborative -> ★ elixir-phoenix | typescript-sveltekit\n")
	sb.WriteString("full-stack JS web/SSR/content -> ★ typescript-sveltekit | typescript-nextjs\n")
	sb.WriteString("CRUD/MVP/admin/content platform -> ★ ruby-rails | python-django\n")
	sb.WriteString("React required/Vercel -> typescript-nextjs\n")
	sb.WriteString("Node.js API/microservice -> typescript-fastify\n")
	sb.WriteString("high-perf API/CLI/infra -> go-service\n")
	sb.WriteString("enterprise API/C# -> dotnet-api\n")
	sb.WriteString("Python API/ML/data -> python-fastapi\n")
	sb.WriteString("Python full-stack/admin/CMS -> python-django\n")
	sb.WriteString("native mobile -> dart-flutter\n")
	sb.WriteString("perf-critical systems -> rust-axum\n")
	sb.WriteString("PHP -> laravel\n\n")

	sb.WriteString("Catalog IDs (for extraction step):\n")
	for _, line := range catalogSummaryLines() {
		sb.WriteString(line)
		sb.WriteByte('\n')
	}

	return sb.String()
}

// scaffoldCommandForProfile returns the CLI scaffold command for a given profile ID.
func scaffoldCommandForProfile(profileID string) string {
	commands := map[string]string{
		"elixir-phoenix":      "mix phx.new {{name}}",
		"typescript-sveltekit": "npm create svelte@latest",
		"ruby-rails":          "rails new {{name}}",
		"typescript-nextjs":   "npx create-next-app@latest",
		"typescript-fastify":  "npm init -y && npm install fastify",
		"go-service":          "go mod init {{module}}",
		"dotnet-api":          "dotnet new webapi -n {{name}}",
		"python-fastapi":      "mkdir {{name}} && cd {{name}} && python -m venv .venv && pip install fastapi uvicorn",
		"python-django":       "django-admin startproject {{name}}",
		"dart-flutter":        "flutter create {{name}}",
		"rust-axum":           "cargo new {{name}}",
		"laravel":             "composer create-project laravel/laravel {{name}}",
	}
	if cmd, ok := commands[profileID]; ok {
		return cmd
	}
	return "(no scaffold command defined)"
}
