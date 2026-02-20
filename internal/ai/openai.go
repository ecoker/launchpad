package ai

import (
	"bytes"
	"context"
	"encoding/json"
	"fmt"
	"io"
	"net/http"
	"strings"
	"time"
)

const (
	openAIResponsesURL = "https://api.openai.com/v1/responses"
	defaultModel       = "gpt-4.1"
)

// OpenAIProvider implements Provider using the OpenAI Responses API.
type OpenAIProvider struct {
	apiKey             string
	model              string
	httpClient         *http.Client
	previousResponseID string
}

// OpenAIOption configures an OpenAIProvider.
type OpenAIOption func(*OpenAIProvider)

// WithModel overrides the default model.
func WithModel(model string) OpenAIOption {
	return func(p *OpenAIProvider) {
		if model != "" {
			p.model = model
		}
	}
}

// WithHTTPClient overrides the default HTTP client.
func WithHTTPClient(c *http.Client) OpenAIOption {
	return func(p *OpenAIProvider) {
		p.httpClient = c
	}
}

// NewOpenAIProvider creates a provider backed by the OpenAI Responses API.
func NewOpenAIProvider(apiKey string, opts ...OpenAIOption) *OpenAIProvider {
	p := &OpenAIProvider{
		apiKey:     strings.TrimSpace(apiKey),
		model:      defaultModel,
		httpClient: &http.Client{Timeout: 120 * time.Second},
	}
	for _, o := range opts {
		o(p)
	}
	return p
}

// Send implements Provider.
func (p *OpenAIProvider) Send(ctx context.Context, message, systemPrompt string) (string, error) {
	type reqBody struct {
		Model              string `json:"model"`
		Instructions       string `json:"instructions,omitempty"`
		PreviousResponseID string `json:"previous_response_id,omitempty"`
		Input              string `json:"input"`
	}
	body := reqBody{
		Model:              p.model,
		Input:              message,
		PreviousResponseID: p.previousResponseID,
		Instructions:       systemPrompt,
	}
	payload, err := json.Marshal(body)
	if err != nil {
		return "", fmt.Errorf("marshal: %w", err)
	}

	for attempt := 1; attempt <= 3; attempt++ {
		req, reqErr := http.NewRequestWithContext(
			ctx, http.MethodPost, openAIResponsesURL, bytes.NewReader(payload),
		)
		if reqErr != nil {
			return "", fmt.Errorf("build request: %w", reqErr)
		}
		req.Header.Set("Authorization", "Bearer "+p.apiKey)
		req.Header.Set("Content-Type", "application/json")

		res, doErr := p.httpClient.Do(req)
		if doErr != nil {
			return "", fmt.Errorf("http: %w", doErr)
		}
		respBytes, readErr := io.ReadAll(res.Body)
		res.Body.Close()
		if readErr != nil {
			return "", fmt.Errorf("read body: %w", readErr)
		}

		if res.StatusCode == http.StatusTooManyRequests {
			time.Sleep(time.Duration(attempt) * 2 * time.Second)
			continue
		}
		if res.StatusCode < 200 || res.StatusCode >= 300 {
			return "", fmt.Errorf(
				"OpenAI API error (HTTP %d) — check your API key and account status",
				res.StatusCode,
			)
		}

		var out responsesAPIResponse
		if jsonErr := json.Unmarshal(respBytes, &out); jsonErr != nil {
			return "", fmt.Errorf("decode response: %w", jsonErr)
		}
		text := out.text()
		if text == "" {
			return "", fmt.Errorf("empty response from API — try again or check your input")
		}
		p.previousResponseID = out.ID
		return text, nil
	}
	return "", fmt.Errorf("rate limited after 3 retries — wait a moment and try again")
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
