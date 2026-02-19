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
	systemPrompt := ""
	if e.previousResponseID == "" {
		systemPrompt = conversationSystemPrompt()
	}
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
		"  \"profile_id\": \"<typescript-react|python-data|elixir-phoenix|dotnet-api|laravel|go-service>\",\n" +
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

	prompt := fmt.Sprintf(
"Generate AI instruction files for the project %q.\n\n"+
"Selected: profile=%s | addons=%s | assets=%s\n\n"+
"Use ONLY the asset content below as your source. Do not invent conventions.\n\n"+
"%s\n"+
"Output ONLY file blocks — no prose before or after:\n"+
"===FILE: relative/path===\n(content)\n===END_FILE===\n\n"+
"Required:\n"+
"1. .github/copilot-instructions.md — always-on standards from core + profile assets\n"+
"2. .github/instructions/*.instructions.md — one per concern, YAML frontmatter applyTo glob required\n"+
"3. AGENTS.md — multi-agent ground rules\n"+
"4. .github/prompts/start.prompt.md — YAML frontmatter: description, agent: \"agent\", tools list; body bootstraps the project\n",
projectName,
sel.ProfileID,
strings.Join(sel.AddonIDs, ", "),
strings.Join(summary, ", "),
contextBlocks.String(),
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
sb.WriteString("You are a senior engineering advisor embedded in Launchpad, a project scaffolding tool.\n\n")
sb.WriteString("Your job: help the developer choose the right tech stack and coding standards for their project.\n\n")
sb.WriteString("How to work:\n")
sb.WriteString("- Listen carefully to what they describe\n")
sb.WriteString("- Ask targeted questions to resolve ambiguity: team size, deploy target, scale, real-time needs\n")
sb.WriteString("- When alternatives are given, explain tradeoffs and give a clear recommendation\n")
sb.WriteString("- Keep responses brief and decision-focused\n")
sb.WriteString("- Do not provide implementation code, file contents, or build tutorials\n")
sb.WriteString("- Once you have a clear decision, explain it in one or two sentences and put READY_TO_GENERATE as the final line\n\n")
sb.WriteString("Rules:\n")
sb.WriteString("- Never ask about coding style or philosophy — Launchpad already has strong opinions on those\n")
sb.WriteString("- Never recommend a stack outside the catalog below\n")
sb.WriteString("- Be direct and opinionated — you are a partner, not a form\n\n")
sb.WriteString("Available stacks and assets:\n")
for _, line := range catalogSummaryLines() {
sb.WriteString(line)
sb.WriteByte('\n')
}
return sb.String()
}
